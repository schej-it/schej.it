<template>
  <div class="tw-p-4 tw-select-none" style="-webkit-touch-callout: none">
    <!-- Day header -->
    <div class="tw-flex">
      <div class="tw-w-12" />
      <div v-for="(day, i) in days" :key="i" class="tw-flex-1">
        <div class="tw-text-center">
          <div class="tw-uppercase tw-font-light tw-text-xs">
            {{ day.dateString }}
          </div>
          <div class="tw-text-lg tw-capitalize">{{ day.dayText }}</div>
        </div>
      </div>
      <div v-if="!calendarOnly" class="sm:tw-w-32" />
    </div>

    <div class="sm:tw-flex">
      <div class="sm:tw-flex-1">
        <!-- Times and grid -->
        <div class="tw-flex">
          <div class="tw-w-12">
            <div
              v-for="(time, i) in times"
              :key="i"
              class="tw-h-5 tw-text-xs tw-pt-1 tw-pr-2 tw-text-right tw-uppercase tw-font-light"
            >
              {{ time.text }}
            </div>
          </div>
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
              >
                <!-- Timeslots -->
                <div v-for="(time, t) in times" :key="t" class="tw-w-full">
                  <div
                    class="timeslot tw-h-5 tw-border-light-gray tw-border-r"
                    :class="timeslotClass(day, time, d, t)"
                    v-on="timeslotVon(d, t)"
                  ></div>
                </div>

                <!-- Calendar events -->
                <template v-if="editing || showCalendarEvents">
                  <div
                    v-for="(event, e) in calendarEventsByDay[d]"
                    :key="`${d}-${e}`"
                    class="tw-absolute tw-w-full tw-p-px tw-select-none"
                    :style="event.style"
                  >
                    <div class="tw-hidden tw-bg-blue/25 tw-bg-dark-gray/25" />
                    <div
                      :class="`tw-bg-${noEventNames ? 'dark-gray' : 'blue'}/25`"
                      class="tw-border-light-gray /*tw-border-solid*/ tw-border tw-w-full tw-h-full tw-text-ellipsis tw-text-xs tw-rounded tw-p-px tw-overflow-hidden"
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
                </template>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="!calendarOnly" class="tw-p-4 sm:tw-py-0 sm:tw-pr-0 sm:tw-w-32">
        <div class="tw-font-medium tw-mb-2 tw-flex tw-items-center">
          <span class="tw-mr-1">Responders</span>
          <div class="tw-bg-black tw-text-white tw-font-bold tw-w-5 tw-h-5 tw-flex tw-justify-center tw-items-center tw-rounded-full tw-text-xs"
            >{{respondents.length}}</div
          >
        </div>
        <div
          class="tw-pl-4 tw-text-sm tw-grid tw-grid-cols-2 tw-gap-x-2 sm:tw-block"
        >
          <div
            v-for="user in respondents"
            :key="user._id"
            :class="
              !curTimeslotAvailability[user._id]
                ? 'tw-line-through tw-text-gray'
                : 'hover:tw-font-bold'
            "
            class="tw-py-1"
            @mouseover="curUser = user._id"
            @mouseleave="curUser = ''"
          >
            {{ user.firstName + " " + user.lastName + " " }}
            <v-tooltip right transition="slide-x-transition">
              <template v-slot:activator="{ on, attrs }">
                <v-icon
                  :class="!isGuest(user) ? 'tw-text-blue' : 'tw-text-dark-gray'"
                  class="-tw-mt-0.5"
                  small
                  v-bind="attrs"
                  v-on="on"
                  >{{
                    !isGuest(user) ? "mdi-check-decagram" : "mdi-account"
                  }}</v-icon
                >
              </template>
              <span>{{ !isGuest(user) ? "Verified account" : "Guest" }}</span>
            </v-tooltip>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  timeIntToTimeText,
  getDateDayOffset,
  dateCompare,
  compareDateDay,
  dateToTimeInt,
  getDateWithTimeInt,
  post,
  onLongPress,
  isBetween,
  clamp,
  isPhone,
} from "@/utils";
import { mapActions, mapState } from "vuex";

export default {
  name: "ScheduleOverlap",

  props: {
    eventId: { type: String, default: "" },
    startDate: { type: Date, required: true },
    endDate: { type: Date, required: true },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    responses: { type: Object, default: () => ({}) },

    loadingCalendarEvents: { type: Boolean, default: false },
    calendarEvents: { type: Array, required: true },
    initialShowCalendarEvents: { type: Boolean, default: false },

    noEventNames: { type: Boolean, default: false },
    calendarOnly: { type: Boolean, default: false },
  },

  data() {
    return {
      max: 0, // The max number of respondents for a given timeslot
      showCalendarEvents: this.initialShowCalendarEvents, // Whether we are showing the current user's calendar events
      availability: new Set(), // Current availability of the current user, an array of dates
      editing: false, // Whether editing the current user's availability
      unsavedChanges: false, // Whether there are unsaved availability changes

      curTimeslotAvailability: {}, // An object containing the people that are available and unavailable for the given timeslot, maps their user id to either true or false
      curTimeslot: { dayIndex: -1, timeIndex: -1 },
      curUser: "", // The id of the current user to show the availability of

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
    };
  },

  computed: {
    ...mapState(["authUser"]),
    availabilityArray() {
      /* Returns the availibility as an array */
      return [...this.availability].map((item) => new Date(item));
    },
    calendarEventsByDay() {
      /* Returns a 2d array of events based on the day they take place. Index 0 = first day */

      // TODO: calendar event spanning two days breaks this (how to fix: split calendar events into two or more events if span multiple days)
      const arr = [];
      for (let i = 0; i < this.days.length; ++i) {
        arr[i] = [];
      }
      for (const calendarEvent of this.calendarEvents) {
        // The number of hours since start time
        // TODO: this doesn't account for if the startTime and the startTime of the event are on different days
        const hoursOffset =
          (calendarEvent.startDate.getTime() -
            getDateWithTimeInt(
              calendarEvent.startDate,
              this.startTime
            ).getTime()) /
          (1000 * 60 * 60);
        // The length of the event in hours
        const hoursLength =
          (calendarEvent.endDate.getTime() -
            calendarEvent.startDate.getTime()) /
          (1000 * 60 * 60);

        for (const d in this.days) {
          const day = this.days[d];
          if (compareDateDay(day.dateObject, calendarEvent.startDate) == 0) {
            arr[d].push({
              ...calendarEvent,
              style: {
                top: `calc(${hoursOffset} * 2 * 1.25rem)`, // 1.25 rem = tw-h-5
                height: `calc(${hoursLength} * 2 * 1.25rem)`,
              },
            });
            break;
          }
        }
      }
      return arr;
    },
    days() {
      /* Return the days that are encompassed by startDate and endDate */
      const days = [];
      const daysOfWeek = ["sun", "mon", "tue", "wed", "thu", "fri", "sat"];
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
      ];
      let curDate = this.startDate;
      while (curDate.getTime() <= this.endDate.getTime()) {
        days.push({
          dayText: daysOfWeek[curDate.getDay()],
          dateString: months[curDate.getMonth()] + " " + curDate.getDate(),
          dateObject: curDate,
        });
        curDate = getDateDayOffset(curDate, 1);
      }

      return days;
    },
    isPhone() {
      return isPhone(this.$vuetify);
    },
    respondents() {
      return Object.values(this.parsedResponses).map((r) => r.user);
    },
    parsedResponses() {
      /* Parses responses so that if _id is null (i.e. guest user), then it is set to the guest user's name */
      const parsed = {};
      for (const k of Object.keys(this.responses)) {
        const newUser = {
          ...this.responses[k].user,
          _id: k,
        };
        parsed[k] = {
          ...this.responses[k],
          user: newUser,
        };
      }
      return parsed;
    },
    responsesFormatted() {
      /* Formats the responses in a map where date/time is mapped to the people that are available then */
      const formatted = new Map();
      for (const day of this.days) {
        for (const time of this.times) {
          const date = getDateWithTimeInt(day.dateObject, time.timeInt);
          formatted.set(date.getTime(), new Set());

          for (const response of Object.values(this.parsedResponses)) {
            const index = response.availability.findIndex(
              (d) => dateCompare(d, date) === 0
            );
            if (index !== -1) {
              // TODO: determine whether I should delete the index??
              formatted.get(date.getTime()).add(response.user._id);
            }
          }

          // Update max
          if (formatted.get(date.getTime()).size > this.max) {
            this.max = formatted.get(date.getTime()).size;
          }
        }
      }
      return formatted;
    },
    times() {
      /* Returns the times that are encompassed by startTime and endTime */
      const times = [];

      let t = this.startTime;
      while (t != this.endTime) {
        times.push({
          timeInt: t,
          text: timeIntToTimeText(t),
        });
        times.push({
          timeInt: t + 0.5,
        });
        t++;
        if (t > 23) t = 0;
      }

      return times;
    },
    userHasResponded() {
      return this.authUser && this.authUser._id in this.parsedResponses;
    },
  },

  methods: {
    ...mapActions(["showInfo"]),
    getRespondentsForDateTime(date, time) {
      /* Returns an array of respondents for the given date/time */
      const d = getDateWithTimeInt(date, time);
      return this.responsesFormatted.get(d.getTime());
    },
    resetCurUserAvailability() {
      /* resets cur user availability to the response stored on the server */
      this.availability = new Set();

      if (this.userHasResponded) {
        this.responses[this.authUser._id].availability.forEach((item) =>
          this.availability.add(new Date(item).getTime())
        );
        this.$nextTick(() => (this.unsavedChanges = false));
      }
    },
    async setAvailability() {
      /* Constructs the availability array using calendarEvents array */
      // This is not a computed property because we should be able to change it manually from what it automatically fills in
      this.availability = new Set();
      for (const d in this.days) {
        const day = this.days[d];
        for (const time of this.times) {
          // Check if there exists a calendar event that overlaps [time, time+0.5]
          const startDate = getDateWithTimeInt(day.dateObject, time.timeInt);
          const endDate = getDateWithTimeInt(
            day.dateObject,
            time.timeInt + 0.5
          );
          const index = this.calendarEventsByDay[d].findIndex((e) => {
            return (
              (dateCompare(e.startDate, startDate) < 0 &&
                dateCompare(e.endDate, startDate) > 0) ||
              (dateCompare(e.startDate, endDate) < 0 &&
                dateCompare(e.endDate, endDate) > 0) ||
              (dateCompare(e.startDate, startDate) == 0 &&
                dateCompare(e.endDate, endDate) == 0)
            );
          });
          if (index === -1) {
            this.availability.add(startDate.getTime());
          }
        }
      }
      this.showInfo(
        "Your availability has been set automatically using your Google Calendar!"
      );
    },
    setTimeslotSize() {
      /* Gets the dimensions of each timeslot and assigns it to the timeslot variable */
      ({ width: this.timeslot.width, height: this.timeslot.height } = document
        .querySelector(".timeslot")
        .getBoundingClientRect());
    },
    showAvailability(d, t) {
      this.curTimeslot = { dayIndex: d, timeIndex: t };
      const available = this.getRespondentsForDateTime(
        this.days[d].dateObject,
        this.times[t].timeInt
      );
      for (const respondent of this.respondents) {
        if (available.has(respondent._id)) {
          this.curTimeslotAvailability[respondent._id] = true;
        } else {
          this.curTimeslotAvailability[respondent._id] = false;
        }
      }
    },
    startEditing() {
      this.editing = true;
    },
    stopEditing() {
      this.editing = false;
    },
    async submitAvailability(name = "") {
      const payload = { availability: this.availabilityArray };
      if (this.authUser) {
        payload.guest = false;
      } else {
        payload.guest = true;
        payload.name = name;
      }
      await post(`/events/${this.eventId}/response`, payload);
      this.$emit("refreshEvent");
      this.unsavedChanges = false;
    },
    timeslotClass(day, time, d, t) {
      /* Returns a class string for the given timeslot div */
      let c = "";

      // Border style
      if (this.curTimeslot.dayIndex === d && this.curTimeslot.timeIndex === t) {
        c += "tw-border tw-border-dashed tw-border-black tw-z-10 ";
      } else {
        if (!("text" in time)) c += "tw-border-b ";
        if (d === 0) c += "tw-border-l tw-border-l-gray ";
        if (d === this.days.length - 1) c += "tw-border-r-gray ";
        if (t === 0) c += "tw-border-t tw-border-t-gray ";
        if (t === this.times.length - 1) c += "tw-border-b-gray ";
      }

      // Fill style
      if (this.editing) {
        // Show only current user availability

        const inDragRange = this.inDragRange(d, t);
        if (inDragRange) {
          // Set style if drag range goes over the current timeslot
          if (this.dragType === this.DRAG_TYPES.ADD) {
            c += "tw-bg-avail-green-300 ";
          } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
          }
        } else {
          // Otherwise just show the current availability
          const date = getDateWithTimeInt(day.dateObject, time.timeInt);
          if (this.availability.has(date.getTime())) {
            c += "tw-bg-avail-green-300 ";
          }
        }
      } else {
        if (this.curUser) {
          const respondents = this.getRespondentsForDateTime(
            day.dateObject,
            time.timeInt
          );
          if (respondents.has(this.curUser)) {
            c += "tw-bg-avail-green-300 ";
          }
        } else {
          // Show everyone's availability
          const numRespondents = this.getRespondentsForDateTime(
            day.dateObject,
            time.timeInt
          ).size;
          if (numRespondents > 0) {
            const frac = numRespondents / this.max;
            const colors = [
              //'tw-bg-avail-green-50',
              "tw-bg-avail-green-100",
              "tw-bg-avail-green-200",
              "tw-bg-avail-green-300",
              "tw-bg-avail-green-400",
              "tw-bg-avail-green-500",
              //'tw-bg-light-blue',
              //'tw-bg-avail-green-600',
            ];
            c += colors[parseInt(frac * colors.length - 1)] + " ";
          }
        }
      }

      return c;
    },
    timeslotVon(d, t) {
      if (!this.editing) {
        return {
          click: () => this.showAvailability(d, t),
          mouseover: () => this.showAvailability(d, t),
        };
      }
      return {};
    },
    resetCurTimeslot() {
      this.curTimeslotAvailability = {};
      for (const respondent of this.respondents) {
        this.curTimeslotAvailability[respondent._id] = true;
      }
      this.curTimeslot = { dayIndex: -1, timeIndex: -1 };
    },

    /* Drag Stuff */
    normalizeXY(e) {
      /* Normalize the touch event to be relative to element */
      let pageX, pageY;
      if ("touches" in e) {
        // is a touch event
        ({ pageX, pageY } = e.touches[0]);
      } else {
        // is a mouse event
        ({ pageX, pageY } = e);
      }
      const { left, top } = e.currentTarget.getBoundingClientRect();
      const x = pageX - left;
      const y = pageY - top - window.scrollY;
      return { x, y };
    },
    getDateFromXY(x, y) {
      /* Returns a date for the timeslot we are currently hovering over given the x and y position */
      const { width, height } = this.timeslot;
      let dayIndex = Math.floor(x / width);
      let timeIndex = Math.floor(y / height);
      dayIndex = clamp(dayIndex, 0, this.days.length - 1);
      timeIndex = clamp(timeIndex, 0, this.times.length - 1);

      return {
        dayIndex,
        timeIndex,
        date: getDateWithTimeInt(
          this.days[dayIndex].dateObject,
          this.times[timeIndex].timeInt
        ),
      };
    },
    endDrag() {
      if (!this.editing) return;

      if (!this.dragStart || !this.dragCur) return;

      // Update availability set based on drag region
      let dayInc =
        (this.dragCur.dayIndex - this.dragStart.dayIndex) /
        Math.abs(this.dragCur.dayIndex - this.dragStart.dayIndex);
      let timeInc =
        (this.dragCur.timeIndex - this.dragStart.timeIndex) /
        Math.abs(this.dragCur.timeIndex - this.dragStart.timeIndex);
      if (isNaN(dayInc)) dayInc = 1;
      if (isNaN(timeInc)) timeInc = 1;
      let d = this.dragStart.dayIndex;
      while (d != this.dragCur.dayIndex + dayInc) {
        let t = this.dragStart.timeIndex;
        while (t != this.dragCur.timeIndex + timeInc) {
          const date = getDateWithTimeInt(
            this.days[d].dateObject,
            this.times[t].timeInt
          );
          if (this.dragType === this.DRAG_TYPES.ADD) {
            this.availability.add(date.getTime());
          } else if (this.dragType === this.DRAG_TYPES.REMOVE) {
            this.availability.delete(date.getTime());
          }
          t += timeInc;
        }
        d += dayInc;
      }
      this.availability = new Set(this.availability);

      // Set dragging defaults
      this.dragging = false;
      this.dragStart = null;
      this.dragCur = null;
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
        );
      }
    },
    moveDrag(e) {
      if (!this.editing) return;

      e.preventDefault();
      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      );
      this.dragCur = { dayIndex, timeIndex };
    },
    startDrag(e) {
      if (!this.editing) return;

      this.dragging = true;
      const { dayIndex, timeIndex, date } = this.getDateFromXY(
        ...Object.values(this.normalizeXY(e))
      );
      this.dragStart = { dayIndex, timeIndex };
      this.dragCur = { dayIndex, timeIndex };

      // Set drag type
      if (this.availability.has(date.getTime())) {
        this.dragType = this.DRAG_TYPES.REMOVE;
      } else {
        this.dragType = this.DRAG_TYPES.ADD;
      }
    },
    isGuest(user) {
      return user._id == user.firstName;
    },
  },

  watch: {
    availability() {
      this.unsavedChanges = true;
    },
    calendarEvents: {
      handler() {
        //if (!this.userHasResponded && !this.calendarOnly) this.setAvailability()
      },
    },
    respondents: {
      immediate: true,
      handler() {
        this.curTimeslotAvailability = {};
        for (const respondent of this.respondents) {
          this.curTimeslotAvailability[respondent._id] = true;
        }
      },
    },
  },

  created() {
    this.resetCurUserAvailability();
  },

  mounted() {
    // Get timeslot size
    this.setTimeslotSize();
    window.addEventListener("resize", this.setTimeslotSize);

    if (!this.calendarOnly) {
      const timesEl = document.getElementById("times");

      if (isPhone(this.$vuetify)) {
        timesEl.addEventListener("touchstart", this.startDrag);
        timesEl.addEventListener("touchmove", this.moveDrag);

        timesEl.addEventListener("touchend", this.endDrag);
        timesEl.addEventListener("touchcancel", this.endDrag);
      } else {
        timesEl.addEventListener("mousedown", this.startDrag);
        timesEl.addEventListener("mousemove", this.moveDrag);
        timesEl.addEventListener("mouseup", this.endDrag);
      }
    }
  },
};
</script>
