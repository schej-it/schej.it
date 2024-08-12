<template>
  <v-card>
    <v-card-title>Choose a calendar provider</v-card-title>
    <v-card-text>
      <div class="tw-flex tw-flex-col tw-gap-2">
        <v-btn block>
          <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
            <v-img
              class="tw-flex-initial"
              width="20"
              height="20"
              src="@/assets/google_logo.svg"
            />
            <v-spacer />
            Google Calendar
            <v-spacer />
          </div>
        </v-btn>
        <v-btn block>
          <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
            <v-img
              class="tw-flex-initial"
              width="20"
              height="20"
              src="@/assets/apple_logo.svg"
            />
            <v-spacer />
            Apple Calendar
            <v-spacer />
          </div>
        </v-btn>
        <v-btn block>
          <div class="tw-flex tw-w-full tw-items-center tw-gap-2">
            <v-img
              class="tw-flex-initial"
              width="20"
              height="20"
              src="@/assets/outlook_logo.svg"
            />
            <v-spacer />
            Outlook Calendar
            <v-spacer />
          </div>
        </v-btn>
      </div>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: "CalendarTypeSelector",
}
</script>
