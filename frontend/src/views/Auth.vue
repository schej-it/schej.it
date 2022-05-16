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

      if (state) {
        switch (state.type) {
          case 'join':
            this.$router.replace({ name: 'event', params: { eventId: state.eventId } })
            break
        }
      }
    })
  },
}
</script>