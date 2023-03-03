<template>
  <div
    :class="`${
      toggled ? 'tw-border-green' : 'tw-border-gray'
    } tw-w-full tw-transition-all tw-flex tw-flex-col tw-text-left tw-border-[1px] tw-p-4 sm:tw-p-6 tw-rounded-md tw-overflow-hidden tw-cursor-pointer`"
    @click="() => (toggled = !toggled)"
  >
    <div
      class="tw-flex tw-flex-row tw-justify-between tw-content-center sm:tw-text-lg"
    >
      <div class="tw-font-medium tw-mr-4">{{ question }}</div>
      <v-icon
        size="x-large"
        :class="`${
          toggled ? 'tw-rotate-45 tw-text-green' : 'tw-rotate-0 tw-text-gray'
        }`"
      >mdi-plus</v-icon>
    </div>

    <v-expand-transition>
      <div v-if="toggled">
        <div class="sm:tw-text-lg tw-pt-4 sm:tw-pt-6">
          <div>{{ answer }}</div>
          <div class="tw-gap-2 tw-flex tw-flex-col">
            <div v-for="(point, index) in points" class="tw-flex tw-items-center">
              <div class="tw-w-5 tw-h-5 tw-mr-2 tw-shrink-0 tw-bg-green tw-text-white tw-rounded-full tw-flex tw-items-center tw-justify-center">{{ index + 1 }}</div>
              <div>{{ point }}</div>
            </div>
          </div>
        </div>
      </div>
    </v-expand-transition>
  </div>
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
