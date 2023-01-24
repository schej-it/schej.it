<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12 tw-mt-5">

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />
    <div class="tw-p-4">
      <EventType v-for="eventType, t in events" :key="t" :eventType="eventType" class="tw-mb-6"></EventType>
    </div>

    <!-- FAB -->
    <BottomFab @click="dialog = true">
      <v-icon>mdi-plus</v-icon>
    </BottomFab>
  </div>
</template>

<script>
import NewEventDialog from '@/components/NewEventDialog.vue'
import EventType from '@/components/EventType.vue'
import BottomFab from '@/components/BottomFab.vue'
import { mapState, mapActions } from 'vuex'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
    EventType,
    BottomFab,
  },

  data: () => ({
    dialog: false,
  }),

  computed: {
    ...mapState([ 'createdEvents', 'joinedEvents' ]),
    events() {
      return [
          {
            header: 'Events I created',
            events: this.createdEvents.reverse(),
          },
          {
            header: 'Events I joined',
            events: this.joinedEvents.reverse(),
          },
        ] 
    }
  },

  methods: {
    ...mapActions( ['getEvents'] )
  },

  created() {
    this.getEvents()
  },
}
</script>
