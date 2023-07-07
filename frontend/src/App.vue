<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="tw-bg-blue" :text="info" />
    <SignInNotSupportedDialog v-model="webviewDialog" />
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
            src="@/assets/schej_logo_with_text.png"
            transition="scale-transition"
            :width="isPhone ? 70 : 90"
          />
        </router-link>

        <v-spacer />

        <AuthUserMenu v-if="authUser" />
        <v-btn v-else id="top-right-sign-in-btn" text @click="signIn"
          >Sign in</v-btn
        >
      </div>
    </div>

    <v-main>
      <div class="tw-flex tw-h-screen tw-flex-col">
        <div
          class="tw-relative tw-flex-1 tw-overscroll-auto"
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

.v-text-field.v-text-field--solo:not(.v-text-field--solo-flat)
  > .v-input__control
  > .v-input__slot {
  filter: drop-shadow(0 1px 2px rgb(0 0 0 / 0.1))
    drop-shadow(0 1px 1px rgb(0 0 0 / 0.06)) !important;
  box-shadow: none !important;
}

.v-menu__content {
  filter: drop-shadow(0 10px 8px rgb(0 0 0 / 0.04))
    drop-shadow(0 4px 3px rgb(0 0 0 / 0.1)) !important;
  box-shadow: none !important;
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

export default {
  name: "App",

  metaInfo: {
    title: "schej - Finding a time to meet, made simple",
    htmlAttrs: {
      lang: "en-US",
    },
    // meta: [
    //   { charset: "utf-8" },
    //   {
    //     name: "description",
    //     content: `schej helps you quickly find the best time for your group to meet. It's like When2meet with Google Calendar integration.`,
    //   },
    //   { name: "viewport", content: "width=device-width, initial-scale=1" },
    // ],
  },

  components: {
    AutoSnackbar,
    AuthUserMenu,
    SignInNotSupportedDialog,
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
        this.$route.name !== "createEvent" &&
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
    handleScroll(e) {
      this.scrollY = window.scrollY
    },
    redirectUser(authenticated) {
      let authRoutes = ["home", "settings"]
      let noAuthRoutes = ["landing", "createEvent"]

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
