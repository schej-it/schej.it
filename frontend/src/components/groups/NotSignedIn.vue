<template>
  <v-fade-transition>
    <div
      v-if="loaded"
      class="tw-flex tw-h-full tw-flex-col tw-items-center tw-justify-center tw-p-2"
    >
      <div class="tw-mb-8 tw-flex tw-max-w-[26rem] tw-flex-col tw-items-center">
        <v-img
          :src="owner.picture"
          :width="90"
          :height="90"
          class="tw-mb-4 tw-rounded-full tw-text-center"
        />
        <h1 class="tw-mb-2 tw-text-center tw-text-xl tw-font-medium">
          {{ owner.firstName ?? "" }} invited you to join "{{ event.name }}"
        </h1>
        <div class="tw-text-center tw-text-dark-gray">
          Join the group now to share your real-time <br v-if="!isPhone" />
          calendar availability with each other!
        </div>
      </div>
      <v-btn @click="join" color="primary" class="tw-mb-8"
        >Join with Google Calendar</v-btn
      >
      <div class="tw-text-center tw-text-dark-gray">
        Already have a Schej account? <a @click="signIn">Sign in instead</a>
      </div>
    </div>
  </v-fade-transition>
</template>

<script>
import { get, isPhone } from "@/utils"

export default {
  name: "NotSignedIn",

  props: {
    event: { type: Object, required: true },
  },

  components: {},

  data() {
    return {
      owner: {},
      loaded: false,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    join() {},
    signIn() {},
  },

  async created() {
    this.owner = await get(`/users/${this.event.ownerId}`)
    this.loaded = true
  },
}
</script>
