<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="info" :text="info" />
    <div
      v-if="showHeader"
      class="tw-h-14 tw-bg-green"
      dark
    >
      <div class="tw-relative tw-px-2 tw-flex tw-items-center tw-justify-center tw-max-w-6xl tw-h-full tw-m-auto">
        <v-img
          @click="navigate(0)"
          alt="Schej.it Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/logo_dark.svg"
          transition="scale-transition"
          width="120"
        />

        <template v-if="!centerHeaderLogo">

          <v-spacer />

          <div
            class="tw-absolute tw-h-full tw-hidden sm:tw-flex"
          >
            <div 
              v-for="{ text, icon, route }, i in tabs"
              :key="text"
              class="tw-w-28 tw-flex tw-flex-col tw-justify-center tw-items-center tw-flex-1 tw-h-full tw-select-none tw-cursor-pointer tw-brightness-95 hover:tw-brightness-150"
              :class="$route.name === route.name ? `tw-border-b-4 tw-border-white tw-border-solid tw-brightness-150` : ''"
              @click="navigate(i)"
            >
              <v-icon class="tw-text-white">{{ icon }}</v-icon>
              <div class="tw-text-white tw-text-sm">{{ text }}</div>
            </div>
          </div>

          <v-spacer />

          <v-avatar v-if="authUser">
            <img referrerpolicy="no-referrer" :src="authUser.picture">
          </v-avatar>

        </template>
      </div>
    </div>

    <v-main class="tw-overflow-y-auto tw-flex tw-justify-center">
      <div class="tw-max-w-6xl tw-mx-auto tw-h-full">
        <router-view v-if="loaded" />
      </div>
    </v-main>

    <div
      v-if="showBottomNavbar"
      class="tw-h-14 tw-bg-green tw-flex"
    >
      <div 
        v-for="{ text, icon, route }, i in tabs"
        :key="text"
        class="tw-flex tw-flex-col tw-justify-center tw-items-center tw-flex-1 tw-h-full tw-select-none tw-cursor-pointer tw-brightness-95 hover:tw-brightness-150"
        :class="$route.name === route.name ? `tw-border-b-4 tw-border-white tw-border-solid tw-brightness-150` : ''"
        @click="navigate(i)"
      >
        <v-icon class="tw-text-white">{{ icon }}</v-icon>
        <div class="tw-text-white tw-text-sm">{{ text }}</div>
      </div>
    </div>
  </v-app>
</template>

<script>
import { mapMutations, mapState } from 'vuex';
import { get, isPhone } from './utils';
import AutoSnackbar from '@/components/AutoSnackbar'

export default {
  name: 'App',

  components: {
    AutoSnackbar,
  },

  data: () => ({
    loaded: false,
    tabs: [
      {
        text: 'Home',
        icon: 'mdi-home',
        route: { name: 'home' },
      },
      {
        text: 'My schedule',
        icon: 'mdi-calendar',
        route: { name: 'schedule' },
      },
      {
        text: 'Friends',
        icon: 'mdi-account-multiple',
        route: { name: 'friends' },
      },
    ],
    tab: 0,
  }),

  computed: {
    ...mapState([ 'authUser', 'error', 'info' ]),
    showHeader() {
      return (
        this.$route.name !== 'sign-in' &&
        this.$route.name !== 'auth' &&
        this.$route.name !== 'privacy-policy'
      )
    },
    showBottomNavbar() {
      return (
        isPhone(this.$vuetify) &&
        this.$route.name !== 'sign-in' &&
        this.$route.name !== 'join' &&
        this.$route.name !== 'auth'
      )
    },
    centerHeaderLogo() {
      return (
        this.$route.name === 'join'
      )
    },
  },

  methods: {
    ...mapMutations([ 'setAuthUser' ]),
    navigate(i) {
      this.tab = i
      this.$router.push(this.tabs[i].route).catch(e => {})
    },
    fixHeight() {
      // Fix height on mobile
      document.querySelector('.v-application--wrap').style.height = window.innerHeight + 'px'

      let items = 0 // Counts the number of fixed height items (header and navbar)
      if (this.showHeader) items++
      if (this.showBottomNavbar) items++

      document.querySelector('.v-main').style.maxHeight = `calc(${window.innerHeight}px - ${items} * 3.5rem)`
    },
    redirectUser(authenticated) {
      let authRoutes = ['home', 'schedule', 'friends', 'event']
      let noAuthRoutes = ['sign-in']

      if (!authenticated) {
        if (authRoutes.includes(this.$route.name)) {
          this.$router.replace({ name: 'sign-in' })
          console.log('redirecting to SIGN IN')
        }
      } else {
        if (noAuthRoutes.includes(this.$route.name)) {
          this.$router.replace({ name: 'home' })
          console.log('redirecting to HOME')
        }
      }
    },
  },

  async created() {
    await get('/user/profile')
      .then(authUser => {
        // console.log(authUser)
        this.setAuthUser(authUser)
      }).catch(() => {
        this.setAuthUser(null)
      })

    this.loaded = true    
  },

  mounted() {
    this.fixHeight()
    window.addEventListener('resize', this.fixHeight)
  },

  watch: {
    authUser: {
      immediate: true,
      handler() {
        /*if (this.authUser) {
          this.redirectAuthUser(true)
        } else {
          this.redirectAuthUser(false)
        }*/
      }
    },
    '$vuetify.breakpoint.name': {
      handler() {
        this.fixHeight()
      },
    },
    $route: {
      immediate: true,
      handler() {
        this.fixHeight()
        get('/auth/status').then(data => {
          this.redirectUser(true)
        }).catch(err => {
          this.redirectUser(false)
        })
      }
    },
  },
};
</script>

<style>
</style>