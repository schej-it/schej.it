<template>
  <div v-if="event" class="tw-mt-8">
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

    <GuestDialog v-model="guestDialog" @submit="saveChangesAsGuest" :respondents="Object.keys(event.responses)"/>

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
          v-if="!authUser && selectedGuestRespondent"
          outlined
          class="tw-text-green tw-bg-white"
          @click="editGuestAvailability"
        >
          {{ `Edit ${selectedGuestRespondent}'s availability` }}
        </v-btn>
        <v-btn
          v-else
          outlined
          class="tw-text-green tw-bg-white"
          :disabled="loading && !userHasResponded"
          @click="addAvailability"
        >
          {{ userHasResponded ? 'Edit availability' : 'Add availability' }}
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
import GuestDialog from '@/components/GuestDialog.vue'
import { errors } from '@/constants'

export default {
  name: 'Event',

  props: {
    eventId: { type: String, required: true },
  },

  components: {
    GuestDialog,
    ScheduleOverlap,
  },

  data: () => ({
    choiceDialog: false,
    guestDialog: false,

    loading: true,
    calendarEvents: [],
    event: null,
    scheduleOverlapComponent: null,

    curGuestId: '', // Id of the current guest being edited
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
    selectedGuestRespondent() {
      return this.scheduleOverlapComponent && this.scheduleOverlapComponent.selectedGuestRespondent
    },
  },

  methods: {
    ...mapActions([ 'showError', 'showInfo' ]),
    addAvailability() {
      /* Show choice dialog if not signed in, otherwise, immediately start editing availability */
      if (!this.scheduleOverlapComponent) return

      if (this.authUser) {
        this.scheduleOverlapComponent.startEditing()
        if (!this.userHasResponded) {
          this.scheduleOverlapComponent.setAvailabilityAutomatically()
        }
      } else {
        this.choiceDialog = true
      }
    },
    cancelEditing() {
      /* Cancels editing and resets availability to previous */
      if (!this.scheduleOverlapComponent) return

      this.scheduleOverlapComponent.resetCurUserAvailability()
      this.scheduleOverlapComponent.stopEditing()
      this.curGuestId = ''
    },
    copyLink() {
      /* Copies event link to clipboard */
      navigator.clipboard.writeText(`${window.location.origin}/e/${this.eventId}`)
      this.showInfo('Link copied to clipboard!')
    },
    async refreshEvent() {
      /* Refresh event details */
      this.event = await get(`/events/${this.eventId}`)
      processEvent(this.event)
    },
    setAvailabilityAutomatically() {
      /* Prompts user to sign in when "set availability automatically" button clicked */
      signInGoogle({ type: 'join', eventId: this.eventId }, true)
      this.choiceDialog = false
    },
    setAvailabilityManually() {
      /* Starts editing after "set availability manually" button clicked */
      if (!this.scheduleOverlapComponent) return

      this.scheduleOverlapComponent.startEditing()
      this.choiceDialog = false
    },
    editGuestAvailability() {
      /* Edits the selected guest's availability */
      if (!this.scheduleOverlapComponent) return

      this.curGuestId = this.selectedGuestRespondent
      this.scheduleOverlapComponent.setAvailability(this.selectedGuestRespondent)
      this.scheduleOverlapComponent.startEditing()
    },
    async saveChanges() {
      /* Shows guest dialog if not signed in, otherwise saves auth user's availability */
      if (!this.scheduleOverlapComponent) return
      
      if (!this.authUser) {
        if (this.curGuestId) {
          this.saveChangesAsGuest(this.curGuestId)
          this.curGuestId = ''
        } else {
          this.guestDialog = true
        }
        return 
      } 

      await this.scheduleOverlapComponent.submitAvailability()

      this.showInfo('Changes saved!')
      this.scheduleOverlapComponent.stopEditing()
    },
    async saveChangesAsGuest(name) {
      /* After guest dialog is submitted, submit availability with the given name */
      if (!this.scheduleOverlapComponent) return
      
      if (name.length > 0) {
        await this.scheduleOverlapComponent.submitAvailability(name)

        this.showInfo('Changes saved!')
        this.scheduleOverlapComponent.resetCurUserAvailability()
        this.scheduleOverlapComponent.stopEditing()
        this.guestDialog = false
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
      this.loading = false
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