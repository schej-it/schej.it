<template>
  <div class="tw-flex tw-flex-col tw-justify-center tw-relative tw-pt-28">
    <v-img
      v-if="image"
      alt="schej character"
      :src="image"
      :height="180"
      transition="fade-transition"
      contain
      class="tw-block tw-absolute tw-top-0 tw-right-4"
    />
    <v-card
      :class="`${classes} ${
        highlighted &&
        'tw-bg-gradient-to-b tw-from-lightest-green tw-to-darkest-green tw-text-white'
      } tw-flex tw-w-[290px] tw-flex-col tw-gap-6 tw-rounded-xl tw-p-6`"
    >
      <div class="tw-flex tw-items-center tw-gap-2 tw-text-base tw-font-medium">
        {{ header }}
        <div
          v-if="specialText"
          class="tw-rounded-md tw-bg-white tw-px-2 tw-py-1 tw-text-xs tw-font-medium tw-text-green"
        >
          {{ specialText }}
        </div>
      </div>
      <div class="tw-flex tw-h-[52px] tw-flex-col tw-justify-center tw-gap-1">
        <div class="tw-flex tw-items-baseline tw-gap-1">
          <div class="tw-text-4xl tw-font-medium">{{ price }}</div>
          <div class="tw-text-sm">{{ priceText }}</div>
        </div>
        <div
          :class="`tw-text-sm ${
            highlighted ? 'tw-text-off-white' : 'tw-text-dark-gray'
          }`"
        >
          {{ subtitle }}
        </div>
      </div>
      <div class="tw-flex tw-flex-col tw-gap-1.5">
        <div
          v-for="feature in features"
          class="tw-flex tw-items-center tw-gap-2 tw-text-sm"
        >
          <v-icon
            small
            :class="highlighted ? 'tw-text-white' : 'tw-text-light-green'"
            >mdi-check</v-icon
          >
          {{ feature }}
        </div>
      </div>
      <slot name="button"></slot>
    </v-card>
  </div>
</template>

<script>
export default {
  name: "PricingCard",

  props: {
    classes: { type: String, default: "" },
    header: { type: String, required: true },
    price: { type: String, required: true },
    subtitle: { type: String, required: true },
    features: { type: Array, required: true },
    specialText: { type: String, required: false },
    priceText: { type: String, required: false },
    highlighted: { type: Boolean, default: false },
    image: { type: String, default: ""}
  },
}
</script>
