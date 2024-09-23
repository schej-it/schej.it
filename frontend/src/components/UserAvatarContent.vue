<template>
  <v-avatar v-if="user" :size="size">
    <img v-if="user.picture" :src="user.picture" referrerpolicy="no-referrer" />
    <v-icon
      class="-tw-mt-1"
      :size="size"
      v-else-if="user.calendarType === calendarTypes.APPLE"
    >
      mdi-apple
    </v-icon>
    <v-icon
      :size="size"
      v-else-if="user.calendarType === calendarTypes.OUTLOOK"
    >
      mdi-microsoft-outlook
    </v-icon>
    <div
      v-else
      :class="`tw-flex tw-size-full tw-items-center tw-justify-center tw-bg-green tw-text-${textSize} tw-text-white`"
    >
      {{ user.firstName?.charAt(0) ?? user.email?.charAt(0) ?? "" }}
    </div>
  </v-avatar>
</template>

<script>
import { calendarTypes } from "@/constants"

export default {
  name: "UserAvatarContent",
  props: {
    user: Object,
    size: { type: Number, default: 48 },
  },

  computed: {
    calendarTypes() {
      return calendarTypes
    },
    textSize() {
      return this.size <= 24 ? "sm" : "lg"
    },
  },
}
</script>
