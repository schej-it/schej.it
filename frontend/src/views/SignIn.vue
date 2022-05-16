<template>
  <div class="tw-flex tw-flex-col tw-h-full">
    <div class="tw-px-16 tw-flex-1 tw-bg-green tw-flex tw-flex-col tw-justify-center tw-items-center">
      <v-img
        class="sm:tw-max-w-md tw-flex-initial"
        alt="Schej.it Logo"
        contain
        src="@/assets/logo_dark.svg"
        transition="scale-transition"
      />
      <div
        class="tw-font-light tw-text-white tw-text-lg"
      >Scheduling made simple</div>
    </div>
    <div class="tw-flex tw-justify-center">
      <v-btn 
        class="tw-my-16 tw-bg-blue" 
        dark
        @click="signIn"
      >Sign in with Google</v-btn>
    </div>
  </div>
</template>

<script>
import { signInGoogle } from '@/utils'
export default {
  name: 'SignIn',

  methods: {
    signIn() {
      signInGoogle()
    },
  },
}
</script>