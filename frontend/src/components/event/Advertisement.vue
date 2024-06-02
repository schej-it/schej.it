<template>
  <div v-if="showAd" class="tw-flex tw-cursor-pointer" @click="navigateToAd">
    <v-img
      alt="tomotime ad"
      :src="
        isPhone
          ? require('@/assets/ads/tomotime_mobile.png')
          : require('@/assets/ads/tomotime.png')
      "
      width="0"
      transition="fade-transition"
      class="tw-relative tw-shadow-md sm:tw-shadow-none"
      ><div
        class="tw-absolute tw-left-0 tw-top-0 tw-ml-0 tw-mt-0 tw-flex tw-h-8 tw-w-8 tw-items-center tw-justify-center tw-rounded-br-lg tw-bg-gray tw-bg-opacity-60 tw-text-xs sm:tw-ml-[6px] sm:tw-mt-[5px] sm:tw-h-10 sm:tw-w-10 sm:tw-rounded-tl-lg"
      >
        AD
      </div></v-img
    >
  </div>
</template>

<script>
import { mapState } from "vuex"
import { isPhone } from "@/utils"
import { get } from "@/utils"
import { guestUserId } from "@/constants"

export default {
  name: "Advertisement",

  props: {
    ownerId: { type: String, default: "" },
  },

  data: () => ({
    link: "https://tomotime.app",
    eduOnly: true,
    owner: null,
  }),

  async mounted() {
    if (this.ownerId !== guestUserId)
      this.owner = await get(`/users/${this.ownerId}`)
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    navigateToAd() {
      this.$posthog?.capture("Clicked ad", {
        link: this.link,
      })
      window.open(this.link, "_blank")
    },
    isEduEmail(email) {
      const split = email.split(".")
      return split[split.length - 1] === "edu"
    },
  },

  watch: {},

  computed: {
    ...mapState(["authUser"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showAd() {
      if (this.eduOnly) {
        return (
          (this.authUser && this.isEduEmail(this.authUser.email)) ||
          (this.owner && this.isEduEmail(this.owner.email))
        )
      } else {
        return true
      }
    },
  },

  components: {},
}
</script>
