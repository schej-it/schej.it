import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    events: {
      'asdf': {
        name: 'Meeting #1',
        startDate: new Date(),
        endDate: new Date(new Date().getTime() + 2 * 24*60*60*1000),
        respondents: 4,
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
