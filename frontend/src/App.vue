<template>
  <v-app>
    <div
      v-if="showHeader"
      class="tw-h-14 tw-bg-green"
      dark
    >
      <div class="tw-flex tw-items-center tw-justify-center tw-w-full tw-h-full">
        <v-img
          alt="Schej.it Logo"
          class="shrink"
          contain
          src="@/assets/logo_dark.svg"
          transition="scale-transition"
          width="120"
        />
      </div>
    </div>

    <v-main>
      <router-view/>
    </v-main>

    <div
      v-if="showNavbar"
      class="tw-h-14 tw-bg-green tw-flex"
    >
      <div 
        v-for="{ text, icon, route }, i in tabs"
        :key="text"
        class="tw-flex tw-flex-col tw-justify-center tw-items-center tw-flex-1 tw-h-full tw-select-none tw-cursor-pointer hover:tw-bg-dark-green"
        :class="$route.name === route.name ? `tw-bg-dark-green` : ''"
        @click="navigate(i)"
      >
        <v-icon class="tw-text-white">{{ icon }}</v-icon>
        <div class="tw-text-white tw-text-sm">{{ text }}</div>
      </div>
    </div>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data: () => ({
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
    showHeader() {
      return this.$route.name !== 'sign-in'
    },
    showNavbar() {
      return (
        this.$route.name !== 'sign-in' &&
        this.$route.name !== 'join'
      )
    },
  },

  methods: {
    navigate(i) {
      this.tab = i
      this.$router.push(this.tabs[i].route).catch(e => {})
    }
  },
};
</script>

<style>
</style>