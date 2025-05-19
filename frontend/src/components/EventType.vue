<!-- Displays an event type (i.e. created or joined) on the home page -->
<template>
  <div class="tw-mb-5">
    <div
      class="tw-flex tw-flex-row tw-items-center tw-justify-between tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
    >
      <div class="tw-flex tw-flex-col">
        {{ eventType.header }}
        <div
          v-if="
            eventType.header === 'Events I created' &&
            enablePaywall &&
            !isPremiumUser
          "
          class="tw-flex tw-items-center tw-gap-2 tw-text-sm tw-font-normal tw-text-very-dark-gray"
        >
          <div>
            {{ authUser?.numEventsCreated }} / {{ numFreeEvents }} free events
          </div>
          <div
            class="tw-text-xs tw-font-medium tw-text-green"
            @click="openUpgradeDialog"
          >
            Upgrade
          </div>
        </div>
      </div>
      <div
        v-if="eventType.header === 'Events I created'"
        class="tw-hidden tw-flex-row tw-items-center tw-gap-2 sm:tw-flex"
      >
        <div
          @click="showW2MDialog = true"
          class="tw-cursor-pointer tw-text-sm tw-font-normal tw-text-dark-gray tw-underline"
        >
          Convert When2meet to Schej
        </div>
        <div
          class="tw-rounded-md tw-bg-green tw-px-2 tw-py-1 tw-text-xs tw-text-white"
        >
          NEW
        </div>
      </div>
      <div
        v-if="eventType.events.length > defaultNumEventsToShow"
        @click="toggleShowAll"
        class="tw-mt-2 tw-cursor-pointer tw-text-sm tw-font-normal tw-text-very-dark-gray sm:tw-hidden"
      >
        Show {{ showAll ? "less" : "more"
        }}<v-icon :class="showAll && 'tw-rotate-180'">mdi-chevron-down</v-icon>
      </div>
    </div>

    <div
      v-if="eventType.events.length === 0"
      class="tw-my-3 tw-text-very-dark-gray"
    >
      {{ emptyText.length > 0 ? emptyText : "No events yet!" }}
    </div>
    <div
      v-else
      class="tw-gr id-cols-1 tw-my-3 tw-grid tw-gap-3 sm:tw-grid-cols-2 lg:tw-grid-cols-3"
    >
      <EventItem
        class="tw-cursor-pointer"
        v-for="(event, i) in sortedEvents.slice(0, defaultNumEventsToShow)"
        :key="i"
        :event="event"
      />
    </div>
    <!-- Show more events sections -->
    <!-- TODO: might want to change for less code repeat -->
    <div v-if="eventType.events.length > defaultNumEventsToShow">
      <v-expand-transition>
        <div
          v-if="showAll"
          class="tw-gr id-cols-1 tw-my-3 tw-grid tw-gap-3 sm:tw-grid-cols-2 lg:tw-grid-cols-3"
        >
          <EventItem
            v-for="(event, i) in sortedEvents.slice(
              defaultNumEventsToShow,
              eventType.events.length
            )"
            :key="i"
            class="tw-cursor-pointer"
            :event="event"
          />
        </div>
      </v-expand-transition>
      <div
        @click="toggleShowAll"
        class="tw-mt-4 tw-hidden tw-cursor-pointer tw-text-sm tw-text-very-dark-gray sm:tw-block"
      >
        Show {{ showAll ? "less" : "more"
        }}<v-icon :class="showAll && 'tw-rotate-180'">mdi-chevron-down</v-icon>
      </div>
    </div>
    <!-- Add the dialog component -->
    <When2meetImportDialog v-model="showW2MDialog" />
  </div>
</template>

<script>
import EventItem from "@/components/EventItem.vue"
import { numFreeEvents } from "@/constants"
import { mapState } from "vuex"
import { isPremiumUser } from "@/utils"
import When2meetImportDialog from "@/components/When2meetImportDialog.vue"

export default {
  name: "EventType",

  components: {
    EventItem,
    When2meetImportDialog,
  },

  props: {
    eventType: { type: Object, required: true },
    emptyText: { type: String, default: "" },
  },

  data: () => ({
    showAll: false,
    showW2MDialog: false,
  }),

  computed: {
    ...mapState(["authUser", "enablePaywall"]),
    defaultNumEventsToShow() {
      return this.$vuetify.breakpoint.lgAndUp ? 6 : 4
    },
    numEventsToShow() {
      return this.showAll
        ? this.eventType.events.length
        : this.defaultNumEventsToShow
    },
    sortedEvents() {
      // Events are sorted serverside, so no need to sort here
      return this.eventType.events
    },
    numFreeEvents() {
      return numFreeEvents
    },
    isPremiumUser() {
      return isPremiumUser(this.authUser)
    },
  },

  methods: {
    toggleShowAll() {
      this.showAll = !this.showAll
    },
    openUpgradeDialog() {
      this.$emit("openUpgradeDialog")
    },
  },
}
</script>
