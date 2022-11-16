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
      <div class="tw-pt-5 tw-px-5 tw-flex tw-justify-between tw-items-center tw-max-w-6xl tw-m-auto">
        <div class="tw-bg-white tw-px-4 tw-py-1 tw-rounded-full">
        <v-img
          alt="schej Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/schej_logo_with_text.svg"
          transition="scale-transition"
          :width="logoWidth"
        />
      </div>

        <v-btn dark text @click="signIn">Sign in</v-btn>
      </div>
      <div class="tw-flex tw-flex-col tw-items-center">
        <div
          class="tw-mt-28 tw-mb-4 lg:tw-mb-10 tw-text-2xl sm:tw-text-5xl lg:tw-text-7xl tw-font-medium tw-text-center tw-text-white"
        >
          Scheduling made simple.
        </div>
        <v-btn rounded class="tw-bg-blue" dark @click="getStarted" :large="$vuetify.breakpoint.smAndUp" :x-large="$vuetify.breakpoint.mdAndUp"
          >Create an event</v-btn
        >
      </div>
    </div>

    <div
      class="-tw-translate-y-1/2 -tw-mb-72 md:tw-grid md:tw-place-content-center"
    >
      <LandingPageCalendar />
    </div>
    
    <div class="tw-py-20 sm:tw-pb-36">
      <div class="tw-flex tw-flex-col tw-items-center tw-mt-10 tw-px-4">
        <p class="sm:tw-text-3xl tw-text-center tw-text-2xl tw-mb-6 tw-leading-10">
          ✍️ Mark your <span class="tw-underline tw-decoration-green tw-decoration-4">availability</span> alongside all your 📅 Google Calendar events.
          <br>
          <!-- It’s like <span class="tw-underline">when2meet</span> but with Google Calendar integration. -->
        </p>
        <v-btn rounded class="tw-bg-blue" dark @click="getStarted" :large="$vuetify.breakpoint.smAndUp" :x-large="$vuetify.breakpoint.mdAndUp"
          >Get started</v-btn
        >
      </div>
    </div>


    <div class="tw-bg-green tw-flex tw-flex-col">
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
    logoWidth() {
      return this.isPhone ? 80 : 120
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
