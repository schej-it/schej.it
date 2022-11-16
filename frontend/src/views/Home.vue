<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12">

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />

    <div class="tw-p-4">
      <div v-for="eventType, t in events" :key="t">
        <div class="tw-text-2xl tw-font-bold">{{ eventType.header }}</div>
        
        <div 
          v-if="eventType.events.length === 0"
          class="tw-my-3"
        >
          No events yet!
        </div>
        <div v-else class="tw-grid tw-grid-cols-1 sm:tw-grid-cols-2 md:tw-grid-cols-3 tw-gap-2 tw-my-3">
          <EventItem  
            class="tw-cursor-pointer"
            v-for="event, i in eventType.events" 
            :key="i"
            :event="event" 
            @click="goToEvent(event._id)"
          />
        </div>
      </div>
    </div>

    <!-- FAB -->
    <BottomFab @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>
  </div>
</template>

<script>
import NewEventDialog from '@/components/NewEventDialog'
import EventItem from '@/components/EventItem'
import BottomFab from '@/components/BottomFab.vue'
import { get } from '@/utils'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
    EventItem,
    BottomFab,
  },

  data: () => ({
    dialog: false,
    events: null,
  }),

  methods: {
    goToEvent(eventId) {
      this.$router.push({ name: 'event', params: { eventId } })
    }
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
