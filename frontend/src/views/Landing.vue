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

    <div class="tw-h-screen tw-flex tw-flex-col">

      <!-- Header -->
      <div class="tw-mb-12">
        <div class="tw-pt-5 tw-px-5 tw-flex tw-justify-between tw-items-center tw-max-w-6xl tw-m-auto">
            <v-img
              alt="schej Logo"
              class="shrink tw-cursor-pointer"
              contain
              src="@/assets/schej_logo_with_text.png"
              transition="scale-transition"
              :width="logoWidth"
            />

          <v-btn text @click="signIn">Sign in</v-btn>
        </div>
        <div class="tw-flex tw-flex-col tw-items-center">
          <div
            class="tw-mt-16 tw-mb-4 lg:tw-mb-10 tw-text-2xl sm:tw-text-5xl lg:tw-text-6xl tw-font-medium tw-text-center"
          >
            Finding a time to meet,<br> made simple.
          </div>
          <div
            class="tw-mb-4 lg:tw-mb-10 sm:tw-text-xl lg:tw-text-2xl tw-text-center"
          >
            It's like when2meet with <br v-if="isPhone"> Google Calendar integration
          </div>
          <v-btn class="tw-bg-green tw-rounded-lg" dark @click="getStarted" :large="$vuetify.breakpoint.smAndUp" :x-large="$vuetify.breakpoint.mdAndUp"
            >Let's schej it</v-btn
          >
        </div>
      </div>

      <!-- Calendar -->
      <div class="tw-flex-1 md:tw-flex md:tw-justify-center tw-relative tw-pb-12">
        <div 
          class="tw-absolute tw-bg-green tw-w-full"
          style="top: 10rem; height: calc(100% - 10rem);"
        >
        </div>
        <div>
          <LandingPageCalendar />
        </div>
      </div>

      <!-- Privacy Policy -->
      <div class="tw-bg-green tw-flex tw-flex-col">
        <div class="tw-text-center tw-text-white tw-m-2">	
          <router-link class="tw-text-white tw-font-medium" :to="{ path: 'privacy-policy' }">Privacy Policy</router-link>
        </div>
      </div>

    </div>
  </div>
</template>

<script>
import LandingPageCalendar from '@/components/LandingPageCalendar.vue'
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
      return 90
    },
  },

  methods: {
    signInGoogle() {
      signInGoogle(null, true)
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
