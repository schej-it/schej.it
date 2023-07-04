<template>
  <div
    class="tw-drop-shadow tw-fixed tw-left-0 tw-bottom-16 tw-bg-white tw-w-full tw-px-4 tw-h-16 sm:tw-w-initial sm:tw-relative sm:tw-bottom-0 sm:tw-drop-shadow-none sm:tw-px-0 sm:tw-h-[unset] sm:tw-flex-1 tw-flex tw-items-center tw-justify-between"
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
