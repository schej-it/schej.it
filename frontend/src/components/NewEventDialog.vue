<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-[28rem] tw-m-0 tw-max-h-full"
  >
    <NewEvent
      :event="event"
      :editEvent="editEvent"
      :allow-notifications="allowNotifications"
      @input="$emit('input', false)"
    />
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
