<template>
  <v-container
    @click="$emit('click')"
    v-ripple
    class="tw-bg-light-gray tw-rounded-md tw-flex tw-text-black tw-justify-between tw-align-middle tw-py-2"
  >
    <div class="tw-flex">
      <div class="tw-mr-3">
        <v-avatar>
          <img v-if="!friend.picture" src="https://cdn.vuetifyjs.com/images/john.jpg" />
          <img v-else :src="friend.picture" />
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
          ><span v-else>
            in 
            <span class="tw-text-light-blue tw-font-bold">
              {{ this.friend.status }}
            </span>
          </span>
        </div>
      </div>
    </div>
    <div v-if="chevron">
      <v-icon class="mt-2">mdi-chevron-right</v-icon>
    </div>
  </v-container>
</template>

<script>
export default {
  name: 'FriendItem',

  props: {
    friend: { type: Object, required: true },
    chevron: { type: Boolean, default: false }, // Whether to show the chevron icon
  },

  data: () => ({}),

  computed: {},

  methods: {},
}
</script>
