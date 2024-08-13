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
            How would you like to add <br v-if="isPhone" />
            your availability?
          </div>
          <div class="tw-flex tw-flex-col tw-gap-2">
            <v-btn block @click="showPermissions" class="tw-bg-white">
              <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
                <v-img
                  class="tw-flex-initial"
                  width="20"
                  height="20"
                  src="@/assets/google_logo.svg"
                />
                <v-spacer />
                Autofill with Google Calendar
                <v-spacer />
              </div>
            </v-btn>
            <v-btn block class="tw-bg-white">
              <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
                <v-img
                  class="tw-flex-initial"
                  width="20"
                  height="20"
                  src="@/assets/apple_logo.svg"
                />
                <v-spacer />
                Autofill with Apple Calendar
                <v-spacer />
              </div>
            </v-btn>
            <v-btn block class="tw-bg-white">
              <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
                <v-img
                  class="tw-flex-initial"
                  width="20"
                  height="20"
                  src="@/assets/outlook_logo.svg"
                />
                <v-spacer />
                Autofill with Outlook Calendar
                <v-spacer />
              </div>
            </v-btn>
            <div class="tw-flex tw-items-center tw-gap-3">
              <v-divider />
              <div
                class="tw-text-center tw-text-xs tw-font-medium tw-text-dark-gray"
              >
                or
              </div>
              <v-divider />
            </div>
            <v-btn @click="setAvailabilityManually" block>Manually</v-btn>
          </div>
        </div>
      </v-expand-transition>
      <v-expand-transition>
        <CalendarPermissionsCard
          v-show="state === states.PERMISSIONS"
          cancelLabel="Back"
          @cancel="showChoices"
          @allow="setAvailabilityAutomatically"
        />
      </v-expand-transition>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"
import CalendarPermissionsCard from "@/components/CalendarPermissionsCard"

export default {
  name: "MarkAvailabilityDialog",

  props: {
    value: { type: Boolean, required: true },
  },

  components: {
    CalendarPermissionsCard,
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
