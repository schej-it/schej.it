<template>
  <div>

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />

    <v-container class="pa-5">
      <div class="headline font-weight-bold tw-mb-3">My events</div>

      <EventItem 
        v-for="event, i in events" 
        :key="i"
        :event="event" 
        @click="goToEvent(event._id)"
        class="tw-mb-2"
      />

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
import { get } from '@/utils'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
    EventItem
  },

  data: () => ({
    dialog: false,
    events: [],
  }),

  methods: {
    goToEvent(eventId) {
      console.log('WHAT')
      this.$router.push({ name: 'event', params: { eventId } })
    }
  },

  created() {
    get('/user/events')
      .then(data => {
        console.log(data)
        this.events = data 
      }).catch(err => {
        console.error(err)
      })
  },
}
</script>
