<template>
  <div
    v-if="show && !isPhone"
    class="tw-relative tw-bg-[#1A1A1E] tw-px-8 tw-py-3 tw-text-center tw-text-sm tw-text-white"
  >
    <div
      class="tw-m-auto tw-flex tw-max-w-6xl tw-flex-col tw-items-center tw-justify-center sm:tw-flex-row"
    >
      <span>
        Join our Discord community for updates, feedback, and support!
      </span>
      <v-btn
        :href="discordUrl"
        target="_blank"
        outlined
        color="white"
        class="tw-mt-3 tw-flex-shrink-0 sm:tw-ml-4 sm:tw-mt-0"
        small
        @click="trackDiscordClick"
      >
        Join Discord
      </v-btn>
    </div>
    <v-btn
      icon
      small
      @click="dismiss"
      class="tw-absolute tw-right-2 tw-top-1/2 -tw-translate-y-1/2"
    >
      <v-icon color="white">mdi-close</v-icon>
    </v-btn>
  </div>
</template>

<script>
import { isPhone } from "@/utils"
export default {
  name: "DiscordBanner",

  data() {
    return {
      discordUrl: "https://discord.gg/v6raNqYxx3",
      show: false,
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    localStorageKey() {
      // Use a versioned key to force re-showing banner if message changes
      return `discordBannerDismissed_v1`
    },
  },

  methods: {
    dismiss() {
      this.show = false
      localStorage.setItem(this.localStorageKey, "true")
      this.$posthog?.capture("discord_banner_dismissed")
    },
    trackDiscordClick() {
      this.$posthog?.capture("discord_banner_clicked", {
        discordUrl: this.discordUrl,
      })
    },
  },

  watch: {
    $route: {
      immediate: true,
      handler() {
        const showOnRoute = this.$route.name === "landing"
        const userHasDismissed =
          localStorage.getItem(this.localStorageKey) === "true"

        this.show = !userHasDismissed && showOnRoute
      },
    },
  },
}
</script>
