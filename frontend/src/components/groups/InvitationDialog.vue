<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-[35rem] tw-m-0 tw-max-h-full"
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
    persistent
  >
    <v-expand-transition>
      <v-card
        class="tw-overflow-none tw-relative tw-flex tw-flex-col tw-rounded-lg tw-px-2 tw-py-5 tw-transition-all"
      >
        <v-card-text>
          <div
            class="tw-mb-5 tw-text-wrap tw-text-xl tw-font-medium tw-text-black"
          >
            Accept invitation to share your calendars with "{{ group.name }}"?
          </div>

          <CalendarAccounts
            :sync-with-backend="false"
            :allow-add-calendar-account="false"
            :toggle-state="true"
            @toggleCalendarAccount="toggleCalendarAccount"
            @toggleSubCalendarAccount="toggleSubCalendarAccount"
          ></CalendarAccounts>

          <div class="tw-mb-2 tw-mt-5 tw-font-medium tw-text-black">
            These calendars will be shared with
          </div>
          <div>
            <UserChip
              v-for="user in group.attendees.filter((u) => !u.declined)"
              :key="user.email"
              :user="user"
            ></UserChip>
          </div>
        </v-card-text>

        <v-card-actions>
          <v-btn
            text
            class="tw-text-dark-gray tw-underline"
            @click="rejectInvitation"
            >Reject invitation</v-btn
          >
          <v-spacer />
          <v-btn
            class="tw-bg-green tw-px-5 tw-text-white tw-transition-opacity"
            @click="acceptInvitation"
            >Accept Invitation</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-expand-transition>
  </v-dialog>
</template>

<script>
import { mapState } from "vuex"
import { isPhone, post } from "@/utils"
import CalendarAccounts from "@/components/settings/CalendarAccounts.vue"
import UserChip from "@/components/general/UserChip.vue"

export default {
  name: "InvitationDialog",

  emits: ["input"],

  props: {
    value: { type: Boolean, required: true },
    group: { type: Object },
  },

  components: {
    CalendarAccounts,
    UserChip,
  },

  data: () => ({
    calendarAccounts: {},
  }),

  mounted() {
    this.calendarAccounts = this.authUser.calendarAccounts
  },

  computed: {
    ...mapState(["authUser"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    rejectInvitation() {
      post(`/events/${this.group._id}/decline`).then((res) => {
        console.log(res)
      })
    },

    acceptInvitation() {
      const payload = { useCalendarAvailability: true, enabledCalendars: {} }

      /** Determine which sub calendars are enabled - same code can be used for submitAvailability in scheduleOverlap.vue */
      for (const email in this.calendarAccounts) {
        if (this.calendarAccounts[email].enabled) {
          payload.enabledCalendars[email] = []
          for (const subCalendarId in this.calendarAccounts[email]
            .subCalendars) {
            if (
              this.calendarAccounts[email].subCalendars[subCalendarId].enabled
            ) {
              payload.enabledCalendars[email].push(subCalendarId)
            }
          }
        }
      }

      post(`/events/${this.group._id}/response`, payload).then((res) => {
        console.log(res)
        this.$emit("input", false)
      })
    },

    toggleCalendarAccount(payload) {
      this.calendarAccounts[payload.email].enabled = payload.enabled
    },

    toggleSubCalendarAccount(payload) {
      this.calendarAccounts[payload.email].subCalendars[
        payload.subCalendarId
      ].enabled = payload.enabled
    },
  },
}
</script>
