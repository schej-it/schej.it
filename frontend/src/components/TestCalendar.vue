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
      :loadingCalendarEvents="loading"
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
  getDateWithTimeNum,
  dateToTimeNum,
  clampDateToTimeNum,
  getCalendarEvents,
  isPhone,
} from '@/utils'

export default {
  props: {
    noEventNames: { type: Boolean, default: false },
  },

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    startDate: getDateWithTimeNum(new Date(), 9),
    endDate: getDateWithTimeNum(getDateDayOffset(new Date(), 2), 24),
    startTime: 9,
    endTime: 0,
    calendarEvents: [],
    loading: false,
  }),

  created() {
    if (this.isPhone) {
      this.startDate = getDateWithTimeNum(new Date(), 9)
      this.endDate = getDateWithTimeNum(getDateDayOffset(new Date(), 2), 24)
    } else {
      this.startDate = getDateWithTimeNum(new Date(), 9)
      this.endDate = getDateWithTimeNum(getDateDayOffset(new Date(), 6), 24)
    }

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
      this.loading = true
      getCalendarEvents({ 
        startDate: this.startDate,
        endDate: this.endDate,
        startTime: this.startTime,
        endTime: this.endTime,
      }).then(data => {
        this.calendarEvents = data
        this.loading = false
      }).catch((err) => {
        console.error(err)
        if (err.error.code === 401 || err.error.code === 403) {
          signInGoogle(null, true)
        }
      })
    },
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },
}
</script>
