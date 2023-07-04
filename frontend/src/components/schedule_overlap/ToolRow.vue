<template>
  <div>
    <div
      class="tw-min-h-[5rem] tw-flex-1 tw-flex tw-items-center sm:tw-mt-0 tw-text-sm tw-justify-center sm:tw-justify-between"
    >
      <div
        :class="
          state === states.EDIT_AVAILABILITY
            ? 'tw-justify-center'
            : 'tw-justify-between'
        "
        class="tw-flex-wrap tw-flex tw-py-4 tw-gap-y-4 tw-gap-x-4 sm:tw-gap-x-8 sm:tw-justify-start tw-flex-1"
      >
        <!-- Select timezone -->
        <TimezoneSelector
          :value="curTimezone"
          @input="(val) => $emit('update:curTimezone', val)"
          :timezones="Object.keys(timezoneMap)"
        />

        <template v-if="state !== states.EDIT_AVAILABILITY">
          <div class="tw-flex tw-justify-center tw-items-center tw-gap-2">
            <div>Show best times</div>
            <v-switch
              id="show-best-times-toggle"
              class="tw-mt-0 tw-py-2.5"
              :input-value="showBestTimes"
              @change="updateShowBestTimes"
              color="#219653"
              hide-details
            />
          </div>
        </template>
        <template v-else-if="isWeekly && !isPhone">
          <GCalWeekSelector />
        </template>
      </div>

      <div
        v-if="isOwner && state !== states.EDIT_AVAILABILITY"
        style="width: 180.16px"
        class="tw-hidden sm:tw-block"
      >
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
          >
            Schedule
          </v-btn>
        </template>
      </div>
    </div>

    <!-- GCal week selector when user is using phone view -->
    <template v-if="isPhone && isWeekly && state === states.EDIT_AVAILABILITY">
      <div class="tw-text-sm tw-h-16">
        <GCalWeekSelector />
      </div>
    </template>
  </div>
</template>

<script>
import TimezoneSelector from "./TimezoneSelector.vue"
import GCalWeekSelector from "./GCalWeekSelector.vue"
import { isPhone } from "@/utils"

export default {
  name: "ToolRow",

  props: {
    state: { type: String, required: true },
    states: { type: Object, required: true },
    curTimezone: { type: String, required: true },
    timezoneMap: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    isOwner: { type: Boolean, required: true },
    isWeekly: { type: Boolean, required: true },
  },

  components: {
    TimezoneSelector,
    GCalWeekSelector,
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },

  methods: {
    updateShowBestTimes(val) {
      this.$emit("update:showBestTimes", !!val)
      this.$emit("onShowBestTimesChange", !!val)
    },
  },
}
</script>
