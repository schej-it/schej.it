<template>
  <v-card
    :flat="dialog"
    :class="{ 'tw-py-4': !dialog }"
    class="tw-overflow-none tw-relative tw-flex tw-max-w-[28rem] tw-flex-col tw-rounded-lg tw-transition-all"
  >
    <v-card-title class="tw-mb-2 tw-flex tw-px-4 sm:tw-px-8">
      <div>
        {{ edit ? "Edit event" : "New event" }}
      </div>
      <v-spacer />
      <template v-if="dialog">
        <v-btn v-if="showHelp" icon @click="helpDialog = true">
          <v-icon>mdi-help-circle</v-icon>
        </v-btn>
        <v-btn v-else @click="$emit('input', false)" icon>
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <HelpDialog
          v-model="helpDialog"
          text="Use events to poll people for their availabilities on certain days"
        />
      </template>
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
            @click="() => toggleAdvancedOptions()"
            ><span class="tw-mr-1">Advanced options</span>
            <v-icon :class="`tw-rotate-${showAdvancedOptions ? '180' : '0'}`"
              >mdi-chevron-down</v-icon
            ></v-btn
          >
          <v-expand-transition>
            <div v-show="showAdvancedOptions">
              <div class="tw-my-2 tw-space-y-4">
                <TimezoneSelector v-model="timezone" label="Timezone" />
                <EmailInput
                  v-show="authUser"
                  ref="emailReminders"
                  @requestContactsAccess="requestContactsAccess"
                  labelColor="tw-text-very-dark-gray"
                  :addedEmails="addedEmails"
                  @update:emails="(newEmails) => (emails = newEmails)"
                >
                  <template v-slot:header>
                    <div class="tw-flex tw-gap-1">
                      <div class="tw-text-very-dark-gray">
                        Set up email reminders
                      </div>

                      <v-tooltip top>
                        <template v-slot:activator="{ on, attrs }">
                          <v-icon small v-bind="attrs" v-on="on"
                            >mdi-information-outline
                          </v-icon>
                        </template>
                        <div>
                          Reminder emails will be sent the day of event
                          creation,<br />one day after, and three days after.
                          You will also receive <br />an email when everybody
                          has filled out the Schej.
                        </div>
                      </v-tooltip>
                    </div>
                  </template>
                </EmailInput>
              </div>
            </div>
          </v-expand-transition>
          <div class="tw-bg-red" ref="advancedOpenScrollTo"></div>
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
        {{ edit ? "Edit" : "Create" }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { eventTypes, dayIndexToDayString, authTypes } from "@/constants"
import {
  post,
  put,
  timeNumToTimeString,
  dateToTimeNum,
  getISODateString,
  isPhone,
  validateEmail,
  signInGoogle,
} from "@/utils"
import { mapActions, mapState } from "vuex"
import TimezoneSelector from "./schedule_overlap/TimezoneSelector.vue"
import HelpDialog from "./HelpDialog.vue"
import EmailInput from "./event/EmailInput.vue"
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
    edit: { type: Boolean, default: false },
    dialog: { type: Boolean, default: true },
    allowNotifications: { type: Boolean, default: true },
    contactsPayload: { type: Object, default: () => ({}) },
    inDialog: { type: Boolean, default: true },
    showHelp: { type: Boolean, default: false },
  },

  components: {
    TimezoneSelector,
    HelpDialog,
    EmailInput,
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
    emails: [], // For email reminders

    helpDialog: false,
  }),

  mounted() {
    if (Object.keys(this.contactsPayload).length > 0)
      this.toggleAdvancedOptions(true)
  },

  computed: {
    ...mapState(["authUser"]),
    formComplete() {
      let emailsValid = true

      for (const email of this.emails) {
        if (!validateEmail(email)) {
          emailsValid = false
          break
        }
      }

      return (
        this.name.length > 0 &&
        (this.selectedDays.length > 0 || this.selectedDaysOfWeek.length > 0) &&
        emailsValid //&&
        // (this.startTime < this.endTime ||
        //   (this.endTime === 0 && this.startTime != 0))
      )
    },
    addedEmails() {
      if (Object.keys(this.contactsPayload).length > 0)
        return this.contactsPayload.emails
      return this.event && this.event.remindees
        ? this.event.remindees.map((r) => r.email)
        : []
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
      if (this.edit) {
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
      this.selectedDaysOfWeek = []
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

      const name = this.name
      const notificationsEnabled = this.notificationsEnabled
      const remindees = this.emails
      if (!this.edit) {
        // Create new event on backend
        post("/events", {
          name,
          duration,
          dates,
          notificationsEnabled,
          remindees,
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
              eventName: name,
              eventDuration: duration,
              eventDates: JSON.stringify(dates),
              eventNotificationsEnabled: notificationsEnabled,
              eventRemindees: remindees,
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
            name,
            duration,
            dates,
            notificationsEnabled,
            remindees,
            type,
          })
            .then(() => {
              this.$posthog?.capture("Event edited", {
                eventId: this.event._id,
                eventName: name,
                eventDuration: duration,
                eventDates: JSON.stringify(dates),
                eventNotificationsEnabled: notificationsEnabled,
                eventRemindees: remindees,
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
    toggleAdvancedOptions(delayed = false) {
      this.showAdvancedOptions = !this.showAdvancedOptions

      const openScrollEl = this.$refs.advancedOpenScrollTo

      if (this.inDialog && openScrollEl && this.showAdvancedOptions) {
        setTimeout(
          () => openScrollEl.scrollIntoView({ behavior: "smooth" }),
          delayed ? 500 : 200
        )
      }
    },

    /** Redirects user to oauth page requesting access to the user's contacts */
    requestContactsAccess({ emails }) {
      const payload = {
        emails,
      }
      signInGoogle({
        state: {
          type: authTypes.EVENT_CONTACTS,
          eventId: this.event ? this.event._id : "",
          payload,
        },
        requestContactsPermission: true,
      })
    },
    /** Update state based on the contactsPayload after granting contacts access */
    contactsAccessGranted({ curScheduledEvent, ...data }) {
      this.curScheduledEvent = curScheduledEvent
      this.$refs.confirmDetailsDialog?.setData(data)
      this.confirmDetailsDialog = true
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
