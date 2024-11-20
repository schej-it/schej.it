<template>
  <router-link
    :to="{
      name: linkTo,
      params: { [identifier]: event.shortId ?? event._id },
    }"
  >
    <v-container
      v-ripple
      class="tw-flex tw-min-h-16 tw-items-center tw-justify-between tw-rounded-lg tw-bg-white tw-px-4 tw-py-2.5 tw-text-black tw-drop-shadow tw-transition-all hover:tw-drop-shadow-md sm:tw-py-3"
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
        <div
          v-if="isGroup && !userHasResponded"
          class="tw-inline-block tw-text-sm tw-italic tw-text-gray"
        >
          Invited
        </div>
        <v-chip
          v-else
          small
          class="tw-m-0.5 tw-bg-off-white tw-text-very-dark-gray"
        >
          <v-icon left small> mdi-account-multiple </v-icon>
          {{ Object.keys(this.event.responses).length }}
        </v-chip>
        <v-menu
          v-if="showOptions"
          v-model="showMenu"
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

          <v-list class="tw-py-1" dense>
            <v-list-item @click="copyLink">
              <v-list-item-content>
                <v-list-item-title>Copy link</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-dialog
              v-if="!isGroup"
              v-model="duplicateDialog"
              width="400"
              persistent
            >
              <template v-slot:activator="{ on, attrs }">
                <v-list-item id="duplicate-event-btn" v-bind="attrs" v-on="on">
                  <v-list-item-content>
                    <v-list-item-title>Duplicate</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
              <v-card>
                <v-card-title>Duplicate {{ typeText }}</v-card-title>
                <v-card-text>
                  <v-text-field
                    v-model="duplicateDialogOptions.name"
                    required
                    placeholder="Name your event..."
                    :disabled="duplicateDialogOptions.loading"
                    hide-details
                    solo
                  />
                  <v-checkbox
                    v-model="duplicateDialogOptions.copyAvailability"
                    label="Copy responses"
                    :disabled="duplicateDialogOptions.loading"
                    hide-details
                    class="tw-mt-2"
                  />
                </v-card-text>
                <v-card-actions>
                  <v-spacer />
                  <v-btn
                    text
                    @click="duplicateDialog = false"
                    :disabled="duplicateDialogOptions.loading"
                    >Cancel</v-btn
                  >
                  <v-btn
                    text
                    color="primary"
                    @click="duplicateEvent"
                    :loading="duplicateDialogOptions.loading"
                    >Confirm</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-dialog>
            <v-dialog v-model="removeDialog" width="400" persistent>
              <template v-slot:activator="{ on, attrs }">
                <v-list-item
                  id="delete-event-btn"
                  class="red--text"
                  v-bind="attrs"
                  v-on="on"
                >
                  <v-list-item-content>
                    <v-list-item-title>Delete {{ typeText }}</v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </template>
              <v-card>
                <v-card-title>Are you sure?</v-card-title>
                <v-card-text
                  >Are you sure you want to delete this
                  {{ typeText }}?</v-card-text
                >
                <v-card-actions>
                  <v-spacer />
                  <v-btn text @click="removeDialog = false">Cancel</v-btn>
                  <v-btn text color="error" @click="removeEvent"
                    >I'm sure</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-list>
        </v-menu>
        <v-icon v-else class="tw-ml-2 tw-mr-1 tw-opacity-75"
          >mdi-chevron-right</v-icon
        >
      </div>
    </v-container>
  </router-link>
</template>

<script>
import { getDateRangeStringForEvent, _delete, isPhone, post } from "@/utils"
import { mapActions, mapState } from "vuex"
import { eventTypes } from "@/constants"

export default {
  name: "EventItem",

  props: {
    event: { type: Object, required: true },
  },

  data: () => ({
    showMenu: false,
    duplicateDialog: false,
    duplicateDialogOptions: {
      name: "",
      copyAvailability: false,
      loading: false,
    },
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
    isGroup() {
      return this.event.type === eventTypes.GROUP
    },
    isSignUp() {
      return this.event.isSignUpForm
    },
    linkTo() {
      if (this.isGroup) {
        return "group"
      } else if (this.isSignUp) {
        return "signUp"
      }

      return "event"
    },
    identifier() {
      if (this.isGroup) {
        return "groupId"
      } else if (this.isSignUp) {
        return "signUpId"
      }
      return "eventId"
    },
    typeText() {
      return this.isGroup ? "group" : "event"
    },
    userHasResponded() {
      return this.authUser?._id in this.event.responses
    },
  },

  methods: {
    ...mapActions(["showError", "showInfo", "getEvents"]),
    copyLink() {
      /* Copies event link to clipboard */
      navigator.clipboard.writeText(
        `${window.location.origin}/e/${this.event.shortId ?? this.event._id}`
      )
      this.showInfo("Link copied to clipboard!")
      this.showMenu = false
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    removeEvent() {
      _delete(`/events/${this.event._id}`)
        .then(() => {
          this.getEvents()
          this.$refs.menu.save() // NOTE: Not sure why but without this line, the menu persists to the next event.

          this.$posthog?.capture("Event removed", {
            eventId: this.event._id,
            eventName: this.event.name,
            eventDuration: this.event.duration,
            eventDates: this.event.dates,
            eventNotificationsEnabled: this.event.notificationsEnabled,
            eventType: this.event.type,
          })
        })
        .catch((err) => {
          this.showError(
            "There was a problem removing that event! Please try again later."
          )
        })
    },
    duplicateEvent() {
      this.duplicateDialogOptions.loading = true
      post(`/events/${this.event._id}/duplicate`, {
        eventName: this.duplicateDialogOptions.name,
        copyAvailability: this.duplicateDialogOptions.copyAvailability,
      })
        .then(({ eventId, shortId }) => {
          this.getEvents()
          this.$refs.menu.save() // NOTE: Not sure why but without this line, the menu persists to the next event.

          this.$posthog?.capture("Event duplicated", {
            eventId: eventId,
            eventName: this.duplicateDialogOptions.name,
            eventDuration: this.event.duration,
            eventDates: this.event.dates,
            eventNotificationsEnabled: this.event.notificationsEnabled,
            eventType: this.event.type,
            copyAvailability: this.duplicateDialogOptions.copyAvailability,
          })
        })
        .catch((err) => {
          this.showError(
            "There was a problem duplicating that event! Please try again later."
          )
        })
        .finally(() => {
          this.duplicateDialogOptions.loading = false
        })
    },
  },

  watch: {
    duplicateDialog: {
      immediate: true,
      handler(val) {
        if (val) {
          this.duplicateDialogOptions.name = `Copy of ${this.event.name}`
        }
      },
    },
  },
}
</script>
