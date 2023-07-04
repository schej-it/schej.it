<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card class="tw-p-4 sm:tw-p-6">
      <v-expand-transition>
        <div v-show="state === states.CHOICES">
          <div class="tw-text-md tw-pb-4 tw-text-center">
            How would you like to mark <br v-if="isPhone" />
            your availability?
          </div>
          <div class="">
            <v-btn
              @click="showPermissions"
              class="tw-mb-2 tw-bg-green"
              dark
              block
            >
              <div class="-tw-mx-4 tw-text-sm">
                Automatically with Google Calendar
              </div>
            </v-btn>
            <v-btn @click="setAvailabilityManually" block>Manually</v-btn>
          </div>
        </div>
      </v-expand-transition>
      <v-expand-transition>
        <div v-show="state === states.PERMISSIONS">
          <div class="tw-text-md tw-mb-4 tw-font-medium">
            We need the following permissions to access your Google Calendar
            events
          </div>

          <div class="tw-mb-8 tw-ml-4 tw-flex tw-flex-col tw-gap-4">
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

          <div class="tw-text-md tw-mb-4 tw-font-medium">
            These permissions will
          </div>
          <div
            class="tw-mb-8 tw-flex tw-flex-col tw-gap-4 tw-text-sm tw-text-very-dark-gray"
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
            <v-btn @click="showChoices" text class="tw-mr-2 tw-flex-1">
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
