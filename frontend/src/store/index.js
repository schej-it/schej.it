import Vue from 'vue'
import Vuex from 'vuex'
import { getDateWithTime, getDateDayOffset } from '@/utils'

Vue.use(Vuex)

const date1 = new Date()
const date2 = getDateDayOffset(new Date(), 1)
const date3 = getDateDayOffset(new Date(), 2)

export default new Vuex.Store({
  state: {
    error: '',
    info: '',

    authUser: null,
  },
  getters: {
  },
  mutations: {
    setError(state, error) {
      state.error = error
    },
    setInfo(state, info) {
      state.info = info
    },

    setAuthUser(state, authUser) {
      state.authUser = authUser
    },
  },
  actions: {
    // Error & info
    showError({ commit }, error) {
      commit('setError', '')
      setTimeout(() => commit('setError', error), 0)
    },
    showInfo({ commit }, info) {
      commit('setInfo', '')
      setTimeout(() => commit('setInfo', info), 0)
    },
  },
  modules: {
  }
})
