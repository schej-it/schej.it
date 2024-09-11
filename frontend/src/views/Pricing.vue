<template>
  <div class="tw-h-full tw-pt-24">
    <div class="tw-flex tw-justify-center">
      <div
        class="tw-w-fit tw-bg-gradient-to-r tw-from-lightest-green tw-to-darkest-green tw-bg-clip-text tw-pb-2 tw-text-center tw-text-4xl tw-font-medium tw-text-transparent sm:tw-text-5xl"
      >
        Free to Schej,
        <br tw="tw-block" />
        upgrade as you go!
      </div>
    </div>
    <div
      class="tw-text-md tw-mt-5 tw-px-10 tw-text-center tw-text-base tw-text-very-dark-gray"
    >
      Support the Schej team and unlock more scheduling features.
    </div>
    <div
      class="tw-mt-10 tw-flex tw-flex-col tw-items-center tw-justify-center tw-gap-6 tw-p-4 md:tw-flex-row md:tw-items-end"
    >
      <PricingCard v-bind="pricingPlans.FREE">
        <template #button>
          <v-btn
            href="/"
            color="tw-bg-white"
            class="tw-w-full"
            style="box-shadow: 0 0 #0000 !important"
            >Create an event</v-btn
          >
        </template>
      </PricingCard>
      <PricingCard v-bind="pricingPlans.PRO">
        <template #button>
          <v-btn
            href="/upgrade-to-pro"
            color="tw-bg-white"
            class="tw-w-full tw-text-green"
            >Upgrade to Pro</v-btn
          >
        </template>
      </PricingCard>
      <PricingCard v-bind="pricingPlans.ENTERPRISE">
        <template #button>
          <v-btn color="tw-bg-white" class="tw-w-full">Contact us</v-btn>
        </template>
      </PricingCard>
    </div>
  </div>
</template>

<script>
import PricingCard from "@/components/pricing/PricingCard.vue"
import { pricingPlans } from "@/constants"

export default {
  name: "Pricing",

  metaInfo: {
    title: "Pricing - Schej",
  },

  components: {
    PricingCard,
  },

  methods: {
    createNew() {
      this.$emit("setNewDialogOptions", {
        show: true,
        contactsPayload: {},
        openNewGroup: false,
      })
    },
  },

  computed: {
    pricingPlans() {
      return pricingPlans
    },
  },
}
</script>
