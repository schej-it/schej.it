<template>
  <div class="tw-mx-auto tw-mb-12 tw-mt-4 tw-max-w-6xl sm:tw-mt-7">
    <!-- Dialog -->
    <NewEventDialog v-model="dialog" :contactsPayload="contactsPayload" />
    <v-fade-transition>
      <div
        class="tw-grid tw-gap-4 tw-p-4 sm:tw-gap-8"
        v-if="!loading || events"
      >
        <EventType
          v-for="(eventType, t) in events"
          :key="t"
          :eventType="eventType"
        ></EventType></div
    ></v-fade-transition>

    <!-- FAB -->
    <BottomFab id="create-event-btn" @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>
  </div>
</template>

<script>
import NewEventDialog from "@/components/NewEventDialog.vue"
import EventType from "@/components/EventType.vue"
import BottomFab from "@/components/BottomFab.vue"
import { mapState, mapActions } from "vuex"
import { eventTypes } from "@/constants"

export default {
  name: "Home",

  components: {
    NewEventDialog,
    EventType,
    BottomFab,
  },

  props: {
    contactsPayload: {
      type: Object,
      default: () => ({}),
    },
  },

  data: () => ({
    dialog: false,
    loading: true,
  }),

  mounted() {
    // If coming from enabling contacts, show the dialog. Checks if contactsPayload is not an Observer.
    this.dialog = Object.keys(this.contactsPayload).length > 0
  },

  computed: {
    ...mapState(["createdEvents", "joinedEvents"]),
    events() {
      return [
        ...(this.weeklyEvents.length > 0
          ? [
              {
                header: "Weekly events",
                events: this.weeklyEvents,
              },
            ]
          : []),
        {
          header: "Events I created",
          events: this.createdEventsWithSpecificDates,
        },
        {
          header: "Events I joined",
          events: this.joinedEvents,
        },
      ]
    },
    createdEventsWithSpecificDates() {
      return this.createdEvents.filter((e) => e.type !== eventTypes.DOW)
    },
    weeklyEvents() {
      return this.createdEvents
        .filter((e) => e.type === eventTypes.DOW)
        .concat(this.joinedEvents.filter((e) => e.type === eventTypes.DOW))
    },
  },

  methods: {
    ...mapActions(["getEvents"]),
  },

  created() {
    this.getEvents().then(() => {
      this.loading = false
    })
  },
}
</script>
