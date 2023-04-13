<template>
  <div v-if="event" class="tw-mt-8">
    <!-- Mark availability option dialog -->
    <v-dialog v-model="choiceDialog" width="400" content-class="tw-m-0">
      <v-card class="tw-text-center sm:tw-p-6 tw-p-4">
        <div class="tw-text-md tw-pb-4">
          How would you like to mark <br v-if="isPhone" />
          your availability?
        </div>
        <div class="">
          <v-btn
            @click="setAvailabilityAutomatically"
            class="tw-bg-blue tw-mb-2"
            dark
            block
          >
            <div class="tw-text-sm -tw-mx-4">
              Automatically with Google Calendar
            </div>
          </v-btn>
          <v-btn @click="setAvailabilityManually" block>Manually</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Google sign in not supported dialog -->
    <SignInNotSupportedDialog v-model="webviewDialog" />

    <!-- Guest Dialog -->
    <GuestDialog
      v-model="guestDialog"
      @submit="saveChangesAsGuest"
      :respondents="Object.keys(event.responses)"
    />

    <!-- Edit event dialog -->
    <NewEventDialog v-model="editEventDialog" :event="event" edit-event />

    <div class="tw-max-w-5xl tw-mx-auto tw-mt-4">
      <div class="tw-mx-8">
        <!-- Title and copy link -->

        <div class="tw-text-black tw-flex tw-items-center">
          <div>
            <div class="tw-text-xl sm:tw-text-3xl">{{ event.name }}</div>
            <div class="tw-flex tw-items-baseline tw-gap-1">
              <div class="tw-text-sm sm:tw-text-base tw-font-normal">
                {{ dateString }}
              </div>
              <v-btn
                v-if="isOwner"
                id="edit-event-btn"
                @click="editEvent"
                icon
                dense
                class="tw-text-green tw-min-w-0 tw-px-2 tw-text-sm sm:tw-text-base"
                ><v-icon v-if="!isPhone" small>mdi-pencil</v-icon
                ><span v-else class="tw-underline">Edit</span></v-btn
              >
            </div>
          </div>
          <v-spacer />
          <div class="tw-flex tw-flex-row tw-items-center tw-gap-2.5">
            <div>
              <v-btn
                :icon="isPhone"
                :outlined="!isPhone"
                class="tw-text-green"
                @click="copyLink"
              >
                <span v-if="!isPhone" class="tw-text-green tw-mr-2"
                  >Copy link</span
                >
                <v-icon class="tw-text-green" v-if="!isPhone"
                  >mdi-content-copy</v-icon
                >
                <v-icon class="tw-text-green" v-else>mdi-share</v-icon>
              </v-btn>
            </div>
            <div v-if="!isPhone">
              <template v-if="!isEditing">
                <v-btn
                  v-if="!authUser && selectedGuestRespondent"
                  min-width="10.25rem"
                  class="tw-text-white tw-bg-green"
                  @click="editGuestAvailability"
                >
                  {{ `Edit ${selectedGuestRespondent}'s availability` }}
                </v-btn>
                <v-btn
                  v-else
                  width="10.25rem"
                  class="tw-text-white tw-bg-green"
                  :disabled="loading && !userHasResponded"
                  @click="addAvailability"
                >
                  {{
                    userHasResponded ? "Edit availability" : "Add availability"
                  }}
                </v-btn>
              </template>
              <template v-else>
                <v-btn
                  class="tw-text-red tw-mr-1 tw-w-20"
                  @click="cancelEditing"
                  outlined
                >
                  Cancel
                </v-btn>
                <v-btn
                  class="tw-text-white tw-bg-green tw-w-20"
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
        :eventId="eventId"
        v-bind="event"
        :loadingCalendarEvents="loading"
        :calendarEventsByDay="calendarEventsByDay"
        @refreshEvent="refreshEvent"
        :selectTimezone="true"
      />
    </div>
    <div class="tw-h-16"></div>

    <!-- Bottom bar for phones -->
    <div
      v-if="isPhone"
      class="tw-flex tw-items-center tw-fixed tw-bottom-0 tw-bg-green tw-w-full tw-px-4 tw-h-16"
    >
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
          {{ userHasResponded ? "Edit availability" : "Add availability" }}
        </v-btn>
      </template>
      <template v-else>
        <v-btn text class="tw-text-white" @click="cancelEditing">
          Cancel
        </v-btn>
        <v-spacer />
        <v-btn class="tw-text-green tw-bg-white" @click="saveChanges">
          Save
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
  getCalendarEventsByDay,
  getDateRangeStringForEvent,
} from "@/utils";
import { mapActions, mapState } from "vuex";

import NewEventDialog from "@/components/NewEventDialog.vue";
import ScheduleOverlap from "@/components/ScheduleOverlap.vue";
import GuestDialog from "@/components/GuestDialog.vue";
import { errors } from "@/constants";
import isWebview from "is-ua-webview";
import SignInNotSupportedDialog from "@/components/SignInNotSupportedDialog.vue";

export default {
  name: "Event",

  props: {
    eventId: { type: String, required: true },
    fromSignIn: { type: Boolean, default: false },
  },

  components: {
    GuestDialog,
    ScheduleOverlap,
    NewEventDialog,
    SignInNotSupportedDialog,
  },

  data: () => ({
    choiceDialog: false,
    webviewDialog: false,
    guestDialog: false,
    editEventDialog: false,
    loading: true,
    calendarEventsByDay: [],
    event: null,
    scheduleOverlapComponent: null,
    scheduleOverlapComponentLoaded: false,

    curGuestId: "", // Id of the current guest being edited
    calendarPermissionGranted: false,
  }),

  computed: {
    ...mapState(["authUser", "events"]),
    dateString() {
      return getDateRangeStringForEvent(this.event);
    },
    isEditing() {
      return this.scheduleOverlapComponent?.editing;
    },
    isOwner() {
      return this.authUser?._id === this.event.ownerId;
    },
    isPhone() {
      return isPhone(this.$vuetify);
    },
    areUnsavedChanges() {
      return this.scheduleOverlapComponent?.unsavedChanges;
    },
    userHasResponded() {
      return this.authUser?._id in this.event.responses;
    },
    selectedGuestRespondent() {
      return this.scheduleOverlapComponent?.selectedGuestRespondent;
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo"]),
    addAvailability() {
      /* Show choice dialog if not signed in, otherwise, immediately start editing availability */
      if (!this.scheduleOverlapComponent) return;

      if (
        (this.authUser && this.calendarPermissionGranted) ||
        this.userHasResponded
      ) {
        this.scheduleOverlapComponent.startEditing();
        if (!this.userHasResponded) {
          this.scheduleOverlapComponent.setAvailabilityAutomatically();
        }
      } else {
        this.choiceDialog = true;
      }
    },
    cancelEditing() {
      /* Cancels editing and resets availability to previous */
      if (!this.scheduleOverlapComponent) return;

      this.scheduleOverlapComponent.resetCurUserAvailability();
      this.scheduleOverlapComponent.stopEditing();
      this.curGuestId = "";
    },
    copyLink() {
      /* Copies event link to clipboard */
      navigator.clipboard.writeText(
        `${window.location.origin}/e/${this.eventId}`
      );
      this.showInfo("Link copied to clipboard!");
    },
    editEvent() {
      /* Show edit event dialog */
      this.editEventDialog = true;
    },
    async refreshEvent() {
      /* Refresh event details */
      this.event = await get(`/events/${this.eventId}`);
      processEvent(this.event);
    },
    setAvailabilityAutomatically() {
      /* Prompts user to sign in when "set availability automatically" button clicked */
      if (isWebview(navigator.userAgent)) {
        // Show dialog prompting user to user a real browser
        this.webviewDialog = true;
      } else {
        // Or sign in if user is already using a real browser
        if (this.authUser) {
          // Request permission if calendar permissions not yet granted
          signInGoogle(
            { type: "event-add-availability", eventId: this.eventId },
            false,
            true
          );
        } else {
          // Ask the user to select the account they want to sign in with if not logged in yet
          signInGoogle(
            { type: "event-add-availability", eventId: this.eventId },
            true,
            true
          );
        }
      }
      this.choiceDialog = false;
    },
    setAvailabilityManually() {
      /* Starts editing after "set availability manually" button clicked */
      if (!this.scheduleOverlapComponent) return;

      this.scheduleOverlapComponent.startEditing();
      this.choiceDialog = false;
    },
    editGuestAvailability() {
      /* Edits the selected guest's availability */
      if (!this.scheduleOverlapComponent) return;

      this.curGuestId = this.selectedGuestRespondent;
      this.scheduleOverlapComponent.populateUserAvailability(
        this.selectedGuestRespondent
      );
      this.scheduleOverlapComponent.startEditing();
    },
    async saveChanges() {
      /* Shows guest dialog if not signed in, otherwise saves auth user's availability */
      if (!this.scheduleOverlapComponent) return;

      if (!this.authUser) {
        if (this.curGuestId) {
          this.saveChangesAsGuest(this.curGuestId);
          this.curGuestId = "";
        } else {
          this.guestDialog = true;
        }
        return;
      }

      await this.scheduleOverlapComponent.submitAvailability();

      this.showInfo("Changes saved!");
      this.scheduleOverlapComponent.stopEditing();
    },
    async saveChangesAsGuest(name) {
      /* After guest dialog is submitted, submit availability with the given name */
      if (!this.scheduleOverlapComponent) return;

      if (name.length > 0) {
        await this.scheduleOverlapComponent.submitAvailability(name);

        this.showInfo("Changes saved!");
        this.scheduleOverlapComponent.resetCurUserAvailability();
        this.scheduleOverlapComponent.stopEditing();
        this.guestDialog = false;
      }
    },

    onBeforeUnload(e) {
      if (this.isEditing) {
        e.preventDefault();
        e.returnValue = "";
        return;
      }

      delete e["returnValue"];
    },
  },

  async created() {
    window.addEventListener("beforeunload", this.onBeforeUnload);

    // Get event details
    try {
      await this.refreshEvent();
    } catch (err) {
      switch (err.error) {
        case errors.EventNotFound:
          this.showError("The specified event does not exist!");
          this.$router.replace({ name: "home" });
          return;
      }
    }

    // Get user's calendar
    getCalendarEventsByDay(this.event)
      .then((events) => {
        this.calendarEventsByDay = events;
        this.loading = false;

        // Set user availability automatically if we're in editing mode and they haven't responded
        if (
          this.authUser &&
          this.isEditing &&
          !this.userHasResponded &&
          this.scheduleOverlapComponent
        ) {
          this.$nextTick(() => {
            this.scheduleOverlapComponent.setAvailabilityAutomatically();
          });
        }

        this.calendarPermissionGranted = true;
      })
      .catch((err) => {
        this.loading = false;
        console.error(err);
        if (err.error.code === 401 || err.error.code === 403) {
          this.calendarPermissionGranted = false;
        }
      });
  },

  watch: {
    event() {
      if (this.event) {
        this.$nextTick(() => {
          this.scheduleOverlapComponent = this.$refs.scheduleOverlap;
        });
      }
    },
    scheduleOverlapComponent() {
      if (!this.scheduleOverlapComponentLoaded) {
        this.scheduleOverlapComponentLoaded;

        // Put into editing mode if just signed in
        if (this.fromSignIn) {
          this.scheduleOverlapComponent.startEditing();
        }
      }
    },
  },
};
</script>
