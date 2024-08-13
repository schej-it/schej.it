<template>
  <div class="tw-mt-1 tw-max-w-full sm:tw-mt-2 sm:tw-max-w-[calc(100%-236px)]">
    <v-dialog v-model="showFullDescription" width="500" content-class="tw-m-0">
      <v-card>
        <v-card-title class="tw-flex tw-items-center tw-justify-between">
          Description
          <v-btn icon @click="showFullDescription = false" class="-tw-mr-2">
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>{{ event.description }}</v-card-text>
      </v-card>
    </v-dialog>

    <div
      v-if="event.description && !isEditing"
      class="tw-inline-flex tw-max-w-full tw-cursor-pointer tw-items-center tw-gap-2 tw-rounded tw-border tw-border-light-gray-stroke tw-bg-light-gray tw-p-2 tw-text-xs tw-font-normal tw-text-very-dark-gray hover:tw-bg-[#eeeeee] sm:tw-text-sm"
      @click="showFullDescription = true"
    >
      <div class="tw-grow tw-truncate">
        {{ event.description }}
      </div>
      <v-btn
        v-if="canEdit"
        key="edit-description-btn"
        class="-tw-my-1"
        icon
        small
        @click="isEditing = true"
      >
        <v-icon small>mdi-pencil</v-icon>
      </v-btn>
    </div>
    <template v-else-if="canEdit">
      <v-btn
        v-if="!isEditing"
        text
        class="-tw-ml-2 tw-mt-0 tw-w-min tw-px-2 tw-text-dark-gray"
        @click="isEditing = true"
      >
        + Add description
      </v-btn>
      <div v-else class="tw-flex tw-w-full tw-items-center tw-gap-2">
        <v-textarea
          v-model="newDescription"
          placeholder="Enter a description..."
          class="-tw-mt-2 tw-flex-grow tw-text-xs sm:tw-text-sm"
          autofocus
          :rows="1"
          auto-grow
          hide-details
        ></v-textarea>
        <v-btn
          icon
          :small="isPhone"
          @click="
            newDescription = event.description
            isEditing = false
          "
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <v-btn icon :small="isPhone" color="primary" @click="saveDescription"
          ><v-icon>mdi-check</v-icon></v-btn
        >
      </div>
    </template>
  </div>
</template>

<script>
import { mapActions } from "vuex"
import { isPhone, put } from "@/utils"

export default {
  name: "EventDescription",

  props: {
    event: {
      type: Object,
      required: true,
    },
    canEdit: {
      type: Boolean,
      required: true,
    },
  },

  data() {
    return {
      isEditing: false,
      newDescription: this.event.description ?? "",
      showFullDescription: false,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    ...mapActions(["showError"]),
    saveDescription() {
      const oldEvent = { ...this.event }
      const newEvent = {
        ...this.event,
        description: this.newDescription,
      }
      this.$emit("update:event", newEvent)
      this.isEditing = false
      put(`/events/${this.event._id}`, newEvent).catch((err) => {
        console.error(err)
        this.showError("Failed to save description! Please try again later.")
        this.$emit("update:event", {
          ...oldEvent,
        })
      })
    },
  },
}
</script>
