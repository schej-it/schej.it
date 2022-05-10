<template>
  <div>

    <!-- Dialog -->
    <NewEventDialog 
      v-model="dialog"
    />

    <!-- FAB -->
    <v-scale-transition appear origin="center">
      <v-btn 
        fab
        absolute
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-4"
        @click="dialog = true"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
import NewEventDialog from '@/components/NewEventDialog'

export default {
  name: 'Home',

  components: {
    NewEventDialog,
  },

  data: () => ({
    dialog: false,
  }),
}
</script>
