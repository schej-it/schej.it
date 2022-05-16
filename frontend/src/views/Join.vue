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
        ><v-icon class="tw-text-white tw-mr-2 tw-text-base">mdi-account-multiple</v-icon>{{ event.responses.length }} respondents</div>
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
import { getDateRangeString, signInGoogle } from '@/utils'
import { mapState } from 'vuex'

export default {
  name: 'Join',

  props: {
    eventId: { type: String, required: true },
  },

  computed: {
    ...mapState([ 'events' ]),
    dateString() {
      return getDateRangeString(this.event.startDate, this.event.endDate)
    },
    event() {
      return this.events[this.eventId]
    },
  },

  methods: {
    signIn() {
      signInGoogle(JSON.stringify({ type: 'join', eventId: this.eventId }))
    },
  },

  mounted() {
    google.accounts.id.initialize({
      client_id: '523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com',
      callback: this.handleCredentialResponse
    });
    google.accounts.id.renderButton(
      document.getElementById('sign-in-google'),
      { theme: 'filled_blue', size: 'large', text: 'continue_with' } 
    )
  },
}
</script>