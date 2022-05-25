<template>
  <v-avatar v-if="user">
    <img :src="user.picture" />
  </v-avatar>
</template>

<script>
export default {
  name: 'UserAvatarContent',
  props: {
    user: Object,
  },
}
</script>
