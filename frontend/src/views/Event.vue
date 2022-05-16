<template>
  <div v-if="event">
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
import { getDateRangeString, get, signInGoogle, dateCompare, dateToTimeInt, getDateDayOffset, clampDateToTimeInt, post } from '@/utils'
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
    event: null,
  }),

  computed: {
    ...mapState([ 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
  },

  methods: {
    test() {
      console.log(document.querySelector('.v-main').scrollTop)
    },
  },

  async created() {
    this.event = await get(`/events/${this.eventId}`)
    this.event.startDate = new Date(this.event.startDate)
    this.event.endDate = new Date(this.event.endDate)
    
    get(`/user/availability?timeMin=${this.event.startDate.toISOString()}&timeMax=${getDateDayOffset(this.event.endDate, 1).toISOString()}`).then(data => {
      this.calendarEvents = data.items
        .filter(event => {
          return (
            'start' in event &&
            'end' in event &&
            'dateTime' in event.end &&
            'dateTime' in event.start &&
            dateToTimeInt(event.end.dateTime) > this.event.startTime &&
            dateToTimeInt(event.start.dateTime) < this.event.endTime
          )
        })
        .map(event => ({ 
          summary: event.summary,
          startDate: clampDateToTimeInt(new Date(event.start.dateTime), this.event.startTime, 'upper'),
          endDate: clampDateToTimeInt(new Date(event.end.dateTime), this.event.endTime, 'lower'),
        }))
    }).catch(err => {
      console.error(err)
      if (err.code === 401) {
        signInGoogle({ type: 'join', eventId: this.eventId })
      }
    })

    /*post(`/events/${this.eventId}/response`, {
      availability: ['lmao', 'what', 'cool'],
    }).then(console.log).catch(console.error)*/
    //document.querySelector('.v-main').addEventListener('scroll', this.test)
  },
}
</script>