<template>
  <div>
    <div class="tw-bg-green tw-h-screen">
      <div class="tw-pt-5 tw-px-5 tw-flex tw-justify-between">
        <v-img
          alt="Schej.it Logo"
          class="shrink tw-cursor-pointer"
          contain
          src="@/assets/logo_dark.svg"
          transition="scale-transition"
          width="120"
        />

        <v-btn dark text @click="$router.push('/sign-in')">Sign in</v-btn>
      </div>
      <div class="tw-flex tw-flex-col tw-items-center">
        <div
          class="tw-mt-24 tw-mb-4 tw-text-3xl sm:tw-text-5xl lg:tw-text-7xl tw-font-medium tw-text-center tw-text-white"
        >
          Scheduling made simple.
        </div>
        <v-btn rounded class="tw-bg-blue" dark @click="$router.push('/sign-in')"
          >Get started</v-btn
        >
      </div>
    </div>
    <div
      class="-tw-translate-y-1/2 -tw-mb-72 md:tw-grid md:tw-place-content-center"
    >
      <LandingPageCalendar />
    </div>
    <div class="tw-pt-20 tw-pb-56">
      <div class="tw-flex tw-flex-row tw-justify-center">
        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-10">
          <h1 class="tw-font-bold tw-text-2xl">Schedule meetings with ease</h1>
          <p class="tw-text-gray">
            Mark your availability alongside all your Google Calendar events.
            It’s like when2meet but with Google Calendar integration.
          </p>
        </div>
        <div class="tw-flex tw-flex tw-w-80">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://download.logo.wine/logo/Google_Calendar/Google_Calendar-Logo.wine.png"
            transition="scroll-y-reverse-transition"
            width="600"
          />
        </div>
      </div>

      <div class="tw-flex tw-flex-row tw-justify-center tw-mt-20">
        <div class="tw-flex tw-flex tw-w-80 tw-justify-center">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://cdn-icons-png.flaticon.com/512/3721/3721924.png"
            transition="scale-transition"
            width="180"
          />
        </div>

        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-7">
          <h1 class="tw-font-bold tw-text-2xl">Send your availability</h1>
          <p class="tw-text-gray">
            Send your friends a screenshot of the times you are available. No
            more hassle of manually typing out time ranges of when you're free.
          </p>
        </div>
      </div>

      <div class="tw-flex tw-flex-row tw-justify-center tw-mt-20">
        <div class="tw-flex tw-flex-col tw-w-80 tw-mt-7">
          <h1 class="tw-font-bold tw-text-2xl">Peep your friends’ schedules</h1>
          <p class="tw-text-gray">
            Add your friends to see what they’ve been up to. Control who can see
            what.
          </p>
        </div>
        <div class="tw-flex tw-flex tw-w-80 tw-justify-center">
          <v-img
            alt="Schej.it Logo"
            class="shrink tw-cursor-pointer"
            contain
            src="https://icon-library.com/images/friends-icon-transparent/friends-icon-transparent-1.jpg"
            transition="scale-transition"
            width="200"
          />
        </div>
      </div>
    </div>

    <div
      class="-tw-translate-y-1/2 -tw-mb-52 md:tw-grid md:tw-place-content-center"
    >
      <v-img
        alt="Schej.it Logo"
        contain
        src="@/assets/logo_square_512_512.png"
        width="200"
      />
    </div>

    <div class="tw-bg-green tw-h-96 tw-flex tw-justify-center tw-items-center">
      <v-btn rounded class="tw-bg-blue" dark @click="$router.push('/sign-in')"
        >Get started</v-btn
      >
    </div>
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
