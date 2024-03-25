<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-[28rem] tw-m-0 tw-max-h-full"
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
  >
    <v-card class="tw-py-4">
      <div v-if="!noTabs" class="-tw-mt-4 tw-flex tw-rounded sm:tw-px-8">
        <v-tabs v-model="tab">
          <v-tab v-for="t in tabs" :key="t.type" class="tw-text-xs">{{ t.title }}</v-tab>
        </v-tabs>
        <v-spacer />
        <v-btn
          absolute
          @click="$emit('input', false)"
          icon
          class="tw-right-0 tw-mr-2 tw-self-center"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </div>

      <v-expand-transition>
        <div v-show="createEvent">
          <NewEvent
            :event="event"
            :edit="edit"
            :allow-notifications="allowNotifications"
            @input="$emit('input', false)"
            :contactsPayload="contactsPayload"
            :show-help="!noTabs"
          />
        </div>
      </v-expand-transition>
      <v-expand-transition>
        <div v-show="createGroup">
          <NewGroup
            :event="event"
            :edit="edit"
            @input="$emit('input', false)"
            :show-help="!noTabs"
          />
        </div>
      </v-expand-transition>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"
import NewEvent from "@/components/NewEvent.vue"
import NewGroup from "./NewGroup.vue"

export default {
  name: "NewDialog",

  emits: ["input"],

  props: {
    value: { type: Boolean, required: true },
    type: { type: String, default: "event" }, // Either "event" or "group"
    event: { type: Object },
    edit: { type: Boolean, default: false },
    allowNotifications: { type: Boolean, default: true },
    contactsPayload: { type: Object, default: () => ({}) },
    noTabs: { type: Boolean, default: false },
  },

  components: {
    NewEvent,
    NewGroup,
  },

  data: () => ({
    tab: 0,
    tabs: [
      { title: "Event", type: "event" },
      { title: "Availability group", type: "group" },
    ],
    TAB_TYPES: {
      EVENT: "event",
      GROUP: "group",
    },
  }),

  created() {},

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    createEvent() {
      if (this.noTabs) {
        return this.type === this.TAB_TYPES.EVENT
      } else {
        return this.tabs[this.tab].type === this.TAB_TYPES.EVENT
      }
    },
    createGroup() {
      if (this.noTabs) {
        return this.type === this.TAB_TYPES.GROUP
      } else {
        return this.tabs[this.tab].type === this.TAB_TYPES.GROUP
      }
    },
  },

  methods: {},
}
</script>
