import Vue from "vue"
import Vuex from "vuex"
import { get } from "@/utils"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    error: "",
    info: "",

    authUser: null,

    createdEvents: [],
    joinedEvents: [],

    // Feature flags
    groupsEnabled: true,
    signUpFormEnabled: false,
    daysOnlyEnabled: false,
    overlayAvailabilitiesEnabled: false,
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

    setGroupsEnabled(state, enabled) {
      state.groupsEnabled = enabled
    },
    setSignUpFormEnabled(state, enabled) {
      state.signUpFormEnabled = enabled
    },
    setDaysOnlyEnabled(state, enabled) {
      state.daysOnlyEnabled = enabled
    },
    setOverlayAvailabilitiesEnabled(state, enabled) {
      state.overlayAvailabilitiesEnabled = enabled
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

    async refreshAuthUser({ commit }) {
      const authUser = await get("/user/profile")
      commit("setAuthUser", authUser)
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
