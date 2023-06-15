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

    <div class="tw-min-h-screen tw-flex tw-flex-col">
      <!-- Header -->
      <div class="tw-mb-12">
        <div
          class="tw-pt-5 tw-px-5 tw-flex tw-items-center tw-max-w-6xl tw-m-auto"
        >
          <v-img
            alt="schej Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="@/assets/schej_logo_with_text.png"
            transition="scale-transition"
            :width="logoWidth"
          />

          <v-spacer />

          <v-btn text href="landing#how-it-works">How it works</v-btn>
          <v-btn text @click="signIn">Sign in</v-btn>
        </div>

        <div class="tw-flex tw-flex-col tw-items-center tw-mt-8">
          <NewEvent :dialog="false"/>
        </div>
      </div>
    </div>

    <!-- Privacy Policy -->
    <div class="tw-bg-green tw-flex tw-flex-col">
      <div
        class="tw-flex tw-flex-row tw-justify-around tw-m-2 tw-py-4 tw-font-medium"
      >
        <a class="tw-text-white" href="#how-it-works">How it works</a>
        <router-link class="tw-text-white" :to="{ path: 'privacy-policy' }"
          >Privacy Policy</router-link
        >
        <a class="tw-text-white" @click="signIn">Sign in</a>
      </div>
    </div>
  </div>
</template>

<script>
import LandingPageCalendar from "@/components/LandingPageCalendar.vue"
import { isPhone, signInGoogle } from "@/utils"
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import FAQ from "@/components/FAQ.vue"
import Header from "@/components/Header.vue"
import NewEvent from "@/components/NewEvent.vue"
import NumberBullet from "@/components/NumberBullet.vue"
import { confetti } from "tsparticles-confetti"

export default {
  name: "CreateEvent",

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
    FAQ,
    Header,
    NumberBullet,
    NewEvent,
  },

  data: () => ({
    dialog: false,
    somedialog: true,
    dialogType: 0,
    DIALOG_TYPES: {
      SIGN_IN: 0,
      SIGN_UP: 1,
    },
    howItWorksSteps: [
      "Create a schej",
      "",
      "Share the schej and find the best time to meet!",
    ],
    faqs: [
      {
        question: "Do I need to sign in with Google in order to use schej?",
        answer:
          "Signing in with Google is required to create events, but anybody can add their availability once an event is created, whether or not they are signed in with Google!",
      },
      {
        question: "How is schej different from lettucemeet or when2meet?",
        points: [
          "Much better UI (web and mobile)",
          "Seemless and working calendar integration",
          "No ads :)",
        ],
      },
      {
        question: "Is Google Calendar access required in order to use schej?",
        answer:
          "Nope! You can manually input your availability, but we highly recommend allowing Google Calendar access in order to view your calendar events while doing so.",
      },
      {
        question: "Will other people be able to see my calendar events?",
        answer:
          "Nope! All other users will be able to see is the availability that you enter for an event.",
      },
      {
        question: "How do I edit my availability?",
        answer:
          'If you are signed in, simply click the "Edit availability" button. If you entered your availability as a guest, click on your name first and then "Edit availability".',
      },
    ],
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
    confetti() {
      confetti({
        spread: 360,
        ticks: 1000,
        // count: 100,
        shapes: ["image"],
        scalar: 5,
        shapeOptions: {
          image: [
            {
              src: require("@/assets/schej_logo.png"),
              width: 32,
              height: 32,
            },
          ],
        },
      })
    },
    signInGoogle() {
      signInGoogle({ state: null, selectAccount: true })
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
