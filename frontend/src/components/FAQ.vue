<template>
  <v-container
    :class="`${
      toggled ? 'tw-h-40 tw-border-green' : 'tw-h-20 tw-border-gray'
    } tw-w-full tw-transition-all tw-flex tw-flex-col tw-text-left tw-border-[1px] tw-p-6 tw-rounded-md tw-overflow-hidden tw-cursor-pointer`"
    @click="() => toggled = !toggled"
  >
    <div class="tw-flex tw-flex-row tw-justify-between tw-content-center sm:tw-text-xl lg:tw-text-2xl">
      <div class="tw-font-bold">{{ question }}</div>
              <v-icon size="x-large" :class="`${toggled ? 'tw-rotate-45 tw-text-green' : 'tw-rotate-0 tw-text-gray'}`">mdi-plus</v-icon>
    </div>
    <div class="tw-mt-5 sm:tw-text-xl lg:tw-text-2xl">{{ answer }}</div>
  </v-container>
</template>

<script>
export default {
  name: "FAQ",

  props: {
    question: { type: String, required: true },
    answer: { type: String, required: true },
  },

  data: () => ({
    toggled: false,
  }),

  computed: {},

  methods: {},
}
</script>
