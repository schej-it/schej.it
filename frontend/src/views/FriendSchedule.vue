<template>
  <div class="tw-p-4">
    <FriendItem :friend="exampleFriend" />

    <TestCalendar />
  </div>
</template>

<script>
import FriendItem from '@/components/FriendItem'
import TestCalendar from '@/components/TestCalendar'

export default {
  name: 'FriendSchedule',

  components: {
    FriendItem,
    TestCalendar,
  },

  data: () => ({
    exampleFriend: {
      name: 'Tony Xin',
      status: 'free',
      picture: 'https://media-exp1.licdn.com/dms/image/C5603AQHPyWrhwOLQ_g/profile-displayphoto-shrink_200_200/0/1632279309883?e=1658361600&v=beta&t=SF_4_hjgwSN4dlQt_EXTKVkNuuqaty5_HQRF_E7jVNo',
    },
  }),
}
</script>