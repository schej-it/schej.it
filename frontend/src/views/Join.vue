<template>
  <div v-if="event" class="tw-flex tw-flex-col tw-h-full">
    <div class="tw-flex-1 tw-bg-green tw-text-white tw-flex tw-justify-center tw-items-center">
      <div class="-tw-mt-8 tw-flex tw-flex-col tw-items-center tw-space-y-3">
        <div 
          class="tw-font-bold tw-text-5xl"
        >{{ event.name }}</div>
        <div
          class="tw-font-light tw-text-2xl"
        >{{ dateString }}</div>
        <div
          class="tw-font-light tw-text-sm tw-flex tw-items-center"
        ><v-icon class="tw-text-white tw-mr-2 tw-text-base">mdi-account-multiple</v-icon>{{ Object.keys(event.responses).length }} respondents</div>
      </div>
    </div>
    <div class="tw-flex tw-flex-col tw-items-center tw-justify-center tw-py-12">
      <v-btn
        v-if="authUser"
        class="tw-bg-blue"
        dark
        @click="join"
      >Join event</v-btn>
      <template v-else>
        <v-btn 
          class="tw-bg-blue tw-mb-2" 
          dark
          @click="signIn"
        >Sign in with Google</v-btn>
        <div 
          class="tw-text-xs tw-mx-10 tw-text-black tw-text-center"
        >Schej.it automatically inputs your availability <br> using your google calendar</div>
      </template>
    </div>
  </div>
</template>

<script>
import { getDateRangeString, signInGoogle, get } from '@/utils'
import { mapState } from 'vuex'

export default {
  name: 'Join',

  props: {
    eventId: { type: String, required: true },
  },

  data: () => ({
    event: null,
  }),

  computed: {
    ...mapState([ 'authUser', 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
  },

  methods: {
    join() {
      this.$router.replace({ name: 'event', params: { eventId: this.eventId } })
    },
    signIn() {
      signInGoogle({ type: 'join', eventId: this.eventId })
    },
  },

  async created() {
    // Get event details
    this.event = await get(`/events/${this.eventId}`)
    this.event.startDate = new Date(this.event.startDate)
    this.event.endDate = new Date(this.event.endDate)
  }
}
</script>