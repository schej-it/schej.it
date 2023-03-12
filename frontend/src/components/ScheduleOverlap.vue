<template>
  <div class="tw-p-4 tw-select-none" style="-webkit-touch-callout: none">
    <div class="tw-flex tw-flex-wrap">
      <!-- Times -->
      <div class="tw-w-12 tw-mt-12">
        <div
          v-for="(time, i) in times"
          :key="i"
          class="tw-h-5 tw-text-xs tw-pt-1 tw-pr-2 tw-text-right tw-uppercase tw-font-light"
        >
          {{ time.text }}
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
                v-if="showCalendarEvents && loadingCalendarEvents"
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
                <div v-if="editing || showCalendarEvents">
                  <v-fade-transition
                    v-for="(event, e) in calendarEventsByDay[d]"
                    :key="`${d}-${e}`"
                    appear
                  >
                    <div
                      class="tw-absolute tw-w-full tw-p-px tw-select-none"
                      :style="event.style"
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
              </div>
            </div>
          </div>

          
        </div>

      </div>

      <div class="break" v-if="isPhone"></div>

      <!-- Respondents -->
      <div
        v-if="!calendarOnly"
        class="tw-p-4 sm:tw-pl-8 sm:tw-py-0 sm:tw-pr-0 sm:tw-pt-12 sm:tw-w-48"
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
            v-for="user in respondents"
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
              class="tw-mr-1 tw-break-all tw-transition-all"
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

    <div class="tw-flex">
      <div class="sm:tw-w-12"></div>

      <!-- Select timezone -->
      <div
        class="tw-flex-1 tw-flex tw-justify-center tw-items-center tw-mt-4 tw-text-sm"
        v-if="selectTimezone"
      >
        <div class="tw-mt-px tw-mr-2">
          Shown in
        </div>
        <v-select
          id="timezone-select"
          v-model="curTimezone"
          class="tw-text-sm -tw-mt-px tw-flex-none tw-min-w-min tw-max-w-xl"
          :items="Object.keys(timezoneMap)"
          dense
          color="#219653"
          item-color="green"
          hide-details
        ></v-select>
      </div>

      <div class="sm:tw-w-48"></div>
    </div>

  </div>
</template>

<style scoped>
.animate-bg-color {
  transition: background-color 0.25s ease-in-out;
}

.break {
  flex-basis: 100%;
  height: 0;
}

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
  getDateDayOffset,
  dateCompare,
  compareDateDay,
  getDateWithTimeNum,
  post,
  isBetween,
  clamp,
  isPhone,
  utcTimeToLocalTime,
} from "@/utils"
import { mapActions, mapState } from "vuex"
import UserAvatarContent from "./UserAvatarContent.vue"
import ZigZag from "./ZigZag.vue"
import timezoneData from "@/data/timezones.json"

export default {
  name: "ScheduleOverlap",
  props: {
    eventId: { type: String, default: "" },
    startDate: { type: Date, required: false },
    endDate: { type: Date, required: false },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    dates: { type: Array, required: false },
    responses: { type: Object, default: () => ({}) },
    loadingCalendarEvents: { type: Boolean, default: false },
    calendarEvents: { type: Array, required: true },
    initialShowCalendarEvents: { type: Boolean, default: false },
    noEventNames: { type: Boolean, default: false },
    calendarOnly: { type: Boolean, default: false },
    selectTimezone: { type: Boolean, default: false },
    interactable: { type: Boolean, default: true },
  },
  data() {
    return {
      max: 0,
      showCalendarEvents: this.initialShowCalendarEvents,
      availability: new Set(),
      availabilityAnimTimeouts: [], // Timeouts for availability animation
      availabilityAnimEnabled: false,
      maxAnimTime: 1200,
      editing: false,
      unsavedChanges: false,
      curTimeslotAvailability: {},
      curTimeslot: { dayIndex: -1, timeIndex: -1 },
      curRespondent: "", // Id of the active respondent (set on hover)
      curRespondentSelected: false,
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

      /* Variables for timezone */
      curTimezone: this.getLocalTimezone(),
    }
  },
  computed: {
    ...mapState(["authUser"]),
    availabilityArray() {
      /* Returns the availibility as an array */
      return [...this.availability].map((item) => new Date(item))
    },
    calendarEventsByDay() {
      /* Returns a 2d array of events based on the day they take place. Index 0 = first day */
      // TODO: calendar event spanning two days breaks this (how to fix: split calendar events into two or more events if span multiple days)
      const arr = []
      for (let i = 0; i < this.days.length; ++i) {
        arr[i] = []
      }

      let startTime
      if (this.startDate) {
        // Legacy date representation
        startTime = this.startTime
      } else {
        // New date representation
        // Do not need to specify this.timezoneOffset because we want to use the local timezoneoffset for event display
        startTime = utcTimeToLocalTime(this.startTime)
      }

      for (const calendarEvent of this.calendarEvents) {
        // calendarEventDayStart is a date representation of the event start time for the day the calendar event takes place
        const calendarEventDayStart = getDateWithTimeNum(calendarEvent.startDate, startTime)
        if (calendarEventDayStart.getTime() > calendarEvent.startDate.getTime()) {
          // Go back a day if calendarEventDayStart is past the calendarEvent start time
          calendarEventDayStart.setDate(calendarEventDayStart.getDate() - 1);
        }

        // The number of hours since start time
        const hoursOffset =
          (calendarEvent.startDate.getTime() - calendarEventDayStart.getTime()) /
          (1000 * 60 * 60)

        // The length of the event in hours
        const hoursLength =
          (calendarEvent.endDate.getTime() -
            calendarEvent.startDate.getTime()) /
          (1000 * 60 * 60)

        // Don't display event if the event is 0 hours long
        if (hoursLength == 0) continue

        for (const d in this.days) {
          const day = this.days[d]
          if (compareDateDay(day.dateObject, calendarEventDayStart) == 0) {
            arr[d].push({
              ...calendarEvent,
              style: {
                top: `calc(${hoursOffset} * 2 * 1.25rem)`,
                height: `calc(${hoursLength} * 2 * 1.25rem)`,
              },
            })
            break
          }
        }
      }
      return arr
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

      if (this.startDate) {
        // Legacy date representation method
        let curDate = this.startDate
        while (curDate.getTime() <= this.endDate.getTime()) {
          days.push({
            dayText: daysOfWeek[curDate.getDay()],
            dateString: `${months[curDate.getMonth()]} ${curDate.getDate()}`,
            dateObject: curDate,
          })
          curDate = getDateDayOffset(curDate, 1)
        }
      } else {
        // New date representation method
        for (const date of this.dates) {
          const hours = Math.floor(this.startTime)
          const minutes = Math.floor((this.startTime - hours) * 60)
          const paddedHours = String(hours).padStart(2, "0")
          const paddedMinutes = String(minutes).padStart(2, "0")

          // dateObject stores the original date object without timezone manipulation
          const dateObject = new Date(`${date}T${paddedHours}:${paddedMinutes}:00Z`)

          // Offset curDate by the correct timezone offset
          const curDate = new Date(dateObject)
          curDate.setHours(curDate.getHours() - this.timezoneOffset / 60)

          days.push({
            dayText: daysOfWeek[curDate.getUTCDay()],
            dateString: `${
              months[curDate.getUTCMonth()]
            } ${curDate.getUTCDate()}`,
            dateObject: dateObject,
          })
        }
      }
      return days
    },
    isPhone() {
      return isPhone(this.$vuetify)
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
          const date = this.getDateTime(day.dateObject, time.timeNum)
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
      let startTime
      let endTime

      if (this.startDate) {
        // Legacy date representation method
        startTime = this.startTime
        endTime = this.endTime
      } else {
        // New date representation method
        startTime = utcTimeToLocalTime(this.startTime, this.timezoneOffset)
        endTime = utcTimeToLocalTime(this.endTime, this.timezoneOffset)
      }

      let t = startTime
      while (t != endTime) {
        // Convert timeNum back to local time if using new date representation
        let timeNum = this.startDate
          ? t
          : utcTimeToLocalTime(utcTimeToLocalTime(t, -this.timezoneOffset))

        times.push({
          timeNum: timeNum,
          text: timeNumToTimeText(t),
        })
        times.push({
          timeNum: timeNum + 0.5,
        })
        t++
        t %= 24
      }
      return times
    },
    timezoneOffset() {
      return this.timezoneMap[this.curTimezone] * -1 // Multiplying by -1 because offset is flipped
    },
    timezoneMap() {
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
      return this.curTimeslot.dayIndex !== -1 && this.curTimeslot.timeIndex !== -1
    },
  },
  methods: {
    ...mapActions(["showInfo"]),

    /*
      Dates
    */
    getDateTime(calendarDate, timeNum) {
      /* 
        Returns a date object given the date and time, where calendarDate represents the start date of each individual day on the calendar
        and timeNum represents the time row for that day
      */
      let startTime

      if (this.startDate) {
        // Legacy date representation method
        startTime = this.startTime
      } else {
        // New date representation method
        startTime = utcTimeToLocalTime(this.startTime)
      }

      const date = getDateWithTimeNum(calendarDate, timeNum)
      if (timeNum < startTime) {
        // Go to the next day if timeNum is less than start time
        date.setDate(date.getDate() + 1)
      }

      return date
    },

    /*
      Respondent
    */
    mouseOverRespondent(e, id) {
      if (!this.curRespondentSelected) this.curRespondent = id
    },
    mouseLeaveRespondent(e) {
      if (!this.curRespondentSelected) this.curRespondent = ""
    },
    clickRespondent(e, id) {
      this.curRespondentSelected = true
      this.curRespondent = id
      e.stopPropagation()
    },
    deselectRespondent(e) {
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

    /*
      Aggregate user availability
    */
    getRespondentsForDateTime(date, time) {
      /* Returns an array of respondents for the given date/time */
      const d = this.getDateTime(date, time)
      return this.responsesFormatted.get(d.getTime())
    },
    showAvailability(d, t) {
      if (this.editing && this.isPhone) {
        // Don't show currently selected timeslot when on phone and editing
        return
      }

      this.curTimeslot = { dayIndex: d, timeIndex: t }

      if (this.editing || this.curRespondent) {
        // Don't show availability when editing or when respondent is selected
        return
      }

      const available = this.getRespondentsForDateTime(
        this.days[d].dateObject,
        this.times[t].timeNum
      )
      for (const respondent of this.respondents) {
        if (available.has(respondent._id)) {
          this.curTimeslotAvailability[respondent._id] = true
        } else {
          this.curTimeslotAvailability[respondent._id] = false
        }
      }
    },

    /*
      Current user availability
    */
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
          const startDate = this.getDateTime(day.dateObject, time.timeNum)
          const endDate = this.getDateTime(day.dateObject, time.timeNum + 0.5)
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
              this.showInfo(
                "Your availability has been set automatically using your Google Calendar!"
              )
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

    /*
      Timeslot
    */
    setTimeslotSize() {
      /* Gets the dimensions of each timeslot and assigns it to the timeslot variable */
      ({ width: this.timeslot.width, height: this.timeslot.height } = document
        .querySelector(".timeslot")
        .getBoundingClientRect())
    },
    timeslotClassStyle(day, time, d, t) {
      /* Returns a class string for the given timeslot div */
      let c = ""
      const s = {}
      // Animation
      if (this.availabilityAnimEnabled) {
        c += "animate-bg-color "
      }

      // Border style
      if (this.curTimeslot.dayIndex === d && this.curTimeslot.timeIndex === t) {
        c += "tw-border tw-border-dashed tw-border-black tw-z-10 "
      } else {
        if (!("text" in time)) c += "tw-border-b "
        if (d === 0) c += "tw-border-l tw-border-l-gray "
        if (d === this.days.length - 1) c += "tw-border-r-gray "
        if (t === 0) c += "tw-border-t tw-border-t-gray "
        if (t === this.times.length - 1) c += "tw-border-b-gray "
      }
      // Fill style
      if (this.editing) {
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
          const date = this.getDateTime(day.dateObject, time.timeNum)
          if (this.availability.has(date.getTime())) {
            c += "tw-bg-avail-green-300 "
          }
        }
      } else {
        if (this.curRespondent) {
          // Show the currently selected respondent's availability
          const respondent = this.curRespondent
          const respondents = this.getRespondentsForDateTime(
            day.dateObject,
            time.timeNum
          )
          if (respondents.has(respondent)) {
            c += "tw-bg-avail-green-300 "
          }
        } else {
          // Show everyone's availability
          const numRespondents = this.getRespondentsForDateTime(
            day.dateObject,
            time.timeNum
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

    /* 
      Editing
    */
    startEditing() {
      this.editing = true
    },
    stopEditing() {
      this.editing = false
      this.stopAvailabilityAnim()
    },

    /* 
      Drag Stuff 
    */
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
        date: this.getDateTime(
          this.days[dayIndex].dateObject,
          this.times[timeIndex].timeNum
        ),
      }
    },
    endDrag() {
      if (!this.editing) return
      if (!this.dragStart || !this.dragCur) return
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
          const date = this.getDateTime(
            this.days[d].dateObject,
            this.times[t].timeNum
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
      if (!this.editing) return
      e.preventDefault()
      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      )
      this.dragCur = { dayIndex, timeIndex }
    },
    startDrag(e) {
      if (!this.editing) return
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
    
    /*
      Timezone
    */
    getLocalTimezone() {
      const split = new Date()
        .toLocaleTimeString("en-us", { timeZoneName: "short" })
        .split(" ")
      const localTimezone = split[split.length - 1]
      
      return localTimezone
    }
  },
  watch: {
    availability() {
      this.unsavedChanges = true
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

    addEventListener("click", this.deselectRespondent)
  },
  mounted() {
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
  components: { UserAvatarContent, ZigZag },
}
</script>
