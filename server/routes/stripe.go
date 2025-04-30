package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/price"
	"schej.it/server/utils"
)

func InitStripe(router *gin.RouterGroup) {
	stripeRouter := router.Group("/stripe")

	stripeRouter.POST("/create-checkout-session", createCheckoutSession)
	stripeRouter.GET("/price", getPrice)
}

type CheckoutSessionPayload struct {
	PriceID string `json:"priceId" binding:"required"`
}

func createCheckoutSession(c *gin.Context) {
	fmt.Println("Creating checkout session")

	var payload CheckoutSessionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	domain := utils.GetBaseUrl()
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				// Provide the exact Price ID (for example, price_1234) of the product you want to sell
				Price:    stripe.String(payload.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:             stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:       stripe.String(domain + "/upgrade/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:        stripe.String(domain + "/upgrade?canceled=true"),
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
	fmt.Println("Fetching products and prices")

	// Get the experiment query parameter
	exp := c.Query("exp")
	priceId := ""

	switch exp {
	case "test2":
		priceId = os.Getenv("STRIPE_PRICE_ID_2")
	case "test3":
		priceId = os.Getenv("STRIPE_PRICE_ID_3")
	default:
		priceId = os.Getenv("STRIPE_PRICE_ID_1")
	}

	params := &stripe.PriceParams{}
	result, err := price.Get(priceId, params)
	if err != nil {
		log.Printf("price.Get error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": result})
}
