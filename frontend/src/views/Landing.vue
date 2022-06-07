<template>
  <div>
    <v-dialog 
      v-model="dialog"
      :width="400"
    >
      <v-card>
        <v-card-title>
          <span v-if="dialogType === DIALOG_TYPES.SIGN_IN">Sign in</span>
          <span v-else-if="dialogType === DIALOG_TYPES.SIGN_UP">Sign up</span>
        </v-card-title>
        <v-card-text class="tw-flex tw-flex-col tw-items-center">
          <SignInGoogleBtn
            @click="signInGoogle"
            dark
          />
          <div class="tw-text-xs tw-text-center">By continuing, you agree to our <router-link :to="{ name: 'privacy-policy' }">privacy policy</router-link></div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <div class="tw-bg-green tw-h-screen">
      <div class="tw-pt-5 tw-px-5 tw-flex tw-justify-between tw-max-w-6xl tw-m-auto">
        <v-img
          alt="schej.it Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/logo_dark.svg"
          transition="scale-transition"
          width="120"
        />

        <v-btn dark text @click="signIn">Sign in</v-btn>
      </div>
      <div class="tw-flex tw-flex-col tw-items-center">
        <div
          class="tw-mt-28 tw-mb-4 lg:tw-mb-10 tw-text-2xl sm:tw-text-5xl lg:tw-text-7xl tw-font-medium tw-text-center tw-text-white"
        >
          Scheduling made simple.
        </div>
        <v-btn rounded class="tw-bg-blue" dark @click="getStarted" :large="$vuetify.breakpoint.smAndUp" :x-large="$vuetify.breakpoint.mdAndUp"
          >Get started</v-btn
        >
      </div>
    </div>

    <div
      class="-tw-translate-y-1/2 -tw-mb-72 md:tw-grid md:tw-place-content-center"
    >
      <LandingPageCalendar />
    </div>
    
    <div class="tw-pt-20 tw-pb-56">

      <div class="tw-flex tw-flex-col tw-items-center sm:tw-flex-row-reverse sm:tw-justify-center">
        <div class="tw-flex tw-flex tw-w-80">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://download.logo.wine/logo/Google_Calendar/Google_Calendar-Logo.wine.png"
            transition="scroll-y-reverse-transition"
            width="600"
          />
        </div>
        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-10 tw-px-4 sm:tw-px-0">
          <h1 class="tw-font-bold sm:tw-text-2xl tw-text-center sm:tw-text-left tw-text-xl tw-mb-2">Schedule meetings with ease</h1>
          <p class="tw-text-dark-gray">
            Mark your availability alongside all your Google Calendar events.
            It’s like when2meet but with Google Calendar integration.
          </p>
        </div>
      </div>

      <div class="tw-flex tw-flex-col tw-items-center sm:tw-flex-row sm:tw-justify-center tw-mt-20">
        <div class="tw-flex tw-flex tw-w-80 tw-justify-center">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://cdn-icons-png.flaticon.com/512/3721/3721924.png"
            transition="scroll-y-reverse-transition"
            width="180"
          />
        </div>

        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-7 tw-px-4 sm:tw-px-0">
          <h1 class="tw-font-bold sm:tw-text-2xl tw-text-center sm:tw-text-left tw-text-xl tw-mb-2">Send your availability</h1>
          <p class="tw-text-dark-gray">
            Send your friends a screenshot of the times you are available. No
            more hassle of manually typing out time ranges of when you're free.
          </p>
        </div>
      </div>

      <div class="tw-flex tw-flex-col tw-items-center sm:tw-flex-row-reverse sm:tw-justify-center tw-mt-20">
        <div class="tw-flex tw-flex tw-w-80 tw-justify-center">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://icon-library.com/images/friends-icon-transparent/friends-icon-transparent-1.jpg"
            transition="scroll-y-reverse-transition"
            width="200"
          />
        </div>
        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-7 tw-px-4 sm:tw-px-0">
          <h1 class="tw-font-bold sm:tw-text-2xl tw-text-center sm:tw-text-left tw-text-xl tw-mb-2">Peep your friends’ schedules</h1>
          <p class="tw-text-dark-gray">
            Add your friends to see what they’ve been up to. Control who can see
            what.
          </p>
        </div>
      </div>
    </div>

    <div
      class="-tw-translate-y-[41%] -tw-mb-52 tw-grid tw-place-content-center"
    >
      <img
        class="tw-select-none"
        style="user-drag: none;"
        alt="Schej.it Logo"
        contain
        src="@/assets/logo_square_512_512.png"
        width="200"
      >
    </div>

    <div class="tw-bg-green tw-h-80 tw-flex tw-flex-col">
      <div class="tw-flex-1 tw-flex tw-justify-center tw-items-center tw-mt-14">
        <v-btn rounded class="tw-bg-blue" dark @click="getStarted" :large="$vuetify.breakpoint.smAndUp" :x-large="$vuetify.breakpoint.mdAndUp"
          >Get started</v-btn
        >
      </div>
      <div class="tw-text-center tw-text-white tw-m-2">	
        <router-link class="tw-text-white tw-font-medium" :to="{ path: 'privacy-policy' }">Privacy Policy</router-link>
      </div>
    </div>
  </div>
</template>

<script>
import LandingPageCalendar from '@/components/LandingPageCalendar'
import { isPhone, signInGoogle } from '@/utils'
import SignInGoogleBtn from '@/components/SignInGoogleBtn.vue'

export default {
  name: 'Landing',

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
  },

  data: () => ({
    dialog: false,
    dialogType: 0,
    DIALOG_TYPES: {
      SIGN_IN: 0,
      SIGN_UP: 1,
    },
  }),

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    signInGoogle() {
      signInGoogle()
    },
    signIn() {
      this.dialog = true
      this.dialogType = this.DIALOG_TYPES.SIGN_IN
    },
    getStarted() {
      this.dialog = true
      this.dialogType = this.DIALOG_TYPES.SIGN_UP
    },
  },
}
</script>
