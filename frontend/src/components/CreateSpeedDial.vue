<template>
  <v-scale-transition appear origin="center">
    <v-speed-dial
      v-model="speedDial"
      fixed
      transition="slide-y-reverse-transition"
      class="tw-bottom-4 tw-left-0 tw-right-0 tw-mx-auto tw-w-fit sm:tw-bottom-8"
    >
      <template v-slot:activator>
        <v-btn v-model="speedDial" fab color="primary" dark>
          <v-icon v-if="fab"> mdi-close </v-icon>
          <v-icon v-else> mdi-plus </v-icon>
        </v-btn>
      </template>
      <v-btn fab small dark class="tw-bg-dark-gray">
        <v-icon small>mdi-account-group</v-icon>
      </v-btn>
      <v-btn fab small dark class="tw-bg-dark-gray">
        <v-icon small>mdi-calendar-cursor</v-icon>
      </v-btn>
    </v-speed-dial>
  </v-scale-transition>
</template>

<script>
export default {
  name: "CreateSpeedDial",

  emits: ["createGroup", "createEvent"],

  data() {
    return {
      speedDial: false,
    }
  },
}
</script>
