<template>
  <v-btn
    :class="dark ? 'tw-bg-blue' : 'tw-bg-white'"
    class="tw-mb-2 tw-h-10 tw-p-0"
    @click="(e) => $emit('click', e)"
  >
    <div class="tw-flex tw-items-center">
      <div
        class="tw-ml-0.5 tw-grid tw-h-9 tw-w-9 tw-place-content-center tw-rounded tw-bg-white"
      >
        <v-img
          contain
          width="1.25rem"
          height="1.25rem"
          src="@/assets/google_logo.svg"
        ></v-img>
      </div>
      <div class="tw-mx-3" :class="dark ? 'tw-text-white' : 'tw-text-black'">
        {{ text }}
      </div>
    </div>
  </v-btn>
</template>

<script>
export default {
  name: "SignInGoogleBtn",

  emits: ["click"],

  props: {
    text: { type: String, default: "Continue with Google" },
    dark: { type: Boolean, default: false },
  },
}
</script>
