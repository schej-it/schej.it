<template>
  
</template>

<script>
import { get, post } from '@/utils'
import { mapMutations } from 'vuex'

export default {
  name: 'Auth',

  methods: {
    ...mapMutations([ 'setAuthUser' ]),
  },

  async created() {
    let { error, code, state } = this.$route.query
    if (error) return
    
    if (state) state = JSON.parse(state)

    // Sign in and set auth user
    await post('/auth/sign-in', { code })
    const authUser = await get('/user/profile')
    this.setAuthUser(authUser)

    // Redirect to the correct place based on "state", otherwise, just redirect to home
    if (state) {
      switch (state.type) {
        case 'join':
          this.$router.replace({ name: 'event', params: { eventId: state.eventId } })
          break
      }
    } else {
      this.$router.replace({ name: "home" })
    }
  },
}
</script>