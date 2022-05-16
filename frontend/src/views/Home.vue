<template>
  <div>

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />

    <v-container class="pa-5">
      <div class="headline font-weight-bold tw-mb-3">My events</div>

      <EventItem :event="exampleEvent"></EventItem>

    </v-container>

    <!-- FAB -->
    <v-scale-transition appear origin="center">
      <v-btn 
        fab
        absolute
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-4"
        @click="dialog = true"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
import NewEventDialog from '@/components/NewEventDialog'
import EventItem from '@/components/EventItem'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
    EventItem
  },

  data: () => ({
    dialog: false,
    exampleEvent: {
      name: "Meeting #1",
      startDate: new Date(),
      endDate: new Date(),
      responses: ["arjun", "tony"]
    }
  }),
}
</script>
