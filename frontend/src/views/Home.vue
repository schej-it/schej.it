<template>
  <div
    class="tw-mx-auto tw-mb-24 tw-mt-4 tw-max-w-6xl tw-space-y-4 sm:tw-mb-12 sm:tw-mt-7"
  >
    <!-- Preload images -->
    <div class="tw-hidden">
      <img src="@/assets/doodles/boba/0.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/1.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/2.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/3.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/4.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/5.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/6.jpg" alt="preload" />
      <img src="@/assets/doodles/boba/7.jpg" alt="preload" />
    </div>
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

    <div
      class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]"
      v-if="!loading || eventsNotEmpty"
    >
      <div
        class="tw-mb-3 tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
      >
        Tools
      </div>
      <div class="tw-flex tw-flex-row tw-items-center tw-gap-2">
        <div
          @click="showW2MDialog = true"
          class="tw-cursor-pointer tw-text-sm tw-font-normal tw-text-dark-gray tw-underline"
        >
          Convert When2meet to Schej
        </div>
      </div>
    </div>

    <div v-if="!loading || eventsNotEmpty" class="tw-flex tw-justify-center">
      <div
        class="animate-boba tw-size-48 tw-bg-contain tw-bg-no-repeat sm:tw-size-48"
      ></div>
    </div>

    <div class="tw-flex tw-flex-col tw-items-center tw-justify-between">
      <router-link
        class="tw-text-xs tw-font-medium tw-text-gray"
        :to="{ name: 'privacy-policy' }"
      >
        Privacy Policy
      </router-link>
    </div>

    <!-- FAB -->
    <BottomFab v-if="isPhone" id="create-event-btn" @click="createNew">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>

    <!-- When2meet Import Dialog -->
    <When2meetImportDialog v-model="showW2MDialog" />
  </div>
</template>

<script>
import EventType from "@/components/EventType.vue"
import BottomFab from "@/components/BottomFab.vue"
import CreateSpeedDial from "@/components/CreateSpeedDial.vue"
import When2meetImportDialog from "@/components/When2meetImportDialog.vue"
import { mapState, mapActions, mapMutations } from "vuex"
import { eventTypes } from "@/constants"
import { isPhone, get } from "@/utils"

export default {
  name: "Home",

  metaInfo: {
    title: "Home - Schej",
  },

  components: {
    EventType,
    BottomFab,
    CreateSpeedDial,
    When2meetImportDialog,
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
    showW2MDialog: false,
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
    ...mapMutations(["setAuthUser"]),
    ...mapActions(["getEvents"]),
    userRespondedToEvent(event) {
      return event.hasResponded ?? false
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
    get("/user/profile")
      .then((authUser) => {
        this.setAuthUser(authUser)
      })
      .catch(() => {
        this.setAuthUser(null)
      })
  },
}
</script>

<style>
@keyframes boba {
  0% {
    background-image: url("@/assets/doodles/boba/0.jpg");
  }
  12.5% {
    background-image: url("@/assets/doodles/boba/1.jpg");
  }
  25% {
    background-image: url("@/assets/doodles/boba/2.jpg");
  }
  37.5% {
    background-image: url("@/assets/doodles/boba/3.jpg");
  }
  50% {
    background-image: url("@/assets/doodles/boba/4.jpg");
  }
  62.5% {
    background-image: url("@/assets/doodles/boba/5.jpg");
  }
  75% {
    background-image: url("@/assets/doodles/boba/6.jpg");
  }
  87.5% {
    background-image: url("@/assets/doodles/boba/7.jpg");
  }
  100% {
    background-image: url("@/assets/doodles/boba/0.jpg");
  }
}

.animate-boba {
  animation: boba 1.04s steps(1) infinite;
  animation-play-state: paused;
  transition: animation-play-state 0s 1.04s;
}

.animate-boba:hover {
  animation-play-state: running;
  transition: animation-play-state 0s;
}
</style>
