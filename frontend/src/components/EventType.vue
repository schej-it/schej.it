<!-- Displays an event type (i.e. created or joined) on the home page -->
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
      class="tw-grid tw-gr id-cols-1 sm:tw-grid-cols-2 md:tw-grid-cols-2 tw-gap-2 sm:tw-gap-4 tw-my-3"
    >
      <EventItem
        class="tw-cursor-pointer"
        v-for="(event, i) in sortedEvents.slice(
          0,
          DEFAULT_NUM_EVENTS_TO_SHOW
        )"
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
          class="tw-grid tw-gr id-cols-1 sm:tw-grid-cols-2 md:tw-grid-cols-2 tw-gap-2 sm:tw-gap-4"
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
      <div @click="toggleShowAll" class="tw-cursor-pointer tw-text-very-dark-gray tw-text-sm" :class="showAll && 'tw-mt-2'">Show {{showAll ? 'less' : 'more'}}<v-icon :class="showAll && 'tw-rotate-180'">mdi-chevron-down</v-icon></div>
    </div>
  </div>
</template>

<script>
import EventItem from "@/components/EventItem.vue";

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
        : this.DEFAULT_NUM_EVENTS_TO_SHOW;
    },
    sortedEvents() {
      const sorted = [...this.eventType.events]
      sorted.sort((a, b) => {
        let aStartDate, bStartDate
        if (a.dates && a.dates.length > 0) {
          aStartDate = a.dates[0]
        } else {
          aStartDate = a.startDate
        }
        if (b.dates && b.dates.length > 0) {
          bStartDate = b.dates[0]
        } else {
          bStartDate = b.startDate
        }

        return new Date(bStartDate) - new Date(aStartDate)
      })
      return sorted
    },
  },

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
