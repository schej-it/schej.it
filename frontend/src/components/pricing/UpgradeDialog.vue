<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card class="tw-rounded-lg tw-p-6 sm:tw-p-8">
      <div class="tw-mb-8 tw-flex tw-flex-col tw-items-start tw-gap-4">
        <h2 class="tw-text-xl tw-font-medium">
          Upgrade to <span class="tw-text-green-600">Schej Premium</span>
        </h2>
        <div class="tw-mb-4 tw-text-sm tw-font-medium tw-text-dark-gray">
          You've run out of
          <v-tooltip
            top
            content-class="tw-bg-very-dark-gray tw-shadow-lg tw-opacity-100 tw-py-4"
          >
            <template v-slot:activator="{ on, attrs }">
              <span v-bind="attrs" v-on="on">
                <span class="tw-cursor-pointer tw-underline">free events</span>.
              </span>
            </template>
            <div>
              Free users can have up to {{ numFreeEvents }} events associated
              <br />
              with their account at a time.
            </div>
          </v-tooltip>
          Delete old events or upgrade to create unlimited events.
        </div>
        <div
          class="tw-inline-block tw-rounded tw-bg-light-green/10 tw-px-2 tw-py-1 tw-text-sm tw-font-medium tw-text-light-green"
        >
          One-time payment
        </div>
        <div class="tw-font-medium">
          <span class="tw-mr-1 tw-text-4xl">{{ formattedPrice }}</span>
          <span class="tw-text-base">for lifetime access</span>
        </div>
        <ul
          class="tw-inline-block tw-space-y-0.5 tw-p-0 tw-text-sm tw-font-medium tw-text-very-dark-gray"
        >
          <li class="tw-flex tw-items-center">
            <v-icon class="tw-mr-2 tw-text-light-green">mdi-check</v-icon>
            <span>Unlimited events</span>
          </li>
          <li class="tw-flex tw-items-center">
            <v-icon class="tw-mr-2 tw-text-light-green">mdi-check</v-icon>
            <span>Unlimited availability groups</span>
          </li>
          <li class="tw-flex tw-items-center">
            <v-icon class="tw-mr-2 tw-text-light-green">mdi-check</v-icon>
            <span>Any new premium features we add</span>
          </li>
        </ul>
      </div>
      <v-btn
        class="tw-mb-0.5"
        color="primary"
        dark
        block
        @click="handleUpgrade"
      >
        Upgrade
      </v-btn>
      <v-btn text block @click="$emit('input', false)"> Cancel </v-btn>
    </v-card>
  </v-dialog>
</template>

<script>
import { get, post } from "@/utils"
import { mapState } from "vuex"
import { numFreeEvents } from "@/constants"

export default {
  name: "UpgradeDialog",

  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      loaded: false,
      price: null,
      checkoutUrl: null,
    }
  },

  computed: {
    ...mapState(["featureFlagsLoaded", "pricingPageConversion", "authUser"]),
    formattedPrice() {
      if (!this.price) return ""
      return "$" + (this.price.unit_amount / 100).toFixed(2)
    },
    numFreeEvents() {
      return numFreeEvents
    },
  },

  methods: {
    async fetchPrice() {
      const res = await get("/stripe/price?exp=" + this.pricingPageConversion)
      this.price = res.price
      await this.createCheckoutSession()
      this.loaded = true
    },
    handleUpgrade() {
      window.location.href = this.checkoutUrl
    },
    async createCheckoutSession() {
      const res = await post("/stripe/create-checkout-session", {
        priceId: this.price.id,
        userId: this.authUser._id,
        originUrl: window.location.href,
      })
      this.checkoutUrl = res.url
    },
  },

  watch: {
    featureFlagsLoaded: {
      handler(newVal) {
        if (newVal) {
          this.fetchPrice()
        }
      },
      immediate: true,
    },
  },
}
</script>
