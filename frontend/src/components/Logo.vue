<template>
  <v-img
    :alt="alt"
    class="shrink tw-cursor-pointer"
    contain
    :src="src"
    transition="scale-transition"
    :width="width"
  />
</template>

<script>
import { isPhone } from "@/utils"

export default {
  name: "Logo",

  props: {
    type: {
      type: "schej" | "betterwhen2meet" | "aprilfools",
      default: "schej",
    },
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    alt() {
      if (this.type === "betterwhen2meet") {
        return "Betterwhen2meet Logo"
      }

      return "Schej Logo"
    },
    src() {
      switch (this.type) {
        case "schej":
          return require("@/assets/schej_logo_with_text.png")
        case "betterwhen2meet":
          return require("@/assets/april_fools_logo.png")
        case "aprilfools":
          return require("@/assets/april_fools_logo.png")
      }
    },
    width() {
      switch (this.type) {
        case "schej":
          return this.isPhone ? 70 : 90
        case "betterwhen2meet":
          return this.isPhone ? 200 : 300
        case "aprilfools":
          return this.isPhone ? 200 : 300
      }
    },
  },
}
</script>
