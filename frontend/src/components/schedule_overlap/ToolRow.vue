<template>
  <div>
    <div
      class="tw-flex tw-min-h-[5rem] tw-flex-1 tw-items-center tw-justify-center tw-text-sm sm:tw-mt-0 sm:tw-justify-between"
    >
      <div
        :class="`tw-justify-${
          state === states.EDIT_AVAILABILITY ? 'center' : 'between'
        } 
        `"
        class="tw-flex tw-flex-1 tw-flex-wrap tw-gap-x-4 tw-gap-y-4 tw-py-4 sm:tw-justify-start sm:tw-gap-x-8"
      >
        <!-- Select timezone -->
        <TimezoneSelector
          :value="curTimezone"
          @input="(val) => $emit('update:curTimezone', val)"
        />

        <template v-if="state !== states.EDIT_AVAILABILITY">
          <div
            v-if="numResponses > 1"
            class="tw-flex tw-items-center tw-justify-center tw-gap-2"
          >
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
          <div
            v-if="isPhone"
            class="tw-flex tw-basis-full tw-items-center tw-gap-x-2"
          >
            Show
            <v-select
              :value="mobileNumDays"
              @input="$emit('update:mobileNumDays', $event)"
              :items="mobileNumDaysOptions"
              :menu-props="{ auto: true }"
              item-text="label"
              item-value="value"
              class="-tw-mt-px tw-flex-none tw-shrink tw-basis-24 tw-text-sm"
              dense
              hide-details
            />
            at a time
          </div>
        </template>
        <template v-else-if="isWeekly && !isPhone">
          <v-spacer />
          <div class="tw-min-w-fit">
            <GCalWeekSelector
              :calendar-permission-granted="calendarPermissionGranted"
              :week-offset="weekOffset"
              @update:weekOffset="(val) => $emit('update:weekOffset', val)"
            />
          </div>
        </template>
      </div>

      <div
        v-if="numResponses > 0 && state !== states.EDIT_AVAILABILITY"
        style="width: 181.5px"
        class="tw-hidden sm:tw-flex"
      >
        <template v-if="state !== states.SCHEDULE_EVENT">
          <v-btn
            outlined
            class="tw-w-full tw-text-green"
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
            class="tw-mr-1 tw-text-red"
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
      <div class="tw-h-16 tw-text-sm">
        <GCalWeekSelector
          :calendar-permission-granted="calendarPermissionGranted"
          :week-offset="weekOffset"
          @update:weekOffset="(val) => $emit('update:weekOffset', val)"
        />
      </div>
    </template>

    <!-- Instructions when user is using phone view -->
    <v-slide-y-reverse-transition>
      <template v-if="isPhone && hintText != '' && showHintText">
        <div
          :class="`tw-fixed tw-left-0 tw-bottom-${
            isWeekly &&
            calendarPermissionGranted &&
            state === states.EDIT_AVAILABILITY
              ? 32
              : 16
          } tw-z-10 tw-flex tw-w-full tw-items-center tw-justify-between tw-gap-1 tw-bg-light-gray tw-px-2 tw-py-2 tw-text-sm tw-text-dark-gray`"
        >
          <div :class="`tw-flex tw-gap-${hintText.length > 60 ? 2 : 1}`">
            <v-icon small>mdi-information-outline</v-icon>
            <div>
              {{ hintText }}
            </div>
          </div>
          <v-icon small @click="closeHint">mdi-close</v-icon>
        </div>
      </template>
    </v-slide-y-reverse-transition>

    <!-- force tailwind classes to compile -->
    <div class="tw-bottom-16 tw-bottom-32"></div>
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
    curTimezone: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    isWeekly: { type: Boolean, required: true },
    isGroup: { type: Boolean, required: true },
    calendarPermissionGranted: { type: Boolean, required: true },
    weekOffset: { type: Number, required: true },
    numResponses: { type: Number, required: true },
    mobileNumDays: { type: Number, default: 3 }, // The number of days to show at a time on mobile
  },

  components: {
    TimezoneSelector,
    GCalWeekSelector,
  },

  data: () => ({
    hintTextState: true,
    mobileNumDaysOptions: [
      { label: "3 days", value: 3 },
      { label: "7 days", value: 7 },
    ],
  }),

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    hintStateLocalStorageKey() {
      return `closedHintText${this.state}` + ("&isGroup" ? this.isGroup : "")
    },
    showHintText() {
      return this.hintTextState && !localStorage[this.hintStateLocalStorageKey]
    },
    hintText() {
      switch (this.state) {
        case this.isGroup && this.states.EDIT_AVAILABILITY:
          return "Toggle which calendars are used. Editing availability is disabled and determined from your calendar events each week."
        case this.states.EDIT_AVAILABILITY:
          return "Tap and drag to add your available times in green"
        case this.states.SCHEDULE_EVENT:
          return "Tap and drag on the calendar to schedule a Google Calendar event during those times"
        default:
          return ""
      }
    },
  },

  methods: {
    updateShowBestTimes(val) {
      this.$emit("update:showBestTimes", !!val)
      this.$emit("onShowBestTimesChange", !!val)
    },
    closeHint() {
      this.hintTextState = false
      localStorage[this.hintStateLocalStorageKey] = true
    },
  },
}
</script>
