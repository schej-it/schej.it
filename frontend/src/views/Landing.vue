<template>
  <div class="tw-relative">
    <div class="tw-min-h-screen tw-flex tw-flex-col tw-max-w-6xl tw-m-auto">
      <div
        style="top: 40rem; height: calc(100% - 40rem + 2rem)"
        class="tw-bg-green tw-absolute tw-w-screen tw-left-1/2 -tw-translate-x-1/2 tw-pb-4"
      ></div>

      <!-- Header -->
      <div class="tw-mb-16">
        <div class="tw-pt-5 tw-px-5 tw-flex tw-items-center">
          <v-img
            alt="schej Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="@/assets/schej_logo_with_text.png"
            transition="scale-transition"
            :width="logoWidth"
          />

          <v-spacer />

          <v-btn text @click="howItWorksDialog = true">How it works</v-btn>
          <v-btn text @click="signIn">Sign in</v-btn>
        </div>
      </div>

      <div class="tw-flex tw-gap-24">
        <!-- Left side -->
        <div>
          <!-- Hero -->
          <div class="tw-mb-32 tw-mx-4">
            <div
              class="tw-text-2xl sm:tw-text-3xl lg:tw-text-5xl tw-font-medium tw-mb-4"
            >
              <div class="tw-leading-tight">Finding a time to meet,</div>
              <div class="tw-leading-snug">made simple.</div>
            </div>
            <div class="tw-text-lg tw-text-very-dark-gray">
              <b>Automatically</b> fill out your availability with Google
              Calendar—<br />
              it’s like When2Meet with Google Calendar integration!
            </div>
          </div>

          <!-- Calendar -->
          <LandingPageCalendar />
        </div>

        <!-- Right side -->
        <div class="mr-4">
          <NewEvent :dialog="false" />
        </div>
      </div>
    </div>

    <!-- Privacy Policy -->
    <!-- <div class="tw-bg-green tw-flex tw-flex-col">
      <div
        class="tw-flex tw-flex-row tw-justify-around tw-m-2 tw-py-4 tw-font-medium"
      >
        <router-link class="tw-text-white" :to="{ path: 'privacy-policy' }"
          >Privacy Policy</router-link
        >
      </div>
    </div> -->

    <!-- How it works -->
    <v-dialog
      v-model="howItWorksDialog"
      id="how-it-works"
      :width="400"
      class="tw-m-0"
    >
      <v-card>
        <v-card-title>How it works</v-card-title>
        <v-card-text class="tw-text-black tw-flex tw-flex-col tw-gap-2 tw-mt-3">
          <div
            v-for="(step, i) in howItWorksSteps"
            :key="i"
            class="tw-flex tw-items-center tw-gap-3"
          >
            <NumberBullet>{{ i + 1 }}</NumberBullet>
            <div v-html="step"></div>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- Sign in dialog -->
    <v-dialog v-model="signInDialog" :width="400">
      <v-card>
        <v-card-title>Sign in</v-card-title>
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
  </div>
</template>

<script>
import LandingPageCalendar from "@/components/LandingPageCalendar.vue"
import { isPhone, signInGoogle } from "@/utils"
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import FAQ from "@/components/FAQ.vue"
import Header from "@/components/Header.vue"
import NumberBullet from "@/components/NumberBullet.vue"
import NewEvent from "@/components/NewEvent.vue"
import { confetti } from "tsparticles-confetti"

export default {
  name: "Landing",

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
    FAQ,
    Header,
    NumberBullet,
    NewEvent,
  },

  data: () => ({
    howItWorksDialog: false,
    howItWorksSteps: [
      "Create a schej event",
      "Automatically fill out your availability with Google Calendar",
      "Share the schej and find the best time to meet!",
    ],
    signInDialog: false,
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
      this.signInDialog = true
    },
  },
}
</script>
