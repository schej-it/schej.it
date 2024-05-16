<template>
  <div v-if="event" class="tw-mt-8 tw-h-full">
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

    <!-- Pages Not Visited dialog -->
    <v-dialog
      v-model="pagesNotVisitedDialog"
      max-width="400"
      content-class="tw-m-0"
    >
      <v-card>
        <v-card-title>Are you sure?</v-card-title>
        <v-card-text
          ><span class="tw-font-medium"
            >You're about to add your availability without filling out all pages
            of this Schej.</span
          >
          Click the left and right arrows at the top to switch between
          pages.</v-card-text
        >
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="pagesNotVisitedDialog = false">Cancel</v-btn>
          <v-btn
            text
            color="primary"
            @click="
              () => {
                saveChanges(true)
                this.pagesNotVisitedDialog = false
              }
            "
            >Add anyways</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <div class="tw-mx-auto tw-mt-4 tw-max-w-5xl">
      <div class="tw-mx-4">
        <!-- Title and copy link -->
        <div class="tw-flex tw-items-center tw-text-black">
          <div>
            <div
              class="sm:mb-2 tw-flex tw-flex-wrap tw-items-center tw-gap-x-4 tw-gap-y-2"
            >
              <div class="tw-text-xl sm:tw-text-3xl">{{ event.name }}</div>
              <template v-if="isGroup">
                <div class="">
                  <v-chip
                    :small="isPhone"
                    class="tw-cursor-pointer tw-select-none tw-rounded tw-bg-light-gray tw-px-2 tw-font-medium sm:tw-px-3"
                    @click="helpDialog = true"
                    >Availability group</v-chip
                  >
                </div>
                <HelpDialog v-model="helpDialog">
                  <template v-slot:header>Availability group</template>
                  <div class="mb-4">
                    Use availability groups to see group members' weekly
                    calendar availabilities. Your actual calendar events are NOT
                    visible to others.
                  </div>
                </HelpDialog>
              </template>
            </div>
            <div class="tw-flex tw-items-baseline tw-gap-1">
              <div
                class="tw-text-sm tw-font-normal tw-text-very-dark-gray sm:tw-text-base"
              >
                {{ dateString }}
              </div>
              <template v-if="canEdit">
                <v-btn
                  v-if="isPhone"
                  id="edit-event-btn"
                  @click="editEvent"
                  class="tw-px-2 tw-text-sm tw-text-green tw-underline"
                  text
                  >Edit event</v-btn
                >
                <v-btn
                  v-else
                  id="edit-event-btn"
                  @click="editEvent"
                  icon
                  dense
                  class="tw-min-w-0 tw-px-2 tw-text-sm tw-text-green sm:tw-text-base"
                >
                  <v-icon small>mdi-pencil</v-icon>
                </v-btn>
              </template>
            </div>
          </div>
          <v-spacer />
          <div class="tw-flex tw-flex-row tw-items-center tw-gap-2.5">
            <div v-if="isGroup">
              <v-btn
                v-if="weekOffset != 0"
                :icon="isPhone"
                text
                class="tw-mr-1 tw-text-very-dark-gray sm:tw-mr-2.5"
                @click="resetWeekOffset"
              >
                <v-icon class="sm:tw-mr-2">mdi-calendar-today</v-icon>
                <span v-if="!isPhone">Today</span>
              </v-btn>
              <v-btn
                :icon="isPhone"
                :outlined="!isPhone"
                class="tw-text-green"
                @click="refreshCalendar"
                :loading="loading"
              >
                <v-icon class="tw-mr-1" v-if="!isPhone">mdi-refresh</v-icon>
                <span v-if="!isPhone" class="tw-mr-2">Refresh</span>
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
                  class="tw-text-white tw-transition-opacity"
                  :class="'tw-bg-green'"
                  :disabled="loading && !userHasResponded"
                  :style="{ opacity: availabilityBtnOpacity }"
                  @click="() => addAvailability()"
                >
                  {{
                    userHasResponded || isGroup || showGuestEditAvailability
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
                  class="tw-w-20 tw-text-white"
                  :class="'tw-bg-green'"
                  @click="() => saveChanges()"
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
        :showGuestEditAvailability="showGuestEditAvailability"
        @addAvailability="addAvailability"
        @refreshEvent="refreshEvent"
        @highlightAvailabilityBtn="highlightAvailabilityBtn"
        @deleteAvailability="deleteAvailability"
        @setCurGuestId="(id) => (curGuestId = id)"
      />
    </div>

    <template v-if="showFeedbackBtn">
      <div class="tw-w-full tw-border-t tw-border-solid tw-border-gray"></div>

      <div class="tw-flex tw-flex-col tw-items-center" v-if="showFeedbackBtn">
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
        <div class="tw-w-full tw-border-t tw-border-solid tw-border-gray"></div>
        <v-btn
          class="tw-h-16"
          block
          text
          href="https://www.paypal.com/donate/?hosted_button_id=KWCH6LGJCP6E6"
          target="_blank"
        >
          Donate
        </v-btn>
      </div>
    </template>

    <div class="tw-h-8"></div>

    <!-- Bottom bar for phones -->
    <div
      v-if="isPhone"
      class="tw-fixed tw-bottom-0 tw-z-20 tw-flex tw-h-16 tw-w-full tw-items-center tw-px-4"
      :class="isScheduling ? 'tw-bg-blue' : 'tw-bg-green'"
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
          class="tw-bg-white tw-text-green tw-transition-opacity"
          :style="{ opacity: availabilityBtnOpacity }"
          @click="editGuestAvailability"
        >
          {{ `Edit ${selectedGuestRespondent}'s availability` }}
        </v-btn>
        <v-btn
          v-else
          class="tw-bg-white tw-text-green tw-transition-opacity"
          :disabled="loading && !userHasResponded"
          :style="{ opacity: availabilityBtnOpacity }"
          @click="() => addAvailability()"
        >
          {{
            userHasResponded || showGuestEditAvailability
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
        <v-btn class="tw-bg-white tw-text-green" @click="() => saveChanges()">
          Save
        </v-btn>
      </template>
      <template v-else-if="isScheduling">
        <v-btn text class="tw-text-white" @click="cancelScheduleEvent">
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn
          :disabled="!allowScheduleEvent"
          class="tw-bg-white tw-text-blue"
          @click="confirmScheduleEvent"
        >
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
import HelpDialog from "@/components/HelpDialog.vue"

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
    HelpDialog,
  },

  data: () => ({
    choiceDialog: false,
    webviewDialog: false,
    guestDialog: false,
    editEventDialog: false,
    invitationDialog: false,
    pagesNotVisitedDialog: false,
    helpDialog: false,
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
    allowScheduleEvent() {
      return this.scheduleOverlapComponent?.allowScheduleEvent
    },
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

    /** Whether to show "Edit availability" for a guest, and allow them to edit their availability with the main button */
    showGuestEditAvailability() {
      const c = this.scheduleOverlapComponent
      return (
        !this.authUser &&
        // this.event?.blindAvailabilityEnabled &&
        // !c?.isOwner &&
        c?.guestName?.length > 0 &&
        c?.guestName in c?.parsedResponses
      )
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo"]),
    /** Show choice dialog if not signed in, otherwise, immediately start editing availability */
    addAvailability(ignoreGuestEdit = false) {
      if (!this.scheduleOverlapComponent) return

      // Edit guest availability if guest edit availability enabled
      if (!ignoreGuestEdit && this.showGuestEditAvailability) {
        this.curGuestId = this.scheduleOverlapComponent?.guestName
        this.scheduleOverlapComponent.populateUserAvailability(this.curGuestId)
        this.scheduleOverlapComponent?.startEditing()
        return
      }

      // Start editing if calendar permission granted or user has responded, otherwise show choice dialog
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
              type: this.isGroup
                ? authTypes.GROUP_ADD_AVAILABILITY
                : authTypes.EVENT_ADD_AVAILABILITY,
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
    async saveChanges(ignorePagesNotVisited = false) {
      /* Shows guest dialog if not signed in, otherwise saves auth user's availability */
      if (!this.scheduleOverlapComponent) return

      // If user hasn't responded and they haven't gone to the next page, show pages not visited dialog
      if (
        !this.userHasResponded &&
        this.curGuestId.length === 0 &&
        !this.scheduleOverlapComponent.pageHasChanged &&
        !ignorePagesNotVisited &&
        this.scheduleOverlapComponent.numPages > 1
      ) {
        this.pagesNotVisitedDialog = true
        return
      }

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
          // Check if the primary calendar has an error
          // We don't care if other calendars have an error, because if they do we just dont show them
          if (eventsMap[this.authUser.email].error) {
            this.calendarPermissionGranted = false
            return
          }

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

    /** Resets week offset to 0 */
    resetWeekOffset() {
      this.weekOffset = 0
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

      // Redirect if we're at the wrong route
      if (this.event.type === eventTypes.GROUP) {
        if (this.$route.name === "event") {
          this.$router.replace({
            name: "group",
            params: {
              groupId: this.eventId,
              initialTimezone: this.initialTimezone,
              fromSignIn: this.fromSignIn,
              contactsPayload: this.contactsPayload,
            },
          })
        }
      } else {
        if (this.$route.name === "group") {
          this.$router.replace({
            name: "event",
            params: {
              eventId: this.eventId,
              initialTimezone: this.initialTimezone,
              fromSignIn: this.fromSignIn,
              contactsPayload: this.contactsPayload,
            },
          })
        }
      }
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
