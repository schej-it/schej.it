<template>
  <v-dialog
    :value="value"
    @click:outside="handleDialogInput"
    no-click-animation
    persistent
    content-class="tw-max-w-[28rem]"
    :fullscreen="isPhone"
    scrollable
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
  >
    <UnsavedChangesDialog v-model="unsavedChangesDialog" @leave="exitDialog">
    </UnsavedChangesDialog>
    <v-card class="tw-pt-4">
      <div v-if="!_noTabs" class="tw-flex tw-rounded sm:-tw-mt-4 sm:tw-px-8">
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
          ref="event"
          key="event"
          :event="event"
          :edit="edit"
          :allow-notifications="allowNotifications"
          @input="$emit('input', false)"
          @update:formEmpty="(val) => (formEmpty = val)"
          :contactsPayload="this.type == 'event' ? contactsPayload : {}"
          :show-help="!_noTabs"
        />
        <NewGroup
          v-else-if="tab === 'group'"
          ref="group"
          key="group"
          :event="event"
          :edit="edit"
          @input="$emit('input', false)"
          @update:formEmpty="(val) => (formEmpty = val)"
          :show-help="!_noTabs"
          :calendarPermissionGranted="calendarPermissionGranted"
          :contactsPayload="this.type == 'group' ? contactsPayload : {}"
        />
      </template>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"
import NewEvent from "@/components/NewEvent.vue"
import UnsavedChangesDialog from "@/components/general/UnsavedChangesDialog.vue"
import NewGroup from "./NewGroup.vue"
import { mapState } from "vuex"

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
    calendarPermissionGranted: { type: Boolean, default: true },
  },

  components: {
    NewEvent,
    NewGroup,
    UnsavedChangesDialog,
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

      unsavedChangesDialog: false,
      formEmpty: true,
    }
  },

  computed: {
    ...mapState(["groupsEnabled"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    _noTabs() {
      if (!this.groupsEnabled) {
        return true
      } else {
        return this.noTabs
      }
    },
  },

  methods: {
    handleDialogInput() {
      if (this.formEmpty) {
        this.$emit("input", false)
      } else {
        this.unsavedChangesDialog = true
      }
    },
    exitDialog() {
      this.$emit("input", false)
      this.$refs[this.tab].reset()
    },
  },

  watch: {
    groupsEnabled: {
      immediate: true,
      handler() {
        if (this.groupsEnabled) {
          this.tabs = [
            { title: "Event", type: "event" },
            { title: "Availability group", type: "group" },
          ]
        } else {
          this.tabs = [{ title: "Event", type: "event" }]
        }
      },
    },
    tab() {
      this.formEmpty = true
    },
  },
}
</script>
