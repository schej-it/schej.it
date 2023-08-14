<template></template>

<script>
import { get, post } from "@/utils"
import { mapMutations } from "vuex"
import { authTypes } from "@/constants"

export default {
  name: "Auth",

  methods: {
    ...mapMutations(["setAuthUser"]),
  },

  async created() {
    let { error, code, state } = this.$route.query
    if (error) this.$router.replace({ name: "home" })

    if (state) state = JSON.parse(state)

    // Sign in and set auth user
    try {
      if (process.env.NODE_ENV === "development")
        console.log({ code, timezoneOffset: new Date().getTimezoneOffset() })

      if (state?.type === authTypes.ADD_CALENDAR_ACCOUNT || state?.type === authTypes.ADD_CALENDAR_ACCOUNT_FROM_EDIT) {
        await post("/user/add-calendar-account", { code })
      } else {
        await post("/auth/sign-in", {
          code,
          timezoneOffset: new Date().getTimezoneOffset(),
        })
        const authUser = await get("/user/profile")
        this.setAuthUser(authUser)
      }

      // Redirect to the correct place based on "state", otherwise, just redirect to home
      if (state) {
        let authUser
        switch (state.type) {
          case authTypes.EVENT_ADD_AVAILABILITY:
            this.$router.replace({
              name: "event",
              params: { eventId: state.eventId, fromSignIn: true },
            })
            break
          case authTypes.EVENT_SIGN_IN:
            this.$router.replace({
              name: "event",
              params: { eventId: state.eventId },
            })
            break
          case authTypes.ADD_CALENDAR_ACCOUNT:
            this.$router.replace({
              name: "settings",
            })
            authUser = await get("/user/profile")
            this.setAuthUser(authUser)
            break
          case authTypes.ADD_CALENDAR_ACCOUNT_FROM_EDIT:
            this.$router.replace({
              name: "event",
              params: { eventId: state.eventId, fromSignIn: true },
            })
            authUser = await get("/user/profile")
            this.setAuthUser(authUser)
            break
        }
      } else {
        this.$router.replace({ name: "home" })
      }
    } catch (err) {
      console.error(err)
    }
  },
}
</script>
