package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/price"
	"github.com/stripe/stripe-go/v82/webhook"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"schej.it/server/db"
	"schej.it/server/logger"
	"schej.it/server/slackbot"
	"schej.it/server/utils"
)

func InitStripe(router *gin.RouterGroup) {
	stripeRouter := router.Group("/stripe")

	stripeRouter.POST("/create-checkout-session", createCheckoutSession)
	stripeRouter.GET("/price", getPrice)
	stripeRouter.POST("/fulfill-checkout", fulfillCheckout)
	stripeRouter.POST("/webhook", stripeWebhook)
}

type CheckoutSessionPayload struct {
	PriceID   string `json:"priceId" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
	OriginURL string `json:"originUrl" binding:"required"`
}

func createCheckoutSession(c *gin.Context) {
	var payload CheckoutSessionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	originURL := payload.OriginURL
	finalRedirectURL := originURL // This is where the user should end up AFTER the /stripe-redirect page

	// Get the base URL for constructing the intermediate redirect path
	baseURL := utils.GetBaseUrl()
	intermediateRedirectBase, err := url.Parse(baseURL)
	if err != nil {
		logger.StdErr.Printf("Error parsing Base URL '%s': %v. Cannot construct redirect URLs.", baseURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error configuring redirect"})
		return
	}

	// Create success URL (points to /stripe-redirect)
	successURL := *intermediateRedirectBase // Start with base URL
	successURL.Path = "/stripe-redirect"    // Set path
	successQuery := url.Values{}
	successQuery.Set("upgrade", "success")
	successQuery.Set("redirect_url", finalRedirectURL) // Add the final destination
	successURL.RawQuery = successQuery.Encode()
	successURLStr := successURL.String()

	// Create cancel URL (points to /stripe-redirect)
	cancelURL := *intermediateRedirectBase // Start with base URL
	cancelURL.Path = "/stripe-redirect"    // Set path
	cancelQuery := url.Values{}
	cancelQuery.Set("upgrade", "cancel")
	cancelQuery.Set("redirect_url", finalRedirectURL) // Add the final destination
	cancelURL.RawQuery = cancelQuery.Encode()
	cancelURLStr := cancelURL.String()

	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(payload.UserID),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, price_1234) of the product you want to sell
				Price:    stripe.String(payload.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:             stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:       stripe.String(successURLStr + "&session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:        stripe.String(cancelURLStr),
		AutomaticTax:     &stripe.CheckoutSessionAutomaticTaxParams{Enabled: stripe.Bool(true)},
		CustomerCreation: stripe.String(string(stripe.CheckoutSessionCustomerCreationAlways)),
		// Provide the Customer ID (for example, cus_1234) for an existing customer to associate it with this session
		// Customer: "cus_RnhPlBnbBbXapY",
	}

	s, err := session.New(params)

	if err != nil {
		log.Printf("session.New: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create checkout session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}

func getPrice(c *gin.Context) {
	// Get the experiment query parameter
	// exp := c.Query("exp")
	oneMonthPriceId := os.Getenv("STRIPE_ONE_MONTH_PRICE_ID")
	lifetimePriceId := os.Getenv("STRIPE_LIFETIME_PRICE_ID")

	// switch exp {
	// case "test2":
	// 	priceId = os.Getenv("STRIPE_PRICE_ID_2")
	// case "test3":
	// 	priceId = os.Getenv("STRIPE_PRICE_ID_3")
	// default:
	// 	priceId = os.Getenv("STRIPE_PRICE_ID_1")
	// }

	params := &stripe.PriceParams{}
	oneMonthResult, err := price.Get(oneMonthPriceId, params)
	if err != nil {
		log.Printf("price.Get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price"})
		return
	}

	lifetimeResult, err := price.Get(lifetimePriceId, params)
	if err != nil {
		log.Printf("price.Get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"oneMonth": oneMonthResult, "lifetime": lifetimeResult})
}

type FulfillCheckoutPayload struct {
	SessionID string `json:"sessionId" binding:"required"`
}

func fulfillCheckout(c *gin.Context) {
	var payload FulfillCheckoutPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	_fulfillCheckout(payload.SessionID)
}

func _fulfillCheckout(sessionId string) {
	// TODO: Make this function safe to run multiple times,
	// even concurrently, with the same session ID

	// TODO: Make sure fulfillment hasn't already been
	// performed for this Checkout Session

	// Retrieve the Checkout Session from the API with line_items expanded
	params := &stripe.CheckoutSessionParams{}
	params.AddExpand("line_items")

	cs, _ := session.Get(sessionId, params)

	// Check the Checkout Session's payment_status property
	// to determine if fulfillment should be performed
	if cs.PaymentStatus != stripe.CheckoutSessionPaymentStatusUnpaid {
		logger.StdOut.Println("Fulfilling Checkout Session " + sessionId)
		if cs.Customer != nil {
			logger.StdOut.Println("Setting stripe customer ID", cs.Customer.ID)

			// Fetch user from database
			userId := cs.ClientReferenceID
			userIdObj, err := primitive.ObjectIDFromHex(userId)
			if err != nil {
				logger.StdErr.Printf("Error parsing user ID: %v", err)
				return
			}
			user := db.GetUserById(userId)
			if user == nil {
				logger.StdErr.Printf("Error getting user: %v", err)
				return
			}

			// Only upgrade the user if customer ID is different
			if user.StripeCustomerId == nil || *user.StripeCustomerId != cs.Customer.ID {
				var planExpiration primitive.DateTime
				if cs.LineItems != nil && len(cs.LineItems.Data) > 0 {
					priceId := cs.LineItems.Data[0].Price.ID
					priceDescription := ""
					if priceId == os.Getenv("STRIPE_ONE_MONTH_PRICE_ID") {
						priceDescription = "1-month"
						planExpiration = primitive.NewDateTimeFromTime(time.Now().AddDate(0, 1, 0))
					} else if priceId == os.Getenv("STRIPE_LIFETIME_PRICE_ID") {
						priceDescription = "lifetime"
						planExpiration = primitive.NewDateTimeFromTime(time.Now().AddDate(999, 0, 0))
					}
					amountTotal := float32(cs.LineItems.Data[0].AmountTotal) / 100.0

					message := fmt.Sprintf(":moneybag: %s %s (%s) paid for Schej ($%.2f, %s) :moneybag:", user.FirstName, user.LastName, user.Email, amountTotal, priceDescription)
					slackbot.SendTextMessageWithType(message, slackbot.MONETIZATION)
				}

				user.PlanExpiration = &planExpiration
				user.StripeCustomerId = &cs.Customer.ID
				db.UsersCollection.UpdateOne(context.Background(), bson.M{"_id": userIdObj}, bson.M{"$set": user})
			}
		}
	}
}

func stripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.StdErr.Printf("Error reading request body: %v", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	// Pass the request body and Stripe-Signature header to ConstructEvent, along with the webhook signing key.
	// Use the secret provided by your webhook endpoint settings or Stripe CLI.
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		logger.StdErr.Println("STRIPE_WEBHOOK_SECRET not set")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), endpointSecret)

	if err != nil {
		logger.StdErr.Printf("Error verifying webhook signature: %v", err)
		c.AbortWithStatus(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Handle the event
	if event.Type == stripe.EventTypeCheckoutSessionCompleted || event.Type == stripe.EventTypeCheckoutSessionAsyncPaymentSucceeded {
		var cs stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &cs)
		if err != nil {
			logger.StdErr.Printf("Error parsing webhook JSON: %v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		logger.StdOut.Printf("Checkout Session %s completed!\n", cs.ID)
		_fulfillCheckout(cs.ID) // Call fulfillCheckout when session is completed
	}

	c.Status(http.StatusOK) // Return 200 OK to acknowledge receipt of the event
}
