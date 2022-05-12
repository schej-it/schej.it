<template>
  <div>
    <v-container class="pa-5">
      <div class="tw-justify-between tw-flex tw-mb-3">
        <div class="headline font-weight-bold">My friends</div>
        <div class="headline tw-text-gray">
          Requests
          <v-chip color="green" text-color="white">3</v-chip>
        </div>
      </div>

      <v-text-field label="Search" outlined dense></v-text-field>

      <FriendItem :friend="exampleFriend"></FriendItem>
    </v-container>

    <v-scale-transition appear origin="center">
      <v-btn
        fab
        absolute
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-4"
      >
        <v-icon>mdi-account-plus</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
import FriendItem from '@/components/FriendItem'

export default {
  name: 'Friends',

  components: {
    FriendItem,
  },

  data: () => ({
    dialog: false,
    exampleFriend: {
      name: 'Tony Xin',
      status: 'free',
    },
  }),
}
</script>
