<template>
  <div class="tw-flex tw-h-screen tw-items-center tw-justify-center tw-p-4">
    <div class="tw-text-center">
      <v-progress-circular
        indeterminate
        color="primary"
        size="32"
      ></v-progress-circular>
      <!-- <p v-if="message" class="tw-mt-4 tw-text-lg">{{ message }}</p>
      <p v-if="error" class="tw-text-red-600 tw-mt-4 tw-text-lg">{{ error }}</p>
      <p class="tw-text-gray-600 tw-mt-2 tw-text-sm">Redirecting shortly...</p> -->
    </div>
  </div>
</template>

<script>
import { get, post } from "@/utils"
import { mapMutations } from "vuex"

export default {
  name: "StripeRedirect",
  data() {
    return {
      message: "Processing your payment information...",
      error: null,
    }
  },
  methods: {
    ...mapMutations(["setAuthUser"]), // Assuming setAuthUser might be needed if fulfillment updates user state
    async handleRedirect() {
      const urlParams = new URLSearchParams(window.location.search)
      const upgradeStatus = urlParams.get("upgrade")
      const sessionId = urlParams.get("session_id")
      const redirectUrl = urlParams.get("redirect_url")

      if (!redirectUrl) {
        this.error = "Missing redirect information. Cannot proceed."
        // Consider redirecting to a default safe page like home after a delay
        setTimeout(() => {
          this.$router.replace({ name: "home" })
        }, 3000)
        return
      }

      try {
        if (upgradeStatus === "success" && sessionId) {
          this.message = "Payment successful! Finalizing your upgrade..."
          await post("/stripe/fulfill-checkout", { sessionId })
          this.message = "Upgrade complete! Redirecting..."
          // Optionally, refresh user data if fulfillment changes it
          const user = await get("/user/profile")
          this.setAuthUser(user)
        } else if (upgradeStatus === "cancel") {
          this.message = "Upgrade cancelled. Redirecting..."
        } else {
          // If neither success nor cancel, maybe it's an unexpected state
          this.error = "Invalid status received. Redirecting..."
        }
      } catch (err) {
        console.error("Error during Stripe redirect handling:", err)
        // this.error =
        //   "An error occurred while processing your upgrade. Please contact support if the problem persists."
        // Continue to redirect even on error
      } finally {
        // Redirect after a short delay to allow user to see message
        window.location.replace(redirectUrl)
      }
    },
  },
  mounted() {
    this.handleRedirect()
  },
}
</script>
