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
    })
  },
}
</script>