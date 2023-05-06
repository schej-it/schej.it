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

          <v-btn text href="#how-it-works">How it works</v-btn>
          <v-btn text @click="signIn">Sign in</v-btn>
        </div>

        <!-- Hero -->
        <div class="tw-flex tw-flex-col tw-items-center">
          <div
            class="tw-mt-16 tw-mb-4 sm:tw-mb-8 lg:tw-mb-10 tw-text-2xl sm:tw-text-5xl lg:tw-text-6xl tw-font-medium tw-text-center"
          >
            <div class="tw-leading-normal">Finding a time to meet,</div>
            <div class="tw-leading-normal">made simple.</div>
          </div>
          <v-btn
            id="lets-schej-it-btn"
            class="tw-bg-green tw-rounded-lg tw-px-6 sm:tw-px-10 lg:tw-px-12"
            dark
            @click="getStarted"
            :large="$vuetify.breakpoint.smAndUp"
            :x-large="$vuetify.breakpoint.mdAndUp"
          >
            Let's schej it
          </v-btn>
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
    </div>

    <!-- How it works -->
    <div
      id="how-it-works"
      class="tw-pt-12 tw-grid tw-place-content-center tw-px-4"
    >
      <div class="tw-flex tw-flex-col tw-gap-4 tw-mx-auto">
        <div
          class="tw-text-center tw-text-2xl sm:tw-text-3xl lg:tw-text-4xl tw-font-medium tw-mb-4"
        >
          How it works
        </div>
        <div
          v-for="(step, i) in howItWorksSteps"
          :key="i"
          class="tw-flex tw-items-center tw-gap-2"
        >
          <NumberBullet>{{ i + 1 }}</NumberBullet>
          <div class="md:tw-text-xl tw-text-base tw-font-medium">
            <div v-if="i == 1">
              <span
                class="tw-cursor-pointer tw-underline tw-decoration-pale-green hover:tw-decoration-green tw-decoration-4"
                style="text-underline-position: under"
                @click="confetti"
                >Automatically</span
              >
              fill out your availability with Google Calendar,
            </div>
            <div v-else v-html="step"></div>
          </div>
        </div>
      </div>
      <div
        class="tw-mt-10 tw-mb-6 md:tw-mt-20 md:tw-mb-12 tw-text-3xl md:tw-text-6xl tw-font-medium tw-text-center"
      >
        It's that simple.
      </div>
      <v-img
        alt="schej character"
        src="@/assets/schej_character.png"
        :max-height="isPhone ? 200 : 300"
        transition="scale-transition"
        contain
        class="-tw-mb-12"
      />
    </div>

    <!-- Video -->
    <div
      class="tw-flex tw-bg-green tw-px-4 tw-pt-24 tw-pb-12 md:tw-pb-12 tw-justify-center"
    >
      <div class="md:tw-h-96 sm:tw-h-80 tw-h-64 tw-max-w-3xl tw-flex-1">
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
    <div class="tw-flex tw-justify-center tw-pt-12">
      <div class="tw-flex-1 tw-mx-4 sm:tw-mx-16 tw-mb-12 tw-max-w-3xl">
        <div id="faq-section" class="lg:tw-pt-3 tw-text-center">
          <Header> Frequently Asked Questions </Header>
          <div
            class="tw-grid tw-grid-cols-1 tw-gap-3 sm:tw-text-xl lg:tw-text-2xl"
          >
            <FAQ v-for="faq in faqs" :key="faq.question" v-bind="faq" />
          </div>
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
import NumberBullet from "@/components/NumberBullet.vue"

export default {
  name: "Landing",

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
    FAQ,
    Header,
    NumberBullet,
  },

  data: () => ({
    dialog: false,
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
      console.log("confetti!!!")
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
