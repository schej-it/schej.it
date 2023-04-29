<template>
  <span>
    <!-- Confirm emails dialog -->
    <ConfirmEmailsDialog
      v-model="confirmEmailsDialog"
      :respondents="respondents"
      @confirm="createCalendarInvite"
    />
    
    <div class="tw-p-4 tw-select-none" style="-webkit-touch-callout: none">
      <div class="tw-flex tw-flex-wrap">
        <!-- Times -->
        <div class="tw-w-12 tw-mt-12">
          <div class="-tw-mt-[8px]">
            <div
              v-for="(time, i) in times"
              :key="i"
              class="tw-h-5 tw-text-xs tw-pr-2 tw-text-right tw-uppercase tw-font-light"
            >
              {{ time.text }}
            </div>
          </div>
        </div>

        <div
          class="tw-flex-1 tw-flex tw-flex-col tw-overflow-x-auto tw-overflow-y-hidden tw-relative"
        >
          <!-- Days -->
          <div class="tw-flex tw-h-12">
            <div
              v-for="(day, i) in days"
              :key="i"
              class="tw-flex-1"
              style="min-width: 50px"
            >
              <div class="tw-text-center">
                <div class="tw-capitalize tw-font-light tw-text-xs">
                  {{ day.dateString }}
                </div>
                <div class="tw-text-lg tw-capitalize">
                  {{ day.dayText }}
                </div>
              </div>
            </div>
          </div>

          <!-- Calendar -->
          <div class="tw-flex tw-flex-col">
            <div class="tw-flex-1">
              <div
                id="times"
                data-long-press-delay="500"
                class="tw-flex tw-relative"
                @mouseleave="resetCurTimeslot"
              >
                <!-- Loader -->
                <div
                  v-if="
                    (alwaysShowCalendarEvents || editing) &&
                    loadingCalendarEvents
                  "
                  class="tw-absolute tw-grid tw-place-content-center tw-w-full tw-h-full tw-z-10"
                >
                  <v-progress-circular class="tw-text-blue" indeterminate />
                </div>

                <div
                  v-for="(day, d) in days"
                  :key="d"
                  class="tw-flex-1 tw-relative"
                  style="min-width: 50px"
                >
                  <!-- Timeslots -->
                  <div v-for="(time, t) in times" :key="t" class="tw-w-full">
                    <div
                      class="timeslot tw-h-5 tw-border-light-gray tw-border-r"
                      :class="timeslotClassStyle(day, time, d, t).class"
                      :style="timeslotClassStyle(day, time, d, t).style"
                      v-on="timeslotVon(d, t)"
                    ></div>
                  </div>

                  <!-- Calendar events -->
                  <div v-if="editing || alwaysShowCalendarEvents">
                    <v-fade-transition
                      v-for="(event, e) in calendarEventsByDay[d]"
                      :key="`${d}-${e}`"
                      appear
                    >
                      <div
                        class="tw-absolute tw-w-full tw-p-px tw-select-none"
                        :style="{
                          top: `calc(${event.hoursOffset} * 2 * 1.25rem)`,
                          height: `calc(${event.hoursLength} * 2 * 1.25rem)`,
                        }"
                        style="pointer-events: none"
                      >
                        <div
                          class="tw-border-blue tw-border-solid tw-border tw-w-full tw-h-full tw-text-ellipsis tw-text-xs tw-rounded tw-p-px tw-overflow-hidden"
                        >
                          <div
                            :class="`tw-text-${
                              noEventNames ? 'dark-gray' : 'blue'
                            }`"
                            class="tw-font-medium"
                          >
                            {{ noEventNames ? "BUSY" : event.summary }}
                          </div>
                        </div>
                      </div>
                    </v-fade-transition>
                  </div>

                  <!-- Scheduled event -->
                  <div
                    v-if="
                      state !== states.EDIT_AVAILABILITY &&
                      state !== states.SINGLE_AVAILABILITY &&
                      (state === states.SCHEDULE_EVENT || scheduled)
                    "
                  >
                    <div
                      v-if="
                        (dragStart && dragStart.dayIndex === d) ||
                        (!dragStart &&
                          curScheduledEvent &&
                          curScheduledEvent.dayIndex === d)
                      "
                      class="tw-absolute tw-w-full tw-p-px tw-select-none"
                      :style="scheduledEventStyle"
                      style="pointer-events: none"
                    >
                      <div
                        class="tw-border-blue tw-bg-blue tw-border-solid tw-border tw-w-full tw-h-full tw-text-ellipsis tw-text-xs tw-rounded tw-p-px tw-overflow-hidden"
                      >
                        <div class="tw-text-white tw-font-medium">
                          {{ name }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="break" v-if="isPhone"></div>

        <!-- Hint text (mobile) -->
        <div v-if="isPhone" class="tw-flex">
          <div class="tw-w-12"></div>
          <div
            class="tw-text-dark-gray tw-text-xs tw-mt-2"
            style="min-height: 1rem"
          >
            {{ hintText.mobile }}
          </div>
        </div>

        <!-- Respondents -->
        <div
          v-if="!calendarOnly"
          class="tw-py-4 tw-w-full sm:tw-pl-8 sm:tw-py-0 sm:tw-pr-0 sm:tw-pt-12 sm:tw-w-48"
        >
          <div class="tw-font-medium tw-mb-2 tw-flex tw-items-center">
            <div class="tw-mr-1 tw-text-lg">Responses</div>
            <div v-if="isCurTimeslotSelected" class="">
              {{ `(${numUsersAvailable}/${respondents.length})` }}
            </div>
            <div
              v-else
              class="tw-bg-black tw-text-white tw-font-bold tw-w-5 tw-h-5 tw-flex tw-justify-center tw-items-center tw-rounded-full tw-text-xs"
            >
              {{ respondents.length }}
            </div>
          </div>
          <div
            class="/*tw-pl-4*/ tw-text-sm tw-grid tw-grid-cols-2 tw-gap-x-2 sm:tw-block"
          >
            <div
              v-for="(user, i) in respondents"
              :key="user._id"
              class="tw-py-1 tw-flex tw-items-center tw-cursor-pointer"
              :class="respondentClass(user._id)"
              @mouseover="(e) => mouseOverRespondent(e, user._id)"
              @mouseleave="mouseLeaveRespondent"
              @click="(e) => clickRespondent(e, user._id)"
            >
              <UserAvatarContent
                v-if="!isGuest(user)"
                :user="user"
                class="tw-w-4 tw-h-4 -tw-ml-3 -tw-mr-1"
              ></UserAvatarContent>
              <v-icon v-else class="tw-ml-1 tw-mr-3" small>mdi-account</v-icon>

              <div
                class="tw-mr-1 tw-transition-all"
                :class="
                  !curTimeslotAvailability[user._id] &&
                  'tw-line-through tw-text-gray'
                "
              >
                {{ user.firstName + " " + user.lastName }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Hint text (desktop) -->
      <div v-if="!isPhone" class="tw-flex">
        <div class="tw-w-12"></div>
        <div
          class="tw-text-dark-gray tw-text-sm tw-mt-2"
          style="min-height: 1.4rem"
        >
          {{ hintText.desktop }}
        </div>
      </div>

      <div class="tw-flex" v-if="!calendarOnly">
        <div class="sm:tw-w-12"></div>

        <div
          class="tw-flex-1 tw-flex tw-items-center tw-mt-4 sm:tw-mt-0 tw-text-sm tw-justify-center sm:tw-justify-between"
        >
          <div
            class="tw-flex tw-gap-4 sm:tw-gap-8 tw-flex-row tw-justify-between tw-flex-1 sm:tw-flex-none"
          >
            <!-- Select timezone -->
            <TimezoneSelector
              v-model="curTimezone"
              :timezones="Object.keys(timezoneMap)"
            />

            <div class="tw-flex tw-justify-center tw-items-center tw-gap-1">
              <div>Show best times</div>
              <v-switch
                class="-tw-mb-1"
                v-model="showBestTimes"
                color="#219653"
                @change="onShowBestTimesChange"
              />
            </div>
          </div>

          <div
            v-if="authUser && isOwner"
            style="width: 180.16px"
            class="tw-hidden sm:tw-block"
          >
            <template v-if="state !== states.SCHEDULE_EVENT">
              <v-btn
                outlined
                class="tw-text-green tw-w-full"
                @click="scheduleEvent"
              >
                <span class="tw-mr-2">Schedule event</span>
                <v-img
                  src="@/assets/gcal_logo.png"
                  class="tw-flex-none"
                  height="20"
                  width="20"
                />
              </v-btn>
            </template>
            <template v-else>
              <v-btn
                outlined
                class="tw-text-red tw-mr-1"
                @click="cancelScheduleEvent"
              >
                Cancel
              </v-btn>
              <v-btn
                color="primary"
                @click="confirmScheduleEvent"
                :disabled="!curScheduledEvent"
              >
                Schedule
              </v-btn>
            </template>
          </div>
        </div>

        <div class="sm:tw-w-48"></div>
      </div>
    </div>
  </span>
</template>

<style scoped>
.animate-bg-color {
  transition: background-color 0.25s ease-in-out;
}

.break {
  flex-basis: 100%;
  height: 0;
}

@media only screen and (max-width: 600px) {
  ::-webkit-scrollbar {
    -webkit-appearance: none;
  }

  /* ::-webkit-scrollbar:vertical {
    width: 28px;
  } */

  ::-webkit-scrollbar:horizontal {
    height: 18px;
  }

  ::-webkit-scrollbar-thumb {
    background-color: theme("colors.gray");
    /* border-radius: 0px 0px 5px 5px; */
    /* border-radius: 10px; */
    border-top: 10px solid white;
    /* border-bottom: 10px solid white; */
  }

  ::-webkit-scrollbar-track {
    /* background-color: ; */
    border: 1px solid theme("colors.off-white");
    /* border-radius: 0px 0px 5px 5px; */
  }
}
</style>

<style>
/* Make timezone select element the same width as content */
#timezone-select {
  width: 5px;
}
</style>

<script>
import {
  timeNumToTimeText,
  dateCompare,
  getDateHoursOffset,
  post,
  isBetween,
  clamp,
  isPhone,
  utcTimeToLocalTime,
  signInGoogle,
  processCalendarEvents,
} from "@/utils"
import { mapActions, mapState } from "vuex"
import UserAvatarContent from "./UserAvatarContent.vue"
import ZigZag from "./ZigZag.vue"
import timezoneData from "@/data/timezones.json"
import TimezoneSelector from "./TimezoneSelector.vue"
import ConfirmEmailsDialog from "./ConfirmEmailsDialog.vue"
import { authTypes } from "@/constants"

export default {
  name: "ScheduleOverlap",
  props: {
    eventId: { type: String, default: "" }, // ID of event
    ownerId: { type: String, default: "" }, // ID of the owner of the event
    name: { type: String, default: "" }, // Name of event
    startTime: { type: Number, required: true }, // Start time of event
    endTime: { type: Number, required: true }, // End time of event
    duration: { type: Number, required: true }, // Duration of event
    dates: { type: Array, required: true }, // Dates of the event
    responses: { type: Object, default: () => ({}) }, // Map of user id to array of times they are available
    scheduledEvent: { type: Object, default: null }, // The scheduled event if event has already been scheduled

    loadingCalendarEvents: { type: Boolean, default: false }, // Whether we are currently loading the calendar events
    calendarEventsByDay: { type: Array, default: () => [] }, // Array of arrays of calendar events
    alwaysShowCalendarEvents: { type: Boolean, default: false }, // Whether to show calendar events all the time
    noEventNames: { type: Boolean, default: false }, // Whether to show "busy" instead of the event name
    calendarOnly: { type: Boolean, default: false }, // Whether to only show calendar and not respondents or any other controls
    interactable: { type: Boolean, default: true }, // Whether to allow user to interact with component
    showSnackbar: { type: Boolean, default: true }, // Whether to show snackbar when availability is automatically filled in
    animateTimeslotAlways: { type: Boolean, default: false }, // Whether to animate timeslots all the time
  },
  data() {
    return {
      states: {
        HEATMAP: "heatmap", // Display heatmap of availabilities
        SINGLE_AVAILABILITY: "single_availability", // Show one person's availability
        BEST_TIMES: "best_times", // Show only the times that work for most people
        EDIT_AVAILABILITY: "edit_availability", // Edit current user's availability
        SCHEDULE_EVENT: "schedule_event", // Schedule event on gcal
      },
      state: "best_times",

      max: 0, // The max amount of people available at any given time
      availability: new Set(), // The current user's availability
      availabilityAnimTimeouts: [], // Timeouts for availability animation
      availabilityAnimEnabled: false, // Whether to animate timeslots changing colors
      maxAnimTime: 1200, // Max amount of time for availability animation
      unsavedChanges: false, // If there are unsaved availability changes
      curTimeslot: { dayIndex: -1, timeIndex: -1 }, // The currently highlighted timeslot
      curTimeslotAvailability: {}, // The users available for the current timeslot
      curRespondent: "", // Id of the active respondent (set on hover)
      curRespondentSelected: false, // Whether a respondent has been selected (clicked)

      /* Variables for drag stuff */
      DRAG_TYPES: {
        ADD: "add",
        REMOVE: "remove",
      },
      timeslot: {
        width: 0,
        height: 0,
      },
      dragging: false,
      dragType: "add",
      dragStart: null,
      dragCur: null,

      /* Variables for options */
      curTimezone: this.getLocalTimezone(),

      curScheduledEvent: null, // The scheduled event represented in the form {hoursOffset, hoursLength, dayIndex}
      prevScheduledEvent: null, // The scheduled event before making changes
      scheduled: false, // Whether event has been scheduled or not
      showBestTimes: localStorage["showBestTimes"] == "true",
      confirmEmailsDialog: false,
    }
  },
  computed: {
    ...mapState(["authUser"]),
    availabilityArray() {
      /* Returns the availibility as an array */
      return [...this.availability].map((item) => new Date(item))
    },
    allowDrag() {
      return (
        this.state === this.states.EDIT_AVAILABILITY ||
        this.state === this.states.SCHEDULE_EVENT
      )
    },
    days() {
      /* Return the days that are encompassed by startDate and endDate */
      const days = []
      const daysOfWeek = ["sun", "mon", "tue", "wed", "thu", "fri", "sat"]
      const months = [
        "jan",
        "feb",
        "mar",
        "apr",
        "may",
        "jun",
        "jul",
        "aug",
        "sep",
        "oct",
        "nov",
        "dec",
      ]

      // New date representation method
      for (let date of this.dates) {
        date = new Date(date)

        days.push({
          dayText: daysOfWeek[date.getDay()],
          dateString: `${months[date.getMonth()]} ${date.getDate()}`,
          dateObject: date,
        })
      }
      return days
    },
    defaultState() {
      // Either the heatmap or the best_times state, depending on the toggle
      return this.showBestTimes ? this.states.BEST_TIMES : this.states.HEATMAP
    },
    editing() {
      // Returns whether currently in the editing state
      return this.state === this.states.EDIT_AVAILABILITY
    },
    hintText() {
      switch (this.state) {
        case this.states.EDIT_AVAILABILITY:
          return {
            desktop: "Click and drag on the calendar to edit your availability",
            mobile: "Tap and drag on the calendar to edit your availability",
          }
        case this.states.SCHEDULE_EVENT:
          return {
            desktop:
              "Click and drag on the calendar to schedule a Google Calendar event during those times",
            mobile:
              "Tap and drag on the calendar to schedule a Google Calendar event during those times",
          }
        default:
          return { desktop: "", mobile: "" }
      }
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    isOwner() {
      return this.authUser?._id === this.ownerId
    },
    respondents() {
      return Object.values(this.parsedResponses).map((r) => r.user)
    },
    selectedGuestRespondent() {
      if (!this.curRespondentSelected || !this.curRespondent) return

      const user = this.parsedResponses[this.curRespondent].user
      return this.curRespondentSelected && this.isGuest(user)
        ? this.curRespondent
        : ""
    },
    scheduledEventStyle() {
      const style = {}
      let top, height
      if (this.dragging) {
        top = this.dragStart.timeIndex
        height = this.dragCur.timeIndex - this.dragStart.timeIndex + 1
      } else {
        top = this.curScheduledEvent.hoursOffset * 2
        height = this.curScheduledEvent.hoursLength * 2
      }
      style.top = `calc(${top} * 1.25rem)`
      style.height = `calc(${height} * 1.25rem)`
      return style
    },
    parsedResponses() {
      /* Parses responses so that if _id is null (i.e. guest user), then it is set to the guest user's name */
      const parsed = {}
      for (const k of Object.keys(this.responses)) {
        const newUser = {
          ...this.responses[k].user,
          _id: k,
        }
        parsed[k] = {
          ...this.responses[k],
          user: newUser,
        }
      }
      return parsed
    },
    responsesFormatted() {
      /* Formats the responses in a map where date/time is mapped to the people that are available then */
      const formatted = new Map()
      for (const day of this.days) {
        for (const time of this.times) {
          const date = getDateHoursOffset(day.dateObject, time.hoursOffset)
          formatted.set(date.getTime(), new Set())
          for (const response of Object.values(this.parsedResponses)) {
            const index = response.availability.findIndex(
              (d) => dateCompare(d, date) === 0
            )
            if (index !== -1) {
              // TODO: determine whether I should delete the index??
              formatted.get(date.getTime()).add(response.user._id)
            }
          }
          // Update max
          if (formatted.get(date.getTime()).size > this.max) {
            this.max = formatted.get(date.getTime()).size
          }
        }
      }
      return formatted
    },
    times() {
      /* Returns the times that are encompassed by startTime and endTime */
      const times = []

      for (let i = 0; i < this.duration; ++i) {
        const utcTimeNum = this.startTime + i
        const localTimeNum = utcTimeToLocalTime(utcTimeNum, this.timezoneOffset)

        times.push({
          hoursOffset: i,
          text: timeNumToTimeText(localTimeNum),
        })
        times.push({
          hoursOffset: i + 0.5,
        })
      }

      return times
    },
    timezoneOffset() {
      return this.timezoneMap[this.curTimezone] * -1 // Multiplying by -1 because offset is flipped
    },
    timezoneMap() {
      /* Maps timezone name to the timezone offset */
      const map = timezoneData.reduce(function (map, obj) {
        map[obj.name] = obj.offset
        return map
      }, {})

      /* Adds current timezone to map if not in map */
      const localTimezone = this.getLocalTimezone()
      if (!map.hasOwnProperty(localTimezone)) {
        map[localTimezone] = new Date().getTimezoneOffset() * -1 // Multiplying by -1 because offset is flipped
      }
      return map
    },
    userHasResponded() {
      return this.authUser && this.authUser._id in this.parsedResponses
    },
    numUsersAvailable() {
      this.curTimeslot
      let numUsers = 0
      for (const key in this.curTimeslotAvailability) {
        if (this.curTimeslotAvailability[key]) numUsers++
      }
      return numUsers
    },
    isCurTimeslotSelected() {
      return (
        this.curTimeslot.dayIndex !== -1 && this.curTimeslot.timeIndex !== -1
      )
    },
  },
  methods: {
    ...mapActions(["showInfo"]),

    // -----------------------------------
    //#region Date
    // -----------------------------------

    /** Returns a date object from the dayindex and timeindex given */
    getDateFromDayTimeIndex(dayIndex, timeIndex) {
      return getDateHoursOffset(
        this.days[dayIndex].dateObject,
        this.times[timeIndex].hoursOffset
      )
    },

    /** Returns a date object from the dayindex and hoursoffset given */
    getDateFromDayHoursOffset(dayIndex, hoursOffset) {
      return getDateHoursOffset(this.days[dayIndex].dateObject, hoursOffset)
    },
    //#endregion

    // -----------------------------------
    //#region Respondent
    // -----------------------------------
    mouseOverRespondent(e, id) {
      if (!this.curRespondentSelected) {
        if (this.state === this.defaultState) {
          this.state = this.states.SINGLE_AVAILABILITY
        }

        this.curRespondent = id
      }
    },
    mouseLeaveRespondent(e) {
      if (!this.curRespondentSelected) {
        if (this.state === this.states.SINGLE_AVAILABILITY) {
          this.state = this.defaultState
        }

        this.curRespondent = ""
      }
    },
    clickRespondent(e, id) {
      if (this.state === this.defaultState) {
        this.state = this.states.SINGLE_AVAILABILITY
      }

      this.curRespondentSelected = true
      this.curRespondent = id
      e.stopPropagation()
    },
    deselectRespondent(e) {
      if (this.state === this.states.SINGLE_AVAILABILITY) {
        this.state = this.defaultState
      }

      this.curRespondentSelected = false
      this.curRespondent = ""
    },
    respondentClass(id) {
      const c = []
      if (this.curRespondent == id) {
        c.push("tw-font-bold")
      }
      return c
    },

    isGuest(user) {
      return user._id == user.firstName
    },
    //#endregion

    // -----------------------------------
    //#region Aggregate user availability
    // -----------------------------------

    /** Returns an array of respondents for the given date/time */
    getRespondentsForHoursOffset(date, hoursOffset) {
      const d = getDateHoursOffset(date, hoursOffset)
      return this.responsesFormatted.get(d.getTime())
    },
    showAvailability(d, t) {
      if (this.state === this.states.EDIT_AVAILABILITY && this.isPhone) {
        // Don't show currently selected timeslot when on phone and editing
        return
      }

      // Update current timeslot (the timeslot that has a dotted border around it)
      this.curTimeslot = { dayIndex: d, timeIndex: t }

      if (this.state === this.states.EDIT_AVAILABILITY || this.curRespondent) {
        // Don't show availability when editing or when respondent is selected
        return
      }

      // Update current timeslot availability to show who is available for the given timeslot
      const available = this.getRespondentsForHoursOffset(
        this.days[d].dateObject,
        this.times[t].hoursOffset
      )
      for (const respondent of this.respondents) {
        if (available.has(respondent._id)) {
          this.curTimeslotAvailability[respondent._id] = true
        } else {
          this.curTimeslotAvailability[respondent._id] = false
        }
      }
    },
    //#endregion

    // -----------------------------------
    //#region Current user availability
    // -----------------------------------
    resetCurUserAvailability() {
      /* resets cur user availability to the response stored on the server */
      this.availability = new Set()
      if (this.userHasResponded) {
        this.populateUserAvailability(this.authUser._id)
      }
    },
    populateUserAvailability(id) {
      /* Populates the availability set for the auth user from the responses object stored on the server */
      this.responses[id].availability.forEach((item) =>
        this.availability.add(new Date(item).getTime())
      )
      this.$nextTick(() => (this.unsavedChanges = false))
    },
    setAvailabilityAutomatically() {
      /* Constructs the availability array using calendarEvents array */
      // This is not a computed property because we should be able to change it manually from what it automatically fills in
      const tmpAvailability = new Set()
      for (const d in this.days) {
        const day = this.days[d]
        for (const time of this.times) {
          // Check if there exists a calendar event that overlaps [time, time+0.5]
          const startDate = getDateHoursOffset(day.dateObject, time.hoursOffset)
          const endDate = getDateHoursOffset(
            day.dateObject,
            time.hoursOffset + 0.5
          )
          const index = this.calendarEventsByDay[d].findIndex((e) => {
            return (
              (dateCompare(e.startDate, startDate) < 0 &&
                dateCompare(e.endDate, startDate) > 0) ||
              (dateCompare(e.startDate, endDate) < 0 &&
                dateCompare(e.endDate, endDate) > 0) ||
              (dateCompare(e.startDate, startDate) == 0 &&
                dateCompare(e.endDate, endDate) == 0)
            )
          })
          if (index === -1) {
            tmpAvailability.add(startDate.getTime())
          }
        }
      }
      this.animateAvailability(tmpAvailability)
    },
    animateAvailability(availability) {
      /* Animate the filling out of availability using setTimeout */

      this.availabilityAnimEnabled = true
      this.availabilityAnimTimeouts = []

      let msPerBlock = 25
      if (availability.size * msPerBlock > this.maxAnimTime) {
        msPerBlock = this.maxAnimTime / availability.size
      }

      let i = 0
      for (const a of availability) {
        const index = i
        const timeout = setTimeout(() => {
          this.availability.add(a)
          this.availability = new Set(this.availability)

          if (index == availability.size - 1) {
            setTimeout(() => {
              this.availabilityAnimEnabled = false
              if (this.showSnackbar) {
                this.showInfo(
                  "Your availability has been set automatically using your Google Calendar!"
                )
              }
            }, 500)
          }
        }, i * msPerBlock)

        this.availabilityAnimTimeouts.push(timeout)
        i++
      }
    },
    stopAvailabilityAnim() {
      for (const timeout of this.availabilityAnimTimeouts) {
        clearTimeout(timeout)
      }
      this.availabilityAnimEnabled = false
    },
    async submitAvailability(name = "") {
      const payload = { availability: this.availabilityArray }
      if (this.authUser) {
        payload.guest = false
      } else {
        payload.guest = true
        payload.name = name
      }
      await post(`/events/${this.eventId}/response`, payload)
      this.$emit("refreshEvent")
      this.unsavedChanges = false
    },
    //#endregion

    // -----------------------------------
    //#region Timeslot
    // -----------------------------------
    setTimeslotSize() {
      /* Gets the dimensions of each timeslot and assigns it to the timeslot variable */
      ;({ width: this.timeslot.width, height: this.timeslot.height } = document
        .querySelector(".timeslot")
        .getBoundingClientRect())
    },
    timeslotClassStyle(day, time, d, t) {
      /* Returns a class string and style object for the given timeslot div */
      let c = ""
      const s = {}
      // Animation
      if (this.animateTimeslotAlways || this.availabilityAnimEnabled) {
        c += "animate-bg-color "
      }

      // Border style
      if (this.curTimeslot.dayIndex === d && this.curTimeslot.timeIndex === t) {
        // Dashed border for currently selected timeslot
        c += "tw-border tw-border-dashed tw-border-black tw-z-10 "
      } else {
        // Normal border
        if (!("text" in time)) c += "tw-border-b "
        if (d === 0) c += "tw-border-l tw-border-l-gray "
        if (d === this.days.length - 1) c += "tw-border-r-gray "
        if (t === 0) c += "tw-border-t tw-border-t-gray "
        if (t === this.times.length - 1) c += "tw-border-b-gray "
      }

      // Fill style
      if (this.state === this.states.EDIT_AVAILABILITY) {
        // Show only current user availability
        const inDragRange = this.inDragRange(d, t)
        if (inDragRange) {
          // Set style if drag range goes over the current timeslot
          if (this.dragType === this.DRAG_TYPES.ADD) {
            c += "tw-bg-avail-green-300 "
          } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
          }
        } else {
          // Otherwise just show the current availability
          const date = getDateHoursOffset(day.dateObject, time.hoursOffset)
          if (this.availability.has(date.getTime())) {
            c += "tw-bg-avail-green-300 "
          }
        }
      }
      if (this.state === this.states.SINGLE_AVAILABILITY) {
        // Show only the currently selected respondent's availability
        const respondent = this.curRespondent
        const respondents = this.getRespondentsForHoursOffset(
          day.dateObject,
          time.hoursOffset
        )
        if (respondents.has(respondent)) {
          c += "tw-bg-avail-green-300 "
        }
      }
      if (
        this.state === this.states.BEST_TIMES ||
        (this.state === this.states.SCHEDULE_EVENT &&
          this.defaultState === this.states.BEST_TIMES)
      ) {
        const numRespondents = this.getRespondentsForHoursOffset(
          day.dateObject,
          time.hoursOffset
        ).size
        if (this.max > 0 && numRespondents === this.max) {
          // Only set timeslot to green for the times that most people are available
          const green = "#12B981"

          s.backgroundColor = green
        }
      }
      if (
        this.state === this.states.HEATMAP ||
        (this.state === this.states.SCHEDULE_EVENT &&
          this.defaultState === this.states.HEATMAP)
      ) {
        // Show everyone's availability
        const numRespondents = this.getRespondentsForHoursOffset(
          day.dateObject,
          time.hoursOffset
        ).size
        if (numRespondents > 0) {
          // Determine color of timeslot based on number of people available
          const frac = numRespondents / this.max
          const green = "#12B981"
          let alpha = (frac * (255 - 30))
            .toString(16)
            .toUpperCase()
            .substring(0, 2)
          if (frac == 1) alpha = "FF"

          s.backgroundColor = green + alpha
        }
      }
      return { class: c, style: s }
    },
    timeslotVon(d, t) {
      if (this.interactable) {
        return {
          click: () => this.showAvailability(d, t),
          mouseover: () => this.showAvailability(d, t),
        }
      }
      return {}
    },
    resetCurTimeslot() {
      this.curTimeslotAvailability = {}
      for (const respondent of this.respondents) {
        this.curTimeslotAvailability[respondent._id] = true
      }
      this.curTimeslot = { dayIndex: -1, timeIndex: -1 }

      // End drag if mouse left time grid
      this.endDrag()
    },
    //#endregion

    // -----------------------------------
    //#region Editing
    // -----------------------------------
    startEditing() {
      this.state = this.states.EDIT_AVAILABILITY
      // console.log("start editing!!!", this.state)
    },
    stopEditing() {
      this.state = this.defaultState
      this.stopAvailabilityAnim()
    },
    //#endregion

    // -----------------------------------
    //#region Schedule event
    // -----------------------------------
    scheduleEvent() {
      this.state = this.states.SCHEDULE_EVENT
    },
    cancelScheduleEvent() {
      this.state = this.defaultState
    },
    confirmScheduleEvent() {
      this.confirmEmailsDialog = true
    },

    /** Creates a google calendar invite and officially schedules the event on the server */
    createCalendarInvite(emails) {
      const { dayIndex, hoursOffset, hoursLength } = this.curScheduledEvent
      const payload = {
        startDate: this.getDateFromDayHoursOffset(dayIndex, hoursOffset),
        endDate: this.getDateFromDayHoursOffset(
          dayIndex,
          hoursOffset + hoursLength
        ),
        attendeeEmails: emails.filter(
          (email) => email.length > 0 && email !== this.authUser.email
        ),
        curScheduledEvent: this.curScheduledEvent,
      }

      // Schedule event on backend
      post(`/events/${this.eventId}/schedule`, payload)
        .then(() => {
          this.confirmEmailsDialog = false
          this.prevScheduledEvent = this.curScheduledEvent // Needed so the scheduled event stays there after exiting scheduling state
          this.scheduled = true
          this.state = this.defaultState
          this.showInfo("Event has been scheduled!")
        })
        .catch((err) => {
          console.error(err)
          // If calendar edit permission not granted, ask for it
          if (err.error.code === 401 || err.error.code === 403) {
            signInGoogle({
              state: {
                type: authTypes.EVENT_SCHEDULE,
                eventId: this.eventId,
                payload,
              },
              requestEditCalendarPermission: true,
            })
          }
        })
    },

    /** Creates a calendar invite with the given payload (used right after enabling calendar permissions) */
    createCalendarInviteFromPayload(payload) {
      post(`/events/${this.eventId}/schedule`, payload).then(() => {
        this.curScheduledEvent = payload.curScheduledEvent
        this.prevScheduledEvent = this.curScheduledEvent // Needed so the scheduled event stays there after exiting scheduling state
        this.scheduled = true
        this.state = this.defaultState
        this.showInfo("Event has been scheduled!")
      })
    },

    /** Sets curScheduledEvent by reformatting the scheduledEvent stored in the server */
    processScheduledEvent() {
      const eventsByDay = processCalendarEvents(this.dates, this.duration, [
        this.scheduledEvent,
      ])
      for (const d in eventsByDay) {
        if (eventsByDay[d].length > 0) {
          const event = eventsByDay[d][0]
          this.curScheduledEvent = {
            dayIndex: parseInt(d),
            hoursOffset: event.hoursOffset,
            hoursLength: event.hoursLength,
          }
          this.prevScheduledEvent = this.curScheduledEvent
          this.scheduled = true

          break
        }
      }
    },
    //#endregion

    // -----------------------------------
    //#region Drag Stuff
    // -----------------------------------
    normalizeXY(e) {
      /* Normalize the touch event to be relative to element */
      let pageX, pageY
      if ("touches" in e) {
        // is a touch event
        ;({ pageX, pageY } = e.touches[0])
      } else {
        // is a mouse event
        ;({ pageX, pageY } = e)
      }
      const { left, top } = e.currentTarget.getBoundingClientRect()
      const x = pageX - left
      const y = pageY - top - window.scrollY
      return { x, y }
    },
    getDateFromXY(x, y) {
      /* Returns a date for the timeslot we are currently hovering over given the x and y position */
      const { width, height } = this.timeslot
      let dayIndex = Math.floor(x / width)
      let timeIndex = Math.floor(y / height)
      dayIndex = clamp(dayIndex, 0, this.days.length - 1)
      timeIndex = clamp(timeIndex, 0, this.times.length - 1)
      return {
        dayIndex,
        timeIndex,
        date: getDateHoursOffset(
          this.days[dayIndex].dateObject,
          this.times[timeIndex].hoursOffset
        ),
      }
    },
    endDrag() {
      if (!this.allowDrag) return

      if (!this.dragStart || !this.dragCur) return

      if (this.state === this.states.EDIT_AVAILABILITY) {
        // Update availability set based on drag region
        let dayInc =
          (this.dragCur.dayIndex - this.dragStart.dayIndex) /
          Math.abs(this.dragCur.dayIndex - this.dragStart.dayIndex)
        let timeInc =
          (this.dragCur.timeIndex - this.dragStart.timeIndex) /
          Math.abs(this.dragCur.timeIndex - this.dragStart.timeIndex)
        if (isNaN(dayInc)) dayInc = 1
        if (isNaN(timeInc)) timeInc = 1
        let d = this.dragStart.dayIndex
        while (d != this.dragCur.dayIndex + dayInc) {
          let t = this.dragStart.timeIndex
          while (t != this.dragCur.timeIndex + timeInc) {
            const date = getDateHoursOffset(
              this.days[d].dateObject,
              this.times[t].hoursOffset
            )
            if (this.dragType === this.DRAG_TYPES.ADD) {
              this.availability.add(date.getTime())
            } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
              this.availability.delete(date.getTime())
            }
            t += timeInc
          }
          d += dayInc
        }
        this.availability = new Set(this.availability)
      } else if (this.state === this.states.SCHEDULE_EVENT) {
        // Update scheduled event
        const dayIndex = this.dragStart.dayIndex
        const hoursOffset = this.dragStart.timeIndex / 2
        const hoursLength =
          (this.dragCur.timeIndex - this.dragStart.timeIndex + 1) / 2

        if (hoursLength > 0) {
          this.curScheduledEvent = { dayIndex, hoursOffset, hoursLength }
        } else {
          this.curScheduledEvent = null
        }
      }

      // Set dragging defaults
      this.dragging = false
      this.dragStart = null
      this.dragCur = null
    },
    inDragRange(dayIndex, timeIndex) {
      /* Returns whether the given day and time index is within the drag range */
      if (this.dragging) {
        return (
          (isBetween(
            dayIndex,
            this.dragStart.dayIndex,
            this.dragCur.dayIndex
          ) ||
            isBetween(
              dayIndex,
              this.dragCur.dayIndex,
              this.dragStart.dayIndex
            )) &&
          (isBetween(
            timeIndex,
            this.dragStart.timeIndex,
            this.dragCur.timeIndex
          ) ||
            isBetween(
              timeIndex,
              this.dragCur.timeIndex,
              this.dragStart.timeIndex
            ))
        )
      }
    },
    moveDrag(e) {
      if (!this.allowDrag) return

      e.preventDefault()
      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      )
      this.dragCur = { dayIndex, timeIndex }
    },
    startDrag(e) {
      if (!this.allowDrag) return

      this.dragging = true

      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      )
      this.dragStart = { dayIndex, timeIndex }
      this.dragCur = { dayIndex, timeIndex }
      // Set drag type
      if (this.availability.has(date.getTime())) {
        this.dragType = this.DRAG_TYPES.REMOVE
      } else {
        this.dragType = this.DRAG_TYPES.ADD
      }
    },
    //#endregion

    // -----------------------------------
    //#region Options
    // -----------------------------------
    getLocalTimezone() {
      const split = new Date(this.dates[0])
        .toLocaleTimeString("en-us", { timeZoneName: "short" })
        .split(" ")
      const localTimezone = split[split.length - 1]

      return localTimezone
    },
    onShowBestTimesChange() {
      localStorage["showBestTimes"] = this.showBestTimes
      if (
        this.state == this.states.BEST_TIMES ||
        this.state == this.states.HEATMAP
      )
        this.state = this.defaultState
    },
    //#endregion
  },
  watch: {
    availability() {
      this.unsavedChanges = true
    },
    state(nextState, prevState) {
      if (prevState === this.states.SCHEDULE_EVENT) {
        this.curScheduledEvent = this.prevScheduledEvent
      }

      if (nextState === this.states.SCHEDULE_EVENT) {
        this.prevScheduledEvent = this.curScheduledEvent
      }
    },
    calendarEvents: {
      handler() {
        //if (!this.userHasResponded && !this.calendarOnly) this.setAvailability()
      },
    },
    respondents: {
      immediate: true,
      handler() {
        this.curTimeslotAvailability = {}
        for (const respondent of this.respondents) {
          this.curTimeslotAvailability[respondent._id] = true
        }
      },
    },
  },
  created() {
    this.resetCurUserAvailability()
    this.processScheduledEvent()

    addEventListener("click", this.deselectRespondent)
  },
  mounted() {
    // Set initial state to best_times or heatmap depending on show best times toggle.
    this.state = this.showBestTimes ? "best_times" : "heatmap"

    // Get timeslot size
    this.setTimeslotSize()
    window.addEventListener("resize", this.setTimeslotSize)
    if (!this.calendarOnly) {
      const timesEl = document.getElementById("times")
      if (isPhone(this.$vuetify)) {
        timesEl.addEventListener("touchstart", this.startDrag)
        timesEl.addEventListener("touchmove", this.moveDrag)
        timesEl.addEventListener("touchend", this.endDrag)
        timesEl.addEventListener("touchcancel", this.endDrag)
      } else {
        timesEl.addEventListener("mousedown", this.startDrag)
        timesEl.addEventListener("mousemove", this.moveDrag)
        timesEl.addEventListener("mouseup", this.endDrag)
      }
    }
  },
  beforeDestroy() {
    removeEventListener("click", this.deselectRespondent)
  },
  components: {
    UserAvatarContent,
    ZigZag,
    TimezoneSelector,
    ConfirmEmailsDialog,
  },
}
</script>
