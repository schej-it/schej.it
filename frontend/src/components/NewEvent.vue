<template>
  <v-card
    :flat="dialog"
    :class="{ 'tw-py-4': !dialog, 'tw-flex-1': dialog }"
    class="tw-relative tw-flex tw-max-w-[28rem] tw-flex-col tw-overflow-hidden tw-rounded-lg tw-transition-all"
  >
    <v-card-title class="tw-mb-2 tw-flex tw-gap-2 tw-px-4 sm:tw-px-8">
      <div>
        <div class="tw-mb-1">
          {{ edit ? "Edit event" : "New event" }}
        </div>
        <div
          v-if="dialog && showHelp"
          class="tw-text-xs tw-font-normal tw-italic tw-text-dark-gray"
        >
          Ideal for one-time / recurring meetings
        </div>
      </div>
      <v-spacer />
      <template v-if="dialog">
        <v-btn v-if="showHelp" icon @click="helpDialog = true">
          <v-icon>mdi-information-outline</v-icon>
        </v-btn>
        <v-btn v-else @click="$emit('input', false)" icon>
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <HelpDialog v-model="helpDialog">
          <template v-slot:header>Events</template>
          <div class="tw-mb-4">
            Use events to collect people's availabilities and compare them
            across certain days.
          </div>
        </HelpDialog>
      </template>
    </v-card-title>
    <v-card-text class="tw-flex-1 tw-overflow-auto tw-px-4 tw-py-1 sm:tw-px-8">
      <v-form
        ref="form"
        v-model="formValid"
        lazy-validation
        class="tw-flex tw-flex-col tw-gap-y-6"
        :disabled="loading"
      >
        <v-text-field
          ref="name-field"
          v-model="name"
          placeholder="Name your event..."
          hide-details="auto"
          solo
          @keyup.enter="blurNameField"
          :rules="nameRules"
          required
        />

        <SlideToggle
          v-if="daysOnlyEnabled && !edit"
          class="tw-w-full"
          v-model="daysOnly"
          :options="daysOnlyOptions"
        />

        <div>
          <v-expand-transition>
            <div v-if="!daysOnly">
              <div class="tw-mb-2 tw-text-lg tw-text-black">
                What times might work?
              </div>
              <div
                class="tw-mb-6 tw-flex tw-items-baseline tw-justify-center tw-space-x-2"
              >
                <v-select
                  v-model="startTime"
                  menu-props="auto"
                  :items="times"
                  hide-details
                  solo
                ></v-select>
                <div>to</div>
                <v-select
                  v-model="endTime"
                  menu-props="auto"
                  :items="times"
                  hide-details
                  solo
                ></v-select>
              </div>
            </div>
          </v-expand-transition>

          <div class="tw-mb-2 tw-text-lg tw-text-black">
            What
            {{ selectedDateOption === dateOptions.SPECIFIC ? "dates" : "days" }}
            might work?
          </div>
          <v-select
            v-if="!edit && !daysOnly"
            v-model="selectedDateOption"
            :items="Object.values(dateOptions)"
            solo
            hide-details
            class="tw-mb-4"
          />

          <v-expand-transition>
            <v-input
              v-if="selectedDateOption === dateOptions.SPECIFIC || daysOnly"
              v-model="selectedDays"
              hint="Drag to select multiple dates"
              persistent-hint
              hide-details="auto"
              :rules="selectedDaysRules"
              key="date-picker"
            >
              <DatePicker
                v-model="selectedDays"
                :minCalendarDate="minCalendarDate"
              />
            </v-input>
            <v-input
              v-else-if="selectedDateOption === dateOptions.DOW"
              v-model="selectedDaysOfWeek"
              hide-details="auto"
              :rules="selectedDaysRules"
              key="days-of-week"
              class="tw-w-fit"
            >
              <v-btn-toggle
                v-model="selectedDaysOfWeek"
                multiple
                solo
                color="primary"
              >
                <v-btn depressed> S </v-btn>
                <v-btn depressed> M </v-btn>
                <v-btn depressed> T </v-btn>
                <v-btn depressed> W </v-btn>
                <v-btn depressed> T </v-btn>
                <v-btn depressed> F </v-btn>
                <v-btn depressed> S </v-btn>
              </v-btn-toggle>
            </v-input>
          </v-expand-transition>
        </div>

        <v-checkbox
          v-if="allowNotifications"
          v-model="notificationsEnabled"
          hide-details
          class="tw-mt-2"
        >
          <template v-slot:label>
            <span class="tw-text-sm tw-text-very-dark-gray"
              >Email me each time someone joins my event</span
            >
          </template>
        </v-checkbox>

        <div v-if="authUser">
          <v-btn
            class="tw-flex tw-items-end tw-justify-start tw-p-1 tw-text-base"
            block
            text
            @click="() => toggleEmailReminders()"
            ><span class="tw-mr-1">Email reminders</span>
            <v-icon
              :class="`tw-rotate-${showEmailReminders ? '180' : '0'}`"
              :size="30"
              >mdi-chevron-down</v-icon
            ></v-btn
          >
          <v-expand-transition>
            <div v-show="showEmailReminders">
              <div class="tw-my-2 tw-space-y-5">
                <EmailInput
                  v-show="authUser"
                  ref="emailInput"
                  @requestContactsAccess="requestContactsAccess"
                  labelColor="tw-text-very-dark-gray"
                  :addedEmails="addedEmails"
                  @update:emails="(newEmails) => (emails = newEmails)"
                >
                  <template v-slot:header>
                    <div class="tw-flex tw-gap-1">
                      <div class="tw-text-very-dark-gray">
                        Remind people to fill out the Schej
                      </div>

                      <v-tooltip top color="#4F4F4F">
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
          <div ref="emailRemindersOpenScrollTo"></div>
        </div>

        <div>
          <v-btn
            class="-tw-mt-2 tw-flex tw-items-end tw-justify-start tw-p-1 tw-text-base"
            block
            text
            @click="() => toggleAdvancedOptions()"
            ><span class="tw-mr-1">Advanced options</span>
            <v-icon
              :class="`tw-rotate-${showAdvancedOptions ? '180' : '0'}`"
              :size="30"
              >mdi-chevron-down</v-icon
            ></v-btn
          >
          <v-expand-transition>
            <div v-show="showAdvancedOptions">
              <div class="tw-my-2 tw-space-y-5">
                <v-checkbox
                  v-if="authUser"
                  v-model="blindAvailabilityEnabled"
                  hide-details
                >
                  <template v-slot:label>
                    <span class="tw-text-sm tw-text-very-dark-gray"
                      >Only show responses to event creator</span
                    >
                  </template>
                </v-checkbox>
                <v-checkbox
                  v-if="authUser"
                  v-model="sendEmailAfterXResponsesEnabled"
                  hide-details
                >
                  <template v-slot:label>
                    <div
                      :class="
                        !sendEmailAfterXResponsesEnabled && 'tw-opacity-50'
                      "
                      class="tw-flex tw-items-center tw-gap-x-2 tw-text-sm tw-text-very-dark-gray"
                    >
                      <div>Email me after</div>
                      <v-text-field
                        v-model="sendEmailAfterXResponses"
                        @click="
                          (e) => {
                            e.preventDefault()
                            e.stopPropagation()
                          }
                        "
                        :disabled="!sendEmailAfterXResponsesEnabled"
                        dense
                        class="email-me-after-text-field -tw-mt-[2px] tw-w-10"
                        menu-props="auto"
                        hide-details
                        type="number"
                        min="1"
                      ></v-text-field>
                      <div>responses</div>
                    </div>
                  </template>
                </v-checkbox>
                <TimezoneSelector v-model="timezone" label="Timezone" />
              </div>
            </div>
          </v-expand-transition>
          <div ref="advancedOpenScrollTo"></div>
        </div>
      </v-form>
    </v-card-text>
    <v-card-actions class="tw-relative tw-px-4 sm:tw-px-8">
      <div class="tw-relative tw-w-full">
        <v-btn
          :disabled="!formValid"
          block
          :loading="loading"
          color="primary"
          class="tw-mt-4 tw-bg-green"
          @click="submit"
        >
          {{ edit ? "Save edits" : "Create event" }}
        </v-btn>
        <div
          :class="formValid ? 'tw-invisible' : 'tw-visible'"
          class="tw-mt-1 tw-text-xs tw-text-red"
        >
          Please fix form errors before continuing
        </div>
      </div>
    </v-card-actions>
  </v-card>
</template>

<style>
.email-me-after-text-field input {
  padding: 0px !important;
}
</style>

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
  getDateWithTimezone,
} from "@/utils"
import { mapActions, mapState } from "vuex"
import TimezoneSelector from "./schedule_overlap/TimezoneSelector.vue"
import HelpDialog from "./HelpDialog.vue"
import EmailInput from "./event/EmailInput.vue"
import DatePicker from "@/components/DatePicker.vue"
import SlideToggle from "./SlideToggle.vue"

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
    showHelp: { type: Boolean, default: false },
  },

  components: {
    TimezoneSelector,
    HelpDialog,
    EmailInput,
    DatePicker,
    SlideToggle,
  },

  data: () => ({
    formValid: true,
    name: "",
    startTime: 9,
    endTime: 17,
    loading: false,
    selectedDays: [],
    selectedDaysOfWeek: [],
    notificationsEnabled: false,

    daysOnly: false,
    daysOnlyOptions: Object.freeze([
      { text: "Dates and times", value: false },
      { text: "Dates only", value: true },
    ]),

    // Date options
    dateOptions: Object.freeze({
      SPECIFIC: "Specific dates",
      DOW: "Days of the week",
    }),
    selectedDateOption: "Specific dates",

    // Email reminders
    showEmailReminders: false,
    emails: [], // For email reminders

    // Advanced options
    showAdvancedOptions: false,
    blindAvailabilityEnabled: false,
    timezone: {},
    sendEmailAfterXResponsesEnabled: false,
    sendEmailAfterXResponses: 3,

    helpDialog: false,

    // Unsaved changes
    initialEventData: {},
  }),

  mounted() {
    if (Object.keys(this.contactsPayload).length > 0) {
      this.toggleEmailReminders(true)

      /** Get previously filled out data after enabling contacts  */
      this.name = this.contactsPayload.name
      this.startTime = this.contactsPayload.startTime
      this.endTime = this.contactsPayload.endTime
      this.daysOnly = this.contactsPayload.daysOnly
      this.selectedDateOption = this.contactsPayload.selectedDateOption
      this.selectedDaysOfWeek = this.contactsPayload.selectedDaysOfWeek
      this.selectedDays = this.contactsPayload.selectedDays
      this.notificationsEnabled = this.contactsPayload.notificationsEnabled
      this.timezone = this.contactsPayload.timezone

      this.$refs.form.resetValidation()
    }
  },

  computed: {
    ...mapState(["authUser", "daysOnlyEnabled"]),
    nameRules() {
      return [(v) => !!v || "Event name is required"]
    },
    selectedDaysRules() {
      return [
        (selectedDays) =>
          selectedDays.length > 0 || "Please select at least one day",
      ]
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
      this.notificationsEnabled = false
      this.daysOnly = false
      this.selectedDateOption = "Specific dates",
      this.emails = []
      this.showAdvancedOptions = false
      this.blindAvailabilityEnabled = false
      this.sendEmailAfterXResponsesEnabled = false
      this.sendEmailAfterXResponses = 3

      this.$refs.form.resetValidation()
    },
    submit() {
      if (!this.$refs.form.validate()) return

      this.selectedDays.sort()

      // Get duration of event
      let duration = this.endTime - this.startTime
      if (duration < 0) duration += 24

      // Get date objects for each selected day
      let dates = []
      let type = ""
      if (this.daysOnly) {
        duration = 0
        type = eventTypes.SPECIFIC_DATES

        for (const day of this.selectedDays) {
          const date = new Date(`${day} 00:00:00Z`)
          dates.push(date)
        }
      } else {
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
      }

      this.loading = true

      const payload = {
        name: this.name,
        duration: duration,
        dates: dates,
        notificationsEnabled: this.notificationsEnabled,
        blindAvailabilityEnabled: this.blindAvailabilityEnabled,
        daysOnly: this.daysOnly,
        remindees: this.emails,
        type: type,
        sendEmailAfterXResponses: this.sendEmailAfterXResponsesEnabled
          ? parseInt(this.sendEmailAfterXResponses)
          : -1,
      }
      const posthogPayload = {
        eventName: this.name,
        eventDuration: duration,
        eventDates: JSON.stringify(dates),
        eventNotificationsEnabled: this.notificationsEnabled,
        eventBlindAvailabilityEnabled: this.blindAvailabilityEnabled,
        eventDaysOnly: this.daysOnly,
        eventRemindees: this.emails,
        eventType: type,
        eventSendEmailAfterXResponses: this.sendEmailAfterXResponsesEnabled
          ? parseInt(this.sendEmailAfterXResponses)
          : -1,
      }

      if (!this.edit) {
        // Create new event on backend
        post("/events", payload)
          .then(({ eventId, shortId }) => {
            this.$router.push({
              name: "event",
              params: {
                eventId: shortId ?? eventId,
                initialTimezone: this.timezone,
              },
            })

            this.$emit("input", false)
            this.reset()

            posthogPayload.eventId = eventId
            this.$posthog?.capture("Event created", posthogPayload)
          })
          .catch((err) => {
            this.showError(
              "There was a problem creating that event! Please try again later."
            )
          })
          .finally(() => {
            this.loading = false
          })
      } else {
        // Edit event on backend
        if (this.event) {
          put(`/events/${this.event._id}`, payload)
            .then(() => {
              posthogPayload.eventId = this.event._id
              this.$posthog?.capture("Event edited", posthogPayload)

              this.$emit("input", false)
              this.reset()
              window.location.reload()
            })
            .catch((err) => {
              this.showError(
                "There was a problem editing this event! Please try again later."
              )
            })
            .finally(() => {
              this.loading = false
            })
        }
      }
    },
    toggleAdvancedOptions(delayed = false) {
      this.showAdvancedOptions = !this.showAdvancedOptions
      if (this.showAdvancedOptions)
        this.scrollToElement(this.$refs.advancedOpenScrollTo, delayed)
    },

    toggleEmailReminders(delayed = false) {
      this.showEmailReminders = !this.showEmailReminders
      if (this.showEmailReminders)
        this.scrollToElement(this.$refs.emailRemindersOpenScrollTo, delayed)
    },

    scrollToElement(element, delayed = false) {
      const openScrollEl = element
      if (this.dialog && openScrollEl) {
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
        name: this.name,
        startTime: this.startTime,
        endTime: this.endTime,
        daysOnly: this.daysOnly,
        selectedDays: this.selectedDays,
        selectedDaysOfWeek: this.selectedDaysOfWeek,
        selectedDateOption: this.selectedDateOption,
        notificationsEnabled: this.notificationsEnabled,
        timezone: this.timezone,
      }
      signInGoogle({
        state: {
          type: authTypes.EVENT_CONTACTS,
          eventId: this.event ? this.event.shortId ?? this.event._id : "",
          openNewGroup: false,
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

    /** Populates the form fields based on this.event */
    updateFieldsFromEvent() {
      if (this.event) {
        this.name = this.event.name

        // Set start time, accounting for the timezone
        this.startTime = Math.floor(
          dateToTimeNum(getDateWithTimezone(this.event.dates[0]), true)
        )
        this.startTime %= 24

        this.endTime = (this.startTime + this.event.duration) % 24
        this.notificationsEnabled = this.event.notificationsEnabled
        this.blindAvailabilityEnabled = this.event.blindAvailabilityEnabled
        this.daysOnly = this.event.daysOnly

        if (
          this.event.sendEmailAfterXResponses !== null &&
          this.event.sendEmailAfterXResponses > 0
        ) {
          this.sendEmailAfterXResponsesEnabled = true
          this.sendEmailAfterXResponses = this.event.sendEmailAfterXResponses
        }

        if (this.event.daysOnly) {
          this.selectedDateOption = this.dateOptions.SPECIFIC
          const selectedDays = []
          for (let date of this.event.dates) {
            selectedDays.push(getISODateString(date, true))
          }
          this.selectedDays = selectedDays
        } else {
          if (this.event.type === eventTypes.SPECIFIC_DATES) {
            this.selectedDateOption = this.dateOptions.SPECIFIC
            const selectedDays = []
            for (let date of this.event.dates) {
              date = getDateWithTimezone(date)

              selectedDays.push(getISODateString(date, true))
            }
            this.selectedDays = selectedDays
          } else if (this.event.type === eventTypes.DOW) {
            this.selectedDateOption = this.dateOptions.DOW
            const selectedDaysOfWeek = []
            for (let date of this.event.dates) {
              date = getDateWithTimezone(date)

              selectedDaysOfWeek.push(date.getUTCDay())
            }
            this.selectedDaysOfWeek = selectedDaysOfWeek
          }
        }
      }
    },
    resetToEventData() {
      this.updateFieldsFromEvent()
      this.$refs.emailInput.reset()
    },
    setInitialEventData() {
      this.initialEventData = {
        name: this.name,
        startTime: this.startTime,
        endTime: this.endTime,
        daysOnly: this.daysOnly,
        selectedDays: this.selectedDays,
        selectedDaysOfWeek: this.selectedDaysOfWeek,
        selectedDateOption: this.selectedDateOption,
        notificationsEnabled: this.notificationsEnabled,
        emails: [...this.emails],
        blindAvailabilityEnabled: this.blindAvailabilityEnabled,
        sendEmailAfterXResponsesEnabled: this.sendEmailAfterXResponsesEnabled,
        sendEmailAfterXResponses: this.sendEmailAfterXResponses,
      }
    },
    hasEventBeenEdited() {
      return (
        this.name !== this.initialEventData.name ||
        this.startTime !== this.initialEventData.startTime ||
        this.endTime !== this.initialEventData.endTime ||
        this.selectedDateOption !== this.initialEventData.selectedDateOption ||
        JSON.stringify(this.selectedDays) !== JSON.stringify(this.initialEventData.selectedDays) ||
        JSON.stringify(this.selectedDaysOfWeek) !== JSON.stringify(this.initialEventData.selectedDaysOfWeek) ||
        this.daysOnly !== this.initialEventData.daysOnly ||
        this.notificationsEnabled !== this.initialEventData.notificationsEnabled ||
        JSON.stringify(this.emails) !== JSON.stringify(this.initialEventData.emails) ||
        this.blindAvailabilityEnabled !== this.initialEventData.blindAvailabilityEnabled ||
        this.sendEmailAfterXResponsesEnabled !== this.initialEventData.sendEmailAfterXResponsesEnabled ||
        this.sendEmailAfterXResponses !== this.initialEventData.sendEmailAfterXResponses 
      )
    }
  },

  watch: {
    event: {
      immediate: true,
      handler() {
        this.updateFieldsFromEvent()
        this.setInitialEventData()
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
