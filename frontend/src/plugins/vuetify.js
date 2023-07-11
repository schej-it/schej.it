import Vue from "vue"
import Vuetify from "vuetify/lib"
import tailwind from "../../tailwind.config"

Vue.use(Vuetify)

export default new Vuetify({
  theme: {
    themes: {
      light: {
        primary: tailwind.theme.colors.green,
      },
    },
  },
})
