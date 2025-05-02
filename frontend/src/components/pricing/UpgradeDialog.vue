<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="600"
    content-class="tw-m-0"
  >
    <v-card class="tw-rounded-lg tw-p-4 tw-pb-2 sm:tw-p-8 sm:tw-pb-4">
      <div class="tw-mb-4 tw-flex tw-flex-col tw-items-start tw-gap-4">
        <h2 class="tw-text-xl tw-font-medium">
          Upgrade to
          <span
            class="tw-bg-gradient-to-r tw-from-darkest-green tw-to-light-green tw-bg-clip-text tw-text-transparent"
            >Schej Premium</span
          >
        </h2>
        <div class="tw-text-sm tw-font-medium tw-text-dark-gray">
          You've run out of free events. Upgrade to create unlimited events.
        </div>
        <!-- <ul
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
        </ul> -->
      </div>
      <div
        class="tw-mb-8 tw-flex tw-flex-col tw-gap-1 sm:tw-flex-row sm:tw-gap-4"
      >
        <div
          class="tw-flex tw-flex-1 tw-flex-col tw-items-center tw-gap-2 tw-rounded-lg tw-border tw-border-light-green/20 tw-p-4"
        >
          <div
            class="tw-inline-block tw-w-fit tw-rounded tw-px-2 tw-py-1 tw-text-sm tw-font-medium"
          >
            1-month pass
          </div>
          <div class="tw-relative">
            <div class="tw-font-medium">
              <span class="tw-mr-1 tw-text-dark-gray tw-line-through">$15</span>
              <span class="tw-mr-1 tw-text-4xl">{{
                formattedPrice(oneMonthPrice)
              }}</span>
              <span class="tw-text-dark-gray">USD</span>
            </div>
            <v-fade-transition>
              <div
                v-if="oneMonthPrice === null"
                class="tw-absolute tw-left-0 tw-top-0 tw-h-full tw-w-full tw-bg-white"
              ></div>
            </v-fade-transition>
          </div>
          <div class="tw-mb-4 tw-text-center tw-text-sm tw-text-very-dark-gray">
            One-time payment.<br />No subscription.
          </div>
          <v-btn
            class="tw-mb-0.5"
            color="primary"
            outlined
            block
            :dark="!loadingCheckoutUrl[oneMonthPrice?.id]"
            :disabled="loadingCheckoutUrl[oneMonthPrice?.id]"
            :loading="loadingCheckoutUrl[oneMonthPrice?.id]"
            @click="handleUpgrade(oneMonthPrice)"
          >
            Upgrade
          </v-btn>
        </div>
        <div
          class="tw-flex tw-flex-1 tw-flex-col tw-items-center tw-gap-2 tw-rounded-lg tw-border tw-border-light-green/20 tw-bg-white tw-p-4 tw-shadow-lg"
        >
          <div
            class="tw-inline-block tw-w-fit tw-rounded tw-bg-light-green/10 tw-px-2 tw-py-1 tw-text-sm tw-font-medium tw-text-light-green"
          >
            Lifetime access
          </div>
          <div class="tw-relative">
            <div class="tw-font-medium">
              <span class="tw-mr-1 tw-text-dark-gray tw-line-through">$40</span>
              <span class="tw-mr-1 tw-text-4xl">{{
                formattedPrice(lifetimePrice)
              }}</span>
              <span class="tw-text-dark-gray">USD</span>
            </div>
            <v-fade-transition>
              <div
                v-if="lifetimePrice === null"
                class="tw-absolute tw-left-0 tw-top-0 tw-h-full tw-w-full tw-bg-white"
              ></div>
            </v-fade-transition>
          </div>
          <div class="tw-mb-4 tw-text-center tw-text-sm tw-text-very-dark-gray">
            One-time payment.<br />No subscription.
          </div>
          <v-btn
            class="tw-mb-0.5"
            color="primary"
            block
            :dark="!loadingCheckoutUrl[lifetimePrice?.id]"
            :disabled="loadingCheckoutUrl[lifetimePrice?.id]"
            :loading="loadingCheckoutUrl[lifetimePrice?.id]"
            @click="handleUpgrade(lifetimePrice)"
          >
            Upgrade
          </v-btn>
        </div>
      </div>
      <div class="tw-mb-2 tw-text-center">
        <div
          class="tw-cursor-pointer tw-text-xs tw-font-medium tw-text-dark-gray hover:tw-underline"
          @click="showDonatedDialog = true"
        >
          I already donated :)
        </div>
      </div>
      <v-btn text block @click="$emit('input', false)"> Cancel </v-btn>
    </v-card>

    <AlreadyDonatedDialog v-model="showDonatedDialog" />
  </v-dialog>
</template>

<script>
import { get, post } from "@/utils"
import { mapState, mapActions } from "vuex"
import AlreadyDonatedDialog from "./AlreadyDonatedDialog.vue"

export default {
  name: "UpgradeDialog",
  components: {
    AlreadyDonatedDialog,
  },
  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      oneMonthPrice: null,
      lifetimePrice: null,
      loadingCheckoutUrl: {},
      showDonatedDialog: false,
    }
  },

  computed: {
    ...mapState(["featureFlagsLoaded", "pricingPageConversion", "authUser"]),
  },

  methods: {
    ...mapActions(["showError"]),
    formattedPrice(price) {
      if (!price) return ""
      return "$" + Math.floor(price.unit_amount / 100)
    },
    async init() {
      if (this.featureFlagsLoaded && this.authUser) {
        if (!this.oneMonthPrice || !this.lifetimePrice) {
          await this.fetchPrice()
        }
      }
    },
    async fetchPrice() {
      const res = await get("/stripe/price?exp=" + this.pricingPageConversion)
      const { oneMonth, lifetime } = res
      this.oneMonthPrice = oneMonth
      this.lifetimePrice = lifetime
    },
    async handleUpgrade(price) {
      this.$posthog.capture("upgrade_clicked", {
        price: this.formattedPrice(price),
      })
      this.$set(this.loadingCheckoutUrl, price.id, true)
      try {
        const res = await post("/stripe/create-checkout-session", {
          priceId: price.id,
          userId: this.authUser._id,
          originUrl: window.location.href,
        })
        window.location.href = res.url
      } catch (e) {
        console.error(e)
        this.showError(
          "There was an error generating a checkout url. Please try again later."
        )
      } finally {
        this.$set(this.loadingCheckoutUrl, price.id, false)
      }
    },
  },

  watch: {
    featureFlagsLoaded: {
      handler() {
        this.init()
      },
      immediate: true,
    },
    authUser: {
      handler() {
        this.init()
      },
    },
    value: {
      handler() {
        if (this.value) {
          post("/analytics/upgrade-dialog-viewed", {
            userId: this.authUser._id,
            price: `${this.formattedPrice(
              this.oneMonthPrice
            )}, ${this.formattedPrice(this.lifetimePrice)}`,
          })
          this.$posthog.capture("upgrade_dialog_viewed", {
            price: `${this.formattedPrice(
              this.oneMonthPrice
            )}, ${this.formattedPrice(this.lifetimePrice)}`,
          })
        }
      },
    },
  },
}
</script>
