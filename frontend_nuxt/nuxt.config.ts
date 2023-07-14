// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: [
    "@nuxtjs/tailwindcss",
    "@invictus.codes/nuxt-vuetify",
    "@pinia/nuxt",
  ],
  //@ts-ignore
  vuetify: {
    moduleOptions: {
      treeshaking: true,
      autoImport: true,
    },
  },
})
