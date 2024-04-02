<template>
  <div v-if="event" class="tw-h-full">
    <AccessDenied v-if="accessDenied" />
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
import { errors } from "@/constants"
import AccessDenied from "@/components/groups/AccessDenied.vue"

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
  },

  data() {
    return {
      event: null,
    }
  },

  computed: {
    ...mapState(["authUser"]),
    accessDenied() {
      if (this.event.ownerId === this.authUser._id) {
        return false
      }

      const attendees = this.event?.attendees
      if (!attendees) return true

      let found = false
      for (const attendee of attendees) {
        if (attendee.email === this.authUser.email) {
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
