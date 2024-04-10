<template>
  <div v-if="event" class="tw-mt-8">
    <!-- Mark availability option dialog -->
    <MarkAvailabilityDialog
      v-model="choiceDialog"
      @setAvailabilityAutomatically="setAvailabilityAutomatically"
      @setAvailabilityManually="setAvailabilityManually"
    />

    <!-- Google sign in not supported dialog -->
    <SignInNotSupportedDialog v-model="webviewDialog" />

    <!-- Guest Dialog -->
    <GuestDialog
      v-model="guestDialog"
      @submit="saveChangesAsGuest"
      :respondents="Object.keys(event.responses)"
    />

    <!-- Edit event dialog -->
    <NewDialog
      v-model="editEventDialog"
      :type="isGroup ? 'group' : 'event'"
      :event="event"
      :contactsPayload="contactsPayload"
      edit
      no-tabs
    />

    <!-- Group invitation dialog -->
    <InvitationDialog
      v-if="isGroup"
      v-model="invitationDialog"
      :group="event"
      :calendarPermissionGranted="calendarPermissionGranted"
      @refreshEvent="refreshEvent"
      @setAvailabilityAutomatically="setAvailabilityAutomatically"
    ></InvitationDialog>

    <div class="tw-mx-auto tw-mt-4 tw-max-w-5xl">
      <div class="tw-mx-4">
        <!-- Title and copy link -->
        <div class="tw-flex tw-items-center tw-text-black">
          <div>
            <div class="tw-text-xl sm:tw-text-3xl">{{ event.name }}</div>
            <div class="tw-flex tw-items-baseline tw-gap-1">
              <div
                class="tw-text-sm tw-font-normal tw-text-very-dark-gray sm:tw-text-base"
              >
                {{ dateString }}
              </div>
              <v-btn
                v-if="canEdit"
                id="edit-event-btn"
                @click="editEvent"
                icon
                dense
                class="tw-min-w-0 tw-px-2 tw-text-sm tw-text-green sm:tw-text-base"
              >
                <v-icon v-if="!isPhone" small>mdi-pencil</v-icon>
                <span v-else class="tw-underline">Edit</span>
              </v-btn>
            </div>
          </div>
          <v-spacer />
          <div class="tw-flex tw-flex-row tw-items-center tw-gap-2.5">
            <div v-if="isGroup">
              <v-btn
                :icon="isPhone"
                :outlined="!isPhone"
                class="tw-text-green"
                @click="refreshCalendar"
                :loading="loading"
              >
                <v-icon class="tw-mr-1 tw-text-green" v-if="!isPhone"
                  >mdi-refresh</v-icon
                >
                <span v-if="!isPhone" class="tw-mr-2 tw-text-green"
                  >Refresh</span
                >
                <v-icon class="tw-text-green" v-else>mdi-refresh</v-icon>
              </v-btn>
            </div>
            <div v-else>
              <v-btn
                :icon="isPhone"
                :outlined="!isPhone"
                class="tw-text-green"
                @click="copyLink"
              >
                <span v-if="!isPhone" class="tw-mr-2 tw-text-green"
                  >Copy link</span
                >
                <v-icon class="tw-text-green" v-if="!isPhone"
                  >mdi-content-copy</v-icon
                >
                <v-icon class="tw-text-green" v-else>mdi-share</v-icon>
              </v-btn>
            </div>
            <div v-if="!isPhone" class="tw-flex tw-w-40">
              <template v-if="!isEditing">
                <v-btn
                  v-if="!isGroup && !authUser && selectedGuestRespondent"
                  min-width="10.25rem"
                  class="tw-bg-green tw-text-white tw-transition-opacity"
                  :style="{ opacity: availabilityBtnOpacity }"
                  @click="editGuestAvailability"
                >
                  {{ `Edit ${selectedGuestRespondent}'s availability` }}
                </v-btn>
                <v-btn
                  v-else
                  width="10.25rem"
                  class="tw-bg-green tw-text-white tw-transition-opacity"
                  :disabled="loading && !userHasResponded"
                  :style="{ opacity: availabilityBtnOpacity }"
                  @click="addAvailability"
                >
                  {{
                    isGroup
                      ? "Edit calendars"
                      : userHasResponded
                      ? "Edit availability"
                      : "Add availability"
                  }}
                </v-btn>
              </template>
              <template v-else>
                <v-btn
                  class="tw-mr-1 tw-w-20 tw-text-red"
                  @click="cancelEditing"
                  outlined
                >
                  Cancel
                </v-btn>
                <v-btn
                  class="tw-w-20 tw-bg-green tw-text-white"
                  @click="saveChanges"
                >
                  Save
                </v-btn></template
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Calendar -->

      <ScheduleOverlap
        ref="scheduleOverlap"
        :event="event"
        :loadingCalendarEvents="loading"
        :calendarEventsMap="calendarEventsMap"
        :calendarPermissionGranted="calendarPermissionGranted"
        :calendar-availabilities="calendarAvailabilities"
        :weekOffset.sync="weekOffset"
        :curGuestId="curGuestId"
        :initial-timezone="initialTimezone"
        @refreshEvent="refreshEvent"
        @highlightAvailabilityBtn="highlightAvailabilityBtn"
        @deleteAvailability="deleteAvailability"
        @setCurGuestId="(id) => (curGuestId = id)"
      />
    </div>

    <template v-if="showFeedbackBtn">
      <v-divider />

      <div class="tw-flex tw-justify-center" v-if="showFeedbackBtn">
        <v-btn
          class="tw-h-16"
          block
          id="feedback-btn"
          text
          href="https://forms.gle/9AgRy4PQfWfVuBnw8"
          target="_blank"
        >
          Give feedback
        </v-btn>
      </div>
    </template>

    <div class="tw-h-16"></div>

    <!-- Bottom bar for phones -->
    <div
      v-if="isPhone"
      class="tw-fixed tw-bottom-0 tw-z-20 tw-flex tw-h-16 tw-w-full tw-items-center tw-bg-green tw-px-4"
    >
      <template v-if="!isEditing && !isScheduling">
        <v-btn
          v-if="numResponses > 0"
          text
          class="tw-text-white"
          @click="scheduleEvent"
          >Schedule</v-btn
        >
        <v-spacer />
        <v-btn
          v-if="!isGroup && !authUser && selectedGuestRespondent"
          outlined
          class="tw-bg-white tw-text-green tw-transition-opacity"
          :style="{ opacity: availabilityBtnOpacity }"
          @click="editGuestAvailability"
        >
          {{ `Edit ${selectedGuestRespondent}'s availability` }}
        </v-btn>
        <v-btn
          v-else
          outlined
          class="tw-bg-white tw-text-green tw-transition-opacity"
          :disabled="loading && !userHasResponded"
          :style="{ opacity: availabilityBtnOpacity }"
          @click="addAvailability"
        >
          {{
            isGroup
              ? "Edit calendars"
              : userHasResponded
              ? "Edit availability"
              : "Add availability"
          }}
        </v-btn>
      </template>
      <template v-else-if="isEditing">
        <v-btn text class="tw-text-white" @click="cancelEditing">
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn class="tw-bg-white tw-text-green" @click="saveChanges">
          Save
        </v-btn>
      </template>
      <template v-else-if="isScheduling">
        <v-btn text class="tw-text-white" @click="cancelScheduleEvent">
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn class="tw-bg-white tw-text-green" @click="confirmScheduleEvent">
          Schedule
        </v-btn>
      </template>
    </div>
  </div>
</template>

<script>
import {
  get,
  signInGoogle,
  isPhone,
  processEvent,
  getCalendarEventsMap,
  getDateRangeStringForEvent,
} from "@/utils"
import { mapActions, mapState } from "vuex"

import NewDialog from "@/components/NewDialog.vue"
import ScheduleOverlap from "@/components/schedule_overlap/ScheduleOverlap.vue"
import GuestDialog from "@/components/GuestDialog.vue"
import { errors, authTypes, eventTypes } from "@/constants"
import isWebview from "is-ua-webview"
import SignInNotSupportedDialog from "@/components/SignInNotSupportedDialog.vue"
import MarkAvailabilityDialog from "@/components/MarkAvailabilityDialog.vue"
import InvitationDialog from "@/components/groups/InvitationDialog.vue"

export default {
  name: "Event",

  props: {
    eventId: { type: String, required: true },
    fromSignIn: { type: Boolean, default: false },
    initialTimezone: { type: Object, default: () => ({}) },
    contactsPayload: { type: Object, default: () => ({}) },
  },

  components: {
    GuestDialog,
    ScheduleOverlap,
    NewDialog,
    SignInNotSupportedDialog,
    MarkAvailabilityDialog,
    InvitationDialog,
  },

  data: () => ({
    choiceDialog: false,
    webviewDialog: false,
    guestDialog: false,
    editEventDialog: false,
    invitationDialog: false,
    loading: true,
    calendarEventsMap: {},
    event: null,
    scheduleOverlapComponent: null,
    scheduleOverlapComponentLoaded: false,

    curGuestId: "", // Id of the current guest being edited
    calendarPermissionGranted: true,

    weekOffset: 0,

    availabilityBtnOpacity: 1,

    // Availability Groups
    calendarAvailabilities: {}, // maps userId to their calendar events
  }),

  mounted() {
    // If coming from enabling contacts, show the dialog. Checks if contactsPayload is not an Observer.
    this.editEventDialog = Object.keys(this.contactsPayload).length > 0
  },

  computed: {
    ...mapState(["authUser", "events"]),
    dateString() {
      return getDateRangeStringForEvent(this.event)
    },
    isEditing() {
      return this.scheduleOverlapComponent?.editing
    },
    isScheduling() {
      return this.scheduleOverlapComponent?.scheduling
    },
    canEdit() {
      return (
        this.event.ownerId == 0 || this.authUser?._id === this.event.ownerId
      )
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    isSpecificDates() {
      return this.event?.type === eventTypes.SPECIFIC_DATES || !this.event?.type
    },
    isWeekly() {
      return this.event?.type === eventTypes.DOW
    },
    isGroup() {
      return this.event?.type === eventTypes.GROUP
    },
    areUnsavedChanges() {
      return this.scheduleOverlapComponent?.unsavedChanges
    },
    userHasResponded() {
      return this.authUser?._id in this.event.responses
    },
    selectedGuestRespondent() {
      return this.scheduleOverlapComponent?.selectedGuestRespondent
    },
    showFeedbackBtn() {
      return this.isPhone
    },
    numResponses() {
      return this.scheduleOverlapComponent?.respondents.length
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo"]),
    addAvailability() {
      /* Show choice dialog if not signed in, otherwise, immediately start editing availability */
      if (!this.scheduleOverlapComponent) return

      if (
        (this.authUser && this.calendarPermissionGranted) ||
        this.userHasResponded
      ) {
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
      this.curGuestId = ""
    },
    copyLink() {
      /* Copies event link to clipboard */
      navigator.clipboard.writeText(
        `${window.location.origin}/e/${this.event.shortId ?? this.event._id}`
      )
      this.showInfo("Link copied to clipboard!")
    },
    async deleteAvailability() {
      if (!this.scheduleOverlapComponent) return

      if (!this.authUser) {
        if (this.curGuestId) {
          await this.scheduleOverlapComponent.deleteAvailability(
            this.curGuestId
          )
          this.curGuestId = ""
        }
      } else {
        await this.scheduleOverlapComponent.deleteAvailability()
      }

      this.showInfo(this.isGroup ? "Left group!" : "Availability deleted!")
      this.scheduleOverlapComponent.stopEditing()
    },
    editEvent() {
      /* Show edit event dialog */
      this.editEventDialog = true
    },
    /** Refresh event details */
    async refreshEvent() {
      this.event = await get(`/events/${this.eventId}`)
      processEvent(this.event)
    },
    setAvailabilityAutomatically() {
      /* Prompts user to sign in when "set availability automatically" button clicked */
      if (isWebview(navigator.userAgent)) {
        // Show dialog prompting user to use a real browser
        this.webviewDialog = true
      } else {
        // Or sign in if user is already using a real browser
        if (this.authUser) {
          // Request permission if calendar permissions not yet granted
          signInGoogle({
            state: {
              type: authTypes.EVENT_ADD_AVAILABILITY,
              eventId: this.eventId,
            },
            selectAccount: false,
            requestCalendarPermission: true,
          })
        } else {
          // Ask the user to select the account they want to sign in with if not logged in yet
          signInGoogle({
            state: {
              type: authTypes.EVENT_ADD_AVAILABILITY,
              eventId: this.eventId,
            },
            selectAccount: true,
            requestCalendarPermission: true,
          })
        }
      }
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
      this.scheduleOverlapComponent.populateUserAvailability(
        this.selectedGuestRespondent
      )
      this.scheduleOverlapComponent.startEditing()
    },
    async saveChanges() {
      /* Shows guest dialog if not signed in, otherwise saves auth user's availability */
      if (!this.scheduleOverlapComponent) return

      if (!this.authUser) {
        if (this.curGuestId) {
          this.saveChangesAsGuest(this.curGuestId)
          this.curGuestId = ""
        } else {
          this.guestDialog = true
        }
        return
      }

      await this.scheduleOverlapComponent.submitAvailability()

      this.showInfo("Changes saved!")
      this.scheduleOverlapComponent.stopEditing()
    },
    async saveChangesAsGuest(name) {
      /* After guest dialog is submitted, submit availability with the given name */
      if (!this.scheduleOverlapComponent) return

      if (name.length > 0) {
        await this.scheduleOverlapComponent.submitAvailability(name)

        this.showInfo("Changes saved!")
        this.scheduleOverlapComponent.resetCurUserAvailability()
        this.scheduleOverlapComponent.stopEditing()
        this.guestDialog = false
      }
    },
    scheduleEvent() {
      this.scheduleOverlapComponent?.scheduleEvent()
    },
    cancelScheduleEvent() {
      this.scheduleOverlapComponent?.cancelScheduleEvent()
    },
    confirmScheduleEvent() {
      this.scheduleOverlapComponent?.confirmScheduleEvent()
    },

    highlightAvailabilityBtn() {
      if (!this.isPhone) {
        window.scrollTo({ top: 0, behavior: "instant" })
      }

      this.availabilityBtnOpacity = 0.1
      setTimeout(() => {
        this.availabilityBtnOpacity = 1

        setTimeout(() => {
          this.availabilityBtnOpacity = 0.1
          setTimeout(() => {
            this.availabilityBtnOpacity = 1
          }, 100)
        }, 100)
      }, 100)
    },

    /** Refresh calendar availabilities of everybody in the group */
    async fetchCalendarAvailabilities() {
      if (this.event.type !== eventTypes.GROUP) return

      // this.calendarAvailabilities = {}
      const curWeekOffset = this.weekOffset
      return getCalendarEventsMap(this.event, {
        weekOffset: curWeekOffset,
        eventId: this.event._id,
      })
        .then((calendarAvailabilities) => {
          // Don't update calendar availabilities if user
          // selected a different weekoffset by the time these calendar events load
          if (curWeekOffset !== this.weekOffset) return

          this.calendarAvailabilities = calendarAvailabilities
        })
        .catch((err) => {
          console.error(err)
        })
    },

    /** Fetch current user's calendar events */
    async fetchAuthUserCalendarEvents() {
      if (!this.authUser) {
        this.calendarPermissionGranted = false
        return
      }

      // this.calendarEventsMap = {}
      const curWeekOffset = this.weekOffset
      return getCalendarEventsMap(this.event, { weekOffset: curWeekOffset })
        .then((eventsMap) => {
          // Don't set calendar events / set availability if user has already
          // selected a different weekoffset by the time these calendar events load
          if (curWeekOffset !== this.weekOffset) return

          this.calendarEventsMap = eventsMap

          // Set user availability automatically if we're in editing mode and they haven't responded
          if (
            this.authUser &&
            this.isEditing &&
            !this.userHasResponded &&
            !this.areUnsavedChanges &&
            this.scheduleOverlapComponent
          ) {
            this.$nextTick(() => {
              this.scheduleOverlapComponent?.setAvailabilityAutomatically()
            })
          }

          // calendar permission granted is false when every calendar in the calendar map has an error, true otherwise
          this.calendarPermissionGranted = !Object.values(
            this.calendarEventsMap
          ).every((c) => Boolean(c.error))
        })
        .catch((err) => {
          console.error(err)
          this.calendarPermissionGranted = false
        })
    },

    /** Refreshes calendar avaliabilities and fetches current user calendar events */
    refreshCalendar() {
      const promises = []
      promises.push(this.fetchCalendarAvailabilities())
      promises.push(this.fetchAuthUserCalendarEvents())

      const curWeekOffset = this.weekOffset
      this.loading = true
      Promise.allSettled(promises).then(() => {
        // Only set loading to false if promises resolved at the same week offset they were fetched at
        // i.e. no new promises are currently being run
        if (curWeekOffset === this.weekOffset) {
          this.loading = false
        }
      })
    },

    onBeforeUnload(e) {
      if (this.areUnsavedChanges) {
        e.preventDefault()
        e.returnValue = ""
        return
      }

      delete e["returnValue"]
    },
  },

  async created() {
    window.addEventListener("beforeunload", this.onBeforeUnload)

    // Get event details
    try {
      await this.refreshEvent()
    } catch (err) {
      switch (err.error) {
        case errors.EventNotFound:
          this.showError("The specified event does not exist!")
          this.$router.replace({ name: "home" })
          return
      }
    }

    const promises = []
    promises.push(this.fetchCalendarAvailabilities())
    promises.push(this.fetchAuthUserCalendarEvents())

    this.loading = true
    Promise.allSettled(promises).then(() => {
      this.loading = false
    })
  },

  beforeDestroy() {
    window.removeEventListener("beforeunload", this.onBeforeUnload)
  },

  watch: {
    event() {
      if (this.event) {
        this.$nextTick(() => {
          this.scheduleOverlapComponent = this.$refs.scheduleOverlap
        })
      }
    },
    scheduleOverlapComponent() {
      if (!this.scheduleOverlapComponentLoaded) {
        this.scheduleOverlapComponentLoaded = true

        // Put into editing mode if just signed in
        if (this.fromSignIn && !this.isGroup) {
          this.scheduleOverlapComponent.startEditing()
        }

        if (this.isGroup && !this.userHasResponded) {
          this.invitationDialog = true
        }
      }
    },
    weekOffset() {
      this.refreshCalendar()
    },
  },
}
</script>
