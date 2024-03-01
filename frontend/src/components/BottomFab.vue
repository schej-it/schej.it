<template>
  <v-scale-transition appear origin="center">
    <v-btn
      fab
      fixed
      class="tw-bottom-4 tw-left-0 tw-right-0 tw-mx-auto tw-bg-green tw-text-white sm:tw-bottom-8"
      @click="(e) => $emit('click', e)"
      v-bind="$attrs"
    >
      <slot />
    </v-btn>
  </v-scale-transition>
</template>

<script>
export default {
  name: "BottomFab",

  emits: ["click"],
}
</script>
