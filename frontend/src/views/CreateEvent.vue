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
import { isPhone, signInGoogle } from "@/utils"
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import NewEvent from "@/components/NewEvent.vue"
import { confetti } from "tsparticles-confetti"

export default {
  name: "CreateEvent",

  components: {
    SignInGoogleBtn,
    NewEvent,
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
