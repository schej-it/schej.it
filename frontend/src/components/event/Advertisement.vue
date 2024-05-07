<template>
  <div v-if="showAd" class="tw-cursor-pointer" @click="navigateToAd">
    <v-img
      alt="tomotime ad"
      :src="
        isPhone
          ? require('@/assets/ads/tomotime_mobile.png')
          : require('@/assets/ads/tomotime.png')
      "
      transition="fade-transition"
      class="tw-relative tw-max-w-[44rem] tw-shadow-md sm:tw-shadow-none"
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

export default {
  name: "Advertisement",

  props: {},

  data: () => ({
    link: "https://tomotime.app",
    eduOnly: true,
  }),

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
  },

  watch: {},

  computed: {
    ...mapState(["authUser"]),
    isPhone() {
      return isPhone(this.$vuetify)
    },
    showAd() {
      if (this.eduOnly && this.authUser) {
        const split = this.authUser.email.split(".")
        if (split[split.length - 1] === "edu") {
          return true
        }
      }
      return false
    },
  },

  components: {},
}
</script>
