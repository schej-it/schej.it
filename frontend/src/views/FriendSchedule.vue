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
    },
  }),
}
</script>