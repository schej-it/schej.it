<template>
  <router-link :to="{ name: 'event', params: { eventId: event._id } }">
    <v-container
      v-ripple
      class="tw-flex tw-items-center tw-justify-between tw-rounded-lg tw-bg-white tw-px-4 tw-py-2.5 tw-text-black tw-drop-shadow tw-transition-all hover:tw-drop-shadow-md sm:tw-py-3"
      :data-ph-capture-attribute-event-id="event._id"
      :data-ph-capture-attribute-event-name="event.name"
    >
      <div class="tw-ml-1">
        <div>{{ this.event.name }}</div>
        <div class="tw-text-sm tw-font-light tw-text-very-dark-gray">
          {{ dateString }}
        </div>
      </div>
      <div class="tw-min-w-max">
        <v-chip small class="tw-m-0.5 tw-bg-off-white tw-text-very-dark-gray">
          <v-icon left small> mdi-account-multiple </v-icon>
          {{ Object.keys(this.event.responses).length }}
        </v-chip>
        <v-menu
          v-if="showOptions"
          ref="menu"
          :close-on-content-click="false"
          transition="slide-x-transition"
          right
          offset-x
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn plain icon v-bind="attrs" v-on="on" @click.prevent>
              <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
          </template>

          <v-list justify="center" class="tw-py-1">
            <v-dialog v-model="removeDialog" width="400" persistent>
              <template v-slot:activator="{ on, attrs }">
                <v-btn
                  id="delete-event-btn"
                  text
                  class="red--text tw-px-6"
                  v-bind="attrs"
                  v-on="on"
                  block
                  >Delete event</v-btn
                >
              </template>
              <v-card>
                <v-card-title>Are you sure?</v-card-title>
                <v-card-text
                  >Are you sure you want to delete this event?</v-card-text
                >
                <v-card-actions>
                  <v-spacer />
                  <v-btn text @click="removeDialog = false">Cancel</v-btn>
                  <v-btn text color="error" @click="removeEvent()"
                    >I'm sure</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-list>
        </v-menu>
        <v-icon v-else class="tw-mx-1 tw-opacity-75">mdi-chevron-right</v-icon>
      </div>
    </v-container>
  </router-link>
</template>

<script>
import { getDateRangeStringForEvent, _delete, isPhone } from "@/utils"
import { mapActions, mapState } from "vuex"

export default {
  name: "EventItem",

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({
    removeDialog: false,
  }),

  computed: {
    ...mapState(["authUser"]),
    dateString() {
      return getDateRangeStringForEvent(this.event)
    },
    showOptions() {
      return this.event.ownerId === this.authUser._id
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo", "getEvents"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    removeEvent() {
      _delete(`/events/${this.event._id}`)
        .then(() => {
          this.getEvents()
          this.$refs.menu.save() // NOTE: Not sure why but without this line, the menu persists to the next event.
        })
        .catch((err) => {
          this.showError(
            "There was a problem removing that event! Please try again later."
          )
        })
    },
  },
}
</script>
