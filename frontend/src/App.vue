<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="tw-bg-blue" :text="info" />
    <SignInNotSupportedDialog v-model="webviewDialog" />
    <div
      v-if="showHeader"
      class="tw-h-14 sm:tw-h-16 tw-bg-white tw-fixed tw-w-screen tw-z-40"
      dark
      :class="'tw-drop-shadow'"
    >
      <div
        class="tw-relative tw-px-2 tw-flex tw-items-center tw-justify-center tw-max-w-6xl tw-h-full tw-m-auto"
      >
        <v-img
          @click="goHome"
          alt="Schej Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/schej_logo_with_text.png"
          transition="scale-transition"
          width="90"
        />

        <v-spacer />

        <AuthUserMenu v-if="authUser" />
        <v-btn v-else id="top-right-sign-in-btn" text @click="signIn">Sign in</v-btn>
      </div>
    </div>

    <v-main>
      <div class="tw-h-screen tw-flex tw-flex-col">
        <div
          class="tw-flex-1 tw-relative tw-overscroll-auto"
          :class="routerViewClass"
        >
          <router-view v-if="loaded" />
        </div>
      </div>
    </v-main>
  </v-app>
</template>

<style>
@import url("https://fonts.googleapis.com/css2?family=DM+Sans&display=swap");

html {
  overflow-y: auto !important;
  overscroll-behavior: none;
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
</style>

<script>
import { mapMutations, mapState } from "vuex"
import { get, getLocation, isPhone, post, signInGoogle } from "./utils"
import AutoSnackbar from "@/components/AutoSnackbar"
import AuthUserMenu from "./components/AuthUserMenu.vue"
import SignInNotSupportedDialog from "./components/SignInNotSupportedDialog.vue"
import isWebview from "is-ua-webview"

export default {
  name: "App",

  metaInfo: {
    title: "schej - Finding a time to meet, made simple",
    htmlAttrs: {
      lang: "en-US",
    },
    meta: [
      { charset: "utf-8" },
      {
        name: "description",
        content: `schej is a group scheduling platform that helps you find a time to meet.
Users fill out their availability with the help of Google Calendar, and a heat map is generated showing when everybody is available.
It's completely free to use, and it looks great on mobile.`,
      },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
    ],
  },

  components: {
    AutoSnackbar,
    AuthUserMenu,
    SignInNotSupportedDialog
  },

  data: () => ({
    mounted: false,
    loaded: false,
    scrollY: 0,
    webviewDialog: false,
  }),

  computed: {
    ...mapState(["authUser", "error", "info"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showHeader() {
      return (
        this.$route.name !== "landing" &&
        this.$route.name !== "auth" &&
        this.$route.name !== "privacy-policy"
      )
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
    goHome() {
      if (this.$route.name !== "home") {
        this.$router.push({ name: "home" })
      } else {
        location.reload()
      }
    },
    handleScroll(e) {
      this.scrollY = window.scrollY
    },
    redirectUser(authenticated) {
      let authRoutes = ["home"]
      let noAuthRoutes = ["landing"]

      if (!authenticated) {
        if (authRoutes.includes(this.$route.name)) {
          this.$router.replace({ name: "landing" })
          // console.log('redirecting to SIGN IN')
        }
      } else {
        if (noAuthRoutes.includes(this.$route.name)) {
          this.$router.replace({ name: "home" })
          // console.log('redirecting to HOME')
        }
      }
    },
    signIn() {
      if (this.$route.name === "event") {
        if (isWebview(navigator.userAgent)) {
          this.webviewDialog = true
          return
        }
        signInGoogle(
          { type: "event-sign-in", eventId: this.$route.params.eventId },
          true
        )
      }
    },
  },

  async created() {
    await get("/user/profile")
      .then((authUser) => {
        // console.log(authUser)
        this.setAuthUser(authUser)
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
    authUser: {
      immediate: true,
      handler() {
        if (this.authUser) {
          this.redirectUser(true)
        } else {
          this.redirectUser(false)
        }
      },
    },
    $route: {
      immediate: true,
      async handler() {
        const originalHref = window.location.href

        get("/auth/status")
          .then((data) => {
            this.redirectUser(true)
          })
          .catch((err) => {
            this.redirectUser(false)
          })

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
