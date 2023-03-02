<template>
  <v-container
    :class="`${
      toggled ? 'tw-max-h-96 tw-border-green' : 'tw-max-h-24 sm:tw-max-h-20 tw-border-gray'
    } tw-w-full tw-transition-all tw-flex tw-flex-col tw-text-left tw-border-[1px] tw-p-6 tw-rounded-md tw-overflow-hidden tw-cursor-pointer`"
    @click="() => (toggled = !toggled)"
  >
    <div
      class="tw-flex tw-flex-row tw-mb-12 sm:tw-mb-7 tw-justify-between tw-content-center sm:tw-text-xl lg:tw-text-2xl"
    >
      <div class="tw-font-medium">{{ question }}</div>
      <v-icon
        size="x-large"
        :class="`${
          toggled ? 'tw-rotate-45 tw-text-green' : 'tw-rotate-0 tw-text-gray'
        }`"
        >mdi-plus</v-icon
      >
    </div>
    <div class="sm:tw-text-xl lg:tw-text-2xl">
      <div>{{ answer }}</div>
      <div class="tw-gap-2 tw-flex tw-flex-col">
        <div v-for="(point, index) in points" class="tw-flex tw-items-center">
          <div class="tw-mr-1 tw-bg-green tw-text-white tw-w-8 tw-h-8 tw-rounded-full tw-items-center tw-flex tw-justify-center">{{ index + 1 }}</div>
          <div> {{ point }}</div>
        </div>
      </div>
    </div>
  </v-container>
</template>

<script>
export default {
  name: "FAQ",

  props: {
    question: { type: String, required: true },
    answer: { type: String },
    points: { type: Array },
  },

  data: () => ({
    toggled: false,
  }),

  computed: {},

  methods: {},
}
</script>
