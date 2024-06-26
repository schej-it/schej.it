<template>
  <div
    :class="toggleState ? '' : 'tw-w-fit tw-min-w-[288px] tw-drop-shadow'"
    class="tw-flex tw-flex-col tw-rounded-lg tw-bg-white tw-text-black tw-transition-all"
  >
    <v-btn
      v-if="toggleState"
      class="-tw-ml-2 tw-w-[calc(100%+1rem)] tw-justify-between tw-px-2"
      block
      text
      @click="toggleShowCalendars"
    >
      <span class="tw-mr-1 tw-text-base tw-font-medium">My calendars</span>
      <v-icon :class="`tw-rotate-${showCalendars ? '180' : '0'}`"
        >mdi-chevron-down</v-icon
      ></v-btn
    >
    <div
      v-else
      class="tw-border-b tw-border-off-white tw-px-4 tw-py-3 tw-font-medium"
    >
      My calendars
    </div>
    <v-expand-transition>
      <span v-if="showCalendars">
        <div :class="toggleState ? '' : 'tw-px-4 tw-py-2'">
          <CalendarAccount
            v-for="account in calendarAccounts"
            :key="account.email"
            :syncWithBackend="syncWithBackend"
            :toggleState="toggleState"
            :account="account"
            :eventId="eventId"
            :openRemoveDialog="openRemoveDialog"
            :calendarEventsMap="calendarEventsMap"
            :removeDialog="removeDialog"
            :selectedRemoveEmail="selectedRemoveEmail"
            :fillSpace="fillSpace"
            @toggleCalendarAccount="
              (payload) => $emit('toggleCalendarAccount', payload)
            "
            @toggleSubCalendarAccount="
              (payload) => $emit('toggleSubCalendarAccount', payload)
            "
          ></CalendarAccount>
        </div>
        <v-btn
          v-if="allowAddCalendarAccount"
          text
          color="primary"
          :class="
            toggleState ? '-tw-ml-2 tw-mt-0 tw-w-min tw-px-2' : 'tw-w-full'
          "
          @click="addCalendarAccount"
          >+ Add calendar</v-btn
        >
      </span>
    </v-expand-transition>
    <v-dialog v-model="removeDialog" width="500" persistent>
      <v-card>
        <v-card-title>Are you sure?</v-card-title>
        <v-card-text class="tw-text-sm tw-text-dark-gray"
          >Are you sure you want to remove
          {{ selectedRemoveEmail }}?</v-card-text
        >
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="removeDialog = false">Cancel</v-btn>
          <v-btn text color="error" @click="removeAccount">Remove</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapState, mapActions, mapMutations } from "vuex"
import { authTypes } from "@/constants"
import { _delete, signInGoogle } from "@/utils"
import CalendarAccount from "@/components/settings/CalendarAccount.vue"

export default {
  name: "CalendarAccounts",

  props: {
    toggleState: { type: Boolean, default: false }, // Whether to allow user to toggle calendar accounts
    eventId: { type: String, default: "" },
    calendarEventsMap: { type: Object, default: () => {} }, // Object of different users' calendar events
    syncWithBackend: { type: Boolean, default: true }, // Whether toggling calendar accounts also updates the backend
    allowAddCalendarAccount: { type: Boolean, default: true }, // Whether to allow user to add a new calendar account
    initialCalendarAccountsData: { type: Object, default: () => {} }, // Initial data to display for enabled calendar accounts
    fillSpace: { type: Boolean, default: false }, // Whether to fill the available space up
  },

  data: () => ({
    removeDialog: false,
    selectedRemoveEmail: "",
    calendarAccounts: {},
    showCalendars:
      localStorage["showCalendars"] == undefined
        ? true
        : localStorage["showCalendars"] == "true",
  }),

  computed: {
    ...mapState(["authUser"]),
  },

  mounted() {
    this.calendarAccounts = !this.initialCalendarAccountsData
      ? this.authUser.calendarAccounts
      : this.initialCalendarAccountsData
  },

  methods: {
    ...mapActions(["showError"]),
    ...mapMutations(["setAuthUser"]),
    addCalendarAccount() {
      signInGoogle({
        state: {
          type: this.toggleState
            ? authTypes.ADD_CALENDAR_ACCOUNT_FROM_EDIT
            : authTypes.ADD_CALENDAR_ACCOUNT,
          eventId: this.eventId,
        },
        requestCalendarPermission: true,
        selectAccount: true,
      })
    },
    openRemoveDialog(email) {
      this.removeDialog = true
      this.selectedRemoveEmail = email
    },
    removeAccount() {
      _delete(`/user/remove-calendar-account`, {
        email: this.selectedRemoveEmail,
      })
        .then(async () => {
          // TODO: investigate into what the standard method is best
          delete this.authUser.calendarAccounts[this.selectedRemoveEmail]
          this.setAuthUser(this.authUser)
          // const newAuthUser = await get("/user/profile")
          // this.setAuthUser(newAuthUser)
          this.removeDialog = false
        })
        .catch((err) => {
          console.error(err)
          this.showError(
            "There was a problem removing this account! Please try again later."
          )
        })
    },
    toggleShowCalendars() {
      this.showCalendars = !this.showCalendars
      localStorage["showCalendars"] = this.showCalendars
    },
  },

  components: {
    CalendarAccount,
  },
}
</script>
