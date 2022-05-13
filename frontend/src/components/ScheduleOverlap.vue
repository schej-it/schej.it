<template>
  <div class="tw-p-10">
    <div class="tw-flex">
      <div class="tw-w-10" />
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
      <div class="tw-w-10">
        <div 
          v-for="time, i in times"
          :key="i"
          class="tw-h-6 tw-text-xs tw-pt-1"  
        >
          {{ time.text }}
        </div>
      </div>
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
            class="tw-h-6 tw-border-light-gray tw-border-r tw-border-b" 
            :class="borderClass(d, t)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { timeIntToTimeString } from '@/utils'

export default {
  name: 'ScheduleOverlap',

  props: {
    startDate: { type: Date, required: true },
    endDate: { type: Date, required: true },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    responses: { type: Array, required: true },
  },

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
        })
        curDate = new Date(curDate.getTime() + 24*60*60*1000)
      }

      return days
    },
    times() {
      const times = []
      
      for (let t = this.startTime; t < this.endTime; ++t) {
        times.push({
          time: t,
          text: timeIntToTimeString(t),
        })
        times.push({
          time: t + 0.5,
        })
      }
      times.push({
        time: this.endTime,
        text: timeIntToTimeString(this.endTime),
      })

      return times
    },
  },

  methods: {
    borderClass(day, time) {
      let c = ''
      if (day === 0) c += 'tw-border-l tw-border-l-gray '
      if (day === this.days.length-1) c += 'tw-border-r-gray '
      if (time === 0) c+= 'tw-border-t tw-border-t-gray '
      if (time === this.times.length-1) c += 'tw-border-b-gray '
      return c
    },
  },
}
</script>