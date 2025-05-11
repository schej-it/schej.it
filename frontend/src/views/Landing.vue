<template>
  <div class="tw-bg-light-gray">
    <div
      class="tw-relative tw-m-auto tw-mb-32 tw-flex tw-min-h-screen tw-max-w-6xl tw-flex-col"
    >
      <!-- Green background -->
      <div
        class="tw-absolute -tw-bottom-32 tw-left-1/2 tw-h-[40%] tw-w-screen -tw-translate-x-1/2 tw-bg-green sm:tw-h-[800px]"
      ></div>

      <!-- Header -->
      <div class="tw-mb-16">
        <div class="tw-flex tw-items-center tw-px-5 tw-pt-5">
          <Logo type="schej" />

          <v-spacer />

          <LandingPageHeader>
            <v-btn text @click="openHowItWorksDialog">How it works</v-btn>
            <v-btn text href="/blog">Blog</v-btn>
            <v-btn text @click="signIn">Sign in</v-btn>
          </LandingPageHeader>
        </div>
      </div>

      <div class="tw-z-20 tw-mt-12 tw-flex tw-flex-col tw-items-center">
        <div
          class="tw-mx-4 tw-mb-6 tw-flex tw-max-w-[26rem] tw-flex-col tw-items-center sm:tw-w-[35rem] sm:tw-max-w-none"
        >
          <div
            class="tw-mb-4 tw-flex tw-select-none tw-items-center tw-rounded-full tw-border tw-border-light-gray-stroke tw-bg-white/70 tw-px-2.5 tw-py-1.5 tw-text-sm tw-text-dark-gray"
          >
            We're open source!
            <github-button
              v-once
              class="-tw-mb-1 tw-ml-2"
              href="https://github.com/schej-it/schej.it"
              data-show-count="true"
              aria-label="Star schej-it/schej.it on GitHub"
              >Star</github-button
            >
          </div>
          <div
            id="header"
            class="tw-mb-4 tw-text-center tw-text-2xl tw-font-medium sm:tw-text-4xl lg:tw-text-4xl xl:tw-text-5xl"
          >
            <h1>Find a time to meet</h1>
          </div>

          <div
            class="lg:tw-text-md tw-text-left tw-text-center tw-text-sm tw-text-very-dark-gray sm:tw-text-lg md:tw-text-lg xl:tw-text-lg"
          >
            Coordinate group meetings without the back and forth. <br />
            Integrates with your calendar of choice.
          </div>
        </div>

        <v-btn
          id="lets-schej-it-btn"
          class="tw-mb-12 tw-block tw-self-center tw-rounded-lg tw-bg-green tw-px-10 tw-text-base sm:tw-px-10 lg:tw-px-12"
          dark
          @click="newDialog = true"
          large
          :x-large="$vuetify.breakpoint.mdAndUp"
        >
          Create event
        </v-btn>
        <div
          class="tw-rounded-xl tw-border tw-border-light-gray-stroke tw-bg-white tw-shadow-xl"
        >
          <div class="tw-relative tw-mx-4 tw-h-[800px] tw-w-[800px]">
            <v-img
              class="tw-absolute tw-left-0 tw-top-0 tw-transition-opacity tw-duration-300"
              :class="{ 'tw-opacity-0': isVideoPlaying }"
              src="@/assets/img/hero.jpg"
              :height="800"
              :width="800"
              transition="fade-transition"
              contain
            />
            <vue-vimeo-player
              video-url="https://player.vimeo.com/video/1083205305?h=d58bef862a"
              :player-width="800"
              :player-height="800"
              :options="{
                muted: true,
                playsinline: true,
              }"
              :controls="false"
              :autoplay="true"
              :loop="true"
              @play="onPlay"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- How it works -->
    <div
      id="how-it-works"
      class="tw-grid tw-place-content-center tw-px-4 tw-pt-12"
    >
      <div class="tw-mx-auto tw-flex tw-flex-col tw-gap-4">
        <div
          class="tw-mb-4 tw-text-center tw-text-2xl tw-font-medium sm:tw-text-3xl lg:tw-text-4xl"
        >
          How it works
        </div>
        <div
          v-for="(step, i) in howItWorksSteps"
          :key="i"
          class="tw-flex tw-items-center tw-gap-2"
        >
          <NumberBullet>{{ i + 1 }}</NumberBullet>
          <div class="tw-text-base tw-font-medium md:tw-text-xl">
            <div v-if="i == 1">
              <span
                class="tw-underline tw-decoration-[#29BC6888] tw-decoration-4"
                style="text-underline-position: under"
                >Autofill</span
              >
              your availability with Google Calendar
            </div>
            <div v-else v-html="step"></div>
          </div>
        </div>
      </div>
      <div
        class="tw-mb-6 tw-mt-10 tw-text-center tw-text-3xl tw-font-medium md:tw-mb-12 md:tw-mt-20 md:tw-text-6xl"
      >
        It's that simple.
      </div>
      <v-img
        alt="schej character"
        src="@/assets/schej_character.png"
        :height="isPhone ? 200 : 300"
        transition="fade-transition"
        contain
        class="-tw-mb-12"
      />
    </div>

    <!-- Video -->
    <div
      class="tw-flex tw-justify-center tw-bg-green tw-px-4 tw-pb-12 tw-pt-24 md:tw-pb-16"
    >
      <div class="tw-h-64 tw-max-w-3xl tw-flex-1 sm:tw-h-80 md:tw-h-96">
        <iframe
          class="tw-h-full tw-w-full"
          src="https://www.youtube.com/embed/58UcQnaXnBs?fs=0&color=white&rel=0&controls=0"
          title="schej demo"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          allowfullscreen
        ></iframe>
      </div>
    </div>

    <!-- FAQ -->
    <div class="tw-flex tw-justify-center tw-pt-12">
      <div class="tw-mx-4 tw-mb-12 tw-max-w-3xl tw-flex-1 sm:tw-mx-16">
        <div id="faq-section" class="tw-text-center lg:tw-pt-3">
          <Header> Frequently Asked Questions </Header>
          <div
            class="tw-grid tw-grid-cols-1 tw-gap-3 sm:tw-text-xl lg:tw-text-2xl"
          >
            <FAQ
              v-for="faq in faqs"
              :key="faq.question"
              @signIn="signIn"
              v-bind="faq"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Privacy Policy -->
    <div class="tw-flex tw-flex-col tw-bg-green">
      <div
        class="tw-m-2 tw-flex tw-flex-col tw-items-center tw-gap-4 tw-pb-6 tw-pt-6"
      >
        <router-link
          class="tw-font-bold tw-text-white"
          :to="{ name: 'privacy-policy' }"
        >
          Privacy Policy
        </router-link>
        <div class="tw-text-light-gray">Made with 💚 by the schej team</div>
      </div>
    </div>

    <!-- Sign in dialog -->
    <SignInDialog v-model="signInDialog" @signIn="_signIn" />

    <!-- New event dialog -->
    <NewDialog
      v-model="newDialog"
      :allow-notifications="false"
      no-tabs
      @signIn="signIn"
    />

    <!-- Add the dialog component -->
    <HowItWorksDialog
      v-if="showHowItWorksDialog"
      v-model="showHowItWorksDialog"
    />
  </div>
</template>

<style scoped>
@media screen and (min-width: 375px) and (max-width: 640px) {
  #header {
    font-size: 1.875rem !important; /* 30px */
    line-height: 2.25rem !important; /* 36px */
  }
}
</style>

<script>
import LandingPageCalendar from "@/components/landing/LandingPageCalendar.vue"
import { isPhone, signInGoogle, signInOutlook } from "@/utils"
import FAQ from "@/components/FAQ.vue"
import Header from "@/components/Header.vue"
import NumberBullet from "@/components/NumberBullet.vue"
import NewEvent from "@/components/NewEvent.vue"
import NewDialog from "@/components/NewDialog.vue"
import LandingPageHeader from "@/components/landing/LandingPageHeader.vue"
import Logo from "@/components/Logo.vue"
import GithubButton from "vue-github-button"
import SignInDialog from "@/components/SignInDialog.vue"
import { calendarTypes } from "@/constants"
import HowItWorksDialog from "@/components/HowItWorksDialog.vue"
import { vueVimeoPlayer } from "vue-vimeo-player"

export default {
  name: "Landing",

  metaInfo: {
    title: "Schej - Find a time to meet",
  },

  components: {
    LandingPageCalendar,
    FAQ,
    Header,
    NumberBullet,
    NewEvent,
    NewDialog,
    LandingPageHeader,
    GithubButton,
    Logo,
    SignInDialog,
    HowItWorksDialog,
    vueVimeoPlayer,
  },

  data: () => ({
    signInDialog: false,
    newDialog: false,
    githubSnackbar: true,
    howItWorksSteps: [
      "Create a schej event",
      "Autofill your availability with Google Calendar",
      "Share the schej with your group and find the best time to meet!",
    ],
    faqs: [
      {
        question: "How is schej different from lettucemeet or when2meet?",
        points: [
          "Much better UI (web and mobile)",
          "Seamless and working calendar integration",
          "A slew of other features that we don't have space to list here",
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
      {
        question: `I want it so that only I can see people's responses.`,
        answer: `Just check "Only show responses to event creator" under Advanced Options when creating your event! Other respondees will not be able to see each other's names or availability.`,
        authRequired: true,
      },
      {
        question: `Can I receive emails when someone fills out my event?`,
        answer: `Absolutely! Check "Email me each time someone joins my event" when creating an event. <br><br>To receive email notifications after a specific number (X) of responses are added, check "Email me after X responses" in Advanced Options.`,
        authRequired: true,
      },
      {
        question: `How do I send reminders to people to fill out an event?`,
        answer: `Open the "Email Reminders" section when creating an event and input everybody's email address. Reminder emails will be sent the day of event creation, one day after, and three days after. <br><br>You will also receive an email once everybody has filled out the Schej.`,
        authRequired: true,
      },
    ],
    rive: null,
    showSchejy: false,
    showHowItWorksDialog: false,
    isVideoPlaying: false,
  }),

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    loadRiveAnimation() {
      // if (!this.rive) {
      //   this.rive = new Rive({
      //     src: "/rive/schej.riv",
      //     canvas: document.querySelector("canvas"),
      //     autoplay: false,
      //     stateMachines: "wave",
      //     onLoad: () => {
      //       // r.resizeDrawingSurfaceToCanvas()
      //     },
      //   })
      //   setTimeout(() => {
      //     this.showSchejy = true
      //     setTimeout(() => {
      //       this.rive.play("wave")
      //     }, 1000)
      //   }, 4000)
      // } else {
      //   this.rive.play("wave")
      // }
    },
    _signIn(calendarType) {
      if (calendarType === calendarTypes.GOOGLE) {
        signInGoogle({ state: null, selectAccount: true })
      } else if (calendarType === calendarTypes.OUTLOOK) {
        // NOTE: selectAccount is not supported implemented yet for Outlook, maybe add it later
        signInOutlook({ state: null, selectAccount: true })
      }
    },
    signIn() {
      this.signInDialog = true
    },
    openHowItWorksDialog() {
      this.showHowItWorksDialog = true
      this.$posthog.capture("how_it_works_clicked")
    },
    onPlay() {
      setTimeout(() => {
        this.isVideoPlaying = true
      }, 300)
    },
  },

  beforeDestroy() {
    this.rive?.cleanup()
  },

  watch: {
    [`$vuetify.breakpoint.name`]: {
      immediate: true,
      handler() {
        if (this.$vuetify.breakpoint.mdAndUp) {
          setTimeout(() => {
            this.loadRiveAnimation()
          }, 0)
        }
      },
    },
  },
}
</script>
