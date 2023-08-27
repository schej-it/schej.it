<!-- Displays an event type (i.e. created or joined) on the home page -->
<template>
  <div class="tw-mb-5">
    <div class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl tw-flex tw-flex-row tw-items-center tw-justify-between">
      {{ eventType.header }}
      <div
        @click="toggleShowAll"
        class="tw-mt-2 tw-font-normal tw-cursor-pointer tw-text-sm tw-text-very-dark-gray sm:tw-hidden"
      >
        Show {{ showAll ? "less" : "more"
        }}<v-icon :class="showAll && 'tw-rotate-180'">mdi-chevron-down</v-icon>
      </div>
    </div>

    <div v-if="eventType.events.length === 0" class="tw-my-3">
      No events yet!
    </div>
    <div
      v-else
      class="tw-gr id-cols-1 tw-my-3 tw-grid tw-gap-3 sm:tw-grid-cols-2 md:tw-grid-cols-2"
    >
      <EventItem
        class="tw-cursor-pointer"
        v-for="(event, i) in sortedEvents.slice(0, DEFAULT_NUM_EVENTS_TO_SHOW)"
        :key="i"
        :event="event"
        @click="goToEvent(event._id)"
      />
    </div>
    <!-- Show more events sections -->
    <!-- TODO: might want to change for less code repeat -->
    <div v-if="eventType.events.length > DEFAULT_NUM_EVENTS_TO_SHOW">
      <v-expand-transition>
        <div
          v-if="showAll"
          class="tw-gr id-cols-1 tw-grid tw-gap-2 sm:tw-grid-cols-2 sm:tw-gap-4 md:tw-grid-cols-2"
        >
          <EventItem
            v-for="(event, i) in sortedEvents.slice(
              DEFAULT_NUM_EVENTS_TO_SHOW,
              eventType.events.length
            )"
            :key="i"
            class="tw-cursor-pointer"
            :event="event"
            @click="goToEvent(event._id)"
          />
        </div>
      </v-expand-transition>
      <div
        @click="toggleShowAll"
        class="tw-mt-4 tw-cursor-pointer tw-text-sm tw-text-very-dark-gray tw-hidden sm:tw-block"
      >
        Show {{ showAll ? "less" : "more"
        }}<v-icon :class="showAll && 'tw-rotate-180'">mdi-chevron-down</v-icon>
      </div>
    </div>
  </div>
</template>

<script>
import EventItem from "@/components/EventItem.vue"

export default {
  name: "EventType",

  components: {
    EventItem,
  },

  props: {
    eventType: { type: Object, required: true },
  },

  data: () => ({
    showAll: false,
    DEFAULT_NUM_EVENTS_TO_SHOW: 4,
  }),

  computed: {
    numEventsToShow() {
      return this.showAll
        ? this.eventType.events.length
        : this.DEFAULT_NUM_EVENTS_TO_SHOW
    },
    sortedEvents() {
      // Events are sorted serverside, so no need to sort here
      return this.eventType.events
    },
  },

  methods: {
    toggleShowAll() {
      this.showAll = !this.showAll
    },
    goToEvent(eventId) {
      this.$router.push({ name: "event", params: { eventId } })
    },
  },
}
</script>
