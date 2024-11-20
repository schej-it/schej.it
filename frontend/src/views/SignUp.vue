<template>
  <div v-if="event" class="tw-h-full">
    <Event
      :eventId="signUpId"
      :fromSignIn="fromSignIn"
      :editingMode="editingMode"
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
import NotSignedIn from "@/components/groups/NotSignedIn.vue"

export default {
  name: "SignUp",

  props: {
    signUpId: { type: String, required: true },
    fromSignIn: { type: Boolean, default: false },
    editingMode: { type: Boolean, default: false },
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
  },

  methods: {
    ...mapActions(["showError"]),
  },

  async created() {
    try {
      this.event = await get(`/events/${this.signUpId}`)

      // TODO(tony): Redirect to correct routeif we're at the wrong route
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
