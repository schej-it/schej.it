<template>
  <div
    :class="
      toggleState
        ? 'tw-min-w-[175px]'
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
        ><v-icon color="#4F4F4F">mdi-plus</v-icon></v-btn
      >
    </div>
    <div :class="toggleState ? '' : 'tw-px-4 tw-py-2'" class="">
      <div
        v-for="account in calendarAccounts"
        class="tw-group tw-flex tw-h-10 tw-flex-row tw-items-center tw-justify-between tw-text-black"
      >
        <div :class="`tw-gap-${toggleState ? '0' : '2'}`" class="tw-flex tw-w-full tw-flex-row tw-items-center">
          <v-checkbox
            v-if="toggleState"
            v-model="account.enabled"
            @change="(changed) => toggleCalendarAccount(account.email, changed)"
          />
          <v-avatar v-else size="24">
            <v-img :src="account.picture"></v-img
          ></v-avatar>
          <span
            class="tw-align-text-middle tw-inline-block tw-break-words tw-text-sm"
          >
            {{ account.email }}
          </span>
        </div>
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
    toggleState: { type: Boolean, default: false },
    eventId: { type: String, default: "" },
  },

  data: () => ({
    removeDialog: false,
    selectedRemoveEmail: "",
    selected: [],
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
    toggleCalendarAccount(email, enabled) {
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
          delete this.authUser.calendarAccounts[this.selectedRemoveEmail]
          this.setAuthUser(this.authUser)
          // const newAuthUser = await get("/user/profile")
          // this.setAuthUser(newAuthUser)
          this.removeDialog = false
        })
        .catch((err) => {
          console.log(err)
          this.showError(
            "There was a problem removing this account! Please try again later."
          )
        })
    },
  },
}
</script>
