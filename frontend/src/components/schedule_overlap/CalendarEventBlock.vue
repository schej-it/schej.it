<template>
  <transition :name="transitionName" appear>
    <div
      class="tw-absolute tw-w-full tw-select-none tw-p-px"
      :style="blockStyle"
      style="pointer-events: none"
    >
      <div
        class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-rounded tw-border tw-border-solid tw-p-1 tw-text-xs"
        :class="containerClass"
      >
        <div :class="textColor" class="ph-no-capture tw-font-medium">
          {{ noEventNames ? "BUSY" : calendarEvent.summary }}
        </div>
      </div>
    </div>
  </transition>
</template>

<script>
export default {
  name: "CalendarEventBlock",
  props: {
    blockStyle: {
      type: Object,
      default: () => ({}),
    },
    calendarEvent: {
      type: Object,
      required: true,
    },
    isGroup: {
      type: Boolean,
      required: true,
    },
    isEditingAvailability: {
      type: Boolean,
      required: true,
    },
    noEventNames: {
      type: Boolean,
      required: true,
    },
    transitionName: {
      type: String,
      required: true,
    },
  },
  computed: {
    containerClass() {
      if (this.calendarEvent.free) {
        return this.isGroup && !this.isEditingAvailability
          ? "tw-border-white tw-bg-light-blue tw-opacity-50"
          : "tw-border-dashed tw-border-blue"
      } else {
        return this.isGroup && !this.isEditingAvailability
          ? "tw-border-white tw-bg-light-blue"
          : "tw-border-blue"
      }
    },
    textColor() {
      const color =
        this.isGroup && !this.isEditingAvailability
          ? "white"
          : this.noEventNames
          ? "dark-gray"
          : "blue"
      return `tw-text-${color}`
    },
  },
}
</script>
