<!-- Displays auth user's avatar, which displays a menu when clicked -->
<template>
  <v-menu
    v-if="authUser"
    offset-y
    :close-on-content-click="false"
  >
    <template v-slot:activator="{ on }">
      <v-btn icon x-large v-on="on">
        <v-avatar size="48">
          <UserAvatarContent :user="authUser" />
        </v-avatar>
      </v-btn>
    </template>
    <v-list class="py-0">
      <v-list-item>
        <v-list-item-title><strong>{{ `${authUser.firstName} ${authUser.lastName}` }}</strong></v-list-item-title>
      </v-list-item>
      <v-divider></v-divider>
      <v-list-item @click="signOut">
        <v-list-item-title class="red--text">Sign Out</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import UserAvatarContent from '@/components/UserAvatarContent'
import { mapState, mapMutations } from 'vuex'
import { get, post } from '@/utils'

export default {
  name: 'AuthUserMenu',

  components: {
    UserAvatarContent,
  },

  computed: {
    ...mapState([ 'authUser' ]),
  },

  methods: {
      ...mapMutations([ 'setAuthUser' ]),
      async signOut() {
          await post('/auth/sign-out')
          this.setAuthUser(null)
      }
  },
}
</script>