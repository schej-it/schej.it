<template>
  <div class="tw-h-full tw-bg-light-gray tw-pt-24">
    <h1 class="tw-text-center tw-text-3xl tw-font-medium md:tw-text-4xl">
      Schej is free <br class="tw-block md:tw-hidden" />
      for most people
    </h1>
    <div class="tw-mt-5 tw-text-center tw-text-base md:tw-text-lg">
      Pay us if you want to support us :)
    </div>
    <div
      class="tw-mt-10 tw-flex tw-flex-col tw-items-center tw-justify-center tw-gap-6 tw-p-4 md:tw-flex-row md:tw-items-start"
    >
      <PricingCard v-bind="pricingPlans.FREE">
        <template #button>
          <v-btn color="primary" class="tw-w-full tw-text-white"
            >Create an event</v-btn
          >
        </template>
      </PricingCard>
      <PricingCard v-bind="pricingPlans.PRO">
        <template #button>
          <v-btn
            color="tw-bg-very-light-gray"
            class="tw-w-full"
            style="box-shadow: 0 0 #0000 !important"
            >Continue with Pro</v-btn
          >
        </template>
      </PricingCard>
      <PricingCard v-bind="pricingPlans.ENTERPRISE">
        <template #button>
          <v-btn
            color="tw-bg-very-light-gray"
            class="tw-w-full"
            style="box-shadow: 0 0 #0000 !important"
            >Contact us</v-btn
          >
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
  components: {
    PricingCard,
  },

  computed: {
    pricingPlans() {
      return pricingPlans
    },
  },
}
</script>
