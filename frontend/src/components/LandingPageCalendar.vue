<template>
  <v-card 
    class="tw-mx-4 md:tw-w-[42rem] tw-rounded-lg"
    style="box-shadow: 0px 0px 10px 3px rgba(0, 0, 0, 0.1);"
  >
    <ScheduleOverlap
      :startDate="startDate"
      :endDate="endDate"
      :startTime="startTime"
      :endTime="endTime"
      
      :responses="responses"

      :calendarEvents="calendarEvents"
      :initialShowCalendarEvents="true"

      calendarOnly
      :interactable="false"
    />
  </v-card>
</template>

<script>
import ScheduleOverlap from '@/components/ScheduleOverlap'
import { getDateDayOffset, getDateWithTimeNum, isPhone } from '@/utils'

export default {
  name: 'LandingPageCalendar',

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    startDate: getDateWithTimeNum(getDateDayOffset(new Date(), -1), 9),
    endDate: getDateWithTimeNum(getDateDayOffset(new Date(), 1), 22),
    startTime: 9,

    responses: {},
  }),

  computed: {
    endTime() {
      return isPhone(this.$vuetify) ? 17 : 22
    },
    calendarEvents() {
      const day1 = this.startDate
      const day2 = new Date()
      const day3 = this.endDate
      const events = [
        {
          startDate: getDateWithTimeNum(day1, 9),
          endDate: getDateWithTimeNum(day1, 11),
          summary: 'Coffee with Jen',
        },
        {
          startDate: getDateWithTimeNum(day2, 11),
          endDate: getDateWithTimeNum(day2, 16),
          summary: 'Karaoke with friends',
        },
        {
          startDate: getDateWithTimeNum(day3, 13),
          endDate: getDateWithTimeNum(day3, 17),
          summary: 'Study session',
        },
      ]

      if (!isPhone(this.$vuetify)) { 
        events.push({
          startDate: getDateWithTimeNum(day3, 20.5),
          endDate: getDateWithTimeNum(day3, 22),
          summary: 'Hackathon meeting',
        })
      }

      return events
    },
  },

  created() {
    const day1 = this.startDate
    const day2 = new Date()
    const day3 = this.endDate


    this.responses = {
      "62828fec1bc681fa020632f2": {
        user: {_id: '1', name: '1'},
        "availability": [
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
        ]
      },
      "628292fe6e12d2baa9c01395": {
        user: {_id: '2', name: '2'},
        "availability": [
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
        ]
      },
      "628208870df4418ff4213757": {
        user: {_id: '3', name: '3'},
        "availability": [
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
      }
    }
  },
}
</script>