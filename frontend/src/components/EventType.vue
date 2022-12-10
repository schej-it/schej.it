<!-- Displays an event type (i.e. created or joined) -->
<template>
  <div class="tw-mb-5">
    <div class="tw-text-2xl tw-font-bold tw-text-dark-green">
      {{ eventType.header }}
    </div>

    <div v-if="eventType.events.length === 0" class="tw-my-3">
      No events yet!
    </div>
    <div
      v-else
      class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 md:tw-grid-cols-2 tw-gap-2 tw-my-3"
    >
      <EventItem
        class="tw-cursor-pointer"
        v-for="(event, i) in eventType.events.slice(0, numEventsToShow)"
        :key="i"
        :event="event"
        @click="goToEvent(event._id)"
      />
      <div v-if="eventType.events.length > numEventsToShow">
        <v-expand-transition>
          <div v-show="showAll">
            <EventItem
              class="tw-cursor-pointer"
              v-for="(event, i) in eventType.events.slice(numEventsToShow)"
              :key="i"
              :event="event"
              @click="goToEvent(event._id)"
            />
          </div>
        </v-expand-transition>

        <div @click="toggleShowAll">See all</div>
      </div>
    </div>
  </div>
</template>

<script>
import EventItem from "@/components/EventItem";

export default {
  name: "EventType",

  components: {
    EventItem,
  },

  props: {
    eventType: { type: Object, required: true },
  },

  data: () => ({
    numEventsToShow: 4,
    showAll: false,
  }),

  computed: {},

  methods: {
    toggleShowAll() {
      this.showAll = !this.showAll;
    },
    goToEvent(eventId) {
      this.$router.push({ name: "event", params: { eventId } });
    },
  },
};
</script>
