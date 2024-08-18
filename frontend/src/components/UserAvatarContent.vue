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
  },
}
</script>
