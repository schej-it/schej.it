<template>
  <div v-if="event">
    <div class="tw-bg-green tw-text-white tw-flex tw-flex-col tw-items-center tw-space-y-1 tw-py-2">
      <div 
        class="tw-font-bold tw-text-3xl"
      >{{ event.name }}</div>
      <div
        class="tw-font-light tw-text-lg"
      >{{ dateString }}</div>
    </div>
    <div v-if="isCalendarShown" class="tw-h-8 tw-sticky tw-top-0 tw-bg-light-blue tw-w-full tw-z-10 tw-flex tw-items-center tw-py-1 tw-px-2 tw-drop-shadow">
      <div class="tw-text-white tw-text-sm">
        <span v-if="isEditing">Editing...</span>
        <span v-else>Tap and hold to edit availability</span>
      </div>
      <v-spacer />
      <v-btn v-if="isEditing" @click="resetEditing" small text class="tw-text-white">Reset</v-btn>
      <v-btn v-if="areUnsavedChanges" @click="saveChanges" small class="tw-bg-blue" dark>Save</v-btn>
    </div>
    <div v-else class="tw-h-8" />
    <ScheduleOverlap
      ref="scheduleOverlap"
      :eventId="eventId" 
      v-bind="event"
      :calendarEvents="calendarEvents"
      :initialShowCalendarEvents="initialShowCalendarEvents"
      @refreshEvent="refreshEvent"
    />
  </div>
</template>

<script>
import { getDateRangeString, get, signInGoogle, dateCompare, dateToTimeInt, getDateDayOffset, clampDateToTimeInt, post } from '@/utils'
import { mapState } from 'vuex'

import ScheduleOverlap from '@/components/ScheduleOverlap'

export default {
  name: 'Event',

  props: {
    eventId: { type: String, required: true },
  },

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    calendarEvents: [],
    event: null,
    scheduleOverlapComponent: null
  }),

  computed: {
    ...mapState([ 'authUser', 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
    initialShowCalendarEvents() {
      return !(this.authUser._id in this.event.responses)
    },
    isCalendarShown() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.showCalendarEvents
    },
    isEditing() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.editing
    },
    areUnsavedChanges() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.unsavedChanges
    },
  },

  methods: {
    test() {
      console.log(document.querySelector('.v-main').scrollTop)
    },
    async refreshEvent() {
      // Get event details
      this.event = await get(`/events/${this.eventId}`)
      this.event.startDate = new Date(this.event.startDate)
      this.event.endDate = new Date(this.event.endDate)
    },
    resetEditing() {
      if (this.scheduleOverlapComponent) this.scheduleOverlapComponent.setAvailability()
    },
    saveChanges() {
      if (this.scheduleOverlapComponent) this.scheduleOverlapComponent.submitAvailability()
    },
  },

  async created() {
    // Get event details
    this.event = await get(`/events/${this.eventId}`)
    this.event.startDate = new Date(this.event.startDate)
    this.event.endDate = new Date(this.event.endDate)
    
    // Get user's calendar
    get(`/user/calendar?timeMin=${this.event.startDate.toISOString()}&timeMax=${getDateDayOffset(this.event.endDate, 1).toISOString()}`).then(data => {
      this.calendarEvents = data
        .map(event => ({ 
          summary: event.summary,
          startDate: clampDateToTimeInt(new Date(event.startDate), this.event.startTime, 'upper'),
          endDate: clampDateToTimeInt(new Date(event.endDate), this.event.endTime, 'lower'),
        }))
    }).catch(err => {
      console.error(err)
      if (err.code === 401 || err.code === 403) {
        signInGoogle({ type: 'join', eventId: this.eventId })
      }
    })
  },

  watch: {
    event() {
      if (this.event) {
        this.$nextTick(() => {
          this.scheduleOverlapComponent = this.$refs.scheduleOverlap
        })
      }
    },
  },
}
</script>