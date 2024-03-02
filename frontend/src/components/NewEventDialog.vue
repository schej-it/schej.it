<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-[28rem] tw-m-0 tw-max-h-full"
  >
    <v-expand-transition>
      <v-card v-if="state == states.TYPE" class="tw-p-4 sm:tw-p-6">
        <div class="tw-text-md tw-pb-4 tw-text-center">
          What would you like to create?
        </div>
        <v-btn block class="tw-mb-2" @click="state = states.CREATE_EVENT"
          >Event</v-btn
        >
        <v-btn block @click="state = states.CREATE_GROUP"
          >Availability group</v-btn
        >
      </v-card>
      <NewEvent
        v-else-if="state == states.CREATE_EVENT"
        :event="event"
        :editEvent="editEvent"
        :allow-notifications="allowNotifications"
        @input="$emit('input', false)"
      />
    </v-expand-transition>
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

  data: () => ({
    states: {
      TYPE: "type", // Let user choose whether to create a group or event
      CREATE_EVENT: "create_event", // Event creation screen
      CREATE_GROUP: "create_group", // Group creation screen
    },
    state: "type",
  }),

  created() {},

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {},

  watch: {
    value() {
      if (!this.value) {
        setTimeout(() => {
          this.state = this.states.TYPE
        }, 300)
      }
    },
  },
}
</script>
