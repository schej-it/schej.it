<template>
  <div class="tw-p-4">
    <FriendItem 
      :friend="userItem"
    />

    <TestCalendar noEventNames />

    <v-scale-transition appear origin="center">
      <v-btn 
        fab
        fixed
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-16"
      >
        <v-icon>mdi-share</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
import FriendItem from '@/components/FriendItem'
import TestCalendar from '@/components/TestCalendar'
import { mapState } from 'vuex'

export default {
  name: 'Schedule',

  components: {
    FriendItem,
    TestCalendar,
  },

  computed: {
    ...mapState([ 'authUser' ]),
    userItem() {
      const { firstName, lastName, picture } = this.authUser
      return {
        name: firstName + ' ' + lastName,
        picture,
        status: 'free',
      }
    },
  },
}
</script>