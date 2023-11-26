<template>
      <div
        v-if="showAccount(account)"
        class="tw-group tw-flex tw-h-10 tw-flex-row tw-items-center tw-justify-between tw-text-black"
      >
        <div
          :class="`tw-gap-${toggleState ? '0' : '2'}`"
          class="tw-flex tw-w-full tw-flex-row tw-items-center"
        >
          <div v-if="toggleState" class="tw-flex">
            <v-checkbox
              v-model="account.enabled"
              @change="
                (changed) => toggleCalendarAccount(account.email, changed)
              "
            />

            <v-icon class="-tw-ml-2 tw-text-dark-gray">mdi-chevron-right</v-icon>
          </div>
          <v-avatar v-else size="24">
            <v-img :src="account.picture"></v-img
          ></v-avatar>
          <div
            :class="toggleState ? 'tw-w-[180px]' : ''"
            class="tw-align-text-middle tw-inline-block tw-break-words tw-text-sm"
          >
            {{ account.email }}
            </div>
          <v-tooltip top v-if="accountHasError(account)">
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                icon
                v-bind="attrs"
                v-on="on"
                @click="() => reauthenticateCalendarAccount(account)"
              >
                <v-icon>mdi-alert-circle</v-icon>
              </v-btn>
            </template>
            <span>Sign in again</span>
          </v-tooltip>
        </div>
        <!-- Needed to make sure tailwind classes compile -->
        <span class="tw-hidden tw-opacity-0 tw-opacity-100"></span>
        <v-btn
          icon
          :class="`tw-opacity-${
            account.email == selectedRemoveEmail && removeDialog ? '100' : '0'
          } ${
            account.email == authUser.email || toggleState ? 'tw-hidden' : ''
          }`"
          class="group-hover:tw-opacity-100"
          @click="() => openRemoveDialog(account.email)"
          ><v-icon color="#4F4F4F">mdi-close</v-icon></v-btn
        >
      </div>
</template>

<script>
import { mapState, mapActions, mapMutations } from "vuex"
import { authTypes } from "@/constants"
import { get, post, _delete, signInGoogle } from "@/utils"

export default {
  name: "CalendarAccount",

  props: {
    toggleState: { type: Boolean, default: false },
    account: { type: Object, default: () => {} },
    eventId: { type: String, default: "" },
    openRemoveDialog: { type: Function },
    calendarEventsMap: { type: Object, default: () => {} }, // Object of different users' calendar events
    removeDialog: {type: Boolean, default: false},
    selectedRemoveEmail: {type: String, default: ""},
  },

  data: () => ({
    calendarEventsMapCopy: null,
  }),

  computed: {
    ...mapState(["authUser"]),
  },

  methods: {
    ...mapActions(["showError"]),
    accountHasError(account) {
      return (
        this.calendarEventsMapCopy &&
        this.calendarEventsMapCopy[account.email]?.error
      )
    },
    /** don't show account if in toggle state and account has an error */
    showAccount(account) {
      return !(this.toggleState && this.accountHasError(account))
    },
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
    reauthenticateCalendarAccount(account) {
      signInGoogle({
        state: {
          type: this.toggleState
            ? authTypes.ADD_CALENDAR_ACCOUNT_FROM_EDIT
            : authTypes.ADD_CALENDAR_ACCOUNT,
          eventId: this.eventId,
        },
        requestCalendarPermission: true,
        selectAccount: false,
        loginHint: account.email,
      })
    },
    toggleCalendarAccount(email, enabled) {
      post(`/user/toggle-calendar`, { email, enabled }).catch((err) => {
        this.showError(
          "There was a problem with toggling your calendar account! Please try again later."
        )
      })
    },
  },

  watch: {
    calendarEventsMapCopy: {
      immediate: true,
      async handler() {
        // Do a test request to calendarevents route to check if calendar access is allowed for each account
        if (!this.calendarEventsMapCopy) {
          try {
            this.calendarEventsMapCopy = await get(
              `/user/calendars?timeMin=${new Date().toISOString()}&timeMax=${new Date().toISOString()}`
            )
          } catch (err) {
            console.error(err)
          }
        }
      },
    },
  },

  created() {
    this.calendarEventsMapCopy = this.calendarEventsMap
  },
}
</script>

