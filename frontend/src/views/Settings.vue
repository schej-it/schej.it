<template>
  <div class="tw-mx-auto tw-mb-12 tw-mt-5 tw-max-w-6xl">
    <div class="tw-flex tw-flex-col tw-gap-16 tw-p-4">
      <!-- Calendar Access Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
        >
          Calendar access
        </div>
        <div class="tw-flex tw-flex-col tw-gap-5 sm:tw-flex-row sm:tw-gap-28">
          <div class="tw-text-black">
            We do not store your calendar data anywhere on our servers, and we
            only fetch your calendar events for the time frame you specify in
            order to display your calendar events while you fill out your
            availability.
          </div>
          <v-btn
            outlined
            class="tw-text-bright-red"
            href="https://myaccount.google.com/connections?filters=3,4&hl=en"
            target="_blank"
            >Revoke calendar access</v-btn
          >
        </div>
        <CalendarAccounts></CalendarAccounts>
      </div>

      <!-- Permissions Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
        >
          Permissions
        </div>
        <div
          class="tw-flex tw-flex-col tw-rounded-md tw-border-[1px] tw-border-light-gray-stroke"
        >
          <div
            class="tw-flex tw-w-full tw-flex-row tw-border-b-[1px] tw-border-light-gray-stroke"
          >
            <div
              v-for="(h, i) in heading"
              :class="`tw-border-r-[${i == heading.length - 1 ? '0' : '1'}px]`"
              class="tw-w-1/3 tw-border-light-gray-stroke tw-p-4 tw-font-bold"
            >
              {{ h }}
            </div>
          </div>

          <div
            v-for="(c, j) in content"
            :class="`tw-border-b-[${j == content.length - 1 ? '0' : '1'}px]`"
            class="tw-flex tw-w-full tw-flex-row tw-border-light-gray-stroke"
          >
            <div
              v-for="(text, k) in c"
              :class="`tw-border-r-[${k == c.length - 1 ? '0' : '1'}px]`"
              class="tw-w-1/3 tw-border-light-gray-stroke tw-p-4"
            >
              {{ text }}
            </div>
          </div>
        </div>
      </div>

      <!-- Question Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
        >
          Have a question?
        </div>
        <div class="tw-flex tw-flex-col tw-gap-5 sm:tw-flex-row sm:tw-gap-28">
          <div class="tw-text-black">
            Email us at
            <a
              href="mailto:schej.team@gmail.com"
              class="tw-text-black tw-underline"
              >schej.team@gmail.com</a
            >
            with any questions!
          </div>
        </div>
      </div>

      <!-- Delete Account Section -->
      <div class="tw-mt-28 tw-flex tw-flex-row tw-justify-center">
        <div class="tw-w-64">
          <v-dialog v-model="deleteDialog" width="400" persistent>
            <template v-slot:activator="{ on, attrs }">
              <v-btn
                outlined
                class="tw-text-bright-red"
                block
                v-bind="attrs"
                v-on="on"
                >Delete account</v-btn
              >
            </template>
            <v-card>
              <v-card-title>Are you sure?</v-card-title>
              <v-card-text class="tw-text-sm tw-text-dark-gray"
                >Are you sure you want to delete your account? All your account
                data will be lost.</v-card-text
              >
              <div class="tw-mx-6">
                <div class="tw-text-sm tw-text-dark-gray">
                  Type your email in the box below to confirm:
                </div>
                <v-text-field
                  v-model="deleteValidateEmail"
                  autofocus
                  class="tw-flex-initial tw-text-white"
                  :placeholder="authUser.email"
                />
              </div>
              <v-card-actions>
                <v-spacer />
                <v-btn text @click="deleteDialog = false">Cancel</v-btn>
                <v-btn
                  text
                  color="error"
                  @click="deleteAccount()"
                  :disabled="authUser.email != deleteValidateEmail"
                  >Delete</v-btn
                >
              </v-card-actions>
            </v-card>
          </v-dialog>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex"
import { _delete } from "@/utils"
import CalendarAccounts from "@/components/settings/CalendarAccounts.vue"

export default {
  name: "Settings",

  components: { CalendarAccounts },

  data: () => ({
    dialog: false,
    deleteDialog: false,
    deleteValidateEmail: "",
    heading: ["Permission", "Purpose", "Requested When"],
    content: [
      [
        "View all calendars subscribed to",
        "Allows us to display calendar events on all your calendars instead of just your primary calendar",
        "User tries to input availability automatically with Google Calendar",
      ],
      [
        "View all calendar events",
        "Allows us to view the names/times of your calendar events",
        "User tries to input availability automatically with Google Calendar",
      ],
    ],
  }),

  computed: {
    ...mapState(["authUser"]),
  },

  methods: {
    ...mapActions(["showError"]),
    deleteAccount() {
      _delete(`/user`)
        .then(() => {
          window.location.reload()
        })
        .catch((err) => {
          this.showError(
            "There was a problem deleting your account! Please try again later."
          )
        })
    },
  },

  created() {},
}
</script>
