<template>
  <div v-if="isPhone">
    <v-menu v-model="menuOpen">
      <template v-slot:activator="{ on, attrs }">
        <v-btn icon v-bind="attrs" v-on="on">
          <v-icon>mdi-menu</v-icon>
        </v-btn>
      </template>
      <v-card class="tw-flex tw-flex-col tw-gap-2 tw-p-2 tw-text-left">
        <slot></slot>
      </v-card>
    </v-menu>
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
