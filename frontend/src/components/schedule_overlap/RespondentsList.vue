<template>
  <div>
    <div class="tw-flex tw-items-center tw-font-medium">
      <div class="tw-mr-1 tw-text-lg">Responses</div>
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
      class="tw-mt-2 tw-grid tw-grid-cols-2 tw-gap-x-2 tw-text-sm sm:tw-block"
    >
      <template v-if="respondents.length === 0">
        <div class="tw-text-very-dark-gray">No responses yet!</div>
      </template>
      <template v-else>
        <div
          v-for="(user, i) in respondents"
          :key="user._id"
          class="tw-group tw-relative tw-flex tw-cursor-pointer tw-items-center tw-overflow-hidden tw-py-1"
          @mouseover="(e) => $emit('mouseOverRespondent', e, user._id)"
          @mouseleave="$emit('mouseLeaveRespondent')"
          @click="(e) => $emit('clickRespondent', e, user._id)"
        >
          <UserAvatarContent
            v-if="!isGuest(user)"
            :user="user"
            class="-tw-ml-3 -tw-mr-1 tw-h-4 tw-w-4"
          ></UserAvatarContent>
          <v-icon v-else class="tw-ml-1 tw-mr-3" small>mdi-account</v-icon>

          <div
            class="tw-mr-1 tw-transition-all"
            :class="respondentClass(user._id)"
          >
            {{ user.firstName + " " + user.lastName }}
          </div>

          <v-btn
            v-if="!authUser && isGuest(user)"
            absolute
            small
            icon
            class="tw-right-0 tw-bg-white tw-opacity-0 group-hover:tw-opacity-100"
            @click="$emit('editGuestAvailability', user._id)"
            ><v-icon small color="#4F4F4F">mdi-pencil</v-icon></v-btn
          >
          <v-btn
            v-else-if="isOwner"
            absolute
            small
            icon
            class="tw-right-0 tw-bg-white tw-opacity-0 group-hover:tw-opacity-100"
            @click="() => deleteAvailability(user)"
            ><v-icon small color="#4F4F4F">mdi-close</v-icon></v-btn
          >
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import { _delete } from "@/utils"
import UserAvatarContent from "../UserAvatarContent.vue"
import { mapState, mapActions } from "vuex"

export default {
  name: "RespondentsList",

  components: { UserAvatarContent },

  props: {
    curRespondent: { type: String, required: true },
    curRespondents: { type: Array, required: true },
    curTimeslot: { type: Object, required: true },
    curTimeslotAvailability: { type: Object, required: true },
    eventId: { type: String, required: true },
    respondents: { type: Array, required: true },
    isOwner: { type: Boolean, required: true },
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
  },

  methods: {
    ...mapActions(["showError", "showInfo"]),
    respondentClass(id) {
      const c = []
      if (this.curRespondent == id || this.curRespondentsSet.has(id)) {
        c.push("tw-font-bold")
      } else if (this.curRespondents.length > 0) {
        c.push("tw-text-gray")
      }

      if (!this.curTimeslotAvailability[id]) {
        c.push("tw-line-through")
        c.push("tw-text-gray")
      }
      return c
    },
    isGuest(user) {
      return user._id == user.firstName
    },
    async deleteAvailability(user) {
      try {
        await _delete(`/events/${this.eventId}/response`, {
          guest: this.isGuest(user),
          userId: user._id,
          name: user._id,
        })
        this.$emit("refreshEvent")
        this.showInfo("Availability successfully deleted!")
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
