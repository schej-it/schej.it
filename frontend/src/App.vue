<template>
  <v-app>
    <AutoSnackbar color="error" :text="error" />
    <AutoSnackbar color="info" :text="info" />
    <div
      v-if="showHeader"
      class="tw-h-14 tw-bg-white tw-fixed tw-w-screen tw-z-40"
      dark
    >
      <div class="tw-relative tw-px-2 tw-flex tw-items-center tw-justify-center tw-max-w-6xl tw-h-full tw-m-auto">
        <v-img
          @click="goHome"
          alt="Schej.it Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/logo_dark.svg"
          transition="scale-transition"
          width="120"
        />

        <v-spacer />

        <AuthUserMenu></AuthUserMenu>
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
html {
  overflow-y: auto !important; 
  overscroll-behavior: none;
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
    fixHeight() {
      // // Fix height on mobile
      // document.querySelector('.v-application--wrap').style.height = window.innerHeight + 'px'

      // let items = 0 // Counts the number of fixed height items (header and navbar)
      // if (this.showHeader) items++
      // if (this.showBottomNavbar) items++

      // document.querySelector('.v-main').style.maxHeight = `calc(${window.innerHeight}px - ${items} * 3.5rem)`
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

    this.loaded = true    
  },

  mounted() {
    this.mounted = true
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
        if (this.mounted) this.fixHeight()
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