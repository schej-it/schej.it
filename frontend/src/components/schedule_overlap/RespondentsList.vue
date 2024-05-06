<template>
  <div>
    <div class="tw-flex tw-items-center tw-font-medium">
      <div class="tw-mr-1 tw-text-lg">
        {{ !isGroup ? "Responses" : "Members" }}
      </div>
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
      class="tw-mt-2 tw-grid tw-grid-cols-2 tw-gap-x-2 tw-overflow-hidden tw-text-sm sm:tw-block sm:tw-overflow-visible"
      :style="
        maxHeight
          ? `max-height: ${maxHeight}px !important; overflow-y: auto !important;`
          : ''
      "
    >
      <template v-if="respondents.length === 0">
        <div class="tw-text-very-dark-gray">No responses yet!</div>
      </template>
      <template v-else>
        <div
          v-for="(user, i) in respondents"
          :key="user._id"
          class="tw-group tw-relative tw-flex tw-cursor-pointer tw-items-center tw-overflow-hidden tw-overflow-visible tw-py-1"
          @mouseover="(e) => $emit('mouseOverRespondent', e, user._id)"
          @mouseleave="$emit('mouseLeaveRespondent')"
          @click="(e) => clickRespondent(e, user._id)"
        >
          <div class="tw-relative tw-flex tw-items-center">
            <UserAvatarContent
              v-if="!isGuest(user)"
              :user="user"
              class="-tw-ml-3 -tw-mr-1 tw-h-4 tw-w-4"
            ></UserAvatarContent>
            <v-icon v-else class="tw-ml-1 tw-mr-3" small>mdi-account</v-icon>

            <v-simple-checkbox
              @click="(e) => $emit('clickRespondent', e, user._id)"
              color="primary"
              :value="respondentSelected(user._id)"
              class="tw-absolute tw-left-0 tw-top-0 -tw-translate-y-1 tw-bg-white tw-bg-white tw-opacity-0 group-hover:tw-opacity-100"
              :class="
                respondentSelected(user._id) ? 'tw-opacity-100' : 'tw-opacity-0'
              "
            />
          </div>

          <div
            class="tw-mr-1 tw-transition-all"
            :class="respondentClass(user._id)"
          >
            {{
              user.firstName +
              " " +
              user.lastName +
              (respondentIfNeeded(user._id) ? "*" : "")
            }}
          </div>

          <!-- <div v-if="isGroup" class="tw-ml-1">
            <v-icon small class="tw-text-green">mdi-calendar-check</v-icon>
          </div> -->

          <v-btn
            v-if="!authUser && isGuest(user)"
            absolute
            small
            icon
            class="tw-right-0 tw-bg-white tw-opacity-0 tw-transition-none group-hover:tw-opacity-100"
            @click="$emit('editGuestAvailability', user._id)"
            ><v-icon small color="#4F4F4F">mdi-pencil</v-icon></v-btn
          >
          <v-btn
            v-else-if="isOwner && !isGroup"
            absolute
            small
            icon
            class="tw-right-0 tw-bg-white tw-opacity-0 tw-transition-none group-hover:tw-opacity-100"
            @click="() => showDeleteAvailabilityDialog(user)"
            ><v-icon small class="hover:tw-text-red" color="#4F4F4F"
              >mdi-delete</v-icon
            ></v-btn
          >
        </div>
      </template>
      <div
        class="tw-col-span-full tw-mt-2 tw-text-dark-gray"
        :class="showIfNeededStar ? 'tw-visible' : 'tw-invisible'"
      >
        * if needed
      </div>
    </div>

    <div
      v-if="pendingUsers.length > 0"
      class="tw-mb-2 tw-flex tw-items-center tw-font-medium"
    >
      <div class="tw-mr-1 tw-text-lg">Pending</div>
      <div class="tw-font-normal">({{ pendingUsers.length }})</div>
    </div>

    <div>
      <div v-for="(user, i) in pendingUsers" :key="user.email">
        <div class="tw-relative tw-flex tw-items-center">
          <v-icon class="tw-ml-1 tw-mr-3" small>mdi-account</v-icon>
          <div class="tw-mr-1 tw-text-sm tw-transition-all">
            {{ user.email }}
          </div>
        </div>
      </div>
    </div>

    <v-dialog v-model="deleteAvailabilityDialog" width="500" persistent>
      <v-card>
        <v-card-title>Are you sure?</v-card-title>
        <v-card-text class="tw-text-sm tw-text-dark-gray"
          >Are you sure you want to delete
          <strong>{{ userToDelete?.firstName }}</strong
          >'s availability from this
          {{ isGroup ? "group" : "event" }}?</v-card-text
        >
        <v-card-actions>
          <v-spacer />
          <v-btn text @click="deleteAvailabilityDialog = false">Cancel</v-btn>
          <v-btn
            text
            color="error"
            @click="
              () => {
                deleteAvailability(userToDelete)
                deleteAvailabilityDialog = false
              }
            "
            >Delete</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>

    <div v-if="isGroup">
      <v-switch
        inset
        class="tw-mb-4"
        :input-value="showCalendarEvents"
        @change="(val) => $emit('update:showCalendarEvents', Boolean(val))"
        hide-details
      >
        <template v-slot:label>
          <div class="tw-text-sm tw-text-black">Overlay calendar events</div>
        </template>
      </v-switch>
    </div>
  </div>
</template>

<script>
import { _delete, getDateHoursOffset } from "@/utils"
import UserAvatarContent from "../UserAvatarContent.vue"
import { mapState, mapActions } from "vuex"

export default {
  name: "RespondentsList",

  components: { UserAvatarContent },

  props: {
    day: { type: Object, required: false }, // Day of the current timeslot
    time: { type: Object, required: false }, // Time of the current timeslot
    curRespondent: { type: String, required: true },
    curRespondents: { type: Array, required: true },
    curTimeslot: { type: Object, required: true },
    curTimeslotAvailability: { type: Object, required: true },
    eventId: { type: String, required: true },
    respondents: { type: Array, required: true },
    parsedResponses: { type: Object, required: true },
    isOwner: { type: Boolean, required: true },
    maxHeight: { type: Number },
    isGroup: { type: Boolean, required: true },
    attendees: { type: Array, default: () => [] },
    showCalendarEvents: { type: Boolean, required: true },
  },

  data() {
    return {
      deleteAvailabilityDialog: false,
      userToDelete: null,
    }
  },

  computed: {
    ...mapState(["authUser"]),
    curRespondentsSet() {
      return new Set(this.curRespondents)
    },
    isCurTimeslotSelected() {
      return (
        this.curTimeslot.dayIndex !== -1 && this.curTimeslot.timeIndex !== -1
      )
    },
    numUsersAvailable() {
      this.curTimeslot
      let numUsers = 0
      for (const key in this.curTimeslotAvailability) {
        if (this.curTimeslotAvailability[key]) numUsers++
      }
      return numUsers
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
    pendingUsers() {
      if (!this.isGroup) return []

      const respondentEmailsSet = new Set(this.respondents.map((r) => r.email))

      return this.attendees.filter((a) => {
        if (!a.declined && !respondentEmailsSet.has(a.email)) {
          return true
        }
        return false
      })
    },
    showIfNeededStar() {
      for (const user of this.respondents) {
        if (this.respondentIfNeeded(user._id)) {
          return true
        }
      }
      return false
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo"]),
    /** Emit clickRespondent event */
    clickRespondent(e, userId) {
      e.stopImmediatePropagation()
      this.$emit("clickRespondent", e, userId)
    },
    /** Returns the class of the given respondent */
    respondentClass(id) {
      const c = []
      if (/*this.curRespondent == id ||*/ this.curRespondentsSet.has(id)) {
        // c.push("tw-font-bold")
      } else if (this.curRespondents.length > 0) {
        c.push("tw-text-gray")
      }

      if (!this.curTimeslotAvailability[id]) {
        c.push("tw-line-through")
        c.push("tw-text-gray")
      }
      return c
    },
    /** Returns whether the respondent has "ifNeeded" availability for the current timeslot */
    respondentIfNeeded(id) {
      if (!this.day || !this.time) return false

      const date = getDateHoursOffset(
        this.day.dateObject,
        this.time.hoursOffset
      )
      return Boolean(this.parsedResponses[id]?.ifNeeded?.has(date.getTime()))
    },
    /** Returns whether the current respondent is selected (for subset avail) */
    respondentSelected(id) {
      return this.curRespondentsSet.has(id)
    },
    /** Returns whether the user is a guest */
    isGuest(user) {
      return user._id == user.firstName
    },
    /** Shows the delete availability dialog */
    showDeleteAvailabilityDialog(user) {
      this.deleteAvailabilityDialog = true
      this.userToDelete = user
    },
    /** Deletes the user's availability on the server */
    async deleteAvailability(user) {
      try {
        await _delete(`/events/${this.eventId}/response`, {
          guest: this.isGuest(user),
          userId: user._id,
          name: user._id,
        })
        this.$emit("refreshEvent")
        this.showInfo("Availability successfully deleted!")

        this.$posthog?.capture("Deleted availability of another user", {
          eventId: this.eventId,
          userId: user._id,
        })
      } catch (e) {
        console.error(e)
        this.showError(
          "There was an error deleting that person's availability!"
        )
      }
    },
  },
}
</script>
