<template>
  <div class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]">
    <div class="tw-mb-3 tw-flex tw-items-center tw-justify-between">
      <div class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl">
        Dashboard
      </div>
      <v-btn
        text
        @click="createFolderDialog = true"
        class="tw-hidden tw-text-very-dark-gray sm:tw-block"
      >
        <v-icon class="tw-mr-2 tw-text-lg">mdi-folder-plus</v-icon>
        New folder
      </v-btn>
    </div>

    <div>
      <div v-for="folder in folders" :key="folder._id" class="tw-mb-2">
        <div class="tw-flex tw-items-center">
          <v-btn icon small @click="toggleFolder(folder._id)">
            <v-icon>{{
              folderOpenState[folder._id] ? "mdi-menu-down" : "mdi-menu-right"
            }}</v-icon>
          </v-btn>
          <v-chip
            :color="folder.color || '#D3D3D3'"
            small
            class="tw-mr-2 tw-rounded tw-border tw-border-light-gray-stroke tw-px-2 tw-text-sm tw-font-medium"
          >
            {{ folder.name }}
          </v-chip>
          <v-menu offset-y>
            <template v-slot:activator="{ on, attrs }">
              <v-btn icon small v-bind="attrs" v-on="on" @click.stop.prevent>
                <v-icon small>mdi-dots-horizontal</v-icon>
              </v-btn>
            </template>
            <v-list>
              <v-list-item @click.stop.prevent="openDeleteDialog(folder)">
                <v-list-item-title class="tw-text-red-500"
                  >Delete</v-list-item-title
                >
              </v-list-item>
            </v-list>
          </v-menu>
          <v-btn
            icon
            small
            @click.stop.prevent="createEventInFolder(folder._id)"
          >
            <v-icon small>mdi-plus</v-icon>
          </v-btn>
        </div>
        <div v-show="folderOpenState[folder._id]">
          <div
            v-if="eventsByFolder[folder._id]?.length > 0"
            class="tw-grid tw-grid-cols-1 tw-gap-4 tw-py-4 sm:tw-grid-cols-2"
          >
            <EventItem
              v-for="event in eventsByFolder[folder._id]"
              :key="event._id"
              :event="event"
            />
          </div>
          <div v-else class="tw-ml-8 tw-py-4 tw-text-sm tw-text-very-dark-gray">
            No events in this folder
          </div>
        </div>
      </div>

      <div>
        <div class="tw-flex tw-items-center">
          <v-btn icon small @click="toggleFolder('no-folder')">
            <v-icon>{{
              folderOpenState["no-folder"] ? "mdi-menu-down" : "mdi-menu-right"
            }}</v-icon>
          </v-btn>
          <span class="tw-text-sm">No folder</span>
        </div>
        <div v-show="folderOpenState['no-folder']">
          <div
            class="tw-grid tw-grid-cols-1 tw-gap-4 tw-py-4 sm:tw-grid-cols-2"
          >
            <EventItem
              v-for="event in eventsWithoutFolder"
              :key="event._id"
              :event="event"
            />
          </div>
        </div>
      </div>
    </div>
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title class="headline"
          >Delete {{ folderToDelete.name }}?</v-card-title
        >
        <v-card-text
          >Are you sure you want to delete this folder? This will not delete the
          events in it.</v-card-text
        >
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="red darken-1" text @click="confirmDelete">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-dialog v-model="createFolderDialog" max-width="400">
      <v-card>
        <v-card-title>New folder</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newFolderName"
            label="Folder name"
            placeholder="Untitled folder"
            autofocus
            @keydown.enter="confirmCreateFolder"
            hide-details
          ></v-text-field>
          <div class="tw-mt-4">
            <span class="tw-text-gray-500 tw-text-sm">Color</span>
            <div class="tw-mt-2 tw-flex tw-gap-x-3">
              <div
                v-for="color in folderColors"
                :key="color"
                class="tw-h-6 tw-w-6 tw-cursor-pointer tw-rounded-full tw-border tw-border-light-gray-stroke"
                :style="{ backgroundColor: color }"
                :class="{
                  'tw-ring-2 tw-ring-gray tw-ring-offset-2':
                    newFolderColor === color,
                }"
                @click="newFolderColor = color"
              ></div>
            </div>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="createFolderDialog = false">Cancel</v-btn>
          <v-btn color="primary" text @click="confirmCreateFolder"
            >Create</v-btn
          >
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex"
import { eventTypes, folderColors } from "@/constants"
import EventItem from "@/components/EventItem.vue"

export default {
  name: "Dashboard",
  components: {
    EventItem,
  },
  data() {
    return {
      deleteDialog: false,
      folderToDelete: {},
      createFolderDialog: false,
      newFolderName: "",
      newFolderColor: folderColors[3],
      folderOpenState: {
        "no-folder": true,
      },
    }
  },
  computed: {
    ...mapState(["createdEvents", "joinedEvents", "groupsEnabled", "folders"]),
    folderColors() {
      return folderColors
    },
    allEvents() {
      return [...this.createdEvents, ...this.joinedEvents]
    },
    eventsByFolder() {
      const eventsByFolder = {}
      this.allEvents.forEach((event) => {
        if (event.folderId) {
          if (!eventsByFolder[event.folderId]) {
            eventsByFolder[event.folderId] = []
          }
          eventsByFolder[event.folderId].push(event)
        }
      })
      return eventsByFolder
    },
    eventsWithoutFolder() {
      return this.allEvents.filter(
        (event) =>
          !event.folderId || !this.folders.find((f) => f._id === event.folderId)
      )
    },
    createdEventsNonGroup() {
      return this.createdEvents.filter((e) => e.type !== eventTypes.GROUP)
    },
    joinedEventsNonGroup() {
      return this.joinedEvents.filter((e) => e.type !== eventTypes.GROUP)
    },
    availabilityGroups() {
      return this.createdEvents
        .filter((e) => e.type === eventTypes.GROUP)
        .concat(this.joinedEvents.filter((e) => e.type === eventTypes.GROUP))
        .sort((e1, e2) => (this.userRespondedToEvent(e1) ? 1 : -1))
    },
  },

  methods: {
    ...mapActions(["createFolder"]),
    confirmCreateFolder() {
      if (this.newFolderName.trim()) {
        this.createFolder({
          name: this.newFolderName.trim(),
          color: this.newFolderColor,
        })
      }
      this.newFolderName = ""
      this.newFolderColor = folderColors[3]
      this.createFolderDialog = false
    },
    toggleFolder(folderId) {
      this.$set(this.folderOpenState, folderId, !this.folderOpenState[folderId])
    },
    createEventInFolder(folderId) {
      console.log(`Create event in folder ${folderId}`)
    },
    openDeleteDialog(folder) {
      this.folderToDelete = folder
      this.deleteDialog = true
    },
    confirmDelete() {
      this.$store.dispatch("deleteFolder", this.folderToDelete._id)
      this.deleteDialog = false
    },
  },
  watch: {
    folders: {
      handler(newFolders) {
        if (newFolders) {
          newFolders.forEach((folder) => {
            if (this.folderOpenState[folder._id] === undefined) {
              this.$set(this.folderOpenState, folder._id, true) // default to open
            }
          })
        }
      },
      immediate: true,
    },
  },
}
</script>

<style>
.v-expansion-panel-header {
  padding: 16px 4px !important;
}
</style>
