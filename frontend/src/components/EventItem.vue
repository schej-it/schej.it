<template>
  <v-container
    @click="$emit('click')"
    v-ripple
    class="hover:tw-drop-shadow-lg tw-drop-shadow-md tw-shadow-inner tw-transition-all tw-bg-white tw-rounded-lg tw-flex tw-text-black tw-justify-between tw-items-center tw-py-3"
  >
    <div class="tw-ml-1">
      <div>{{ this.event.name }}</div>
      <div class="tw-text-sm tw-font-light">{{ dateString }}</div>
    </div>
    <div class="tw-min-w-max">
      <v-chip small class="tw-text-very-dark-gray tw-m-2 tw-bg-off-white">
        <v-icon left small> mdi-account-multiple </v-icon>
        {{ Object.keys(this.event.responses).length }}
      </v-chip>
      <v-icon large>mdi-chevron-right</v-icon>
    </div>
  </v-container>
</template>

<script>
import { getDateRangeStringForEvent, } from '@/utils'

export default {
  name: 'EventItem',

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({}),

  computed: {
    dateString() {
      return getDateRangeStringForEvent(this.event)
    },
  },

  methods: {},
}
</script>
