<template>
  <v-snackbar
    v-model="show"
    top
    :color="color"
  >
    {{ text }}
    <v-btn
      text
      @click="show = false"
    >
      Close
    </v-btn>
  </v-snackbar>
</template>

<script>
export default {
  name: 'AutoSnackbar',
  
  props: {
    text: {type: String, default: ''},
    color: {type: String, default: ''}
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
        } else
          this.show = false
      },
    },
  },
}
</script>