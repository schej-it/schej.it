<template>
  <v-container
    @click="$emit('click')"
    v-ripple
    class="tw-bg-light-gray tw-rounded-md tw-flex tw-text-black tw-justify-between tw-align-middle tw-py-2"
  >
    <div>
      <div>{{ this.event.name }}</div>
      <div class="tw-text-sm tw-font-light">{{ dateString }}</div>
    </div>
    <div>
      <v-chip class="ma-2" color="gray lighten-2" text-color="gray darken-2">
        <v-icon left> mdi-account-multiple </v-icon>
        {{ this.event.responses.length }}
      </v-chip>
      <v-icon>mdi-chevron-right</v-icon>
    </div>
  </v-container>
</template>

<script>
import { getDateRangeString } from '@/utils'

export default {
  name: 'EventItem',

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({}),

  computed: {
    dateString() {
      return getDateRangeString(new Date(this.event.startDate), new Date(this.event.endDate))
    },
  },

  methods: {},
}
</script>
