<template>
  <v-card
    class="tw-flex tw-w-[290px] tw-flex-col tw-gap-6 tw-rounded-md tw-p-4"
  >
    <div class="tw-text-base tw-font-medium">{{ header }}</div>
    <div class="tw-flex tw-h-[80px] tw-flex-col tw-justify-center tw-gap-1">
      <div class="tw-text-xs tw-font-medium tw-text-red">
        {{ specialText }}
      </div>
      <div class="tw-flex tw-items-baseline tw-gap-1">
        <div class="tw-text-4xl tw-font-medium">{{ price }}</div>
        <div class="tw-text-sm">{{ priceText }}</div>
      </div>
      <div class="tw-text-sm tw-text-dark-gray">{{ subtitle }}</div>
    </div>
    <slot name="button"></slot>
    <div class="tw-flex tw-flex-col tw-gap-1.5">
      <div
        v-for="feature in features"
        class="tw-flex tw-items-center tw-gap-2 tw-text-sm"
      >
        <v-icon small>mdi-check</v-icon>
        {{ feature }}
      </div>
    </div>
  </v-card>
</template>

<script>
export default {
  name: "PricingCard",

  props: {
    header: { type: String, required: true },
    price: { type: String, required: true },
    subtitle: { type: String, required: true },
    features: { type: Array, required: true },
    specialText: { type: String, required: false },
    priceText: { type: String, required: false },
  },
}
</script>
