<template>
  <div v-if="showAd" class="tw-relative tw-cursor-pointer" @click="navigateToAd">
    <v-img
      alt="tomotime ad"
      src="@/assets/ads/tomotime.png"
      transition="fade-transition"
      class="tw-max-w-[60rem]"
    />
    <div class="tw-bg-gray tw-bg-opacity-50 tw-w-8 tw-h-8 sm:tw-w-10 sm:tw-h-10 tw-absolute tw-left-0 tw-top-0 tw-flex tw-justify-center tw-items-center tw-rounded-tl-lg tw-rounded-br-lg tw-ml-0 tw-mt-0 sm:tw-mt-2 sm:tw-ml-2 tw-text-xs ">
      AD
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from "vuex"
import {
  isPhone,
} from "@/utils"

export default {
  name: "Advertisement",

  props: {
  },

  data: () => ({
    link: "https://tomotime.app",
    eduOnly: true,
  }),

  methods: {
    navigateToAd() {
      this.$posthog?.capture("Clicked ad", {
        link: this.link,
      })
      window.open(this.link, "_blank")
    }
  },

  watch: {
    
  },

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

  components: {  },
}
</script>
