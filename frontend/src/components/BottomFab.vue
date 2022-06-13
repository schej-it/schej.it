<template>
  <v-scale-transition appear origin="center">
    <v-btn 
      fab
      fixed
      class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-text-white"
      :class="isPhone ? 'tw-bottom-16' : 'tw-bottom-2'"
      @click="e => $emit('click', e)"
      v-bind="$attrs"
    >
      <slot />
    </v-btn>
  </v-scale-transition>  
</template>

<script>
import { isPhone } from '@/utils'

export default {
  name: 'BottomFab',

  emits: ['click'],

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    }
  },

  created() {
    console.log('PROPS:', this.$props)
    console.log('ATTRS:', this.$attrs)
  }
}
</script>