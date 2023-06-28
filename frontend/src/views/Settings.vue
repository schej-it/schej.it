<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12 tw-mt-5">
    <div class="tw-p-4 tw-flex tw-flex-col tw-gap-16">
      <!-- Calendar Access Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl sm:tw-text-2xl tw-font-medium tw-text-dark-green"
        >
          Calendar Access
        </div>
        <div class="tw-flex tw-flex-col sm:tw-flex-row tw-gap-5 sm:tw-gap-28">
          <div class="tw-text-black">
            We do not store your calendar data anywhere on our servers, and we
            only fetch your calendar events for the time frame you specify in
            order to display your calendar events while you fill out your
            availability.
          </div>
          <v-btn outlined class="tw-text-bright-red"
            >Revoke calendar access</v-btn
          >
        </div>
      </div>

      <!-- Permissions Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl sm:tw-text-2xl tw-font-medium tw-text-dark-green"
        >
          Permissions
        </div>
        <div
          class="tw-flex tw-flex-col tw-border-gray tw-border-t-[1px] tw-border-l-[1px]"
        >
          <div
            class="tw-flex tw-flex-row tw-w-full tw-border-gray tw-border-b-[1px]"
          >
            <div
              v-for="h in heading"
              class="tw-w-1/3 tw-border-gray tw-border-r-[1px] tw-p-4 tw-font-bold"
            >
              {{ h }}
            </div>
          </div>

          <div
            v-for="c in content"
            class="tw-flex tw-flex-row tw-w-full tw-border-gray tw-border-b-[1px]"
          >
            <div
              v-for="text in c"
              class="tw-w-1/3 tw-border-gray tw-border-r-[1px] tw-p-4"
            >
              {{ text }}
            </div>
          </div>
        </div>
      </div>

      <!-- Question Section -->
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div
          class="tw-text-xl sm:tw-text-2xl tw-font-medium tw-text-dark-green"
        >
          Have a question?
        </div>
        <div class="tw-flex tw-flex-col sm:tw-flex-row tw-gap-5 sm:tw-gap-28">
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
      <div class="tw-flex tw-flex-row tw-justify-center tw-mt-28">
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
              <v-card-text class="tw-text-dark-gray tw-text-sm"
                >Are you sure you want to delete your account? All your account
                data will be lost.</v-card-text
              >
              <div class="tw-mx-6">
                <div class="tw-text-dark-gray tw-text-sm">Type your email in the box below to confirm:</div>
              <v-text-field
                v-model="deleteValidateEmail"
                autofocus
                class="tw-text-white tw-flex-initial"
                :placeholder="authUser.email"
              />
            </div>
              <v-card-actions>
                <v-spacer />
                <v-btn text @click="deleteDialog = false">Cancel</v-btn>
                <v-btn text color="error" @click="deleteAccount()" :disabled="authUser.email != deleteValidateEmail">Delete</v-btn>
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

export default {
  name: "Settings",

  components: {},

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
    ...mapActions([]),
    deleteAccount() {
      _delete(`/user`)
        .catch((err) => {
          this.showError(
            "There was a problem deleting your account! Please try again later."
          )
        })
    }
  },

  created() {},
}
</script>
