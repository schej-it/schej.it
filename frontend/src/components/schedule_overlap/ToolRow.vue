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
        class="tw-flex tw-flex-1 tw-flex-wrap tw-gap-x-4 tw-gap-y-2 tw-py-4 sm:tw-justify-start sm:tw-gap-x-8"
      >
        <!-- Select timezone -->
        <div class="tw-flex tw-items-center tw-gap-2">
          <TimezoneSelector
            v-if="!event.daysOnly"
            class="tw-w-full sm:tw-w-[unset]"
            :value="curTimezone"
            @input="(val) => $emit('update:curTimezone', val)"
          />
          <v-select
            :value="timeType"
            @input="$emit('update:timeType', $event)"
            :items="timeTypeOptions"
            :menu-props="{ auto: true }"
            item-text="label"
            item-value="value"
            class="-tw-mt-px tw-w-16 tw-text-sm"
            dense
            hide-details
          />
        </div>

        <template v-if="state !== states.EDIT_AVAILABILITY && isPhone">
          <ExpandableSection
            label="Options"
            :value="showEventOptions"
            @input="$emit('toggleShowEventOptions')"
            class="tw-mt-2 tw-w-full"
          >
            <div class="tw-flex tw-flex-col tw-gap-2 tw-pt-2">
              <template v-if="numResponses > 1">
                <v-switch
                  inset
                  id="show-best-times-toggle"
                  :input-value="showBestTimes"
                  @change="(val) => $emit('update:showBestTimes', !!val)"
                  hide-details
                >
                  <template v-slot:label>
                    <div class="tw-text-sm tw-text-black">
                      Show best {{ event.daysOnly ? "days" : "times" }}
                    </div>
                  </template>
                </v-switch>
                <v-switch
                  inset
                  id="hide-if-needed-toggle"
                  :input-value="hideIfNeeded"
                  @change="(val) => $emit('update:hideIfNeeded', !!val)"
                  hide-details
                >
                  <template v-slot:label>
                    <div class="tw-text-sm tw-text-black">
                      Hide if needed {{ event.daysOnly ? "days" : "times" }}
                    </div>
                  </template>
                </v-switch>
              </template>
              <div
                v-if="!event.daysOnly"
                class="tw-flex tw-basis-full tw-items-center tw-gap-x-2 tw-py-1"
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
            </div>
          </ExpandableSection>
        </template>
        <template
          v-if="state === states.EDIT_AVAILABILITY && isWeekly && !isPhone"
        >
          <v-spacer />
          <div class="tw-min-w-fit">
            <GCalWeekSelector
              v-if="calendarPermissionGranted"
              :week-offset="weekOffset"
              @update:weekOffset="(val) => $emit('update:weekOffset', val)"
            />
          </div>
        </template>
      </div>

      <div
        v-if="
          !event.daysOnly &&
          numResponses > 0 &&
          state !== states.EDIT_AVAILABILITY
        "
        style="width: 181.5px"
        class="tw-hidden sm:tw-flex"
      >
        <template v-if="state !== states.SCHEDULE_EVENT">
          <v-btn
            outlined
            class="tw-w-full tw-text-blue"
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
            :disabled="!allowScheduleEvent"
            class="tw-bg-blue tw-text-white"
            @click="(e) => $emit('confirmScheduleEvent', e)"
          >
            Schedule
          </v-btn>
        </template>
      </div>
    </div>

    <Advertisement
      class="tw-mt-5 sm:tw-mt-10"
      :ownerId="event.ownerId"
    ></Advertisement>
  </div>
</template>

<script>
import TimezoneSelector from "./TimezoneSelector.vue"
import GCalWeekSelector from "./GCalWeekSelector.vue"
import { isPhone } from "@/utils"
import Advertisement from "../event/Advertisement.vue"
import ExpandableSection from "../ExpandableSection.vue"
import { timeTypes } from "@/constants"

export default {
  name: "ToolRow",

  props: {
    event: { type: Object, required: true },
    state: { type: String, required: true },
    states: { type: Object, required: true },
    curTimezone: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    hideIfNeeded: { type: Boolean, required: true },
    isWeekly: { type: Boolean, required: true },
    calendarPermissionGranted: { type: Boolean, required: true },
    weekOffset: { type: Number, required: true },
    numResponses: { type: Number, required: true },
    mobileNumDays: { type: Number, default: 3 }, // The number of days to show at a time on mobile
    allowScheduleEvent: { type: Boolean, required: true },
    showEventOptions: { type: Boolean, required: true },
    timeType: { type: String, required: true },
  },

  components: {
    TimezoneSelector,
    GCalWeekSelector,
    Advertisement,
    ExpandableSection,
  },

  data: () => ({
    mobileNumDaysOptions: [
      { label: "3 days", value: 3 },
      { label: "7 days", value: 7 },
    ],
    timeTypeOptions: [
      { label: "12h", value: timeTypes.HOUR12 },
      { label: "24h", value: timeTypes.HOUR24 },
    ],
  }),

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
  },
}
</script>
