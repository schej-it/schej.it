<template>
  <v-card
    class="tw-overflow-none tw-relative tw-flex tw-max-w-[28rem] tw-flex-col tw-rounded-lg tw-py-4"
  >
    <v-card-title class="tw-mb-2 tw-flex tw-px-4 sm:tw-px-8">
      <div>
        {{ editEvent ? "Edit event" : "New event" }}
      </div>
      <v-spacer />
      <v-btn v-if="dialog" @click="$emit('input', false)" icon>
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </v-card-title>
    <v-card-text class="tw-flex-1 tw-overflow-auto tw-px-4 tw-py-1 sm:tw-px-8">
      <div class="tw-flex tw-flex-col tw-space-y-6">
        <v-text-field
          ref="name-field"
          v-model="name"
          placeholder="Name your event..."
          autofocus
          :disabled="loading"
          hide-details
          solo
          @keyup.enter="blurNameField"
        />

        <div>
          <div class="tw-mb-2 tw-text-lg tw-text-black">
            What times might work?
          </div>
          <div class="tw-flex tw-items-baseline tw-justify-center tw-space-x-2">
            <v-select
              v-model="startTime"
              :disabled="loading"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
            <div>to</div>
            <v-select
              v-model="endTime"
              :disabled="loading"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
          </div>
        </div>

        <div>
          <div class="tw-mb-2 tw-text-lg tw-text-black">
            What
            {{ selectedDateOption === dateOptions.SPECIFIC ? "dates" : "days" }}
            might work?
          </div>
          <v-select
            v-model="selectedDateOption"
            :items="Object.values(dateOptions)"
            solo
            hide-details
            class="tw-mb-4"
          />

          <v-expand-transition>
            <v-date-picker
              v-if="selectedDateOption === dateOptions.SPECIFIC"
              v-model="selectedDays"
              no-title
              multiple
              color="primary"
              :show-current="false"
              class="tw-min-w-full tw-border-0 tw-drop-shadow sm:tw-min-w-0"
              :min="minCalendarDate"
              full-width
            />
            <div v-else-if="selectedDateOption === dateOptions.DOW">
              <v-btn-toggle
                v-model="selectedDaysOfWeek"
                multiple
                solo
                color="primary"
              >
                <v-btn> S </v-btn>
                <v-btn> M </v-btn>
                <v-btn> T </v-btn>
                <v-btn> W </v-btn>
                <v-btn> T </v-btn>
                <v-btn> F </v-btn>
                <v-btn> S </v-btn>
              </v-btn-toggle>
            </div>
          </v-expand-transition>
        </div>
        <v-checkbox
          v-if="allowNotifications"
          v-model="notificationsEnabled"
          label="Email me each time someone joins my event!"
          hide-details
          class="tw-mt-2"
        />
        <div>
          <v-btn
            class="tw-justify-start tw-pl-0"
            block
            text
            @click="showAdvancedOptions = !showAdvancedOptions"
            ><span class="tw-mr-1">Advanced options</span>
            <v-icon>{{
              showAdvancedOptions ? "mdi-chevron-up" : "mdi-chevron-down"
            }}</v-icon></v-btn
          >
          <v-expand-transition>
            <div v-show="showAdvancedOptions">
              <div class="tw-my-2">
                <TimezoneSelector
                  class="tw-mb-2"
                  v-model="timezone"
                  label="Timezone"
                />
              </div>
            </div>
          </v-expand-transition>
        </div>
      </div>
    </v-card-text>
    <v-card-actions class="tw-relative tw-px-8">
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
import { eventTypes, dayIndexToDayString } from "@/constants"
import {
  post,
  put,
  timeNumToTimeString,
  dateToTimeNum,
  getISODateString,
  isPhone,
} from "@/utils"
import { mapActions } from "vuex"
import TimezoneSelector from "./schedule_overlap/TimezoneSelector.vue"
import dayjs from "dayjs"
import utcPlugin from "dayjs/plugin/utc"
import timezonePlugin from "dayjs/plugin/timezone"
dayjs.extend(utcPlugin)
dayjs.extend(timezonePlugin)

export default {
  name: "NewEvent",

  emits: ["input"],

  props: {
    event: { type: Object },
    editEvent: { type: Boolean, default: false },
    dialog: { type: Boolean, default: true },
    allowNotifications: { type: Boolean, default: true },
  },

  components: {
    TimezoneSelector,
  },

  data: () => ({
    name: "",
    startTime: 9,
    endTime: 17,
    loading: false,
    selectedDays: [],
    selectedDaysOfWeek: [],
    notificationsEnabled: false,

    // Date options
    dateOptions: Object.freeze({
      SPECIFIC: "Specific dates",
      DOW: "Days of the week",
    }),
    selectedDateOption: "Specific dates",

    // Advanced options
    showAdvancedOptions: false,
    timezone: {},
  }),

  computed: {
    formComplete() {
      return (
        this.name.length > 0 &&
        (this.selectedDays.length > 0 || this.selectedDaysOfWeek.length > 0) //&&
        // (this.startTime < this.endTime ||
        //   (this.endTime === 0 && this.startTime != 0))
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
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    ...mapActions(["showError"]),
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

      let dates = []
      let type = ""

      // Get date objects for each selected day
      const startTimeString = timeNumToTimeString(this.startTime)
      if (this.selectedDateOption === this.dateOptions.SPECIFIC) {
        type = eventTypes.SPECIFIC_DATES

        for (const day of this.selectedDays) {
          const date = dayjs.tz(
            `${day} ${startTimeString}`,
            this.timezone.value
          )
          dates.push(date.toDate())
        }
      } else if (this.selectedDateOption === this.dateOptions.DOW) {
        type = eventTypes.DOW

        this.selectedDaysOfWeek.sort((a, b) => a - b)
        for (const dayIndex of this.selectedDaysOfWeek) {
          const day = dayIndexToDayString[dayIndex]
          const date = dayjs.tz(
            `${day} ${startTimeString}`,
            this.timezone.value
          )
          dates.push(date.toDate())
        }
      }

      this.loading = true
      if (!this.editEvent) {
        // Create new event on backend
        post("/events", {
          name: this.name,
          duration,
          dates,
          notificationsEnabled: this.notificationsEnabled,
          type,
        })
          .then(({ eventId }) => {
            this.$router.push({
              name: "event",
              params: { eventId, initialTimezone: this.timezone },
            })
            this.loading = false
            this.$emit("input", false)
            this.reset()

            this.$posthog?.capture("Event created", {
              eventId: eventId,
              eventName: this.name,
              eventDuration: duration,
              eventDates: dates,
              eventNotificationsEnabled: this.notificationsEnabled,
              eventType: type,
            })
          })
          .catch((err) => {
            this.showError(
              "There was a problem creating that event! Please try again later."
            )
          })
      } else {
        // Edit event on backend
        if (this.event) {
          put(`/events/${this.event._id}`, {
            name: this.name,
            duration,
            dates,
            notificationsEnabled: this.notificationsEnabled,
            type,
          })
            .then(() => {
              this.$posthog?.capture("Event edited", {
                eventId: this.event._id,
                eventName: this.name,
                eventDuration: duration,
                eventDates: dates,
                eventNotificationsEnabled: this.notificationsEnabled,
                eventType: type,
              })

              this.$emit("input", false)
              this.reset()
              window.location.reload()
            })
            .catch((err) => {
              this.showError(
                "There was a problem editing this event! Please try again later."
              )
            })
        }
      }
    },
  },

  watch: {
    event: {
      immediate: true,
      handler() {
        // Populate event fields if this.event exists
        if (this.event) {
          this.name = this.event.name
          this.startTime = Math.floor(dateToTimeNum(this.event.dates[0]))
          this.endTime = (this.startTime + this.event.duration) % 24
          this.notificationsEnabled = this.event.notificationsEnabled

          if (this.event.type === eventTypes.SPECIFIC_DATES) {
            this.selectedDateOption = this.dateOptions.SPECIFIC
            const selectedDays = []
            for (const date of this.event.dates) {
              selectedDays.push(getISODateString(date))
            }
            this.selectedDays = selectedDays
          } else if (this.event.type === eventTypes.DOW) {
            this.selectedDateOption = this.dateOptions.DOW
            const selectedDaysOfWeek = []
            for (const date of this.event.dates) {
              selectedDaysOfWeek.push(new Date(date).getDay())
            }
            this.selectedDaysOfWeek = selectedDaysOfWeek
          }
        }
      },
    },
    selectedDateOption() {
      // Reset the other date / day selection when date option is changed
      if (this.selectedDateOption === this.dateOptions.SPECIFIC) {
        this.selectedDaysOfWeek = []
      } else if (this.selectedDateOption === this.dateOptions.DOW) {
        this.selectedDays = []
      }
    },
  },
}
</script>
