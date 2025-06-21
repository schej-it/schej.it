import posthog from "posthog-js"

export default {
  install(Vue, options) {
    Vue.prototype.$posthog = posthog.init(process.env.VUE_APP_POSTHOG_API_KEY, {
      api_host: "https://e.timeful.app",
      capture_pageview: false,
      autocapture: false,
    })
  },
}
