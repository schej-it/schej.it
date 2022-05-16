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
        <div class="tw-flex">
          <div 
            v-for="day, d in days" 
            :key="d"
            class="tw-flex-1 tw-relative"
          >
            <!-- Timeslots -->
            <div 
              v-for="time, t in times"
              :key="t"
              class="tw-w-full"  
            >
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <div 
                    class="tw-h-5 tw-border-light-gray tw-border-r" 
                    :class="timeslotClass(day, time, d, t)"
                    v-bind="attrs"
                    v-on="on"
                  />
                </template>
                <div>
                  <div>Available:</div>
                  <div 
                    v-for="user in getRespondentsForDateTime(day.dateObject, time.timeInt)"
                    :key="user._id"
                  >
                    {{ user.firstName }}
                  </div>
                </div>
              </v-tooltip>
            </div>
            
            <!-- Calendar events -->
            <template v-if="showCalendarEvents">
              <div
                v-for="event, e in calendarEventsByDay[d]" 
                :key="`${d}-${e}`"
                class="tw-absolute tw-w-full tw-p-px"
                :style="event.style"
              >
                <div
                  class="tw-bg-blue/25 tw-border-light-gray /*tw-border-solid*/ tw-border tw-w-full tw-h-full tw-text-ellipsis tw-text-xs tw-rounded tw-p-px tw-overflow-hidden"
                >
                  <div class="tw-text-blue tw-font-medium">
                    {{ noEventNames ? 'BUSY' : event.summary }}
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>

        <template v-if="!calendarOnly">
          <v-btn 
            class="tw-my-2"
            block
            @click="showCalendarEvents = !showCalendarEvents"
          >
            <template v-if="!showCalendarEvents">
              Edit <v-icon small class="tw-ml-1">mdi-pencil</v-icon>
            </template>
            <template v-else>
              View suggested times 
            </template>
          </v-btn>
          <v-btn
            class="tw-mb-2"
            block
            @click="copyLink"
          >
            Copy link <v-icon small class="tw-ml-1">mdi-content-copy</v-icon>
          </v-btn>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import { timeIntToTimeText, getDateDayOffset, dateCompare, compareDateDay, dateToTimeInt, getDateWithTimeInt, post } from '@/utils'
import { mapActions, mapState } from 'vuex'

export default {
  name: 'ScheduleOverlap',

  props: {
    eventId: { type: String, default: '' },
    startDate: { type: Date, required: true },
    endDate: { type: Date, required: true },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    responses: { type: Object, default: () => ({}) },

    calendarEvents: { type: Array, required: true },
    initialShowCalendarEvents: { type: Boolean, default: true },

    noEventNames: { type: Boolean, default: false },
    calendarOnly: { type: Boolean, default: false },
  },

  data() {
    return {
      max: 0, // The max number of respondents for a given timeslot
      showCalendarEvents: this.initialShowCalendarEvents, 
      availability: new Set(), // Current availability of the current user, an array of dates
    }
  },

  computed: {
    ...mapState([ 'authUser' ]),
    availabilityArray() {
      /* Returns the availibility as an array */
      return [...this.availability].map(item => new Date(item))
    },
    calendarEventsByDay() {
      /* Returns a 2d array of events based on the day they take place. Index 0 = first day */
      
      // TODO: calendar event spanning two days breaks this (how to fix: split calendar events into two or more events if span multiple days)
      const arr = []
      for (let i = 0; i < this.days.length; ++i) {
        arr[i] = []
      }
      for (const calendarEvent of this.calendarEvents) {
        const startTime = dateToTimeInt(calendarEvent.startDate)
        const endTime = dateToTimeInt(calendarEvent.endDate)
        for (const d in this.days) {
          const day = this.days[d]
          if (compareDateDay(day.dateObject, calendarEvent.startDate) == 0) {
            arr[d].push({
              ...calendarEvent,
              style: {
                top: `calc(${startTime-this.startTime} * 2 * 1.25rem)`, // 1.25 rem = tw-h-5 
                height: `calc(${endTime-startTime} * 2 * 1.25rem)`
              }
            })
            break
          }
        }
      }
      return arr
    },
    currentResponse() {
      /* Returns a response object for the current user */
      return {
        name: 'jony',
        availability: this.availabilityArray,
      }
    },
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
          const date = getDateWithTimeInt(day.dateObject, time.timeInt)
          formatted.set(date.getTime(), new Set())
          
          for (const response of [...Object.values(this.responses), /*this.currentResponse*/]) {
            const index = response.availability.findIndex(d => dateCompare(d, date) === 0)
            if (index !== -1) {
              // TODO: determine whether I should delete the index??
  
              formatted.get(date.getTime()).add(response.user)
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
          timeInt: t,
          text: timeIntToTimeText(t),
        })
        times.push({
          timeInt: t + 0.5,
        })
      }

      return times
    },
    userHasResponded() {
      return this.authUser._id in this.responses
    }
  },

  methods: {
    ...mapActions([ 'showInfo' ]),
    copyLink() {
      navigator.clipboard.writeText(`http://localhost:8080/j/${this.eventId}`)
    },
    getRespondentsForDateTime(date, time) {
      /* Returns an array of respondents for the given date/time */
      const d = getDateWithTimeInt(date, time)
      return this.responsesFormatted.get(d.getTime())
    },
    async setAvailability() {
      /* Constructs the availability array using calendarEvents array */
      // This is not a computed property because we should be able to change it manually from what it automatically fills in
      // TODO: there is a bug where it still sets your availability if the event length is 30 minutes
      this.availability = new Set()
      for (const d in this.days) {
        const day = this.days[d]
        for (const time of this.times) {
          // Check if there exists a calendar event that overlaps [time, time+0.5]
          const startDate = getDateWithTimeInt(day.dateObject, time.timeInt)
          const endDate = getDateWithTimeInt(day.dateObject, time.timeInt + 0.5)
          const index = this.calendarEventsByDay[d].findIndex(e => {
            return (
              (dateCompare(e.startDate, startDate) < 0 && dateCompare(e.endDate, startDate) > 0) ||
              (dateCompare(e.startDate, endDate) < 0 && dateCompare(e.endDate, endDate) > 0) ||
              (dateCompare(e.startDate, startDate) == 0 && dateCompare(e.endDate, endDate) == 0)
            )
          })
          if (index === -1) {
            this.availability.add(startDate.getTime())
          }
        }
      }
      await post(`/events/${this.eventId}/response`, { availability: this.availabilityArray })
      this.$emit('refreshEvent')
      this.showInfo('All done! Your availability has been set automatically')
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
      if (this.showCalendarEvents) {
        // Show only current user availability
        const date = getDateWithTimeInt(day.dateObject, time.timeInt)
        if (this.availability.has(date.getTime())) {
          c += 'tw-bg-avail-green-300 '
        }
      } else {
        // Show everyone's availability
        const numRespondents = this.getRespondentsForDateTime(day.dateObject, time.timeInt).size
        if (numRespondents > 0) {
          const frac = numRespondents / this.max
          const colors = [
            //'tw-bg-avail-green-50', 
            'tw-bg-avail-green-100', 
            'tw-bg-avail-green-200', 
            'tw-bg-avail-green-300', 
            'tw-bg-avail-green-400', 
            'tw-bg-light-blue', 
            //'tw-bg-avail-green-600',
          ] 
          c += colors[parseInt(frac*colors.length-1)] + ' '
        }
      }

      return c
    },
  },

  watch: {
    calendarEvents: {
      handler() {
        if (!this.userHasResponded && !this.calendarOnly) this.setAvailability()
      },
    },
  },

  created() {
    if (this.userHasResponded) {
      this.availability = new Set()
      this.responses[this.authUser._id].availability.forEach(item => this.availability.add(new Date(item).getTime()))
    }
  },
}
</script>