<template>
  <div class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]">
    <div class="tw-mb-3 tw-flex tw-items-center tw-justify-between">
      <div class="tw-flex tw-flex-col">
        <div
          class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl"
        >
          Dashboard
        </div>
        <div
          v-if="!isPremiumUser"
          class="tw-flex tw-items-baseline tw-gap-2 tw-text-sm tw-font-normal tw-text-very-dark-gray"
        >
          <div>
            {{ authUser?.numEventsCreated }} / {{ numFreeEvents }} free events
            created
          </div>
          <div
            class="tw-cursor-pointer tw-select-none tw-text-xs tw-font-medium tw-text-green tw-underline"
            @click="openUpgradeDialog"
          >
            Upgrade
          </div>
        </div>
      </div>
      <v-btn
        text
        @click="openCreateFolderDialog"
        class="tw-hidden tw-text-very-dark-gray sm:tw-block"
      >
        <v-icon class="tw-mr-2 tw-text-lg">mdi-folder-plus</v-icon>
        New folder
      </v-btn>
    </div>

    <div>
      <div v-for="folder in folders" :key="folder._id" class="tw-group tw-mb-2">
        <div class="tw-flex tw-items-center">
          <v-btn icon small @click="toggleFolder(folder._id)">
            <v-icon>{{
              folderOpenState[folder._id] ? "mdi-menu-down" : "mdi-menu-right"
            }}</v-icon>
          </v-btn>
          <v-chip
            :color="folder.color || '#D3D3D3'"
            small
            class="tw-mr-2 tw-cursor-pointer tw-rounded tw-border tw-border-light-gray-stroke tw-px-2 tw-text-sm tw-font-medium"
            @click="openEditFolderDialog(folder)"
          >
            {{ folder.name }}
          </v-chip>
          <div
            class="tw-invisible tw-flex tw-items-center group-hover:tw-visible"
          >
            <v-menu offset-y>
              <template v-slot:activator="{ on, attrs }">
                <v-btn icon small v-bind="attrs" v-on="on" @click.stop.prevent>
                  <v-icon small>mdi-dots-horizontal</v-icon>
                </v-btn>
              </template>
              <v-list dense class="tw-py-1">
                <v-list-item @click.stop.prevent="openEditFolderDialog(folder)">
                  <v-list-item-title>Edit</v-list-item-title>
                </v-list-item>
                <v-list-item @click.stop.prevent="openDeleteDialog(folder)">
                  <v-list-item-title class="tw-text-red"
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
        </div>
        <div v-show="folderOpenState[folder._id]">
          <draggable
            :list="[
              ...eventsByFolder[folder._id].groups,
              ...eventsByFolder[folder._id].events,
            ]"
            group="events"
            @end="onEnd"
            :data-folder-id="folder._id"
            draggable=".item"
            class="tw-relative tw-grid tw-min-h-[52px] tw-grid-cols-1 tw-gap-4 tw-py-4 sm:tw-grid-cols-2"
          >
            <template v-slot:header>
              <div
                v-if="
                  eventsByFolder[folder._id].groups.length === 0 &&
                  eventsByFolder[folder._id].events.length === 0
                "
                class="tw-absolute tw-left-0 tw-ml-8 tw-py-4 tw-text-sm tw-text-very-dark-gray"
              >
                No events in this folder
              </div>
            </template>
            <EventItem
              v-for="event in eventsByFolder[folder._id].groups"
              :key="event._id"
              :id="event._id"
              :event="event"
              :folder-id="folder._id"
              class="item"
            />
            <div class="tw-col-span-full"></div>
            <EventItem
              v-for="event in eventsByFolder[folder._id].events"
              :key="event._id"
              :id="event._id"
              :event="event"
              :folder-id="folder._id"
              class="item"
            />
          </draggable>
        </div>
      </div>

      <div v-if="allEvents.length > 0">
        <div class="tw-flex tw-items-center">
          <v-btn icon small @click="toggleFolder('no-folder')">
            <v-icon>{{
              folderOpenState["no-folder"] ? "mdi-menu-down" : "mdi-menu-right"
            }}</v-icon>
          </v-btn>
          <span class="tw-text-sm tw-font-medium">No folder</span>
        </div>
        <div v-show="folderOpenState['no-folder']">
          <draggable
            :list="[
              ...eventsWithoutFolder.groups,
              ...eventsWithoutFolder.events,
            ]"
            group="events"
            @end="onEnd"
            data-folder-id="null"
            draggable=".item"
            class="tw-relative tw-grid tw-min-h-[52px] tw-grid-cols-1 tw-gap-4 tw-py-4 sm:tw-grid-cols-2"
          >
            <template v-slot:header>
              <div
                v-if="
                  eventsWithoutFolder.groups.length === 0 &&
                  eventsWithoutFolder.events.length === 0
                "
                class="tw-absolute tw-left-0 tw-ml-7 tw-py-4 tw-text-sm tw-text-very-dark-gray"
              >
                No events
              </div>
            </template>
            <EventItem
              v-for="event in eventsWithoutFolder.groups"
              :key="event._id"
              :id="event._id"
              :event="event"
              :folder-id="null"
              class="item"
            />
            <div class="tw-col-span-full"></div>
            <EventItem
              v-for="event in eventsWithoutFolder.events"
              :key="event._id"
              :id="event._id"
              :event="event"
              :folder-id="null"
              class="item"
            />
          </draggable>
        </div>
      </div>
      <div v-else>
        <div class="tw-py-4 tw-text-sm tw-text-very-dark-gray">
          No events yet! Create one to get started.
        </div>
      </div>
    </div>
    <v-dialog v-model="deleteDialog" max-width="400">
      <v-card>
        <v-card-title>Delete "{{ folderToDelete.name }}"?</v-card-title>
        <v-card-text
          >Are you sure you want to delete this folder? All events you own in
          this folder will be deleted as well.</v-card-text
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
        <v-card-title>{{ folderDialogTitle }}</v-card-title>
        <v-card-text>
          <v-text-field
            v-model="newFolderName"
            label="Folder name"
            placeholder="Untitled folder"
            autofocus
            @keydown.enter="confirmFolderDialog"
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
          <v-btn text @click="closeFolderDialog">Cancel</v-btn>
          <v-btn color="primary" text @click="confirmFolderDialog">{{
            folderDialogConfirmText
          }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { mapState, mapActions, mapGetters } from "vuex"
import draggable from "vuedraggable"
import {
  eventTypes,
  folderColors,
  numFreeEvents,
  upgradeDialogTypes,
} from "@/constants"
import EventItem from "@/components/EventItem.vue"
import ObjectID from "bson-objectid"

export default {
  name: "Dashboard",
  components: {
    EventItem,
    draggable,
  },
  data() {
    return {
      deleteDialog: false,
      folderToDelete: {},
      createFolderDialog: false,
      newFolderName: "",
      newFolderColor: folderColors[3],
      isEditingFolder: false,
      folderToEdit: null,
      folderOpenState: {
        "no-folder": true,
      },
    }
  },
  computed: {
    ...mapGetters(["isPremiumUser"]),
    ...mapState(["authUser", "events", "groupsEnabled", "folders"]),
    numFreeEvents() {
      return numFreeEvents
    },
    folderColors() {
      return folderColors
    },
    allEvents() {
      return this.events
    },
    allEventsMap() {
      return this.allEvents.reduce((acc, event) => {
        acc[event._id] = event
        return acc
      }, {})
    },
    eventsByFolder() {
      const eventsByFolder = {}
      const allEventIds = new Set(this.allEvents.map((e) => e._id))
      this.folders.forEach((folder) => {
        eventsByFolder[folder._id] = { groups: [], events: [] }
        for (const eventId of folder.eventIds) {
          const event = this.allEventsMap[eventId]
          if (event) {
            if (event.type === eventTypes.GROUP) {
              eventsByFolder[folder._id].groups.push(event)
            } else {
              eventsByFolder[folder._id].events.push(event)
            }
            allEventIds.delete(eventId)
          }
        }
        eventsByFolder[folder._id].groups.sort(this.sortEvents)
        eventsByFolder[folder._id].events.sort(this.sortEvents)
      })
      eventsByFolder["no-folder"] = { groups: [], events: [] }
      for (const eventId of allEventIds) {
        const event = this.allEventsMap[eventId]
        if (event) {
          if (event.type === eventTypes.GROUP) {
            eventsByFolder["no-folder"].groups.push(event)
          } else {
            eventsByFolder["no-folder"].events.push(event)
          }
        }
      }
      eventsByFolder["no-folder"].groups.sort(this.sortEvents)
      eventsByFolder["no-folder"].events.sort(this.sortEvents)
      return eventsByFolder
    },
    eventsWithoutFolder() {
      return this.eventsByFolder["no-folder"]
    },
    folderDialogTitle() {
      return this.isEditingFolder ? "Edit folder" : "New folder"
    },
    folderDialogConfirmText() {
      return this.isEditingFolder ? "Save" : "Create"
    },
  },

  methods: {
    ...mapActions([
      "createFolder",
      "showUpgradeDialog",
      "deleteFolder",
      "setEventFolder",
      "updateFolder",
      "createNew",
    ]),
    sortEvents(a, b) {
      if (ObjectID.isValid(a._id) && ObjectID.isValid(b._id)) {
        return ObjectID(b._id).getTimestamp() - ObjectID(a._id).getTimestamp()
      }
      return 0
    },
    onEnd(evt) {
      const eventId = evt.item.id
      let newFolderId = evt.to.dataset.folderId
      if (newFolderId === "null" || newFolderId === undefined) {
        newFolderId = null
      }

      const fromFolderId = evt.from.dataset.folderId

      // if moving within the same folder, do nothing.
      if (fromFolderId === newFolderId) {
        // Here you might want to handle re-ordering within the same folder
        // For now, we do nothing.
        return
      }

      const event = this.allEvents.find((e) => e._id === eventId)

      if (event) {
        this.setEventFolder({
          eventId: event._id,
          folderId: newFolderId,
        })
      }
    },
    confirmFolderDialog() {
      if (!this.newFolderName.trim()) {
        this.closeFolderDialog()
        return
      }
      if (this.isEditingFolder) {
        this.updateFolder({
          folderId: this.folderToEdit._id,
          name: this.newFolderName.trim(),
          color: this.newFolderColor,
        })
      } else {
        this.createFolder({
          name: this.newFolderName.trim(),
          color: this.newFolderColor,
        })
      }
      this.closeFolderDialog()
    },
    closeFolderDialog() {
      this.createFolderDialog = false
      this.isEditingFolder = false
      this.folderToEdit = null
      this.newFolderName = ""
      this.newFolderColor = folderColors[3]
    },
    openCreateFolderDialog() {
      this.isEditingFolder = false
      this.folderToEdit = null
      this.newFolderName = ""
      this.newFolderColor = folderColors[3]
      this.createFolderDialog = true
    },
    openEditFolderDialog(folder) {
      this.isEditingFolder = true
      this.folderToEdit = folder
      this.newFolderName = folder.name
      this.newFolderColor = folder.color || folderColors[3]
      this.createFolderDialog = true
    },
    toggleFolder(folderId) {
      this.$set(this.folderOpenState, folderId, !this.folderOpenState[folderId])
    },
    createEventInFolder(folderId) {
      this.createNew({
        eventOnly: false,
        folderId: folderId,
      })
    },
    openDeleteDialog(folder) {
      this.folderToDelete = folder
      this.deleteDialog = true
    },
    confirmDelete() {
      this.$store.dispatch("deleteFolder", this.folderToDelete._id)
      this.deleteDialog = false
    },
    openUpgradeDialog() {
      this.showUpgradeDialog({
        type: upgradeDialogTypes.UPGRADE_MANUALLY,
      })
    },
  },
  created() {
    try {
      const storedState = localStorage.getItem("folderOpenState")
      if (storedState) {
        this.folderOpenState = JSON.parse(storedState)
      }
    } catch (e) {
      console.error("Error reading folderOpenState from localStorage", e)
      // If corrupted, remove it
      localStorage.removeItem("folderOpenState")
    }
  },
  watch: {
    folderOpenState: {
      handler(newState) {
        try {
          localStorage.setItem("folderOpenState", JSON.stringify(newState))
        } catch (e) {
          console.error("Error saving folderOpenState to localStorage", e)
        }
      },
      deep: true,
    },
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
