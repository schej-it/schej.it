<template>
  
</template>

<script>
import { get, post } from '@/utils'

export default {
  name: 'Auth',

  mounted() {
    let { error, code, state } = this.$route.query
    if (error) return
    
    if (state) state = JSON.parse(state)

    post('/auth/sign-in', { code }).then(data => {
      console.log(data)

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
    })
  },
}
</script>