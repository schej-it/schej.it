<template>
  <div v-if="event">
    <v-dialog
      v-model="choiceDialog"
      width="400"
      content-class="tw-m-0"
    >
      <v-card class="tw-text-center sm:tw-p-6 tw-p-4">
        <div class="tw-text-md tw-font-semibold tw-pb-4">How would you like to mark <br v-if="isPhone"> your availability?</div>
        <div class="">
          <v-btn 
            @click="setAvailabilityAutomatically"
            class="tw-bg-blue tw-mb-2" 
            dark 
            block
          >
            <div class="tw-text-sm -tw-mx-4">Automatically with Google Calendar</div>
          </v-btn>
          <v-btn @click="setAvailabilityManually" block>Manually</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <div class="tw-max-w-5xl tw-mx-auto tw-mt-4">

      <div class="tw-text-black tw-mx-8 tw-flex tw-items-center">
        <div>
          <div class="tw-text-3xl">{{ event.name }}</div>
          <div class="tw-font-normal">{{ dateString }}</div>
        </div>
        <v-spacer/>
        <div>
          <v-btn
            :icon="isPhone"
            :outlined="!isPhone"
            class="tw-text-green"
            @click="copyLink"
          >
            <span v-if="!isPhone" class="tw-text-green tw-mr-2">Copy link</span>
            <v-icon class="tw-text-green">mdi-content-copy</v-icon>
          </v-btn>
        </div>
      </div>

      <ScheduleOverlap
        ref="scheduleOverlap"
        :eventId="eventId" 
        v-bind="event"
        :loadingCalendarEvents="loading"
        :calendarEvents="calendarEvents"
        @refreshEvent="refreshEvent"
      />
    </div>
    <!-- Placeholder for bottom bar -->
    <div class="tw-h-16"></div>

    <div class="tw-flex tw-items-center tw-fixed tw-bottom-0 tw-bg-green tw-w-full tw-px-4 tw-h-16">
      <template v-if="!isEditing">
        <v-spacer />
        <v-btn
          outlined
          class="tw-text-green tw-bg-white"
          @click="choiceDialog = true"
        >
          Add availability
        </v-btn>
      </template>
      <template v-else>
        <v-btn
          text
          class="tw-text-white"
          @click="cancelEditing"
        >
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn
          class="tw-text-green tw-bg-white"
          @click="saveChanges"
        >
          Save
        </v-btn>
      </template>
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
    choiceDialog: false,

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
    isEditing() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.editing
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    areUnsavedChanges() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.unsavedChanges
    },
    userHasResponded() {
      return this.authUser && this.authUser._id in this.event.responses
    },
  },

  methods: {
    ...mapActions([ 'showError', 'showInfo' ]),
    cancelEditing() {
      if (!this.scheduleOverlapComponent) return

      // TODO: discard changes
      this.scheduleOverlapComponent.stopEditing()
    },
    copyLink() {
      navigator.clipboard.writeText(`${window.location.origin}/e/${this.eventId}`)
      this.showInfo('Link copied to clipboard!')
    },
    async refreshEvent() {
      // Get event details
      this.event = await get(`/events/${this.eventId}`)
      processEvent(this.event)
    },
    setAvailabilityAutomatically() {
      // if (this.scheduleOverlapComponent) this.scheduleOverlapComponent.setAvailability()
      this.choiceDialog = false
    },
    setAvailabilityManually() {
      if (!this.scheduleOverlapComponent) return

      this.scheduleOverlapComponent.startEditing()
      this.choiceDialog = false
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

    // Show dialog if user hasn't responded yet
    // this.choiceDialog = !this.userHasResponded
    
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