<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12 tw-mt-5">

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />

    <div class="tw-p-4">
      <EventType v-for="eventType, t in events" :key="t" :eventType="eventType" class="tw-mb-5"></EventType>
      
    </div>

    <!-- FAB -->
    <BottomFab @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>
  </div>
</template>

<script>
import NewEventDialog from '@/components/NewEventDialog'
import EventType from '@/components/EventType'
import BottomFab from '@/components/BottomFab.vue'
import { get } from '@/utils'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
    EventType,
    BottomFab,
  },

  data: () => ({
    dialog: false,
    events: null,
  }),

  methods: {
    
  },

  created() {
    get('/user/events')
      .then(data => {
        this.events = [
          {
            header: 'Events I created',
            events: data.events.reverse(),
          },
          {
            header: 'Events I joined',
            events: data.joinedEvents.reverse(),
          },
        ] 
      }).catch(err => {
        console.error(err)
      })
  },
}
</script>
