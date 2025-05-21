<template>
  <v-dialog v-model="dialog" max-width="400px" content-class="tw-m-0">
    <v-card>
      <v-card-title>
        <span class="tw-text-xl tw-font-medium">Oops! Feature Not Ready</span>
        <v-spacer />
        <v-btn
          absolute
          @click="dialog = false"
          icon
          class="tw-right-0 tw-mr-2 tw-self-center"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text class="tw-pt-4 tw-text-very-dark-gray">
        You've caught us a bit early! This feature isn't quite ready yet. We're
        working on it and will let you know when "Folders" are released.
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" text @click="dialog = false">Got it!</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "FeatureNotReadyDialog",
  props: {
    value: Boolean,
  },
  computed: {
    dialog: {
      get() {
        return this.value
      },
      set(val) {
        this.$emit("input", val)
      },
    },
  },
}
</script>

<style scoped></style>
