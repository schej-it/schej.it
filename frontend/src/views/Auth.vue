<template>
  
</template>

<script>
import { post } from '@/utils'

export default {
  name: 'Auth',

  mounted() {
    const { error, code, state } = this.$route.query
    if (error) return
    
    post('/auth/sign-in', { code }).then(data => {
      console.log(data)

      /*return fetch(`https://oauth2.googleapis.com/token?client_id=523323684219-jfakov2bgsleeb6den4ktpohq4lcnae2.apps.googleusercontent.com&client_secret=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&code=${code}&grant_type=authorization_code&redirect_uri=${encodeURI('http://localhost:8080/auth')}`, {
        method: 'POST',
      }).then(res => res.json())*/
    })
  },
}
</script>