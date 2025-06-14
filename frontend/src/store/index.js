import Vue from "vue"
import Vuex from "vuex"
import { numFreeEvents, upgradeDialogTypes } from "@/constants"
import { get, isPremiumUser } from "@/utils"
import { createFolder, deleteFolder } from "../utils/FolderClient"

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

    // New dialog
    newDialogOptions: {
      show: false,
      contactsPayload: {},
      openNewGroup: false,
      eventOnly: false,
    },
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

    addFolder(state, folder) {
      state.folders.push(folder)
    },
    removeFolder(state, folderId) {
      state.folders = state.folders.filter((f) => f._id !== folderId)
    },

    setNewDialogOptions(
      state,
      {
        show = false,
        contactsPayload = {},
        openNewGroup = false,
        eventOnly = true,
      }
    ) {
      state.newDialogOptions = {
        show,
        contactsPayload,
        openNewGroup,
        eventOnly,
      }
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

    createNew({ state, getters, commit, dispatch }, { eventOnly = false }) {
      if (
        state.enablePaywall &&
        !getters.isPremiumUser &&
        state.authUser?.numEventsCreated >= numFreeEvents
      ) {
        dispatch("showUpgradeDialog", {
          type: upgradeDialogTypes.CREATE_EVENT,
        })
        return
      }

      commit("setNewDialogOptions", {
        show: true,
        contactsPayload: {},
        openNewGroup: false,
        eventOnly: eventOnly,
      })
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
    async createFolder({ commit, dispatch }, { name, color }) {
      try {
        const folder = await createFolder(name, color)
        commit("addFolder", folder)
      } catch (err) {
        dispatch("showError", "There was a problem creating the folder!")
        console.error(err)
      }
    },
    async deleteFolder({ commit, dispatch }, folderId) {
      try {
        await deleteFolder(folderId)
        commit("removeFolder", folderId)
      } catch (err) {
        dispatch("showError", "There was a problem deleting the folder!")
        console.error(err)
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
