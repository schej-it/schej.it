<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12 tw-mt-4 sm:tw-mt-7">
    <!-- Dialog -->
    <NewEventDialog v-model="dialog" />
    <div class="tw-grid tw-p-4 tw-gap-4 sm:tw-gap-8">
      <EventType
        v-for="(eventType, t) in events"
        :key="t"
        :eventType="eventType"
      ></EventType>
    </div>

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

export default {
  name: "Home",

  components: {
    NewEventDialog,
    EventType,
    BottomFab,
  },

  data: () => ({
    dialog: false,
  }),

  computed: {
    ...mapState(["createdEvents", "joinedEvents"]),
    events() {
      return [
        {
          header: "Events I created",
          events: this.createdEvents,
        },
        {
          header: "Events I joined",
          events: this.joinedEvents,
        },
      ]
    },
  },

  methods: {
    ...mapActions(["getEvents"]),
  },

  created() {
    this.getEvents()
  },
}
</script>
