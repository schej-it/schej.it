<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    :width="400"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title>Sign in</v-card-title>
      <v-card-text class="tw-flex tw-flex-col tw-items-center">
        <!-- TODO: Add Sign in with Outlook -->
        <SignInGoogleBtn
          class="tw-mb-4"
          @click="signIn(calendarTypes.GOOGLE)"
        />
        <div class="tw-text-center tw-text-xs">
          By continuing, you agree to our
          <router-link class="tw-text-blue" :to="{ name: 'privacy-policy' }"
            >privacy policy</router-link
          >
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import SignInGoogleBtn from "@/components/SignInGoogleBtn.vue"
import { calendarTypes } from "@/constants"

export default {
  name: "SignInDialog",
  components: {
    SignInGoogleBtn,
  },
  props: {
    value: { type: Boolean, required: true },
  },
  data() {
    return {
      calendarTypes,
    }
  },
  methods: {
    signIn(provider) {
      this.$emit("signIn", provider)
      // this.$emit("input", false) // Close dialog after emitting
    },
  },
}
</script>
