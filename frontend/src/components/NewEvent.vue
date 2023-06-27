<template>
  <v-card
    class="tw-py-4 tw-flex tw-flex-col tw-rounded-lg tw-relative tw-overflow-none"
  >
    <v-card-title class="tw-px-8 tw-flex tw-mb-2">
      <div>
        {{ editEvent ? "Edit event" : "New event" }}
      </div>
      <v-spacer />
      <v-btn v-if="dialog" @click="$emit('input', false)" icon>
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text class="tw-px-8 tw-overflow-auto tw-py-1 tw-flex-1">
      <div class="tw-space-y-10 tw-flex tw-flex-col">
        <v-text-field
          ref="name-field"
          v-model="name"
          autofocus
          :disabled="loading"
          class="tw-text-white tw-flex-initial"
          placeholder="Name your event..."
          hide-details
          solo
          @keyup.enter="blurNameField"
        />

        <div>
          <div class="tw-text-lg tw-text-black tw-mb-4">
            What times might work?
          </div>
          <div class="tw-flex tw-space-x-2 tw-items-baseline tw-justify-center">
            <v-select
              v-model="startTime"
              :disabled="loading"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
            <div>to</div>
            <v-select
              v-model="endTime"
              :disabled="loading"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
          </div>
        </div>

        <div>
          <div class="tw-text-lg tw-text-black tw-mb-4">
            What dates might work?
          </div>
          <v-select
            v-model="selectedDateOption"
            :items="Object.values(dateOptions)"
            solo
            hide-details
            class="tw-mb-2"
          />

          <v-expand-transition>
            <v-date-picker
              v-if="selectedDateOption === dateOptions.SPECIFIC"
              v-model="selectedDays"
              no-title
              multiple
              color="primary"
              elevation="2"
              :show-current="false"
              class="tw-min-w-full sm:tw-min-w-0 tw-border-0"
              :min="minCalendarDate"
              full-width
            />
            <div v-else-if="selectedDateOption === dateOptions.DOW">
              <div class="tw-flex tw-mt-4">
                <v-btn-toggle
                  v-model="selectedDaysOfWeek"
                  multiple
                  solo
                  color="primary"
                >
                  <v-btn v-if="!mondayStart"> S </v-btn>
                  <v-btn> M </v-btn>
                  <v-btn> T </v-btn>
                  <v-btn> W </v-btn>
                  <v-btn> T </v-btn>
                  <v-btn> F </v-btn>
                  <v-btn> S </v-btn>
                  <v-btn v-if="mondayStart"> S </v-btn>
                </v-btn-toggle>
              </div>
              <v-checkbox
                v-model="mondayStart"
                label="Start on Monday"
                hide-details
              />
            </div>
          </v-expand-transition>
        </div>
        <v-checkbox
          v-if="dialog"
          v-model="notificationsEnabled"
          label="Email me each time someone joins my event!"
          hide-details
          class="tw-mt-2"
        />
      </div>
    </v-card-text>
    <v-card-actions class="tw-px-8 tw-relative">
      <v-btn
        block
        :loading="loading"
        :dark="formComplete"
        class="tw-mt-4 tw-bg-green"
        :disabled="!formComplete"
        @click="submit"
      >
        {{ editEvent ? "Edit" : "Create" }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import {
  post,
  put,
  timeNumToTimeString,
  dateToTimeNum,
  getISODateString,
} from "@/utils"

export default {
  name: "NewEvent",

  emits: ["input"],

  props: {
    event: { type: Object },
    editEvent: { type: Boolean, default: false },
    dialog: { type: Boolean, default: true },
  },

  data: () => ({
    name: "",
    startTime: 9,
    endTime: 17,
    loading: false,
    selectedDays: [],
    selectedDaysOfWeek: [],
    mondayStart: false,
    notificationsEnabled: false,

    dateOptions: Object.freeze({
      SPECIFIC: "Specific dates",
      DOW: "Days of the week",
    }),
    selectedDateOption: "Specific dates",
  }),

  created() {
    // Populate event fields if this.event exists
    if (this.event) {
      this.name = this.event.name
      this.startTime = Math.floor(dateToTimeNum(this.event.dates[0]))
      this.endTime = (this.startTime + this.event.duration) % 24
      this.notificationsEnabled = this.event.notificationsEnabled

      const selectedDays = []
      for (const date of this.event.dates) {
        selectedDays.push(getISODateString(date))
      }
      this.selectedDays = selectedDays
    }
  },

  computed: {
    formComplete() {
      return (
        this.name.length > 0 &&
        this.selectedDays.length > 0 &&
        (this.startTime < this.endTime ||
          (this.endTime === 0 && this.startTime != 0))
      )
    },
    times() {
      const times = []

      for (let h = 1; h < 12; ++h) {
        times.push({ text: `${h} am`, value: h })
      }
      for (let h = 0; h < 12; ++h) {
        times.push({ text: `${h == 0 ? 12 : h} pm`, value: h + 12 })
      }
      times.push({ text: "12 am", value: 0 })

      return times
    },
    minCalendarDate() {
      if (this.editEvent) {
        return ""
      }

      let today = new Date()
      let dd = String(today.getDate()).padStart(2, "0")
      let mm = String(today.getMonth() + 1).padStart(2, "0")
      let yyyy = today.getFullYear()

      return yyyy + "-" + mm + "-" + dd
    },
  },

  methods: {
    blurNameField() {
      this.$refs["name-field"].blur()
    },
    reset() {
      this.name = ""
      this.startTime = 9
      this.endTime = 17
      this.selectedDays = []
    },
    submit() {
      this.selectedDays.sort()

      // Get duration of event
      let duration = this.endTime - this.startTime
      if (duration < 0) duration += 24

      // Get date objects for each selected day
      const startTimeString = timeNumToTimeString(this.startTime)
      const dates = []
      for (const day of this.selectedDays) {
        const date = new Date(`${day}T${startTimeString}`)
        dates.push(date)
      }

      this.loading = true
      if (!this.editEvent) {
        // Create new event on backend
        post("/events", {
          name: this.name,
          duration,
          dates,
          notificationsEnabled: this.notificationsEnabled,
        }).then(({ eventId }) => {
          this.$router.push({ name: "event", params: { eventId } })
          this.loading = false
        })
      } else {
        // Edit event on backend
        if (this.event) {
          put(`/events/${this.event._id}`, {
            name: this.name,
            duration,
            dates,
            notificationsEnabled: this.notificationsEnabled,
          }).then(() => {
            window.location.reload()
          })
        }
      }
    },
  },
}
</script>
