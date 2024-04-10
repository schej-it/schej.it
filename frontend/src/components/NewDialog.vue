<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-[28rem]"
    :fullscreen="isPhone"
    scrollable
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
  >
    <v-card class="tw-py-4">
      <div v-if="!noTabs" class="tw-flex tw-rounded sm:-tw-mt-4 sm:tw-px-8">
        <v-tabs v-model="tab">
          <v-tab
            v-for="t in tabs"
            :key="t.type"
            :tab-value="t.type"
            class="tw-text-xs"
            >{{ t.title }}</v-tab
          >
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

      <template>
        <NewEvent
          v-if="tab === 'event'"
          key="event"
          :event="event"
          :edit="edit"
          :allow-notifications="allowNotifications"
          @input="$emit('input', false)"
          :contactsPayload="contactsPayload"
          :show-help="!noTabs"
        />
        <NewGroup
          v-else-if="tab === 'group'"
          key="group"
          :event="event"
          :edit="edit"
          @input="$emit('input', false)"
          :show-help="!noTabs"
        />
      </template>
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

  data() {
    return {
      tab: this.type,
      tabs: [
        { title: "Event", type: "event" },
        { title: "Availability group", type: "group" },
      ],
      TAB_TYPES: {
        EVENT: "event",
        GROUP: "group",
      },
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {},
}
</script>
