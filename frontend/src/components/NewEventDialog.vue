<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    :fullscreen="isPhone"
    :hide-overlay="isPhone"
    content-class="tw-max-w-[28rem]"
    scrollable
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
  >
    <NewEvent
      :event="event"
      :editEvent="editEvent"
      :allow-notifications="allowNotifications"
      @input="$emit('input', false)"
    ></NewEvent>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"
import NewEvent from "@/components/NewEvent.vue"

export default {
  name: "NewEventDialog",

  emits: ["input"],

  props: {
    value: { type: Boolean, required: true },
    event: { type: Object },
    editEvent: { type: Boolean, default: false },
    allowNotifications: { type: Boolean, default: true },
  },

  components: {
    NewEvent,
  },

  data: () => ({}),

  created() {},

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {},
}
</script>
