<template>
  <div>
    <div
      class="tw-flex tw-w-full tw-justify-between tw-pl-12 tw-relative tw-bg-gray tw-h-0 tw-top-5"
    >
      <v-btn icon color="gray darken-2" @click="shiftBack">
        <v-icon large> mdi-chevron-left </v-icon>
      </v-btn>
      <v-btn icon color="gray darken-2" @click="shiftForward">
        <v-icon large> mdi-chevron-right </v-icon>
      </v-btn>
    </div>
    <ScheduleOverlap
      :startDate="startDate"
      :endDate="endDate"
      :startTime="startTime"
      :endTime="endTime"
      :calendarEvents="calendarEvents"
      :noEventNames="noEventNames"
      calendarOnly
    />
  </div>
</template>

<script>
import ScheduleOverlap from '@/components/ScheduleOverlap'

import {
  signInGoogle,
  get,
  getDateDayOffset,
  getDateWithTime,
  dateToTimeInt,
  clampDateToTimeInt,
} from '@/utils'

export default {
  props: {
    noEventNames: { type: Boolean, default: false },
  },

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    startDate: getDateWithTime(new Date(), '0:00'),
    endDate: getDateWithTime(getDateDayOffset(new Date(), 2), '11:59'),
    startTime: 9,
    endTime: 24,
    calendarEvents: [],
  }),

  created() {
    this.retrieveCalendarEvents()
  },

  methods: {
    shiftBack() {
      this.shiftBy(-1)
    },
    shiftForward() {
      this.shiftBy(1)
    },
    shiftBy(numberDays) {
      this.startDate.setDate(this.startDate.getDate() + numberDays)
      this.endDate.setDate(this.endDate.getDate() + numberDays)
      this.startDate = new Date(this.startDate)
      this.retrieveCalendarEvents()
    },
    retrieveCalendarEvents() {
      get(
        `/user/calendar?timeMin=${this.startDate.toISOString()}&timeMax=${getDateDayOffset(
          this.endDate,
          1
        ).toISOString()}`
      )
        .then((data) => {
          this.calendarEvents = data.map((event) => ({
            summary: event.summary,
            startDate: clampDateToTimeInt(
              new Date(event.startDate),
              this.startTime,
              'upper'
            ),
            endDate: clampDateToTimeInt(
              new Date(event.endDate),
              this.endTime,
              'lower'
            ),
          }))
        })
        .catch((err) => {
          console.error(err)
          if (err.error.code === 401) {
            signInGoogle()
          }
        })
    },
  },
}
</script>
