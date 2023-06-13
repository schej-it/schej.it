<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-mb-12 tw-mt-5">
    <div class="tw-p-4 tw-flex tw-flex-col tw-gap-10">
      <div class="tw-text-4xl lg:tw-text-5xl">Settings</div>
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div class="tw-text-2xl">Calendar Access</div>
        <div class="tw-flex tw-flex-row tw-gap-28">
          <div class="tw-text-dark-gray">
            We do not store your calendar data anywhere on our servers, and we
            only fetch your calendar events for the time frame you specify in
            order to display your calendar events while you fill out your
            availability.
          </div>
          <v-btn outlined class="tw-text-dark-green"
            >Enable calendar access</v-btn
          >
        </div>
      </div>
      <div class="tw-flex tw-flex-col tw-gap-5">
        <div class="tw-text-2xl">Permissions</div>
        <div class="tw-flex tw-flex-col tw-border-gray tw-border-t-[1px] tw-border-l-[1px]">
          <div class="tw-flex tw-flex-row tw-w-full tw-border-gray tw-border-b-[1px]">
            <div v-for="h in heading" class="tw-w-1/3 tw-border-gray tw-border-r-[1px] tw-p-4 tw-text-xl">{{ h }}</div>
          </div>

          <div v-for="c in content" class="tw-flex tw-flex-row tw-w-full tw-border-gray tw-border-b-[1px]">
            <div v-for="text in c" class="tw-w-1/3 tw-border-gray tw-border-r-[1px] tw-p-4">{{ text }}</div>
          </div>
        </div>
      </div>
      <div class="tw-flex tw-flex-row tw-justify-center tw-mt-28"><v-btn outlined class="tw-text-bright-red"
            >Delete account</v-btn
          ></div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from "vuex"

export default {
  name: "Settings",

  components: {},

  data: () => ({
    dialog: false,
    heading: ["Permission", "Purpose", "Requested When"],
    content: [["View all calendars subscribed to", "Allows us to display calendar events on all your calendars instead of just your primary calendar", "User tries to input availability automatically with Google Calendar"], ["View all calendar events", "Allows us to view the names/times of your calendar events", "User tries to input availability automatically with Google Calendar"]]
  }),

  computed: {
    ...mapState([]),
  },

  methods: {
    ...mapActions([]),
  },

  created() {},
}
</script>
