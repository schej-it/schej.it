<template>
  <v-container
    @click="$emit('click')"
    class="tw-bg-light-gray tw-rounded-md tw-flex tw-text-black tw-justify-between tw-align-middle tw-py-2"
  >
    <div class="tw-flex tw-mt-2">
      <div class="tw-mr-3">
        <v-avatar>
          <img
            v-if="!user.picture"
            src="https://cdn.vuetifyjs.com/images/john.jpg"
          />
          <img v-else :src="user.picture" />
        </v-avatar>
      </div>
      <div>
        <div class="tw-font-medium">{{ this.user.name }}</div>
        <div class="tw-text-sm">
          Currently
          <span
            v-if="this.user.status == 'free'"
            class="tw-text-green tw-font-bold"
            >free</span
          ><span v-else>
            in
            <span class="tw-text-light-blue tw-font-bold">
              {{ this.user.status }}
            </span>
          </span>
        </div>
      </div>
      
    </div>

    <div >
      <v-switch
        v-model="hideEventNames"
        inset
      ></v-switch>
    </div>
    
  </v-container>
</template>

<script>
export default {
  name: 'UserItem',

  props: {
    user: { type: Object, required: true },
  },

  data: () => ({
      hideEventNames: 0,
  }),

  computed: {
      
  },

  watch: {
      hideEventNames(val) {
          console.log(val)
      }
  },

  methods: {},
}
</script>
