<template>
  <div>
    <div class="tw-bg-green tw-h-screen">
      <div class="tw-flex tw-flex-col tw-items-center">
        <div class="tw-mt-24 tw-mb-4 tw-text-3xl sm:tw-text-5xl lg:tw-text-7xl tw-font-medium tw-text-center tw-text-white">Scheduling made simple.</div>
        <v-btn class="tw-bg-blue" dark>Get started</v-btn>
      </div>
    </div>
    <div class="-tw-translate-y-1/2 -tw-mb-72 md:tw-grid md:tw-place-content-center">
      <LandingPageCalendar />
    </div>
    <div>hello there</div>
  </div>
</template>

<script>
import LandingPageCalendar from '@/components/LandingPageCalendar'

export default {
  name: 'Landing',

  components: {
    LandingPageCalendar,
  },
}
</script>