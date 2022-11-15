<template>
  <v-scale-transition appear origin="center">
    <v-btn 
      fab
      fixed
      class="tw-bg-green tw-mx-auto tw-left-0 tw-right-0 tw-bottom-2 tw-text-white"
      @click="e => $emit('click', e)"
      v-bind="$attrs"
    >
      <slot />
    </v-btn>
  </v-scale-transition>  
</template>

<script>
export default {
  name: 'BottomFab',

  emits: ['click'],
}
</script>