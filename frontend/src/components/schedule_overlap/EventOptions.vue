<template>
  <ExpandableSection
    v-if="numResponses > 1"
    label="Options"
    :value="showEventOptions"
    @input="$emit('toggleShowEventOptions')"
  >
    <div class="tw-flex tw-flex-col tw-gap-4 tw-pt-2">
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
      <v-switch
        v-if="isGroup && !isPhone"
        inset
        :input-value="showCalendarEvents"
        @change="(val) => $emit('update:showCalendarEvents', Boolean(val))"
        hide-details
      >
        <template v-slot:label>
          <div class="tw-text-sm tw-text-black">Overlay calendar events</div>
        </template>
      </v-switch>
    </div>
  </ExpandableSection>
</template>

<script>
import { isPhone } from "@/utils"
import { eventTypes } from "@/constants"
import ExpandableSection from "@/components/ExpandableSection.vue"

export default {
  name: "EventOptions",

  components: {
    ExpandableSection,
  },

  props: {
    event: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    hideIfNeeded: { type: Boolean, required: true },
    numResponses: { type: Number, required: true },
    showEventOptions: { type: Boolean, required: true },
    showCalendarEvents: { type: Boolean, required: true },
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    isGroup() {
      return this.event.type === eventTypes.GROUP
    },
  },
}
</script>
