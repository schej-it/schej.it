<template>
  <div
    class="sm:tw-w-initial tw-fixed tw-left-0 tw-bottom-16 tw-flex tw-h-16 tw-w-full tw-items-center tw-justify-between tw-bg-white tw-px-4 tw-drop-shadow sm:tw-relative sm:tw-bottom-0 sm:tw-h-[unset] sm:tw-flex-1 sm:tw-px-0 sm:tw-drop-shadow-none"
  >
    <v-btn icon><v-icon>mdi-chevron-left</v-icon></v-btn>
    <div>Showing Google Calendar for week of 7/2</div>
    <v-btn icon><v-icon>mdi-chevron-right</v-icon></v-btn>
  </div>
</template>

<script>
import { isPhone } from "@/utils"

export default {
  name: "GCalWeekSelector",

  data() {
    return {}
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },
}
</script>
