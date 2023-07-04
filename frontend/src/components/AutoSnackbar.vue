<template>
  <v-snackbar v-model="show" top :color="color">
    <span class="tw-mr-2 tw-text-sm">{{ text }}</span>

    <template v-slot:action="{ attrs }">
      <v-btn v-bind="attrs" icon @click="show = false">
        <v-icon>mdi-close</v-icon>
      </v-btn>
    </template>
  </v-snackbar>
</template>

<script>
export default {
  name: "AutoSnackbar",

  props: {
    text: { type: String, default: "" },
    color: { type: String, default: "" },
  },

  data() {
    return {
      show: false,
    }
  },

  created() {},

  watch: {
    text: {
      immediate: true,
      handler(text) {
        if (text) {
          this.show = true
        } else this.show = false
      },
    },
  },
}
</script>
