<template>
  <div v-if="isPhone">
    <v-btn icon @click="toggleMenu">
      <v-icon>mdi-menu</v-icon>
    </v-btn>
    <div
      v-if="menuOpen"
      class="tw-absolute tw-left-0 twj-right-0 tw-z-30 tw-flex tw-w-full tw-flex-col tw-gap-2 tw-bg-white tw-p-2"
    >
      <slot></slot>
    </div>
  </div>
  <div v-else class="tw-flex tw-items-center">
    <slot></slot>
  </div>
</template>

<script>
import { isPhone } from "@/utils"

export default {
  name: "LandingPageHeader",

  data() {
    return {
      menuOpen: false,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    toggleMenu() {
      this.menuOpen = !this.menuOpen
    },
  },
}
</script>
