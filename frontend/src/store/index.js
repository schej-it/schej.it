import Vue from 'vue'
import Vuex from 'vuex'
import { getDateWithTime, getDateDayOffset } from '@/utils'

Vue.use(Vuex)

const date1 = new Date()
const date2 = getDateDayOffset(new Date(), 1)
const date3 = getDateDayOffset(new Date(), 2)

export default new Vuex.Store({
  state: {
    events: {
      'asdf': {
        name: 'Meeting #1',
        startDate: getDateWithTime(new Date(), '00:00'),
        endDate: getDateWithTime(getDateDayOffset(new Date(), 2), '00:00'),
        startTime: 9,
        endTime: 20,
        responses: [
          {
            name: 'arjun',
            times: [
              getDateWithTime(date1, '9:00'),
              getDateWithTime(date1, '9:30'),
              getDateWithTime(date1, '10:00'),
              getDateWithTime(date1, '11:00'),
              getDateWithTime(date1, '17:30'),
              getDateWithTime(date1, '18:00'),
              getDateWithTime(date1, '18:30'),
              getDateWithTime(date1, '19:00'),
            ],
          },
        ],
      },
    },
  },
  getters: {
  },
  mutations: {
  },
  actions: {
  },
  modules: {
  }
})
