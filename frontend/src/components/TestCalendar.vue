<template>
  <ScheduleOverlap 
    :startDate="startDate"
    :endDate="endDate"
    :startTime="startTime"
    :endTime="endTime"
    :calendarEvents="calendarEvents"
    :noEventNames="noEventNames"
    
    calendarOnly
  />
</template>

<script>
import ScheduleOverlap from '@/components/ScheduleOverlap'

import { signInGoogle, get, getDateDayOffset, getDateWithTime, dateToTimeInt, clampDateToTimeInt } from '@/utils'

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
    endTime: 22,
    calendarEvents: [],
  }),

  created() {

    get(`/user/calendar?timeMin=${this.startDate.toISOString()}&timeMax=${getDateDayOffset(this.endDate, 1).toISOString()}`).then(data => {
      this.calendarEvents = data
        .map(event => ({ 
          summary: event.summary,
          startDate: clampDateToTimeInt(new Date(event.startDate), this.startTime, 'upper'),
          endDate: clampDateToTimeInt(new Date(event.endDate), this.endTime, 'lower'),
        }))
    }).catch(err => {
      console.error(err)
      if (err.error.code === 401) {
        signInGoogle()
      }
    })
  }
}
</script>