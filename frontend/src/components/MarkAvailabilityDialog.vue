<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card class="sm:tw-p-6 tw-p-4">
      <v-expand-transition>
        <div v-show="state === states.CHOICES">
          <div class="tw-text-md tw-pb-4 tw-text-center">
            How would you like to mark <br v-if="isPhone" />
            your availability?
          </div>
          <div class="">
            <v-btn
              @click="showPermissions"
              class="tw-bg-green tw-mb-2"
              dark
              block
            >
              <div class="tw-text-sm -tw-mx-4">
                Automatically with Google Calendar
              </div>
            </v-btn>
            <v-btn @click="setAvailabilityManually" block>Manually</v-btn>
          </div>
        </div>
      </v-expand-transition>
      <v-expand-transition>
        <div v-show="state === states.PERMISSIONS">
          <div class="tw-text-md tw-font-medium tw-mb-4">
            We need the following permissions to access your Google Calendar
            events
          </div>

          <div class="tw-flex tw-flex-col tw-gap-4 tw-mb-8 tw-ml-4">
            <div class="tw-flex tw-gap-2 tw-text-sm">
              <v-img
                src="@/assets/gcal_logo.png"
                class="tw-flex-none"
                height="20"
                width="20"
              />
              <div>View events on all your calendars.</div>
            </div>
            <div class="tw-flex tw-gap-2 tw-text-sm">
              <v-img
                src="@/assets/gcal_logo.png"
                class="tw-flex-none"
                height="20"
                width="20"
              />
              <div>See the list of Google calendars you’re subscribed to.</div>
            </div>
          </div>

          <div class="tw-text-md tw-font-medium tw-mb-4">
            These permissions will
          </div>
          <div
            class="tw-flex tw-flex-col tw-gap-4 tw-mb-8 tw-text-very-dark-gray tw-text-sm"
          >
            <div>
              Allow us to display the names/times of your calendar events
            </div>
            <div>
              Allow us to display calendar events on all calendars, not just
              your primary calendar
            </div>
          </div>

          <div class="tw-flex">
            <!-- <v-spacer /> -->
            <v-btn @click="showChoices" text class="tw-flex-1 tw-mr-2">
              Back
            </v-btn>
            <v-btn
              @click="setAvailabilityAutomatically"
              class="tw-flex-1 tw-bg-green"
              dark
            >
              Allow
            </v-btn>
          </div>
        </div>
      </v-expand-transition>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"

export default {
  name: "MarkAvailabilityDialog",

  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      states: {
        CHOICES: "choices", // present user with choice of automatic or manual
        PERMISSIONS: "permissions", // present to user the gcal permissions we request
      },
      state: "choices",
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    setAvailabilityAutomatically() {
      this.$emit("setAvailabilityAutomatically")
    },
    setAvailabilityManually() {
      this.$emit("setAvailabilityManually")
    },
    showPermissions() {
      this.state = this.states.PERMISSIONS
    },
    showChoices() {
      this.state = this.states.CHOICES
    },
  },

  watch: {
    value() {
      if (!this.value) setTimeout(() => (this.state = this.states.CHOICES), 100)
    },
  },
}
</script>
