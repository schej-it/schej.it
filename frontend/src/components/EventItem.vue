<template>
  <v-container
    @click="$emit('click')"
    v-ripple
    class="hover:tw-drop-shadow-lg tw-drop-shadow-md tw-shadow-inner tw-transition-all tw-bg-white tw-rounded-lg tw-flex tw-text-black tw-justify-between tw-items-center tw-py-3"
  >
    <div class="tw-ml-1">
      <div>{{ this.event.name }}</div>
      <div class="tw-text-sm tw-font-light">{{ dateString }}</div>
    </div>
    <div class="tw-min-w-max">
      <v-chip small class="tw-text-very-dark-gray tw-m-2 tw-bg-off-white">
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
          <v-btn plain icon v-bind="attrs" v-on="on">
            <v-icon>mdi-dots-vertical</v-icon>
          </v-btn>
        </template>

        <v-list justify="center">
          <v-dialog v-model="removeDialog" width="400" persistent>
            <template v-slot:activator="{ on, attrs }">
              <v-btn id="delete-event-btn" text small class="red--text" v-bind="attrs" v-on="on"
                >Delete</v-btn
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
</template>

<script>
import { getDateRangeStringForEvent, _delete } from "@/utils"
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
    ...mapState([ 'authUser' ]),
    dateString() {
      return getDateRangeStringForEvent(this.event)
    },
    showOptions() {
      return this.event.ownerId === this.authUser._id
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo", "getEvents"]),
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
