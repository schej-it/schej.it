<template>
  <div v-if="event" class="tw-h-full">
    <NotSignedIn v-if="!authUser" :event="event" />
    <AccessDenied v-else-if="accessDenied" />
    <Event
      v-else
      :eventId="groupId"
      :fromSignIn="fromSignIn"
      :initialTimezone="initialTimezone"
      :contactsPayload="contactsPayload"
    ></Event>
  </div>
</template>

<script>
import Event from "./Event.vue"
import { mapActions, mapState } from "vuex"
import { get } from "@/utils"
import { errors, eventTypes } from "@/constants"
import AccessDenied from "@/components/groups/AccessDenied.vue"
import NotSignedIn from "@/components/groups/NotSignedIn.vue"

export default {
  name: "Group",

  props: {
    groupId: { type: String, required: true },
    fromSignIn: { type: Boolean, default: false },
    initialTimezone: { type: Object, default: () => ({}) },
    contactsPayload: { type: Object, default: () => ({}) },
  },

  components: {
    AccessDenied,
    Event,
    NotSignedIn,
  },

  data() {
    return {
      event: null,
    }
  },

  computed: {
    ...mapState(["authUser"]),
    accessDenied() {
      if (this.event.ownerId === this.authUser?._id) {
        return false
      }

      const attendees = this.event?.attendees
      if (!attendees) return true

      let found = false
      for (const attendee of attendees) {
        if (attendee.email.toLowerCase() === this.authUser.email.toLowerCase()) {
          // The line below is commented out because we want attendee to be able to rejoin group after declining
          // if (attendee.declined) return true

          found = true
          break
        }
      }

      return !found
    },
  },

  methods: {
    ...mapActions(["showError"]),
  },

  async created() {
    try {
      this.event = await get(`/events/${this.groupId}`)

      // Redirect if we're at the wrong route
      if (this.event.type !== eventTypes.GROUP) {
        this.$router.replace({
          name: "event",
          params: {
            eventId: this.groupId,
            initialTimezone: this.initialTimezone,
            fromSignIn: this.fromSignIn,
            contactsPayload: this.contactsPayload,
          },
        })
      }
    } catch (err) {
      switch (err.error) {
        case errors.EventNotFound:
          this.showError("The specified event does not exist!")
          this.$router.replace({ name: "home" })
          return
      }
    }
  },
}
</script>
