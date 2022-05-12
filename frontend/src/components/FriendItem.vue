<template>
  <v-container
    class="tw-bg-light-gray tw-rounded-md tw-flex tw-text-black tw-justify-between tw-align-middle tw-py-2"
  >
    <div class="tw-flex">
      <div class="tw-mr-3">
        <v-avatar>
          <img src="https://cdn.vuetifyjs.com/images/john.jpg" alt="John" />
        </v-avatar>
      </div>
      <div>
        <div class="tw-font-medium">{{ this.friend.name }}</div>
        <div class="tw-text-sm">
          Currently
          <span
            v-if="this.friend.status == 'free'"
            class="tw-text-green tw-font-bold"
            >free</span
          ><span v-else>in {{ this.friend.status }}</span>
        </div>
      </div>
    </div>
    <div>
      <v-icon class="mt-2">mdi-chevron-right</v-icon>
    </div>
  </v-container>
</template>

<script>
export default {
  name: 'FriendItem',

  props: {
    friend: { type: Object, required: true },
  },

  data: () => ({}),

  computed: {},

  methods: {},
}
</script>
