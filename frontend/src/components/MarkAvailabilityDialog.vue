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
          <div class="">
            <v-btn
              @click="showPermissions"
              class="tw-mb-2"
              color="primary"
              dark
              block
            >
              <div class="-tw-mx-4 tw-text-sm">
                Autofill with Google Calendar
              </div>
            </v-btn>
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
