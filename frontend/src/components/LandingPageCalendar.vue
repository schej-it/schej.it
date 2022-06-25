<template>
  <v-card class="tw-mx-4 md:tw-w-[42rem] ">
    <ScheduleOverlap
      :startDate="startDate"
      :endDate="endDate"
      :startTime="startTime"
      :endTime="endTime"
      
      :responses="responses"

      :calendarEvents="calendarEvents"

      showResponsesWithCalendarEvents
      calendarOnly
    />
  </v-card>
</template>

<script>
import ScheduleOverlap from '@/components/ScheduleOverlap'
import { getDateDayOffset, getDateWithTimeInt } from '@/utils'

export default {
  name: 'LandingPageCalendar',

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    startDate: getDateWithTimeInt(getDateDayOffset(new Date(), -1), 9),
    endDate: getDateWithTimeInt(getDateDayOffset(new Date(), 1), 22),
    startTime: 9,
    endTime: 22,

    calendarEvents: [],

    responses: {},
  }),

  created() {
    const day1 = this.startDate
    const day2 = new Date()
    const day3 = this.endDate

    this.calendarEvents = [
      {
        startDate: getDateWithTimeInt(day1, 9),
        endDate: getDateWithTimeInt(day1, 11),
        summary: 'Coffee with Jen',
      },
      {
        startDate: getDateWithTimeInt(day2, 11),
        endDate: getDateWithTimeInt(day2, 16),
        summary: 'Karaoke with friends',
      },
      {
        startDate: getDateWithTimeInt(day3, 13),
        endDate: getDateWithTimeInt(day3, 17),
        summary: 'Study session',
      },
      {
        startDate: getDateWithTimeInt(day3, 20.5),
        endDate: getDateWithTimeInt(day3, 22),
        summary: 'Hackathon meeting',
      },
    ]

    this.responses = {
      "62828fec1bc681fa020632f2": {
        user: {_id: '1', name: '1'},
        "availability": [
          getDateWithTimeInt(day1, 11),
          getDateWithTimeInt(day1, 11.5),
          getDateWithTimeInt(day1, 12),
          getDateWithTimeInt(day1, 12.5),
          getDateWithTimeInt(day1, 13),
          getDateWithTimeInt(day1, 14.5),
          getDateWithTimeInt(day1, 15),
          getDateWithTimeInt(day1, 15.5),
          getDateWithTimeInt(day1, 16),
          getDateWithTimeInt(day1, 16.5),
          getDateWithTimeInt(day1, 17),
          getDateWithTimeInt(day1, 17.5),
          getDateWithTimeInt(day1, 18),
          getDateWithTimeInt(day1, 18.5),
          getDateWithTimeInt(day1, 19),
          getDateWithTimeInt(day1, 20.5),
          getDateWithTimeInt(day1, 21),
          getDateWithTimeInt(day1, 21.5),
          
          getDateWithTimeInt(day2, 9),
          getDateWithTimeInt(day2, 9.5),
          getDateWithTimeInt(day2, 12),
          getDateWithTimeInt(day2, 12.5),
          getDateWithTimeInt(day2, 13),
          getDateWithTimeInt(day2, 13.5),
          getDateWithTimeInt(day2, 14),
          getDateWithTimeInt(day2, 14.5),
          getDateWithTimeInt(day2, 15),
          getDateWithTimeInt(day2, 15.5),
          getDateWithTimeInt(day2, 16),
          getDateWithTimeInt(day2, 16.5),
          getDateWithTimeInt(day2, 17),
          getDateWithTimeInt(day2, 17.5),
          getDateWithTimeInt(day2, 18),
          getDateWithTimeInt(day2, 18.5),
          getDateWithTimeInt(day2, 19),
          getDateWithTimeInt(day2, 19.5),
          getDateWithTimeInt(day2, 20),
          getDateWithTimeInt(day2, 21.5),
          
          getDateWithTimeInt(day3, 12.5),
          getDateWithTimeInt(day3, 13),
          getDateWithTimeInt(day3, 13.5),
          getDateWithTimeInt(day3, 14),
          getDateWithTimeInt(day3, 14.5),
          getDateWithTimeInt(day3, 15),
          getDateWithTimeInt(day3, 15.5),
          getDateWithTimeInt(day3, 16),
          getDateWithTimeInt(day3, 16.5),
          getDateWithTimeInt(day3, 17),
          getDateWithTimeInt(day3, 19.5),
        ]
      },
      "628292fe6e12d2baa9c01395": {
        user: {_id: '2', name: '2'},
        "availability": [
          getDateWithTimeInt(day1, 14),
          getDateWithTimeInt(day1, 14.5),
          getDateWithTimeInt(day1, 15),
          getDateWithTimeInt(day1, 15.5),
          getDateWithTimeInt(day1, 16),
          getDateWithTimeInt(day1, 16.5),
          getDateWithTimeInt(day1, 17),
          getDateWithTimeInt(day1, 17.5),
          getDateWithTimeInt(day1, 18),
          getDateWithTimeInt(day1, 18.5),
          getDateWithTimeInt(day1, 19),
          getDateWithTimeInt(day1, 19.5),
          getDateWithTimeInt(day1, 20),
          getDateWithTimeInt(day1, 20.5),
          getDateWithTimeInt(day1, 21),
          getDateWithTimeInt(day1, 21.5),

          getDateWithTimeInt(day2, 9),
          getDateWithTimeInt(day2, 9.5),
          getDateWithTimeInt(day2, 10),
          getDateWithTimeInt(day2, 10.5),
          getDateWithTimeInt(day2, 16),
          getDateWithTimeInt(day2, 16.5),
          getDateWithTimeInt(day2, 17),
          getDateWithTimeInt(day2, 17.5),
          getDateWithTimeInt(day2, 18),
          getDateWithTimeInt(day2, 18.5),
          getDateWithTimeInt(day2, 19),
          getDateWithTimeInt(day2, 19.5),
          getDateWithTimeInt(day2, 20),
          getDateWithTimeInt(day2, 20.5),
          getDateWithTimeInt(day2, 21),
          getDateWithTimeInt(day2, 21.5),

          getDateWithTimeInt(day3, 9),
          getDateWithTimeInt(day3, 9.5),
          getDateWithTimeInt(day3, 10),
          getDateWithTimeInt(day3, 12.5),
          getDateWithTimeInt(day3, 17),
          getDateWithTimeInt(day3, 18),
          getDateWithTimeInt(day3, 18.5),
          getDateWithTimeInt(day3, 19),
          getDateWithTimeInt(day3, 19.5),
          getDateWithTimeInt(day3, 20),
        ]
      },
      "628208870df4418ff4213757": {
        user: {_id: '3', name: '3'},
        "availability": [
          getDateWithTimeInt(day1, 11),
          getDateWithTimeInt(day1, 11.5),
          getDateWithTimeInt(day1, 12),
          getDateWithTimeInt(day1, 12.5),
          getDateWithTimeInt(day1, 13),
          getDateWithTimeInt(day1, 13.5),
          getDateWithTimeInt(day1, 14),
          getDateWithTimeInt(day1, 14.5),
          getDateWithTimeInt(day1, 15),
          getDateWithTimeInt(day1, 15.5),
          getDateWithTimeInt(day1, 16),
          getDateWithTimeInt(day1, 16.5),
          getDateWithTimeInt(day1, 20),
          getDateWithTimeInt(day1, 20.5),
          getDateWithTimeInt(day1, 21),
          getDateWithTimeInt(day1, 21.5),

          getDateWithTimeInt(day2, 9),
          getDateWithTimeInt(day2, 9.5),
          getDateWithTimeInt(day2, 10),
          getDateWithTimeInt(day2, 10.5),
          getDateWithTimeInt(day2, 11),
          getDateWithTimeInt(day2, 11.5),
          getDateWithTimeInt(day2, 19),
          getDateWithTimeInt(day2, 19.5),
          getDateWithTimeInt(day2, 20),
          getDateWithTimeInt(day2, 20.5),
          getDateWithTimeInt(day2, 21),
          getDateWithTimeInt(day2, 21.5),

          getDateWithTimeInt(day3, 9),
          getDateWithTimeInt(day3, 9.5),
          getDateWithTimeInt(day3, 10),
          getDateWithTimeInt(day3, 10.5),
          getDateWithTimeInt(day3, 11),
          getDateWithTimeInt(day3, 11.5),
          getDateWithTimeInt(day3, 12),
          getDateWithTimeInt(day3, 12.5),
          getDateWithTimeInt(day3, 17),
          getDateWithTimeInt(day3, 17.5),
          getDateWithTimeInt(day3, 18),
          getDateWithTimeInt(day3, 18.5),
          getDateWithTimeInt(day3, 19),
          getDateWithTimeInt(day3, 19.5),
          getDateWithTimeInt(day3, 20),

        ],
      }
    }
  },
}
</script>