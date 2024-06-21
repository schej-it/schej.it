<template>
  <div>
    <div
      class="tw-relative tw-m-auto tw-flex tw-min-h-screen tw-max-w-6xl tw-flex-col"
    >
      <!-- Green background -->
      <div
        class="tw-absolute tw-bottom-0 tw-left-1/2 tw-h-[40%] tw-w-screen -tw-translate-x-1/2 tw-bg-green sm:tw-h-[47%]"
      ></div>

      <!-- Header -->
      <div class="tw-mb-16">
        <div class="tw-flex tw-items-center tw-px-5 tw-pt-5">
          <Logo type="schej" />

          <v-spacer />

          <LandingPageHeader>
            <v-btn text href="#how-it-works">How it works</v-btn>
            <v-btn text href="/blog">Blog</v-btn>
            <v-btn text @click="signIn">Sign in</v-btn>
          </LandingPageHeader>
        </div>
      </div>

      <div
        class="tw-relative tw-flex tw-justify-center tw-pb-12 lg:tw-justify-between lg:tw-pb-24"
      >
        <!-- Left side -->
        <div class="tw-flex tw-flex-col">
          <!-- Hero -->
          <div
            class="tw-mx-4 tw-flex tw-max-w-[26rem] tw-flex-col tw-items-center sm:tw-w-[35rem] sm:tw-max-w-none sm:tw-items-start lg:tw-mb-14"
          >
            <div
              id="header"
              class="tw-mb-4 tw-text-center tw-text-2xl tw-font-medium sm:tw-text-left sm:tw-text-4xl lg:tw-text-4xl xl:tw-text-5xl"
            >
              <div
                class="tw-bg-gradient-to-r tw-from-light-green tw-to-darkest-green tw-bg-clip-text tw-pb-1 tw-text-transparent"
              >
                Let's schej it!
              </div>
              <div class="-tw-mt-1 tw-leading-snug">Find a time to meet</div>
            </div>

            <div
              class="lg:tw-text-md tw-mb-4 tw-text-left tw-text-center tw-text-sm tw-text-very-dark-gray sm:tw-text-left sm:tw-text-lg md:tw-text-lg xl:tw-text-lg"
            >
              <b>Automatically</b> fill in your availability—it’s like When2Meet
              with Google Calendar integration!
            </div>
          </div>

          <v-btn
            id="lets-schej-it-btn"
            class="tw-my-6 tw-block tw-self-center tw-rounded-lg tw-bg-green tw-px-10 tw-text-base sm:tw-px-10 lg:tw-hidden lg:tw-px-12"
            dark
            @click="newDialog = true"
            large
            :x-large="$vuetify.breakpoint.mdAndUp"
          >
            Create event
          </v-btn>

          <!-- Calendar -->
          <div>
            <v-img
              alt="schej character"
              src="@/assets/schejie/wave.png"
              :height="isPhone ? 70 : 80"
              transition="fade-transition"
              contain
              class="-tw-mb-4 tw-mt-2 tw-block sm:tw-mt-6 lg:tw-hidden"
            />
            <LandingPageCalendar class="tw-drop-shadow-lg" />
          </div>
        </div>

        <!-- Right side -->
        <div class="tw-ml-6 tw-mr-4 tw-hidden lg:tw-flex lg:tw-flex-col">
          <!-- <v-img
            alt="schej character"
            src="@/assets/schejie/wave.png"
            :height="90"
            transition="fade-transition"
            contain
            class=""
          /> -->
          <!-- Placeholder when schejy is not shown -->
          <div
            v-if="!showSchejy"
            id="canvas"
            class="-tw-mb-36 -tw-mt-24 tw-h-[350px] tw-w-[350px] tw-overflow-hidden"
          ></div>
          <v-slide-y-reverse-transition>
            <div v-show="showSchejy" class="tw-self-center tw-overflow-hidden">
              <canvas
                id="canvas"
                width="700"
                height="700"
                class="-tw-mb-36 -tw-mt-24 tw-h-[350px] tw-w-[350px] tw-overflow-hidden"
              ></canvas>
            </div>
          </v-slide-y-reverse-transition>
          <NewEvent
            class="tw-drop-shadow-lg"
            :dialog="false"
            :allow-notifications="false"
            @signIn="signIn"
          />
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
        class="tw-m-2 tw-flex tw-flex-col tw-items-center tw-gap-4 tw-pb-6 tw-pt-6 sm:tw-pb-20"
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
    <v-dialog v-model="signInDialog" :width="400">
      <v-card>
        <v-card-title>Sign in</v-card-title>
        <v-card-text class="tw-flex tw-flex-col tw-items-center">
          <SignInGoogleBtn class="tw-mb-4" @click="signInGoogle" />
          <div class="tw-text-center tw-text-xs">
            By continuing, you agree to our
            <router-link class="tw-text-blue" :to="{ name: 'privacy-policy' }"
              >privacy policy</router-link
            >
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- New event dialog -->
    <NewDialog
      v-model="newDialog"
      :allow-notifications="false"
      no-tabs
      @signIn="signIn"
    />

    <!-- GitHub button -->
    <v-snackbar
      v-if="!isPhone"
      min-width="unset"
      v-model="githubSnackbar"
      bottom
      :timeout="-1"
      content-class="tw-flex tw-items-center tw-gap-x-4"
      rounded="pill"
    >
      We're fully open source on GitHub!
      <!-- Place this tag where you want the button to render. -->
      <github-button
        class="-tw-mb-1"
        href="https://github.com/schej-it/schej.it"
        data-color-scheme="no-preference: light; light: light; dark: dark;"
        data-size="large"
        data-show-count="true"
        aria-label="Star schej-it/schej.it on GitHub"
        >Star</github-button
      >
      <template v-slot:action="{ attrs }">
        <v-btn v-bind="attrs" icon @click="githubSnackbar = false">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
    </v-snackbar>
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
import { isPhone, signInGoogle } from "@/utils"
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import FAQ from "@/components/FAQ.vue"
import Header from "@/components/Header.vue"
import NumberBullet from "@/components/NumberBullet.vue"
import NewEvent from "@/components/NewEvent.vue"
import NewDialog from "@/components/NewDialog.vue"
import LandingPageHeader from "@/components/landing/LandingPageHeader.vue"
import Logo from "@/components/Logo.vue"
import { Rive } from "@rive-app/canvas"
import GithubButton from "vue-github-button"

export default {
  name: "Landing",

  metaInfo: {
    title: "Schej - Finding a time to meet, made simple",
  },

  components: {
    LandingPageCalendar,
    SignInGoogleBtn,
    FAQ,
    Header,
    NumberBullet,
    NewEvent,
    NewDialog,
    LandingPageHeader,
    GithubButton,
    Logo,
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
  }),

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    loadRiveAnimation() {
      if (!this.rive) {
        this.rive = new Rive({
          src: "/rive/schej.riv",
          canvas: document.querySelector("canvas"),
          autoplay: false,
          stateMachines: "wave",
          onLoad: () => {
            // r.resizeDrawingSurfaceToCanvas()
          },
        })
        setTimeout(() => {
          this.showSchejy = true
          setTimeout(() => {
            this.rive.play("wave")
          }, 1000)
        }, 4000)
      } else {
        this.rive.play("wave")
      }
    },
    signInGoogle() {
      signInGoogle({ state: null, selectAccount: true })
    },
    signIn() {
      this.signInDialog = true
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
