<template>
  <div
    v-if="calendarPermissionGranted"
    class="sm:tw-w-initial tw-fixed tw-left-0 tw-bottom-16 tw-flex tw-h-16 tw-w-full tw-items-center tw-justify-between tw-bg-white tw-px-2 tw-drop-shadow sm:tw-relative sm:tw-bottom-0 sm:tw-h-[unset] sm:tw-flex-1 sm:tw-px-0 sm:tw-drop-shadow-none"
  >
    <v-btn @click="prevWeek" icon><v-icon>mdi-chevron-left</v-icon></v-btn>
    <div class="tw-text-center sm:tw-w-72">
      Showing Google Calendar for week of {{ weekText }}
    </div>
    <v-btn @click="nextWeek" icon><v-icon>mdi-chevron-right</v-icon></v-btn>
  </div>
</template>

<script>
import { isPhone } from "@/utils"
import dayjs from "dayjs"

export default {
  name: "GCalWeekSelector",

  props: {
    calendarPermissionGranted: { type: Boolean, required: true },
    weekOffset: { type: Number, required: true },
  },

  data() {
    return {}
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    weekText() {
      const date = new Date()
      // Set date to the Sunday of the current week
      date.setDate(date.getDate() - date.getDay())
      // Change date by the weekoffset
      date.setDate(date.getDate() + 7 * this.weekOffset)

      return dayjs(date).format("M/D")
    },
  },

  methods: {
    nextWeek() {
      this.$emit("update:weekOffset", this.weekOffset + 1)
    },
    prevWeek() {
      this.$emit("update:weekOffset", this.weekOffset - 1)
    },
  },
}
</script>
