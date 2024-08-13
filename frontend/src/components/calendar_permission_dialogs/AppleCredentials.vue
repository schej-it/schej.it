<template>
  <div class="tw-flex tw-flex-col tw-gap-6">
    <div class="tw-flex tw-flex-col tw-gap-3">
      <div
        class="tw-text-md tw-flex tw-flex-row tw-items-center tw-justify-start tw-gap-2 tw-font-medium"
      >
        Connect your
        <div class="tw-flex tw-items-center tw-gap-1">
          <v-img
            class="-tw-mt-1 tw-flex-initial"
            width="15"
            height="15"
            src="@/assets/apple_logo.svg"
          />
          <span class="tw-font-medium">Apple Calendar</span>
        </div>
      </div>
      <div class="tw-flex tw-flex-col tw-gap-2">
        <div class="tw-text-sm tw-text-very-dark-gray">
          Generate an app specific password to use with Schej at
          <a
            href="https://appleid.apple.com/account/manage"
            target="_blank"
            rel="noopener noreferrer"
            >https://appleid.apple.com/account/manage</a
          >. Copy and paste the generated app password below.
        </div>
        <div class="tw-text-sm tw-text-very-dark-gray">
          Your credentials will be stored and encrypted.
        </div>
      </div>
    </div>
    <div class="tw-flex tw-flex-col tw-gap-3">
      <v-text-field solo placeholder="Apple ID" hide-details />
      <v-text-field solo placeholder="App password" hide-details />
      <div class="tw-flex tw-items-center tw-gap-2">
        <v-btn text class="tw-grow" @click="$emit('back')">Back</v-btn>
        <v-btn color="primary" class="tw-grow" @click="$emit('continue')"
          >Submit</v-btn
        >
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AppleCredentials",
}
</script>
