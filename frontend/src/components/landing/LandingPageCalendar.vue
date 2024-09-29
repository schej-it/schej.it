<template>
  <v-card class="tw-m-4 tw-rounded-lg lg:tw-w-[34rem]">
    <!-- Brendan W was here -->
    <div class="-tw-ml-3 sm:tw-ml-0">
      <ScheduleOverlap
        ref="scheduleOverlap"
        :event="event"
        :sampleCalendarEventsByDay="calendarEventsByDay"
        calendarOnly
        :interactable="false"
        :showSnackbar="false"
        :alwaysShowCalendarEvents="true"
        animateTimeslotAlways
        :showHintText="false"
      />
    </div>
  </v-card>
</template>

<script>
import ScheduleOverlap from "@/components//schedule_overlap/ScheduleOverlap"
import {
  getDateDayOffset,
  getDateWithTimeNum,
  isPhone,
  dateToTimeNum,
  processTimeBlocks,
} from "@/utils"

export default {
  name: "LandingPageCalendar",

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    dates: [
      getDateWithTimeNum(getDateDayOffset(new Date(), -1), 9),
      getDateWithTimeNum(new Date(), 9),
      getDateWithTimeNum(getDateDayOffset(new Date(), 1), 9),
    ],

    responses: {},

    calendarEventsByDay: [],
  }),

  computed: {
    event() {
      return {
        dates: this.dates,
        duration: this.duration,
        startTime: this.startTime,
        endTime: this.endTime,
        responses: this.responses,
      }
    },
    duration() {
      return isPhone(this.$vuetify) ? 6 : 8
    },
    startTime() {
      return dateToTimeNum(new Date(this.dates[0]), true)
    },
    endTime() {
      return (this.startTime + this.duration) % 24
    },
  },

  methods: {
    getCalendarEventsByDay() {
      const [day1, day2, day3] = this.dates
      const events = [
        {
          startDate: getDateWithTimeNum(day1, 9),
          endDate: getDateWithTimeNum(day1, 10),
          summary: "Coffee with Jen",
        },
        {
          startDate: getDateWithTimeNum(day2, 11),
          endDate: getDateWithTimeNum(day2, 14),
          summary: "Karaoke with friends",
        },
        {
          startDate: getDateWithTimeNum(day3, 13),
          endDate: getDateWithTimeNum(day3, 17),
          summary: "Study session",
        },
      ]

      if (!isPhone(this.$vuetify)) {
        events.push({
          startDate: getDateWithTimeNum(day3, 20.5),
          endDate: getDateWithTimeNum(day3, 22),
          summary: "Hackathon meeting",
        })
      }

      this.calendarEventsByDay = processTimeBlocks(
        this.dates,
        this.duration,
        events
      )
    },
    getResponses() {
      const [day1, day2, day3] = this.dates
      this.responses = {
        "62828fec1bc681fa020632f2": {
          user: { _id: "1", name: "1" },
          availability: [
            getDateWithTimeNum(day1, 10),
            getDateWithTimeNum(day1, 10.5),
            getDateWithTimeNum(day1, 11),
            getDateWithTimeNum(day1, 11.5),
            getDateWithTimeNum(day1, 12),
            getDateWithTimeNum(day1, 12.5),
            getDateWithTimeNum(day1, 13),
            getDateWithTimeNum(day1, 14.5),
            getDateWithTimeNum(day1, 15),
            getDateWithTimeNum(day1, 15.5),
            getDateWithTimeNum(day1, 16),
            getDateWithTimeNum(day1, 16.5),
            getDateWithTimeNum(day1, 17),
            getDateWithTimeNum(day1, 17.5),
            getDateWithTimeNum(day1, 18),
            getDateWithTimeNum(day1, 18.5),
            getDateWithTimeNum(day1, 19),
            getDateWithTimeNum(day1, 20.5),
            getDateWithTimeNum(day1, 21),
            getDateWithTimeNum(day1, 21.5),

            getDateWithTimeNum(day2, 9),
            getDateWithTimeNum(day2, 9.5),
            getDateWithTimeNum(day2, 14),
            getDateWithTimeNum(day2, 14.5),
            getDateWithTimeNum(day2, 15),
            getDateWithTimeNum(day2, 15.5),
            getDateWithTimeNum(day2, 16),
            getDateWithTimeNum(day2, 16.5),
            getDateWithTimeNum(day2, 17),
            getDateWithTimeNum(day2, 17.5),
            getDateWithTimeNum(day2, 18),
            getDateWithTimeNum(day2, 18.5),
            getDateWithTimeNum(day2, 19),
            getDateWithTimeNum(day2, 19.5),
            getDateWithTimeNum(day2, 20),
            getDateWithTimeNum(day2, 21.5),

            getDateWithTimeNum(day3, 12.5),
            getDateWithTimeNum(day3, 17),
            getDateWithTimeNum(day3, 19.5),
          ],
        },
        "628292fe6e12d2baa9c01395": {
          user: { _id: "2", name: "2" },
          availability: [
            getDateWithTimeNum(day1, 14),
            getDateWithTimeNum(day1, 14.5),
            getDateWithTimeNum(day1, 15),
            getDateWithTimeNum(day1, 15.5),
            getDateWithTimeNum(day1, 16),
            getDateWithTimeNum(day1, 16.5),
            getDateWithTimeNum(day1, 17),
            getDateWithTimeNum(day1, 17.5),
            getDateWithTimeNum(day1, 18),
            getDateWithTimeNum(day1, 18.5),
            getDateWithTimeNum(day1, 19),
            getDateWithTimeNum(day1, 19.5),
            getDateWithTimeNum(day1, 20),
            getDateWithTimeNum(day1, 20.5),
            getDateWithTimeNum(day1, 21),
            getDateWithTimeNum(day1, 21.5),

            getDateWithTimeNum(day2, 9),
            getDateWithTimeNum(day2, 9.5),
            getDateWithTimeNum(day2, 10),
            getDateWithTimeNum(day2, 10.5),
            getDateWithTimeNum(day2, 16),
            getDateWithTimeNum(day2, 16.5),
            getDateWithTimeNum(day2, 17),
            getDateWithTimeNum(day2, 17.5),
            getDateWithTimeNum(day2, 18),
            getDateWithTimeNum(day2, 18.5),
            getDateWithTimeNum(day2, 19),
            getDateWithTimeNum(day2, 19.5),
            getDateWithTimeNum(day2, 20),
            getDateWithTimeNum(day2, 20.5),
            getDateWithTimeNum(day2, 21),
            getDateWithTimeNum(day2, 21.5),

            getDateWithTimeNum(day3, 9),
            getDateWithTimeNum(day3, 9.5),
            getDateWithTimeNum(day3, 10),
            getDateWithTimeNum(day3, 12.5),
            getDateWithTimeNum(day3, 17),
            getDateWithTimeNum(day3, 18),
            getDateWithTimeNum(day3, 18.5),
            getDateWithTimeNum(day3, 19),
            getDateWithTimeNum(day3, 19.5),
            getDateWithTimeNum(day3, 20),
          ],
        },
        "628208870df4418ff4213757": {
          user: { _id: "3", name: "3" },
          availability: [
            getDateWithTimeNum(day1, 11),
            getDateWithTimeNum(day1, 11.5),
            getDateWithTimeNum(day1, 12),
            getDateWithTimeNum(day1, 12.5),
            getDateWithTimeNum(day1, 13),
            getDateWithTimeNum(day1, 13.5),
            getDateWithTimeNum(day1, 14),
            getDateWithTimeNum(day1, 14.5),
            getDateWithTimeNum(day1, 15),
            getDateWithTimeNum(day1, 15.5),
            getDateWithTimeNum(day1, 16),
            getDateWithTimeNum(day1, 16.5),
            getDateWithTimeNum(day1, 20),
            getDateWithTimeNum(day1, 20.5),
            getDateWithTimeNum(day1, 21),
            getDateWithTimeNum(day1, 21.5),

            getDateWithTimeNum(day2, 9),
            getDateWithTimeNum(day2, 9.5),
            getDateWithTimeNum(day2, 10),
            getDateWithTimeNum(day2, 10.5),
            getDateWithTimeNum(day2, 19),
            getDateWithTimeNum(day2, 19.5),
            getDateWithTimeNum(day2, 20),
            getDateWithTimeNum(day2, 20.5),
            getDateWithTimeNum(day2, 21),
            getDateWithTimeNum(day2, 21.5),

            getDateWithTimeNum(day3, 9),
            getDateWithTimeNum(day3, 9.5),
            getDateWithTimeNum(day3, 10),
            getDateWithTimeNum(day3, 10.5),
            getDateWithTimeNum(day3, 11),
            getDateWithTimeNum(day3, 11.5),
            getDateWithTimeNum(day3, 12),
            getDateWithTimeNum(day3, 12.5),
            getDateWithTimeNum(day3, 17),
            getDateWithTimeNum(day3, 17.5),
            getDateWithTimeNum(day3, 18),
            getDateWithTimeNum(day3, 18.5),
            getDateWithTimeNum(day3, 19),
            getDateWithTimeNum(day3, 19.5),
            getDateWithTimeNum(day3, 20),
          ],
        },
      }

      for (const id of Object.keys(this.responses)) {
        const fixedAvailability = []
        for (const date of this.responses[id].availability) {
          fixedAvailability.push(date)
          const dateCopy = new Date(date)
          dateCopy.setMinutes(date.getMinutes() + 15)
          fixedAvailability.push(dateCopy)
        }
        this.responses[id].availability = fixedAvailability
      }
    },
    reset() {
      this.responses = {}
      this.calendarEventsByDay = []
    },
    playAnimation() {
      this.reset()

      this.$refs.scheduleOverlap.startEditing()
      setTimeout(() => {
        this.getCalendarEventsByDay()
        this.getResponses()
        setTimeout(() => {
          this.$refs.scheduleOverlap.setAvailabilityAutomatically()
          setTimeout(() => {
            this.$refs.scheduleOverlap.stopEditing()
          }, 2000)
        }, 500)
      }, 200)
    },
  },

  mounted() {
    this.playAnimation()
  },
}
</script>
