<template>
  <div v-if="showAccount()" class="tw-flex tw-flex-col">
    <div
      class="tw-group tw-flex tw-h-10 tw-flex-row tw-items-center tw-justify-between tw-text-black"
    >
      <div
        :class="`tw-gap-${toggleState ? '0' : '2'}`"
        class="tw-flex tw-w-full tw-flex-row tw-items-center"
      >
        <div v-if="toggleState" class="tw-flex tw-items-center">
          <v-checkbox
            v-model="account.enabled"
            @change="(enabled) => toggleCalendarAccount(enabled)"
            class="-tw-mb-[3px]"
          />
          <div
            class="-tw-ml-2 tw-h-fit tw-w-fit tw-cursor-pointer"
            @click="
              () => {
                showSubCalendars = !showSubCalendars
              }
            "
          >
            <!-- Make sure tailwind classes are compiled -->
            <div class="tw-rotate-0 tw-rotate-90"></div>

            <v-icon
              :class="`tw-rotate-${showSubCalendars ? 90 : 0}`"
              class="tw-text-dark-gray tw-transition-all"
              >mdi-chevron-right</v-icon
            >
          </div>
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
        <v-tooltip top v-if="accountHasError()">
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              icon
              v-bind="attrs"
              v-on="on"
              @click="reauthenticateCalendarAccount"
            >
              <v-icon>mdi-alert-circle</v-icon>
            </v-btn>
          </template>
          <span>Sign in again</span>
        </v-tooltip>
      </div>
      <!-- Needed to make sure tailwind classes compile -->
      <span class="tw-hidden tw-opacity-0 tw-opacity-100"></span>

      <!-- Delete account button -->
      <v-btn
        icon
        :class="`tw-opacity-${
          account.email == selectedRemoveEmail && removeDialog ? '100' : '0'
        } ${account.email == authUser.email || toggleState ? 'tw-hidden' : ''}`"
        class="group-hover:tw-opacity-100"
        @click="() => openRemoveDialog(account.email)"
        ><v-icon color="#4F4F4F">mdi-close</v-icon></v-btn
      >
    </div>

    <!-- Sub-calendar accounts -->

    <v-expand-transition>
      <div v-if="showSubCalendars" class="tw-bg-[#EBF7EF]">
        <div
          v-for="(subCalendar, id) in account.subCalendars"
          :key="id"
          class="tw-flex tw-flex-row tw-items-center"
        >
          <v-checkbox
            v-model="subCalendar.enabled"
            @change="(enabled) => toggleSubCalendarAccount(enabled, id)"
            class="tw-h-5 tw-items-center"
          />
          <div
            class="tw-align-text-middle tw-ml-8 tw-inline-block tw-w-40 tw-break-words tw-text-sm"
          >
            {{ subCalendar.name }}
          </div>
        </div>
      </div>
    </v-expand-transition>
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
    removeDialog: { type: Boolean, default: false },
    selectedRemoveEmail: { type: String, default: "" },
    syncWithBackend: { type: Boolean, default: true },
  },

  data: () => ({
    showSubCalendars: false,
    calendarEventsMapCopy: null,
  }),

  computed: {
    ...mapState(["authUser"]),
  },

  methods: {
    ...mapActions(["showError"]),
    accountHasError() {
      return (
        this.calendarEventsMapCopy &&
        this.calendarEventsMapCopy[this.account.email]?.error
      )
    },
    /** don't show account if in toggle state and account has an error */
    showAccount() {
      return !(this.toggleState && this.accountHasError(this.account))
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
    reauthenticateCalendarAccount() {
      signInGoogle({
        state: {
          type: this.toggleState
            ? authTypes.ADD_CALENDAR_ACCOUNT_FROM_EDIT
            : authTypes.ADD_CALENDAR_ACCOUNT,
          eventId: this.eventId,
        },
        requestCalendarPermission: true,
        selectAccount: false,
        loginHint: this.account.email,
      })
    },
    toggleSubCalendarAccount(enabled, subCalendarId) {
      if (this.syncWithBackend) {
        post(`/user/toggle-sub-calendar`, {
          email: this.account.email,
          enabled,
          subCalendarId,
        }).catch((err) => {
          this.showError(
            "There was a problem with toggling your calendar account! Please try again later."
          )
        })
      } else {
        this.$emit("toggleSubCalendarAccount", {
          email: this.account.email,
          enabled,
          subCalendarId,
        })
      }
    },
    toggleCalendarAccount(enabled) {
      if (!enabled) this.showSubCalendars = false

      if (this.syncWithBackend) {
        post(`/user/toggle-calendar`, {
          email: this.account.email,
          enabled,
        }).catch((err) => {
          this.showError(
            "There was a problem with toggling your calendar account! Please try again later."
          )
        })
      } else {
        this.$emit("toggleCalendarAccount", {
          email: this.account.email,
          enabled,
        })
      }
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
