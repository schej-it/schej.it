<template>
  <div class="tw-flex tw-flex-col tw-h-full">
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
        ><v-icon class="tw-text-white tw-mr-2 tw-text-base">mdi-account-multiple</v-icon>{{ event.respondents }} respondents</div>
      </div>
    </div>
    <div class="tw-flex tw-flex-col tw-items-center tw-justify-center tw-py-12">
      <v-btn 
        class="tw-bg-blue tw-mb-2" 
        dark
        @click="signIn"
      >Sign in with Google</v-btn>
      <div 
        class="tw-text-xs tw-mx-10 tw-text-black tw-text-center"
      >Schej.it automatically inputs your availability <br> using your google calendar</div>
    </div>
  </div>
</template>

<script>
import { getDateString } from '@/utils'
import { mapState } from 'vuex'

export default {
  name: 'Join',

  computed: {
    ...mapState([ 'events' ]),
    dateString() {
      return getDateString(this.event.startDate) + ' - ' + getDateString(this.event.endDate)
    },
    event() {
      return this.events[this.eventId]
    },
    eventId() {
      return this.$route.params.eventId
    },
  },
}
</script>