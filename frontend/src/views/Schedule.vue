<template>
  <div>
    <v-scale-transition appear origin="center">
      <v-btn 
        fab
        absolute
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-4"
      >
        <v-icon>mdi-share</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
export default {
  name: 'Schedule',
}
</script>