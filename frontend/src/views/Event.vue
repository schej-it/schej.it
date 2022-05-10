<template>
  <div>
    <div class="tw-bg-green tw-text-white tw-flex tw-flex-col tw-items-center tw-space-y-1 tw-py-2">
      <div 
        class="tw-font-bold tw-text-3xl"
      >{{ event.name }}</div>
      <div
        class="tw-font-light tw-text-lg"
      >{{ dateString }}</div>
    </div>
  </div>
</template>

<script>
import { getDateRangeString } from '@/utils'
import { mapState } from 'vuex'

export default {
  name: 'Event',

  computed: {
    ...mapState([ 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
    event() {
      return this.events[this.eventId]
    },
    eventId() {
      return this.$route.params.eventId
    },
  }
}
</script>