import Vue from "vue"
import Vuex from "vuex"
import { get } from "@/utils"
import { jitsuAnalytics } from "@jitsu/js"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    error: "",
    info: "",

    authUser: null,

    createdEvents: [],
    joinedEvents: [],

    analytics: jitsuAnalytics({
      host: "https://data.schej.it",
      // Browser Write Key configured on Jitsu Site entity.
      // If no Browser Write Key is added for Site entity, Site ID value can be used a Write Key.
      // On Jitsu.Cloud can be omitted if Site has explicitly mapped domain name that is used in host parameter
      writeKey: process.env.JITSU_WRITE_KEY,
    }),
  },
  getters: {},
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

    setCreatedEvents(state, createdEvents) {
      state.createdEvents = createdEvents
    },
    setJoinedEvents(state, joinedEvents) {
      state.joinedEvents = joinedEvents
    },
  },
  actions: {
    // Error & info
    showError({ commit }, error) {
      commit("setError", "")
      setTimeout(() => commit("setError", error), 0)
    },
    showInfo({ commit }, info) {
      commit("setInfo", "")
      setTimeout(() => commit("setInfo", info), 0)
    },

    // Events
    getEvents({ commit, dispatch }) {
      if (this.state.authUser) {
        return get("/user/events")
          .then((data) => {
            commit("setCreatedEvents", data.events)
            commit("setJoinedEvents", data.joinedEvents)
          })
          .catch((err) => {
            dispatch("showError", "There was a problem fetching events!")
          })
      } else {
        return null
      }
    },
  },
  modules: {},
})
