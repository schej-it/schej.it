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
    <ScheduleOverlap 
      v-bind="event"
      :calendarEvents="calendarEvents"
    />
  </div>
</template>

<script>
import { getDateRangeString, get, signInGoogle } from '@/utils'
import { mapState } from 'vuex'

import ScheduleOverlap from '@/components/ScheduleOverlap'

export default {
  name: 'Event',

  props: {
    eventId: { type: String, required: true },
  },

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    calendarEvents: [],
  }),

  computed: {
    ...mapState([ 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
    event() {
      return this.events[this.eventId]
    },
  },

  methods: {
    test() {
      console.log(document.querySelector('.v-main').scrollTop)
    },
  },

  created() {
    get(`/user/availability?timeMin=${this.event.startDate.toISOString()}&timeMax=${this.event.endDate.toISOString()}`).then(data => {
      this.calendarEvents = data.items
        .filter(event => ('dateTime' in event.end && 'dateTime' in event.start))
        .map(event => ({ 
          summary: event.summary,
          startDate: new Date(event.start.dateTime),
          endDate: new Date(event.end.dateTime)
        }))
    }).catch(err => {
      if (err.code === 401) {
        signInGoogle({ type: 'join', eventId: this.eventId })
      }
    })
    //document.querySelector('.v-main').addEventListener('scroll', this.test)
  },
}
</script>