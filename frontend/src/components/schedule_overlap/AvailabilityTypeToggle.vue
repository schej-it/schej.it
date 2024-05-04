<template>
  <div
    class="tw-relative tw-flex tw-w-fit tw-items-center tw-rounded tw-border tw-border-light-gray"
  >
    <div
      class="tw-absolute tw-h-full tw-w-[50%] tw-rounded tw-border-2 tw-transition-all"
      :class="tabs[index].borderClass"
      :style="{
        ...tabs[index].borderStyle,
        transform: `translateX(${index * 100}%)`,
      }"
    ></div>
    <template v-for="(tab, i) in tabs">
      <div
        class="tw-flex-1 tw-cursor-pointer tw-px-4 tw-py-2 tw-text-center tw-text-sm tw-font-medium tw-transition-all"
        :class="i === index ? tab.activeClass : inactiveClass"
        @click="$emit('input', tab.value)"
      >
        {{ tab.text }}
      </div>
    </template>
  </div>
</template>

<script>
import { availabilityTypes } from "@/constants"

export default {
  name: "AvailabilityTypeToggle",

  props: {
    value: { type: String, required: true },
  },

  data() {
    return {
      index: 0,
      tabs: [
        {
          text: "Available",
          activeClass: "tw-text-green tw-bg-green/5",
          borderClass: "tw-border-green",
          borderStyle: { boxShadow: "0px 2px 8px 0px #00994C80" },
          value: availabilityTypes.AVAILABLE,
        },
        {
          text: "If needed",
          activeClass: "tw-text-dark-yellow tw-bg-yellow/5",
          borderClass: "tw-border-orange",
          borderStyle: { boxShadow: "0px 2px 8px 0px #FEDB9380" },
          value: availabilityTypes.IF_NEEDED,
        },
      ],
    }
  },

  computed: {
    inactiveClass() {
      return "tw-text-dark-gray tw-bg-light-gray"
    },
  },

  watch: {
    value: {
      immediate: true,
      handler() {
        this.index = this.tabs.findIndex((tab) => tab.value === this.value)
        if (this.index === -1) this.index = 0
      },
    },
  },
}
</script>
