<template>
  <Event
    :eventId="groupId"
    :fromSignIn="fromSignIn"
    :initialTimezone="initialTimezone"
    :contactsPayload="contactsPayload"
  ></Event>
</template>

<script>
import Event from "./Event.vue"

export default {
  name: "Group",

  props: {
    groupId: { type: String, required: true },
    fromSignIn: { type: Boolean, default: false },
    initialTimezone: { type: Object, default: () => ({}) },
    contactsPayload: { type: Object, default: () => ({}) },
  },

  components: {
    Event,
  },
}
</script>
