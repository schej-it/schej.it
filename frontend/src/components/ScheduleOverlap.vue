<template>
  <div class="tw-p-4">
    <div class="tw-flex">
      <div class="tw-w-12" />
      <div 
        v-for="day, i in days"
        :key="i"
        class="tw-flex-1"
      >
        <div class="tw-text-center">
          <div class="tw-uppercase tw-font-light tw-text-xs">{{ day.dayText }}</div>
          <div class="tw-text-lg">{{ day.date }}</div>
        </div>
      </div>
    </div>
    <div class="tw-flex">
      <div class="tw-w-12">
        <div 
          v-for="time, i in times"
          :key="i"
          class="tw-h-5 tw-text-xs tw-pt-1 tw-pr-2 tw-text-right tw-uppercase tw-font-light"  
        >
          {{ time.text }}
        </div>
      </div>
      <div class="tw-flex-1">
        <div class="tw-flex-1 tw-flex">
          <div 
            v-for="day, d in days" 
            :key="d"
            class="tw-flex-1"
          >
            <div 
              v-for="time, t in times"
              :key="t"
              class="tw-w-full"  
            >
              <div 
                class="tw-h-5 tw-border-light-gray tw-border-r" 
                :class="timeslotClass(day, time, d, t)"
              />
            </div>
          </div>
        </div>

        <v-btn 
          class="tw-mt-2"
          block
        >
        Edit <v-icon small class="tw-ml-1">mdi-pencil</v-icon>
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script>
import { timeIntToTimeText, getDateDayOffset, getDateWithTime, dateCompare } from '@/utils'

export default {
  name: 'ScheduleOverlap',

  props: {
    startDate: { type: Date, required: true },
    endDate: { type: Date, required: true },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    responses: { type: Array, required: true },
  },

  data: () => ({
    max: 0 // The max number of respondents for a given timeslot
  }),

  computed: {
    days() {
      /* Return the days that are encompassed by startDate and endDate */
      const days = []
      const daysOfWeek = ['sun', 'mon', 'tue', 'wed', 'thu', 'fri', 'sat']
      let curDate = this.startDate 
      while (curDate.getTime() <= this.endDate.getTime())  {
        days.push({
          dayText: daysOfWeek[curDate.getDay()],
          date: curDate.getDate(),
          dateObject: curDate,
        })
        curDate = getDateDayOffset(curDate, 1)
      }

      return days
    },
    responsesFormatted() {
      /* Formats the responses in a map where date/time is mapped to the people that are available then */
      const formatted = new Map()
      for (const day of this.days) {
        for (const time of this.times) {
          const date = getDateWithTime(day.dateObject, time.timeString)
          formatted.set(date.getTime(), new Set())
          
          for (const response of this.responses) {
            const index = response.times.findIndex(d => dateCompare(d, date) === 0)
            if (index !== -1) {
              // TODO: determine whether I should delete the index??
              //response.times.splice(index, 1)
  
              formatted.get(date.getTime()).add(response.name)
            }
          }

          // Update max
          if (formatted.get(date.getTime()).size > this.max) {
            this.max = formatted.get(date.getTime()).size
          }
        }
      }
      return formatted
    },
    times() {
      /* Returns the times that are encompassed by startTime and endTime */
      const times = []
      
      for (let t = this.startTime; t < this.endTime; ++t) {
        times.push({
          time: t,
          timeString: `${t}:00`,
          text: timeIntToTimeText(t),
        })
        times.push({
          time: t + 0.5,
          timeString: `${t}:30`,
        })
      }

      return times
    },
  },

  methods: {
    getRespondentsForDateTime(date, time) {
      /* Returns an array of respondents for the given date/time */
      const d = getDateWithTime(date, time)
      return this.responsesFormatted.get(d.getTime())
    },
    timeslotClass(day, time, d, t) {
      /* Returns a class string for the given timeslot div */
      let c = ''
      
      // Border style
      if (!('text' in time)) c += 'tw-border-b '
      if (d === 0) c += 'tw-border-l tw-border-l-gray '
      if (d === this.days.length-1) c += 'tw-border-r-gray '
      if (t === 0) c+= 'tw-border-t tw-border-t-gray '
      if (t === this.times.length-1) c += 'tw-border-b-gray '

      // Fill style
      const numRespondents = this.getRespondentsForDateTime(day.dateObject, time.timeString).size
      if (numRespondents > 0) {
        const frac = numRespondents / this.max
        const colors = [
          //'tw-bg-avail-green-50', 
          'tw-bg-avail-green-100', 
          'tw-bg-avail-green-200', 
          'tw-bg-avail-green-300', 
          'tw-bg-avail-green-400', 
          'tw-bg-avail-green-500', 
          //'tw-bg-avail-green-600',
        ] 
        c += colors[parseInt(frac*colors.length-1)] + ' '
      }

      return c
    },
  },
}
</script>