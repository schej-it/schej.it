<template>
  <div
    :class="
      toggleState
        ? 'tw-min-w-[240px]'
        : 'tw-w-fit tw-min-w-[288px] tw-drop-shadow'
    "
    class="tw-flex tw-flex-col tw-rounded-lg tw-bg-white tw-text-black tw-transition-all"
  >
    <div
      :class="toggleState ? '' : 'tw-border-b-[1px] tw-px-4 tw-pt-1'"
      class="tw-flex tw-flex-row tw-items-center tw-justify-between tw-border-off-white tw-pb-1 tw-align-middle"
    >
      <div class="tw-font-medium">My calendars</div>
      <v-btn @click="addCalendarAccount" icon
        ><v-icon class="tw-text-very-dark-gray">mdi-plus</v-icon></v-btn
      >
    </div>
    <div :class="toggleState ? '' : 'tw-px-4 tw-py-2'" class="">
      <CalendarAccount
        v-for="account in calendarAccounts"
        :key="account.email"
        :toggleState="toggleState"
        :account="account"
        :eventId="eventId"
        :openRemoveDialog="openRemoveDialog"
        :calendarEventsMap="calendarEventsMap"
        :removeDialog="removeDialog"
        :selectedRemoveEmail="selectedRemoveEmail"
      ></CalendarAccount>
    </div>
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
import { get, post, _delete, signInGoogle } from "@/utils"
import CalendarAccount from "@/components/settings/CalendarAccount.vue"

export default {
  name: "CalendarAccounts",

  props: {
    toggleState: { type: Boolean, default: false },
    eventId: { type: String, default: "" },
    calendarEventsMap: { type: Object, default: () => {} }, // Object of different users' calendar events
  },

  data: () => ({
    removeDialog: false,
    selectedRemoveEmail: "",
  }),

  computed: {
    ...mapState(["authUser"]),
    calendarAccounts() {
      return this.authUser.calendarAccounts
    },
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
  },

  components: {
    CalendarAccount,
  },
}
</script>
