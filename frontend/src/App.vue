<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="tw-bg-blue" :text="info" />
    <SignInNotSupportedDialog v-model="webviewDialog" />
    <NewDialog v-model="newDialog" type="event" no-tabs />
    <div
      v-if="showHeader"
      class="tw-fixed tw-z-40 tw-h-14 tw-w-screen tw-bg-white sm:tw-h-16"
      dark
    >
      <div
        class="tw-relative tw-m-auto tw-flex tw-h-full tw-max-w-6xl tw-items-center tw-justify-center tw-px-4"
      >
        <router-link :to="{ name: 'home' }">
          <v-img
            alt="Schej Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="@/assets/april_fools_logo.png"
            transition="scale-transition"
            :width="isPhone ? 200 : 300"
          />
        </router-link>

        <v-spacer />

        <v-btn
          v-if="$route.name !== 'home'"
          id="top-right-create-btn"
          text
          @click="createEvent"
        >
          Create an event
        </v-btn>
        <v-btn
          v-if="showFeedbackBtn"
          id="feedback-btn"
          text
          href="https://forms.gle/9AgRy4PQfWfVuBnw8"
          target="_blank"
        >
          Give feedback
        </v-btn>
        <div v-if="authUser" class="sm:tw-ml-4">
          <AuthUserMenu />
        </div>
        <v-btn v-else id="top-right-sign-in-btn" text @click="signIn">
          Sign in
        </v-btn>
      </div>
    </div>

    <v-main>
      <div class="tw-flex tw-h-screen tw-flex-col">
        <div
          class="tw-relative tw-flex-1 tw-overscroll-auto"
          :class="routerViewClass"
        >
          <router-view v-if="loaded" :key="$route.fullPath" />
        </div>
      </div>
    </v-main>
  </v-app>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=DM+Sans&display=swap");

html {
  overflow-y: auto !important;
  /* overscroll-behavior: none; */
  scroll-behavior: smooth;
}

* {
  font-family: "DM Sans", sans-serif;
  /* touch-action: manipulation !important; */
}

.v-btn {
  letter-spacing: unset !important;
  text-transform: unset !important;
}

.v-text-field.v-text-field--solo:not(.v-text-field--solo-flat)
  > .v-input__control
  > .v-input__slot {
  filter: drop-shadow(0 1px 2px rgb(0 0 0 / 0.1))
    drop-shadow(0 1px 1px rgb(0 0 0 / 0.06)) !important;
  box-shadow: none !important;
}

.v-menu__content {
  box-shadow: 0px 5px 5px -1px rgba(0, 0, 0, 0.1),
    0px 8px 10px 0.5px rgba(0, 0, 0, 0.07), 0px 3px 14px 1px rgba(0, 0, 0, 0.06) !important;
}
</style>

<script>
import { mapMutations, mapState } from "vuex"
import { get, getLocation, isPhone, post, signInGoogle } from "./utils"
import { authTypes } from "./constants"
import AutoSnackbar from "@/components/AutoSnackbar"
import AuthUserMenu from "./components/AuthUserMenu.vue"
import SignInNotSupportedDialog from "./components/SignInNotSupportedDialog.vue"
import isWebview from "is-ua-webview"
import NewDialog from "./components/NewDialog.vue"

export default {
  name: "App",

  metaInfo: {
    htmlAttrs: {
      lang: "en-US",
    },
  },

  components: {
    AutoSnackbar,
    AuthUserMenu,
    SignInNotSupportedDialog,
    NewDialog,
  },

  data: () => ({
    mounted: false,
    loaded: false,
    scrollY: 0,
    webviewDialog: false,
    newDialog: false,
  }),

  computed: {
    ...mapState(["authUser", "error", "info"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showHeader() {
      return (
        this.$route.name !== "createEvent" &&
        this.$route.name !== "landing" &&
        this.$route.name !== "auth" &&
        this.$route.name !== "privacy-policy"
      )
    },
    showFeedbackBtn() {
      return !this.isPhone || this.$route.name === "home"
    },
    routerViewClass() {
      let c = ""
      if (this.showHeader) {
        if (this.isPhone) {
          c += "tw-pt-12 "
        } else {
          c += "tw-pt-14 "
        }
      }
      return c
    },
  },

  methods: {
    ...mapMutations(["setAuthUser"]),
    handleScroll(e) {
      this.scrollY = window.scrollY
    },
    createEvent() {
      this.newDialog = true
    },
    signIn() {
      if (this.$route.name === "event") {
        if (isWebview(navigator.userAgent)) {
          this.webviewDialog = true
          return
        }
        signInGoogle({
          state: {
            type: authTypes.EVENT_SIGN_IN,
            eventId: this.$route.params.eventId,
          },
          selectAccount: true,
        })
      }
    },
  },

  async created() {
    await get("/user/profile")
      .then((authUser) => {
        // console.log(authUser)
        this.setAuthUser(authUser)

        this.$posthog?.identify(authUser._id, {
          email: authUser.email,
          firstName: authUser.firstName,
          lastName: authUser.lastName,
        })
      })
      .catch(() => {
        this.setAuthUser(null)
      })

    // Event listeners
    window.addEventListener("scroll", this.handleScroll)

    this.loaded = true
  },

  mounted() {
    this.mounted = true
    this.scrollY = window.scrollY
  },

  beforeDestroy() {
    window.removeEventListener("scroll", this.handleScroll)
  },

  watch: {
    $route: {
      immediate: true,
      async handler() {
        const originalHref = window.location.href
        if (this.$route.name) {
          this.$posthog?.capture("$pageview")
        }

        // Check for poster query parameter
        if (this.$route.query.p) {
          let location = null
          try {
            location = await getLocation()
          } catch (e) {
            // User probably has adblocker
          }

          post("/analytics/scanned-poster", {
            url: originalHref,
            location,
          })
        }
      },
    },
  },
}
</script>
