<template>
  <div
    class="tw-flex tw-w-fit tw-min-w-[288px] tw-flex-col tw-rounded-lg tw-bg-white tw-text-black tw-drop-shadow tw-transition-all"
  >
    <div
      class="tw-flex tw-flex-row tw-items-center tw-justify-between tw-border-b-[1px] tw-border-off-white tw-py-1 tw-px-4 tw-align-middle"
    >
      <div>My calendars</div>
      <v-btn @click="addCalendarAccount" icon
        ><v-icon color="#4F4F4F">mdi-plus</v-icon></v-btn
      >
    </div>
    <div class="tw-py-2 tw-px-4">
      <div
        v-for="account in calendarAccounts"
        class="tw-group tw-flex tw-flex-row tw-items-center tw-justify-between tw-text-black tw-h-10"
      >
        <div class="tw-flex tw-flex-row tw-items-center tw-gap-2">
          <v-avatar size="24"> <v-img :src="account.picture"></v-img></v-avatar>
          <div class="tw-text-sm">
            {{ account.email }}
          </div>
        </div>
        <v-btn
          icon
          :class="`tw-opacity-${
            account.email == selectedRemoveEmail && removeDialog ? '100' : '0'
          } ${account.email == authUser.email ? 'tw-hidden' : ''}`"
          class="group-hover:tw-opacity-100"
          @click="() => openRemoveDialog(account.email)"
          ><v-icon color="#4F4F4F">mdi-close</v-icon></v-btn
        >
      </div>
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
import { post, _delete, signInGoogle } from "@/utils"

export default {
  name: "CalendarAccounts",

  props: {
    answer: { type: String },
    points: { type: Array },
  },

  data: () => ({
    removeDialog: false,
    selectedRemoveEmail: "",
  }),

  computed: {
    ...mapState(["authUser"]),
    calendarAccounts() {
      return [{email: this.authUser.email, picture: this.authUser.picture}].concat(this.authUser.calendarAccounts)
    }
  },

  methods: {
    ...mapActions(["showError"]),
    ...mapMutations(["setAuthUser"]),
    addCalendarAccount() {
      signInGoogle({
        state: { type: authTypes.ADD_CALENDAR_ACCOUNT },
        requestCalendarPermission: true,
        selectAccount: true,
      })
    },
    toggleCalendarAccount(email, enabled) {
      this.authUser.calendarAccounts.find(
        (account) => account.email == email
      ).enabled = enabled

      post(`/user/toggle-calendar`, { email, enabled }).catch((err) => {
        this.showError(
          "There was a problem with toggling your calendar account! Please try again later."
        )
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
          this.authUser.calendarAccounts = this.authUser.calendarAccounts.filter(
            (account) => account.email != this.selectedRemoveEmail
          )
          this.setAuthUser(this.authUser)
          // const newAuthUser = await get("/user/profile")
          // this.setAuthUser(newAuthUser)
          this.removeDialog = false
        })
        .catch((err) => {
          this.showError(
            "There was a problem removing this account! Please try again later."
          )
        })
    },
  },
}
</script>
