<template>
  <ExpandableSection label="Options" v-model="showOptions">
    <div class="tw-pt-2" v-if="numResponses > 1">
      <v-switch
        inset
        id="show-best-times-toggle"
        :input-value="showBestTimes"
        @change="updateShowBestTimes"
        hide-details
      >
        <template v-slot:label>
          <div class="tw-text-sm tw-text-black">
            Show best {{ event.daysOnly ? "days" : "times" }}
          </div>
        </template>
      </v-switch>
    </div>
  </ExpandableSection>
</template>

<script>
import ExpandableSection from "@/components/ExpandableSection.vue"

export default {
  name: "EventOptions",

  components: {
    ExpandableSection,
  },

  props: {
    event: { type: Object, required: true },
    showBestTimes: { type: Boolean, required: true },
    numResponses: { type: Number, required: true },
  },

  data() {
    return {
      showOptions: false,
    }
  },

  methods: {
    updateShowBestTimes(val) {
      this.$emit("update:showBestTimes", !!val)
    },
  },
}
</script>
