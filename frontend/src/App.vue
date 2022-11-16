<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="info" :text="info" />
    <div
      v-if="showHeader"
      class="tw-h-14 tw-bg-white tw-fixed tw-w-screen tw-z-40"
      dark
      :class="scrollY > 0 && 'tw-drop-shadow'"
    >
      <div 
        class="tw-relative tw-px-2 tw-flex tw-items-center tw-justify-center tw-max-w-6xl tw-h-full tw-m-auto"
      >
        <v-img
          @click="goHome"
          alt="Schej Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/schej_logo_with_text.svg"
          transition="scale-transition"
          width="120"
        />

        <v-spacer />

        <AuthUserMenu v-if="authUser" />
        <v-btn 
          v-else
          text
        >Sign in</v-btn>
      </div>
    </div>

    <v-main>
      <div class="tw-h-screen tw-flex tw-flex-col">
        <div class="tw-flex-1 tw-relative tw-overscroll-auto" :class="routerViewClass">
          <router-view v-if="loaded" />
        </div>
      </div>
    </v-main>

  </v-app>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=DM+Sans&display=swap');

html {
  overflow-y: auto !important; 
  overscroll-behavior: none;
}

* {
  font-family: 'DM Sans', sans-serif;
}

.v-btn {
  letter-spacing: unset !important;
  text-transform: unset !important;
}
</style>

<script>
import { mapMutations, mapState } from 'vuex';
import { get, isPhone } from './utils';
import AutoSnackbar from '@/components/AutoSnackbar'
import AuthUserMenu from './components/AuthUserMenu.vue';

export default {
  name: 'App',

  components: {
    AutoSnackbar,
    AuthUserMenu
  },

  data: () => ({
    mounted: false,
    loaded: false,
    scrollY: 0,
  }),

  computed: {
    ...mapState([ 'authUser', 'error', 'info' ]),
    showHeader() {
      return (
        this.$route.name !== 'landing' &&
        this.$route.name !== 'auth' &&
        this.$route.name !== 'privacy-policy'
      )
    },
    routerViewClass() {
      let c = ''
      if (this.showHeader) c += 'tw-pt-14 '
      return c
    },
  },

  methods: {
    ...mapMutations([ 'setAuthUser' ]),
    goHome() {
      this.$router.push({ name: 'home' })
    },
    handleScroll(e) {
      this.scrollY = window.scrollY
    },
    redirectUser(authenticated) {
      let authRoutes = ['home']
      let noAuthRoutes = ['landing']
      
      if (!authenticated) {
        if (authRoutes.includes(this.$route.name)) {
          this.$router.replace({ name: 'landing' })
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

    // Event listeners
    window.addEventListener('scroll', this.handleScroll)

    this.loaded = true    
  },

  mounted() {
    this.mounted = true
    this.scrollY = window.scrollY
  },

  destroyed() {
    window.removeEventListener('scroll', this.handleScroll)
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
    $route: {
      immediate: true,
      handler() {
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