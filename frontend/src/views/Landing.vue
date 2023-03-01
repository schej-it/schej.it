<template>
  <div>
    <v-dialog v-model="dialog" :width="400">
      <v-card>
        <v-card-title>
          <span v-if="dialogType === DIALOG_TYPES.SIGN_IN">Sign in</span>
          <span v-else-if="dialogType === DIALOG_TYPES.SIGN_UP">Sign up</span>
        </v-card-title>
        <v-card-text class="tw-flex tw-flex-col tw-items-center">
          <SignInGoogleBtn @click="signInGoogle" dark />
          <div class="tw-text-xs tw-text-center">
            By continuing, you agree to our
            <router-link class="tw-text-blue" :to="{ name: 'privacy-policy' }"
              >privacy policy</router-link
            >
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <div class="tw-h-screen tw-flex tw-flex-col">
      <!-- Header -->
      <div class="tw-mb-12">
        <div
          class="tw-pt-5 tw-px-5 tw-flex tw-justify-between tw-items-center tw-max-w-6xl tw-m-auto"
        >
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
            Finding a time to meet,<br />
            made simple.
          </div>
          <div
            class="tw-mb-4 lg:tw-mb-10 sm:tw-text-xl lg:tw-text-2xl tw-text-center"
          >
            It's like when2meet with <br v-if="isPhone" />
            Google Calendar integration
          </div>
          <v-btn
            class="tw-bg-green tw-rounded-lg"
            dark
            @click="getStarted"
            :large="$vuetify.breakpoint.smAndUp"
            :x-large="$vuetify.breakpoint.mdAndUp"
            >Let's schej it</v-btn
          >
        </div>
      </div>

      <!-- Calendar -->
      <div
        class="tw-flex-1 md:tw-flex md:tw-justify-center tw-relative tw-pb-12"
      >
        <div
          class="tw-absolute tw-bg-green tw-w-full"
          style="top: 10rem; height: calc(100% - 10rem)"
        ></div>
        <div>
          <LandingPageCalendar />
        </div>
      </div>

      <!-- What is schej? -->
      <div
        id="how-it-works"
        class="tw-grid sm:tw-grid-cols-1 md:tw-grid-cols-2 tw-bg-white tw-pt-28 tw-p-8 sm:tw-px-28 xl:tw-px-[20%]"
      >
        <!-- Text -->
        <div class="tw-flex tw-flex-col tw-pr-10 lg:tw-pt-3">
          <div class="tw-text-2xl sm:tw-text-5xl lg:tw-text-6xl tw-font-medium">
            What is schej?
          </div>
          <div
            class="tw-mt-10 tw-grid tw-grid-cols-1 tw-gap-6 sm:tw-text-xl lg:tw-text-2xl"
          >
            <div>
              schej is a group scheduling platform that helps you find a time to
              meet.
            </div>
            <div>
              Users fill out their availability with the help of Google
              Calendar, and a heat map is generated showing when everybody is
              available.
            </div>
            <div>
              It's <span class="tw-font-bold">completely free</span> to use, and
              it looks great on mobile.
            </div>
          </div>
        </div>

        <!-- Video -->
        <div class="tw-h-96 tw-mt-5 md:tw-mt-5">
          <iframe
            class="tw-w-full tw-h-full"
            src="https://www.youtube.com/embed/Wzth9Ov7bkI?fs=0&color=white&rel=0&controls=0"
            title="Demo of Schej"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
            allowfullscreen
          ></iframe>
        </div>
      </div>

      <!-- FAQ -->
      <div
        class="tw-bg-white tw-p-8 tw-pt-32 tw-pb-32 sm:tw-px-16 xl:tw-px-[20%]"
      >
        <!-- Text -->
        <div class="tw-pr-10 lg:tw-pt-3 tw-text-center">
          <div class="tw-text-2xl sm:tw-text-5xl lg:tw-text-6xl tw-font-medium">
            Frequently Asked Questions
          </div>
          <div
            class="tw-mt-16 tw-grid tw-grid-cols-1 tw-gap-3 sm:tw-text-xl lg:tw-text-2xl"
          >
            <FAQ
              question="Why should I log in with my Google account"
              answer="Although it isn't required, it's highly recommended since it allows you to view all your events while filling out your availbaility. A huge help!"
            ></FAQ>
            <FAQ
              question="Will other people be able to see my calendar events?"
              answer="Nope! All other users will be able to see is the availability that you edit and enter into the event."
            ></FAQ>
            <FAQ
              question="Can I still add my availability without logging in?"
              answer="Yes! You can enter as a guest. To edit your guest availibility, just click on your name again."
            ></FAQ>
            <FAQ
              question="How is schej different from lettucemeet or when2meet?"
              :points="['Much better UI (web and mobile)', 'Seemless and working calendar integration', 'No ads :)']"
            ></FAQ>
          </div>
        </div>
      </div>

      <!-- Privacy Policy -->
      <div class="tw-bg-green tw-flex tw-flex-col">
        <div
          class="tw-flex tw-flex-row tw-justify-around tw-m-2 tw-h-28 tw-px-[15%] tw-pt-10 tw-font-medium "
        >
          <a class="tw-text-white" href="#how-it-works" >How it works</a>
          <router-link class="tw-text-white" :to="{ path: 'privacy-policy' }"
            >Privacy Policy</router-link
          >
          <a class="tw-text-white" @click="signIn">Sign in</a>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import LandingPageCalendar from "@/components/LandingPageCalendar.vue"
import { isPhone, signInGoogle } from "@/utils"
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import FAQ from "@/components/FAQ.vue"

export default {
  name: "Landing",

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
    FAQ,
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
