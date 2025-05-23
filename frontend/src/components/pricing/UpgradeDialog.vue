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
          <template
            v-if="upgradeDialogType === upgradeDialogTypes.CREATE_EVENT"
          >
            You've run out of free events. Upgrade to create unlimited events.
            <br class="tw-hidden sm:tw-block" />
            Your payment helps us keep the site running.
          </template>
          <template
            v-else-if="upgradeDialogType === upgradeDialogTypes.SCHEDULE_EVENT"
          >
            Upgrade to schedule events with Schej. Your payment helps us keep
            the site running.
          </template>
          <template
            v-else-if="
              upgradeDialogType === upgradeDialogTypes.UPGRADE_MANUALLY
            "
          >
            Create unlimited events with Schej Premium. Your payment helps us
            keep the site running.
          </template>
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
            Monthly
          </div>
          <div class="tw-relative">
            <div class="tw-font-medium">
              <span class="tw-mr-1 tw-text-dark-gray tw-line-through">$15</span>
              <span class="tw-mr-1 tw-text-4xl">{{
                isStudent
                  ? formattedPrice(monthlyStudentPrice)
                  : formattedPrice(monthlyPrice)
              }}</span>
              <span class="tw-text-dark-gray">USD</span>
            </div>
            <v-fade-transition>
              <div
                v-if="monthlyPrice === null"
                class="tw-absolute tw-left-0 tw-top-0 tw-h-full tw-w-full tw-bg-white"
              ></div>
            </v-fade-transition>
          </div>
          <div class="tw-mb-4 tw-text-center tw-text-sm tw-text-very-dark-gray">
            Billed monthly.<br />Cancel anytime.
          </div>
          <v-btn
            class="tw-mb-0.5"
            color="primary"
            outlined
            block
            :dark="
              isStudent
                ? !loadingCheckoutUrl[monthlyStudentPrice?.id]
                : !loadingCheckoutUrl[monthlyPrice?.id]
            "
            :disabled="
              isStudent
                ? loadingCheckoutUrl[monthlyStudentPrice?.id]
                : loadingCheckoutUrl[monthlyPrice?.id]
            "
            :loading="
              isStudent
                ? loadingCheckoutUrl[monthlyStudentPrice?.id]
                : loadingCheckoutUrl[monthlyPrice?.id]
            "
            @click="
              isStudent
                ? handleUpgrade(monthlyStudentPrice)
                : handleUpgrade(monthlyPrice)
            "
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
              <span class="tw-mr-1 tw-text-dark-gray tw-line-through"
                >$100</span
              >
              <span class="tw-mr-1 tw-text-4xl">{{
                isStudent
                  ? formattedPrice(lifetimeStudentPrice)
                  : formattedPrice(lifetimePrice)
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
            :dark="
              isStudent
                ? !loadingCheckoutUrl[lifetimeStudentPrice?.id]
                : !loadingCheckoutUrl[lifetimePrice?.id]
            "
            :disabled="
              isStudent
                ? loadingCheckoutUrl[lifetimeStudentPrice?.id]
                : loadingCheckoutUrl[lifetimePrice?.id]
            "
            :loading="
              isStudent
                ? loadingCheckoutUrl[lifetimeStudentPrice?.id]
                : loadingCheckoutUrl[lifetimePrice?.id]
            "
            @click="
              isStudent
                ? handleUpgrade(lifetimeStudentPrice)
                : handleUpgrade(lifetimePrice)
            "
          >
            Upgrade
          </v-btn>
        </div>
      </div>
      <div class="tw-flex tw-w-full tw-items-center tw-justify-center tw-pb-4">
        <v-checkbox
          id="student-checkbox"
          v-model="isStudent"
          dense
          hide-details
        >
        </v-checkbox>
        <label
          for="student-checkbox"
          class="tw-cursor-pointer tw-select-none tw-text-sm tw-text-very-dark-gray"
          >I'm a student</label
        >
      </div>
      <div
        class="tw-flex tw-w-full tw-items-center tw-justify-center tw-gap-4 tw-text-center"
      >
        <a
          class="tw-cursor-pointer tw-py-2 tw-text-xs tw-font-medium tw-text-dark-gray tw-underline"
          target="_blank"
          href="https://forms.gle/aaBzFvoKkHLPDjio7"
        >
          I don't want to pay
        </a>
        <div
          class="tw-cursor-pointer tw-py-2 tw-text-xs tw-font-medium tw-text-dark-gray tw-underline"
          @click="showDonatedDialog = true"
        >
          I already donated :)
        </div>
      </div>
      <v-btn text block @click="$emit('input', false)"> Cancel </v-btn>
    </v-card>

    <AlreadyDonatedDialog v-model="showDonatedDialog" />
    <StudentProofDialog v-model="showStudentProofDialog" />
  </v-dialog>
</template>

<script>
import { get, post } from "@/utils"
import { mapState, mapActions } from "vuex"
import { upgradeDialogTypes } from "@/constants"
import AlreadyDonatedDialog from "./AlreadyDonatedDialog.vue"
import StudentProofDialog from "./StudentProofDialog.vue"

export default {
  name: "UpgradeDialog",
  components: {
    AlreadyDonatedDialog,
    StudentProofDialog,
  },
  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      monthlyPrice: null,
      lifetimePrice: null,
      monthlyStudentPrice: null,
      lifetimeStudentPrice: null,
      loadingCheckoutUrl: {},
      showDonatedDialog: false,
      isStudent: false,
      showStudentProofDialog: false,
    }
  },

  computed: {
    ...mapState([
      "featureFlagsLoaded",
      "pricingPageConversion",
      "authUser",
      "upgradeDialogType",
      "upgradeDialogData",
    ]),
    upgradeDialogTypes() {
      return upgradeDialogTypes
    },
  },

  methods: {
    ...mapActions(["showError"]),
    formattedPrice(price) {
      if (!price) return ""
      return "$" + Math.floor(price.unit_amount / 100)
    },
    async init() {
      if (this.featureFlagsLoaded) {
        if (!this.lifetimePrice || !this.monthlyPrice) {
          await this.fetchPrice()
        }
      }
    },
    async fetchPrice() {
      const res = await get("/stripe/price?exp=" + this.pricingPageConversion)
      const { lifetime, monthly, lifetimeStudent, monthlyStudent } = res
      this.lifetimePrice = lifetime
      this.monthlyPrice = monthly
      this.lifetimeStudentPrice = lifetimeStudent
      this.monthlyStudentPrice = monthlyStudent
    },
    async handleUpgrade(price) {
      // if (this.isStudent) {
      //   this.showStudentProofDialog = true
      //   this.$posthog.capture("student_upgrade_attempt", {
      //     price: price,
      //   })
      //   return
      // }
      this.$posthog.capture("upgrade_clicked", {
        price: this.formattedPrice(price),
      })
      this.$set(this.loadingCheckoutUrl, price.id, true)
      try {
        let originUrl = window.location.href
        if (this.upgradeDialogData) {
          if (this.upgradeDialogType === upgradeDialogTypes.SCHEDULE_EVENT) {
            originUrl = `${originUrl}?scheduled_event=${encodeURIComponent(
              JSON.stringify(this.upgradeDialogData.scheduledEvent)
            )}`
          }
        }
        const res = await post("/stripe/create-checkout-session", {
          priceId: price.id,
          userId: this.authUser._id,
          isSubscription: price.recurring !== null,
          originUrl: originUrl,
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
    isStudent: {
      handler(val) {
        if (val) {
          this.$posthog.capture("student_pricing_viewed", {
            prices: `${this.formattedPrice(
              this.monthlyStudentPrice
            )}, ${this.formattedPrice(this.lifetimeStudentPrice)}`,
          })
        }
      },
    },
    featureFlagsLoaded: {
      handler() {
        this.init()
      },
      immediate: true,
    },
    value: {
      handler() {
        if (this.value) {
          post("/analytics/upgrade-dialog-viewed", {
            userId: this.authUser?._id ?? this.$posthog?.get_distinct_id(),
            price: `${this.formattedPrice(
              this.monthlyPrice
            )}, ${this.formattedPrice(this.lifetimePrice)}`,
            type: this.upgradeDialogType,
          })
          this.$posthog.capture("upgrade_dialog_viewed", {
            price: `${this.formattedPrice(
              this.monthlyPrice
            )}, ${this.formattedPrice(this.lifetimePrice)}`,
            type: this.upgradeDialogType,
          })
        }
      },
    },
  },
}
</script>
