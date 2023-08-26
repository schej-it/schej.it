<template>
  <span>
    <div class="tw-select-none tw-p-4" style="-webkit-touch-callout: none">
      <div class="tw-flex tw-flex-wrap sm:tw-flex-nowrap">
        <div class="tw-flex tw-grow tw-overflow-hidden">
          <!-- Times -->
          <div class="tw-mt-12 tw-w-12 tw-flex-none">
            <div class="-tw-mt-[8px]">
              <div
                v-for="(time, i) in times"
                :key="i"
                class="tw-h-5 tw-pr-2 tw-text-right tw-text-xs tw-font-light tw-uppercase"
              >
                {{ time.text }}
              </div>
            </div>
          </div>

          <div class="tw-grow tw-overflow-hidden">
            <div class="tw-relative tw-overflow-hidden">
              <div
                ref="calendar"
                @scroll="onCalendarScroll"
                class="tw-relative tw-flex tw-flex-col tw-overflow-x-auto tw-overflow-y-hidden"
              >
                <!-- Days -->
                <div
                  class="tw-z-10 tw-flex tw-h-12 tw-items-center tw-bg-white"
                >
                  <div
                    v-for="(day, i) in days"
                    :key="i"
                    class="tw-flex-1 tw-bg-white"
                    style="min-width: 50px"
                  >
                    <div class="tw-text-center">
                      <div
                        v-if="isSpecificDates"
                        class="tw-text-xs tw-font-light tw-capitalize tw-text-very-dark-gray"
                      >
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
                      class="tw-relative tw-flex"
                      @mouseleave="resetCurTimeslot"
                    >
                      <!-- Loader -->
                      <div
                        v-if="
                          (alwaysShowCalendarEvents || editing) &&
                          loadingCalendarEvents
                        "
                        class="tw-absolute tw-z-10 tw-grid tw-h-full tw-w-full tw-place-content-center"
                      >
                        <v-progress-circular
                          class="tw-text-blue"
                          indeterminate
                        />
                      </div>

                      <div
                        v-for="(day, d) in days"
                        :key="d"
                        class="tw-relative tw-flex-1"
                        style="min-width: 50px"
                      >
                        <!-- Timeslots -->
                        <div
                          v-for="(time, t) in times"
                          :key="t"
                          class="tw-w-full"
                        >
                          <div
                            class="timeslot tw-h-5 tw-border-r tw-border-[#DDDDDD88]"
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
                              class="tw-absolute tw-w-full tw-select-none tw-p-px"
                              :style="{
                                top: `calc(${event.hoursOffset} * 2 * 1.25rem)`,
                                height: `calc(${event.hoursLength} * 2 * 1.25rem)`,
                              }"
                              style="pointer-events: none"
                            >
                              <div
                                class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-rounded tw-border tw-border-solid tw-border-blue tw-p-px tw-text-xs"
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
                        <div v-if="state === states.SCHEDULE_EVENT">
                          <div
                            v-if="
                              (dragStart && dragStart.dayIndex === d) ||
                              (!dragStart &&
                                curScheduledEvent &&
                                curScheduledEvent.dayIndex === d)
                            "
                            class="tw-absolute tw-w-full tw-select-none tw-p-px"
                            :style="scheduledEventStyle"
                            style="pointer-events: none"
                          >
                            <div
                              class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-rounded tw-border tw-border-solid tw-border-blue tw-bg-blue tw-p-px tw-text-xs"
                            >
                              <div class="tw-font-medium tw-text-white">
                                {{ event.name }}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <ZigZag
                v-if="showLeftZigZag"
                left
                class="tw-absolute tw-left-0 tw-top-0 tw-h-full tw-w-3"
              />
              <ZigZag
                v-if="showRightZigZag"
                right
                class="tw-absolute tw-right-0 tw-top-0 tw-h-full tw-w-3"
              />
            </div>

            <!-- Hint text (desktop) -->
            <div v-if="!isPhone && showHintText" class="tw-flex">
              <div
                class="tw-mt-2 tw-text-sm tw-text-dark-gray"
                style="min-height: 1.4rem"
              >
                {{ hintText.desktop }}
              </div>
            </div>

            <ToolRow
              v-if="!calendarOnly && !isPhone"
              :state="state"
              :states="states"
              :cur-timezone.sync="curTimezone"
              :timezone-map="timezoneMap"
              :show-best-times.sync="showBestTimes"
              :is-owner="isOwner"
              :cur-scheduled-event="curScheduledEvent"
              :is-weekly="isWeekly"
              :calendar-permission-granted="calendarPermissionGranted"
              :week-offset="weekOffset"
              :num-responses="respondents.length"
              @update:weekOffset="(val) => $emit('update:weekOffset', val)"
              @onShowBestTimesChange="onShowBestTimesChange"
              @scheduleEvent="scheduleEvent"
              @cancelScheduleEvent="cancelScheduleEvent"
              @confirmScheduleEvent="confirmScheduleEvent"
            />
          </div>
        </div>

        <div class="break" v-if="isPhone"></div>

        <!-- Hint text (mobile) -->
        <div v-if="isPhone && showHintText" class="tw-flex">
          <div class="tw-w-12"></div>
          <div
            class="tw-mt-2 tw-text-xs tw-text-dark-gray"
            style="min-height: 2rem"
          >
            {{ hintText.mobile }}
          </div>
        </div>

        <!-- Respondents -->
        <div
          v-if="!calendarOnly"
          class="tw-w-full tw-py-4 sm:tw-w-48 sm:tw-flex-none sm:tw-py-0 sm:tw-pl-8 sm:tw-pr-0 sm:tw-pt-12"
        >
          <div v-if="state == states.EDIT_AVAILABILITY">
            <CalendarAccounts
              v-if="calendarPermissionGranted"
              :toggleState="true"
              :eventId="event._id"
              :calendar-events-map="calendarEventsMap"
            ></CalendarAccounts>
          </div>

          <div v-else>
            <div class="tw-mb-2 tw-flex tw-items-center tw-font-medium">
              <div class="tw-mr-1 tw-text-lg">Responses</div>
              <div class="tw-font-normal">
                <template v-if="curRespondents.length === 0">
                  {{
                    isCurTimeslotSelected
                      ? `(${numUsersAvailable}/${respondents.length})`
                      : `(${respondents.length})`
                  }}
                </template>
                <template v-else>
                  {{
                    isCurTimeslotSelected
                      ? `(${numCurRespondentsAvailable}/${curRespondents.length})`
                      : `(${curRespondents.length})`
                  }}
                </template>
              </div>
            </div>
            <div
              class="tw-grid tw-grid-cols-2 tw-gap-x-2 tw-text-sm sm:tw-block"
            >
              <template v-if="respondents.length === 0">
                <div class="tw-text-very-dark-gray">No responses yet!</div>
              </template>
              <template v-else>
                <div
                  v-for="(user, i) in respondents"
                  :key="user._id"
                  class="tw-flex tw-cursor-pointer tw-items-center tw-py-1"
                  @mouseover="(e) => mouseOverRespondent(e, user._id)"
                  @mouseleave="mouseLeaveRespondent"
                  @click="(e) => clickRespondent(e, user._id)"
                >
                  <UserAvatarContent
                    v-if="!isGuest(user)"
                    :user="user"
                    class="-tw-ml-3 -tw-mr-1 tw-h-4 tw-w-4"
                  ></UserAvatarContent>
                  <v-icon v-else class="tw-ml-1 tw-mr-3" small
                    >mdi-account</v-icon
                  >

                  <div
                    class="tw-mr-1 tw-transition-all"
                    :class="respondentClass(user._id)"
                  >
                    {{ user.firstName + " " + user.lastName }}
                  </div>
                </div>
              </template>
            </div>
          </div>
        </div>
      </div>

      <ToolRow
        v-if="!calendarOnly && isPhone"
        :state="state"
        :states="states"
        :cur-timezone.sync="curTimezone"
        :timezone-map="timezoneMap"
        :show-best-times.sync="showBestTimes"
        :is-owner="isOwner"
        :cur-scheduled-event="curScheduledEvent"
        :is-weekly="isWeekly"
        :calendar-permission-granted="calendarPermissionGranted"
        :week-offset="weekOffset"
        :num-responses="respondents.length"
        @update:weekOffset="(val) => $emit('update:weekOffset', val)"
        @onShowBestTimesChange="onShowBestTimesChange"
        @scheduleEvent="scheduleEvent"
        @cancelScheduleEvent="cancelScheduleEvent"
        @confirmScheduleEvent="confirmScheduleEvent"
      />
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
  splitCalendarEventsByDay,
  dateToDowDate,
} from "@/utils"
import { eventTypes } from "@/constants"
import { mapActions, mapState } from "vuex"
import UserAvatarContent from "@/components/UserAvatarContent.vue"
import CalendarAccounts from "@/components/settings/CalendarAccounts.vue"
import ZigZag from "./ZigZag.vue"
import timezoneData from "@/data/timezones.json"
import TimezoneSelector from "./TimezoneSelector.vue"
import ConfirmDetailsDialog from "./ConfirmDetailsDialog.vue"
import ToolRow from "./ToolRow.vue"

export default {
  name: "ScheduleOverlap",
  props: {
    event: { type: Object, required: true },

    loadingCalendarEvents: { type: Boolean, default: false }, // Whether we are currently loading the calendar events
    calendarEventsMap: { type: Object, default: () => {} }, // Object of different users' calendar events
    sampleCalendarEventsByDay: { type: Array, required: false }, // Sample calendar events to use for example calendars
    calendarPermissionGranted: { type: Boolean, default: false }, // Whether user has granted google calendar permissions

    weekOffset: { type: Number, default: 0 }, // Week offset used for displaying calendar events on weekly schejs

    alwaysShowCalendarEvents: { type: Boolean, default: false }, // Whether to show calendar events all the time
    noEventNames: { type: Boolean, default: false }, // Whether to show "busy" instead of the event name
    calendarOnly: { type: Boolean, default: false }, // Whether to only show calendar and not respondents or any other controls
    interactable: { type: Boolean, default: true }, // Whether to allow user to interact with component
    showSnackbar: { type: Boolean, default: true }, // Whether to show snackbar when availability is automatically filled in
    animateTimeslotAlways: { type: Boolean, default: false }, // Whether to animate timeslots all the time
    showHintText: { type: Boolean, default: true }, // Whether to show the hint text telling user what to do
  },
  data() {
    return {
      states: {
        HEATMAP: "heatmap", // Display heatmap of availabilities
        SINGLE_AVAILABILITY: "single_availability", // Show one person's availability
        SUBSET_AVAILABILITY: "subset_availability", // Show availability for a subset of people
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
      curRespondents: [], // Id of currently selected respondents (set on click)

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
      showBestTimes: localStorage["showBestTimes"] == "true",

      /* Variables for scrolling */
      calendarScrollLeft: 0, // The current scroll position of the calendar
      calendarMaxScroll: 0, // The maximum scroll amount of the calendar, scrolling to this point means we have scrolled to the end
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
    calendarEventsByDay() {
      /** If this is an example calendar */
      if (this.sampleCalendarEventsByDay) return this.sampleCalendarEventsByDay

      /** If the user isn't logged in */
      if (!this.authUser) return []

      let events = []
      /** Adds events from calendar accounts that are enabled */
      for (const id in this.authUser.calendarAccounts) {
        if (this.authUser.calendarAccounts[id].enabled) {
          events = events.concat(
            this.calendarEventsMap.hasOwnProperty(id)
              ? this.calendarEventsMap[id].calendarEvents
              : []
          )
        }
      }

      const eventsCopy = JSON.parse(JSON.stringify(events))
      const calendarEventsByDay = splitCalendarEventsByDay(
        this.event,
        eventsCopy,
        this.weekOffset
      )

      return calendarEventsByDay
    },
    curRespondentsSet() {
      return new Set(this.curRespondents)
    },
    /** Returns the max number of people in the curRespondents array available at any given time */
    curRespondentsMax() {
      let max = 0
      for (const day of this.days) {
        for (const time of this.times) {
          const num = [
            ...this.getRespondentsForHoursOffset(
              day.dateObject,
              time.hoursOffset
            ),
          ].filter((r) => this.curRespondentsSet.has(r)).length

          if (num > max) max = num
        }
      }
      return max
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

      for (let date of this.event.dates) {
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
            desktop:
              "Click and drag on the calendar to edit your availability. Green means available.",
            mobile:
              "Tap and drag on the calendar to edit your availability. Drag the dates at the top to scroll. Green means available.",
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
      return this.authUser?._id === this.event.ownerId
    },
    isSpecificDates() {
      return this.event.type === eventTypes.SPECIFIC_DATES || !this.event.type
    },
    isWeekly() {
      return this.event.type === eventTypes.DOW
    },
    respondents() {
      return Object.values(this.parsedResponses).map((r) => r.user)
    },
    selectedGuestRespondent() {
      if (this.curRespondents.length !== 1) return ""

      const user = this.parsedResponses[this.curRespondents[0]].user
      return this.isGuest(user) ? user._id : ""
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
      for (const k of Object.keys(this.event.responses)) {
        const newUser = {
          ...this.event.responses[k].user,
          _id: k,
        }
        parsed[k] = {
          ...this.event.responses[k],
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

      for (let i = 0; i < this.event.duration; ++i) {
        const utcTimeNum = this.event.startTime + i
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
    numCurRespondentsAvailable() {
      this.curTimeslot
      let numUsers = 0
      for (const key in this.curTimeslotAvailability) {
        if (
          this.curTimeslotAvailability[key] &&
          this.curRespondentsSet.has(key)
        )
          numUsers++
      }
      return numUsers
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

    showLeftZigZag() {
      return this.calendarScrollLeft > 0
    },
    showRightZigZag() {
      return Math.ceil(this.calendarScrollLeft) < this.calendarMaxScroll
    },
  },
  methods: {
    ...mapActions(["showInfo", "showError"]),

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
      if (this.curRespondents.length === 0) {
        if (this.state === this.defaultState) {
          this.state = this.states.SINGLE_AVAILABILITY
        }

        this.curRespondent = id
      }
    },
    mouseLeaveRespondent(e) {
      if (this.curRespondents.length === 0) {
        if (this.state === this.states.SINGLE_AVAILABILITY) {
          this.state = this.defaultState
        }

        this.curRespondent = ""
      }
    },
    clickRespondent(e, id) {
      this.state = this.states.SUBSET_AVAILABILITY
      this.curRespondent = ""

      if (this.curRespondentsSet.has(id)) {
        // Remove id
        this.curRespondents = this.curRespondents.filter((r) => r != id)

        // Go back to default state if all users deselected
        if (this.curRespondents.length === 0) {
          this.state = this.defaultState
        }
      } else {
        // Add id
        this.curRespondents.push(id)
      }

      e.stopPropagation()
    },
    deselectRespondents(e) {
      // Don't deselect respondents if toggled best times
      if (
        e.target?.previousElementSibling?.id === "show-best-times-toggle" ||
        e.target?.firstChild?.firstChild?.id === "show-best-times-toggle"
      )
        return

      if (this.state === this.states.SUBSET_AVAILABILITY) {
        this.state = this.defaultState
      }

      this.curRespondents = []
    },
    respondentClass(id) {
      const c = []
      if (this.curRespondent == id || this.curRespondentsSet.has(id)) {
        c.push("tw-font-bold")
      } else if (this.curRespondents.length > 0) {
        c.push("tw-text-gray")
      }

      if (!this.curTimeslotAvailability[id]) {
        c.push("tw-line-through")
        c.push("tw-text-gray")
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

    /** Returns a set of respondents for the given date/time */
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
      this.event.responses[id].availability.forEach((item) =>
        this.availability.add(new Date(item).getTime())
      )
      this.$nextTick(() => (this.unsavedChanges = false))
    },
    setAvailabilityAutomatically() {
      /* Constructs the availability array using calendarEvents array */
      // This is not a computed property because we should be able to change it manually from what it automatically fills in
      this.availability = new Set()
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
              this.unsavedChanges = false
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
      await post(`/events/${this.event._id}/response`, payload)
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
            s.backgroundColor = "#00994C88"
          } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
          }
        } else {
          // Otherwise just show the current availability
          const date = getDateHoursOffset(day.dateObject, time.hoursOffset)
          if (this.availability.has(date.getTime())) {
            s.backgroundColor = "#00994C88"
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
          s.backgroundColor = "#00994C88"
        }
      }

      if (
        this.state === this.states.BEST_TIMES ||
        this.state === this.states.HEATMAP ||
        this.state === this.states.SCHEDULE_EVENT ||
        this.state === this.states.SUBSET_AVAILABILITY
      ) {
        let numRespondents
        let max

        if (
          this.state === this.states.BEST_TIMES ||
          this.state === this.states.HEATMAP ||
          this.state === this.states.SCHEDULE_EVENT
        ) {
          numRespondents = this.getRespondentsForHoursOffset(
            day.dateObject,
            time.hoursOffset
          ).size
          max = this.max
        } else if (this.state === this.states.SUBSET_AVAILABILITY) {
          numRespondents = [
            ...this.getRespondentsForHoursOffset(
              day.dateObject,
              time.hoursOffset
            ),
          ].filter((r) => this.curRespondentsSet.has(r)).length

          max = this.curRespondentsMax
        }

        if (this.defaultState === this.states.BEST_TIMES) {
          if (max > 0 && numRespondents === max) {
            // Only set timeslot to green for the times that most people are available
            const green = "#00994C"
            s.backgroundColor = green
          }
        } else if (this.defaultState === this.states.HEATMAP) {
          if (numRespondents > 0) {
            // Determine color of timeslot based on number of people available
            const frac = numRespondents / max
            const green = "#00994C"
            let alpha = (frac * (255 - 30))
              .toString(16)
              .toUpperCase()
              .substring(0, 2)
            if (frac == 1) alpha = "FF"

            s.backgroundColor = green + alpha
          }
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

    /** Redirect user to Google Calendar to finish the creation of the event */
    confirmScheduleEvent() {
      // Get start date, and end date from the area that the user has dragged out
      const { dayIndex, hoursOffset, hoursLength } = this.curScheduledEvent
      let startDate = this.getDateFromDayHoursOffset(dayIndex, hoursOffset)
      let endDate = this.getDateFromDayHoursOffset(
        dayIndex,
        hoursOffset + hoursLength
      )

      if (this.isWeekly) {
        // Determine offset based on current day of the week.
        // People expect the event to be scheduled in the future, not the past, which is why this check exists
        let offset = 0
        if (new Date().getDay() > startDate.getDay()) {
          offset = 1
        }

        // Transform startDate and endDate to be the current week offset
        startDate = dateToDowDate(this.event.dates, startDate, offset, true)
        endDate = dateToDowDate(this.event.dates, endDate, offset, true)
      }

      // Format email string separated by commas
      const emails = this.respondents.map((r) => {
        // Return email if they are not a guest, otherwise return their name
        if (r.email.length > 0) {
          return r.email
        } else {
          return `${r.firstName} (no email)`
        }
      })
      const emailsString = encodeURIComponent(emails.join(","))

      // Format start and end date to be in the format required by gcal (remove -, :, and .000)
      const start = startDate.toISOString().replace(/([-:]|\.000)/g, "")
      const end = endDate.toISOString().replace(/([-:]|\.000)/g, "")

      // Construct Google Calendar event creation template url
      const url = `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(
        this.event.name
      )}&dates=${start}/${end}&details=${encodeURIComponent(
        "\n\nThis event was scheduled with schej: https://schej.it/e/"
      )}${this.event._id}&add=${emailsString}`

      // Navigate to url and reset state
      window.open(url, "_blank")
      this.state = this.defaultState
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
      if (e.touches?.length > 1) return // If dragging with more than one finger

      e.preventDefault()
      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      )
      this.dragCur = { dayIndex, timeIndex }
    },
    startDrag(e) {
      if (!this.allowDrag) return
      if (e.touches?.length > 1) return // If dragging with more than one finger

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
      const split = new Date(this.event.dates[0])
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

    // -----------------------------------
    //#region Scroll
    // -----------------------------------
    onCalendarScroll(e) {
      this.calendarMaxScroll = e.target.scrollWidth - e.target.offsetWidth
      this.calendarScrollLeft = e.target.scrollLeft
    },
    //#endregion
  },
  watch: {
    availability() {
      if (this.state === this.states.EDIT_AVAILABILITY) {
        this.unsavedChanges = true
      }
    },
    state(nextState, prevState) {
      // Reset scheduled event when exiting schedule event state
      if (prevState === this.states.SCHEDULE_EVENT) {
        this.curScheduledEvent = null
      } else if (prevState === this.states.EDIT_AVAILABILITY) {
        this.unsavedChanges = false
      }
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
    calendarEventsByDay() {
      if (
        this.state === this.states.EDIT_AVAILABILITY &&
        this.authUser &&
        !(this.authUser?._id in this.event.responses) && // User hasn't responded yet
        !this.loadingCalendarEvents &&
        !this.unsavedChanges
      ) {
        this.setAvailabilityAutomatically()
      }
    },
  },
  created() {
    this.resetCurUserAvailability()

    addEventListener("click", this.deselectRespondents)
  },
  mounted() {
    // Set initial state to best_times or heatmap depending on show best times toggle.
    this.state = this.showBestTimes ? "best_times" : "heatmap"

    // Set initial calendar max scroll
    this.calendarMaxScroll =
      this.$refs.calendar.scrollWidth - this.$refs.calendar.offsetWidth

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
    removeEventListener("click", this.deselectRespondents)
  },
  components: {
    UserAvatarContent,
    ZigZag,
    TimezoneSelector,
    ConfirmDetailsDialog,
    ToolRow,
    CalendarAccounts,
  },
}
</script>
