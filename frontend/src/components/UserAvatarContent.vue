<template>
  <v-avatar v-if="user" :size="size">
    <img :src="user.picture" referrerpolicy="no-referrer" />
  </v-avatar>
</template>

<script>
export default {
  name: 'UserAvatarContent',
  props: {
    user: Object,
    size: { type: Number, default: 48 },
  },
}
</script>
