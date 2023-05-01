<template>
  <div class="tw-flex">
    <div
      class="tw-flex-1 tw-flex tw-items-center tw-mt-4 sm:tw-mt-0 tw-text-sm tw-justify-center sm:tw-justify-between"
    >
      <div
        class="tw-flex tw-gap-4 sm:tw-gap-8 tw-flex-row tw-justify-between sm:tw-justify-start tw-flex-1"
      >
        <!-- Select timezone -->
        <TimezoneSelector
          :value="curTimezone"
          @input="(val) => $emit('update:curTimezone', val)"
          :timezones="Object.keys(timezoneMap)"
        />

        <div class="tw-flex tw-justify-center tw-items-center tw-gap-1">
          <div>Show best times</div>
          <v-switch
            class="-tw-mb-1"
            :input-value="showBestTimes"
            @change="updateShowBestTimes"
            color="#219653"
          />
        </div>
      </div>

      <div v-if="isOwner" style="width: 180.16px" class="tw-hidden sm:tw-block">
        <template v-if="state !== states.SCHEDULE_EVENT">
          <v-btn
            outlined
            class="tw-text-green tw-w-full"
            @click="(e) => $emit('scheduleEvent', e)"
          >
            <span class="tw-mr-2">Schedule event</span>
            <v-img
              src="@/assets/gcal_logo.png"
              class="tw-flex-none"
              height="20"
              width="20"
            />
          </v-btn>
        </template>
        <template v-else>
          <v-btn
            outlined
            class="tw-text-red tw-mr-1"
            @click="(e) => $emit('cancelScheduleEvent', e)"
          >
            Cancel
          </v-btn>
          <v-btn
            color="primary"
            @click="(e) => $emit('confirmScheduleEvent', e)"
            :disabled="!curScheduledEvent"
          >
            Schedule
          </v-btn>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import TimezoneSelector from "./TimezoneSelector.vue"

export default {
  name: "ToolRow",

  props: {
    state: { type: String, required: true },
    states: { type: Object, required: true },
    curTimezone: { type: String, required: true },
    timezoneMap: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    isOwner: { type: Boolean, required: true },
    curScheduledEvent: { type: Object | null, required: true },
  },

  components: {
    TimezoneSelector,
  },

  methods: {
    updateShowBestTimes(val) {
      this.$emit("update:showBestTimes", !!val)
      this.$emit("onShowBestTimesChange", !!val)
    },
  },
}
</script>
