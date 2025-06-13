import Vue from "vue"
import Vuex from "vuex"
import { get, isPremiumUser } from "@/utils"

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    error: "",
    info: "",

    authUser: null,

    createdEvents: [],
    joinedEvents: [],
    folders: [],

    featureFlagsLoaded: false,

    // Feature flags
    groupsEnabled: true,
    signUpFormEnabled: false,
    daysOnlyEnabled: true,
    overlayAvailabilitiesEnabled: true,
    enablePaywall: false,

    // Experiments
    pricingPageConversion: "control",

    // Upgrade dialog
    upgradeDialogVisible: false,
    upgradeDialogType: null,
    upgradeDialogData: null,
  },
  getters: {
    isPremiumUser(state) {
      return isPremiumUser(state.authUser)
    },
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

    setCreatedEvents(state, createdEvents) {
      state.createdEvents = createdEvents
    },
    setJoinedEvents(state, joinedEvents) {
      state.joinedEvents = joinedEvents
    },
    setFolders(state, folders) {
      state.folders = folders
    },

    setFeatureFlagsLoaded(state, loaded) {
      state.featureFlagsLoaded = loaded
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
    setPricingPageConversion(state, conversion) {
      state.pricingPageConversion = conversion
    },
    setEnablePaywall(state, enabled) {
      state.enablePaywall = enabled
    },
    setUpgradeDialogVisible(state, visible) {
      state.upgradeDialogVisible = visible
    },
    setUpgradeDialogType(state, type) {
      state.upgradeDialogType = type
    },
    setUpgradeDialogData(state, data) {
      state.upgradeDialogData = data
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
        return Promise.allSettled([get("/user/folders"), get("/user/events")])
          .then(([folders, events]) => {
            if (
              folders.status === "fulfilled" &&
              events.status === "fulfilled"
            ) {
              commit("setFolders", folders.value)
              commit("setCreatedEvents", events.value.events)
              commit("setJoinedEvents", events.value.joinedEvents)
            } else {
              dispatch("showError", "There was a problem fetching events!")
              console.error(folders.reason, events.reason)
            }
          })
          .catch((err) => {
            dispatch("showError", "There was a problem fetching events!")
            console.error(err)
          })
      } else {
        return null
      }
    },
    showUpgradeDialog({ commit }, { type, data = null }) {
      commit("setUpgradeDialogVisible", true)
      commit("setUpgradeDialogType", type)
      commit("setUpgradeDialogData", data)
    },
    hideUpgradeDialog({ commit }) {
      commit("setUpgradeDialogVisible", false)
      commit("setUpgradeDialogType", null)
      commit("setUpgradeDialogData", null)
    },
  },
  modules: {},
})
