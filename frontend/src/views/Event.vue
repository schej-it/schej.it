<template>
  <div v-if="event">
    <div class="tw-relative tw-bg-green tw-text-white tw-flex tw-justify-center">
      <div class="tw-flex tw-flex-col tw-items-center tw-space-y-1 tw-py-2 tw-z-10">
        <div 
          class="tw-font-bold tw-text-3xl"
        >{{ event.name }}</div>
        <div
          class="tw-font-light tw-text-lg"
        >{{ dateString }}</div>
      </div>
    </div>
    <div class="tw-max-w-6xl tw-mx-auto">
      <div v-if="isCalendarShown" class="tw-relative tw-h-8 tw-sticky tw-top-0 tw-bg-light-blue tw-w-full tw-z-10 tw-flex tw-items-center tw-justify-center tw-py-1 tw-px-2 tw-drop-shadow">
        <div class="tw-text-white tw-text-sm tw-z-10">
          <span v-if="isPhone">
            <span v-if="isEditing">Editing...</span>
            <span v-else>Tap and hold calendar to enable editing</span>
          </span>
          <span v-else>Drag to edit availability</span>
        </div>
        <v-spacer />
        <v-btn v-if="isEditing || !isPhone" @click="resetEditing" small text class="tw-text-white">Clear</v-btn>
        <v-btn v-if="areUnsavedChanges" @click="saveChanges" small class="tw-bg-blue" dark>Save</v-btn>
      </div>
      <ScheduleOverlap
        ref="scheduleOverlap"
        :eventId="eventId" 
        v-bind="event"
        :loadingCalendarEvents="loading"
        :calendarEvents="calendarEvents"
        :initialShowCalendarEvents="initialShowCalendarEvents"
        @refreshEvent="refreshEvent"
      />
    </div>
  </div>
</template>

<script>
import { getDateRangeString, get, signInGoogle, dateCompare, dateToTimeInt, getDateDayOffset, clampDateToTimeInt, post, ERRORS, isPhone, processEvent, getDateWithTimeInt, getCalendarEvents } from '@/utils'
import { mapActions, mapState } from 'vuex'

import ScheduleOverlap from '@/components/ScheduleOverlap'
import { errors } from '@/constants'

export default {
  name: 'Event',

  props: {
    eventId: { type: String, required: true },
  },

  components: {
    ScheduleOverlap,
  },

  data: () => ({
    loading: true,
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
    isPhone() {
      return isPhone(this.$vuetify)
    },
    areUnsavedChanges() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.unsavedChanges
    },
  },

  methods: {
    ...mapActions([ 'showError', 'showInfo' ]),
    test() {
      console.log(document.querySelector('.v-main').scrollTop)
    },
    async refreshEvent() {
      // Get event details
      this.event = await get(`/events/${this.eventId}`)
      processEvent(this.event)
    },
    resetEditing() {
      if (this.scheduleOverlapComponent) this.scheduleOverlapComponent.setAvailability()
    },
    saveChanges() {
      if (this.scheduleOverlapComponent) this.scheduleOverlapComponent.submitAvailability()

      if (!this.isPhone) {
        this.showInfo('Changes saved!')
        this.scheduleOverlapComponent.showCalendarEvents = false
      }
    },
  },

  async created() {
    // Get event details
    try {
      this.event = await get(`/events/${this.eventId}`)
      processEvent(this.event)
    } catch (err) {
      switch (err.error) {
        case errors.EventNotFound:
          this.showError('The specified event does not exist!')
          this.$router.replace({ name: 'home' })
          return
      }
    }
    
    // Get user's calendar
    getCalendarEvents(this.event).then(events => {
      this.calendarEvents = events
      this.loading = false
    }).catch(err => {
      console.error(err)
      if (err.error.code === 401 || err.error.code === 403) {
        signInGoogle({ type: 'join', eventId: this.eventId }, true)
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