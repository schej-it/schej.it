<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    :fullscreen="isPhone"
    :hide-overlay="isPhone"
    content-class="tw-max-w-xl"
    :transition="isPhone ? `dialog-bottom-transition` : `dialog-transition`"
  >
    <v-card tile class="tw-flex tw-flex-col">
      <v-card-title class="tw-flex">
        <div>New Event</div>
        <v-spacer />
        <v-btn icon @click="$emit('input', false)"
          ><v-icon>mdi-close</v-icon></v-btn
        >
      </v-card-title>
      <v-card-text class="tw-space-y-4 tw-flex tw-flex-col tw-flex-1">
        <v-text-field
          ref="name-field"
          v-model="name"
          autofocus
          :disabled="loading"
          class="tw-text-white tw-flex-initial"
          placeholder="Name of event..."
          hide-details
          @keyup.enter="blurNameField"
        />

        <div>
          <div class="tw-flex tw-space-x-2 tw-items-baseline tw-justify-center">
            <v-select
              v-model="startTime"
              :disabled="loading"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              outlined
              hide-details
            ></v-select>
            <div>to</div>
            <v-select
              v-model="endTime"
              :disabled="loading"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              outlined
              hide-details
            ></v-select>
          </div>
        </div>

        <div>
          <div
            class="tw-text-lg tw-text-black tw-text-center tw-font-medium tw-mt-6 tw-mb-2"
          >
            What dates would you like to meet?
          </div>
          <div class="tw-flex tw-flex-col tw-justify-center tw-items-center">
            <v-date-picker
              v-model="selectedDays"
              no-title
              multiple
              color="primary"
              elevation="2"
              :show-current="false"
              class="tw-min-w-full sm:tw-min-w-0 tw-border-0"
            />
          </div>
          
        </div>

        <v-spacer />

        <v-btn
          :loading="loading"
          :dark="formComplete"
          class="tw-bg-green"
          :disabled="!formComplete"
          @click="submit"
          >{{ editEvent ? 'Edit' : 'Create' }}</v-btn
        >
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone, post, put, utcTimeToLocalTime } from "@/utils";

export default {
  name: "NewEventDialog",

  emits: ["input"],

  props: {
    value: { type: Boolean, required: true },
    event: { type: Object, },
    editEvent: { type: Boolean, default: false },
  },

  data: () => ({
    name: "",
    startTime: 9,
    endTime: 17,
    loading: false,
    selectedDays: [],
  }),

  created() {
    if (this.event) {
      this.name = this.event.name
      this.startTime = utcTimeToLocalTime(this.event.startTime)
      this.endTime = utcTimeToLocalTime(this.event.endTime)
      
      // Format dates for the local timezone
      const paddedStartTime = String(this.event.startTime).padStart(2, '0');
      const selectedDays = []
      for (const date of this.event.dates) {
        const localDate = new Date(`${date}T${paddedStartTime}:00:00Z`);
        const dateString = new Date(localDate.getTime() - new Date().getTimezoneOffset() * 1000 * 60).toISOString().substring(0, 10);
        selectedDays.push(dateString);
      }
      this.selectedDays = selectedDays
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify);
    },
    formComplete() {
      return (
        this.name.length > 0 &&
        this.selectedDays.length > 0 &&
        (this.startTime < this.endTime || (this.endTime === 0 && this.startTime != 0))
      );
    },
    times() {
      const times = [];

      for (let h = 1; h < 12; ++h) {
        times.push({ text: `${h} am`, value: h });
      }
      for (let h = 0; h < 12; ++h) {
        times.push({ text: `${h == 0 ? 12 : h} pm`, value: h + 12 });
      }
      times.push({ text: "12 am", value: 0 });

      return times;
    },
    minCalenderDate() {
      let today = new Date();
      let dd = String(today.getDate()).padStart(2, '0');
      let mm = String(today.getMonth() + 1).padStart(2, '0');
      let yyyy = today.getFullYear();
      
      return yyyy + '-' + mm + '-' + dd;
    }
  },

  methods: {
    blurNameField() {
      this.$refs["name-field"].blur();
    },
    reset() {
      this.name = "";
      this.startTime = 9;
      this.endTime = 17;
      this.selectedDays = [];
    },
    submit() {
      // Calculate UTC dates array and UTC start/end times
      const utcDates = [];
      const paddedStartTime = String(this.startTime).padStart(2, '0');
      const timezoneOffset = new Date().getTimezoneOffset();
      let utcStartTime = (this.startTime + timezoneOffset/60) % 24;
      if (utcStartTime < 0) utcStartTime += 24;
      let utcEndTime = (this.endTime + timezoneOffset/60) % 24;
      if (utcEndTime < 0) utcStartTime += 24;
      
      for (const date of this.selectedDays) {
        const utcDate = new Date(`${date}T${paddedStartTime}:00:00`);
        const dateString = utcDate.toISOString().substring(0, 10);
        utcDates.push(dateString);
      }
      utcDates.sort()

      // Create new event on backend
      this.loading = true;
      if (!this.editEvent) {
        post("/events", {
          name: this.name,
          startTime: utcStartTime,
          endTime: utcEndTime,
          dates: utcDates,
        }).then(({ eventId }) => {
          this.$router.push({ name: "event", params: { eventId } });
          this.loading = false;
        });
      } else {
        if (this.event) {
          put(`/events/${this.event._id}`, {
            name: this.name,
            startTime: utcStartTime,
            endTime: utcEndTime,
            dates: utcDates,
          }).then(() => {
            window.location.reload()
          });
        }
      }
    },
  },
};
</script>
