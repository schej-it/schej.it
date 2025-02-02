<template>
  <div class="tw-mx-auto tw-mb-12 tw-mt-4 tw-max-w-6xl tw-space-y-4 sm:tw-mt-7">
    <div
      v-if="loading && !eventsNotEmpty"
      class="tw-flex tw-h-[calc(100vh-10rem)] tw-w-full tw-items-center tw-justify-center"
    >
      <v-progress-circular
        indeterminate
        color="primary"
        :size="20"
        :width="2"
      ></v-progress-circular>
    </div>
    <template v-if="groupsEnabled">
      <v-fade-transition>
        <div
          class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]"
          v-if="!loading || eventsNotEmpty"
        >
          <EventType
            :eventType="availabilityGroups"
            emptyText="You are not part of any availability groups!"
          />
        </div>
      </v-fade-transition>
    </template>
    <v-fade-transition>
      <div
        class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]"
        v-if="!loading || eventsNotEmpty"
      >
        <div class="tw-grid tw-gap-4 sm:tw-gap-8">
          <EventType
            v-for="(eventType, t) in events"
            :key="t"
            :eventType="eventType"
          ></EventType>
        </div>
      </div>
    </v-fade-transition>

    <!-- FAB -->
    <BottomFab v-if="isPhone" id="create-event-btn" @click="createNew">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>
  </div>
</template>

<script>
import EventType from "@/components/EventType.vue"
import BottomFab from "@/components/BottomFab.vue"
import CreateSpeedDial from "@/components/CreateSpeedDial.vue"
import { mapState, mapActions } from "vuex"
import { eventTypes } from "@/constants"
import { isPhone } from "@/utils"

export default {
  name: "Home",

  metaInfo: {
    title: "Home - Schej",
  },

  components: {
    EventType,
    BottomFab,
    CreateSpeedDial,
  },

  props: {
    contactsPayload: {
      type: Object,
      default: () => ({}),
    },
    openNewGroup: { type: Boolean, default: false },
  },

  data: () => ({
    loading: true,
  }),

  mounted() {
    // If coming from enabling contacts, show the dialog. Checks if contactsPayload is not an Observer.
    this.$emit("setNewDialogOptions", {
      show: Object.keys(this.contactsPayload).length > 0 || this.openNewGroup,
      contactsPayload: this.contactsPayload,
      openNewGroup: this.openNewGroup,
    })
  },

  computed: {
    ...mapState(["createdEvents", "joinedEvents", "authUser", "groupsEnabled"]),
    events() {
      return [
        {
          header: "Events I created",
          events: this.createdEventsNonGroup,
        },
        {
          header: "Events I joined",
          events: this.joinedEventsNonGroup,
        },
      ]
    },
    createdEventsNonGroup() {
      return this.createdEvents.filter((e) => e.type !== eventTypes.GROUP)
    },
    joinedEventsNonGroup() {
      return this.joinedEvents.filter((e) => e.type !== eventTypes.GROUP)
    },
    availabilityGroups() {
      return {
        header: "Availability groups",
        events: this.createdEvents
          .filter((e) => e.type === eventTypes.GROUP)
          .concat(this.joinedEvents.filter((e) => e.type === eventTypes.GROUP))
          .sort((e1, e2) => (this.userRespondedToEvent(e1) ? 1 : -1)),
      }
    },
    eventsNotEmpty() {
      return this.createdEvents.length > 0 || this.joinedEvents.length > 0
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    ...mapActions(["getEvents"]),
    userRespondedToEvent(event) {
      return this.authUser._id in event.responses
    },
    createNew() {
      this.$emit("setNewDialogOptions", {
        show: true,
        contactsPayload: {},
        openNewGroup: false,
      })
    },
  },

  created() {
    this.getEvents().then(() => {
      this.loading = false
    })
  },
}
</script>
