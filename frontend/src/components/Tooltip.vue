<template>
  <div class="tw-relative" ref="tooltipTrigger">
    <slot></slot>
    <div
      v-if="isVisible && content"
      class="tw-pointer-events-none tw-fixed tw-z-50 tw-rounded-lg tw-bg-dark-gray tw-px-1.5 tw-py-1 tw-text-xs tw-text-white tw-shadow-lg tw-transition-opacity tw-duration-200"
      :style="{
        left: `${position.x}px`,
        top: `${position.y}px`,
        transform: 'translate(-50%, -50%)',
      }"
    >
      {{ content }}
    </div>
  </div>
</template>

<script>
export default {
  name: "Tooltip",
  props: {
    content: {
      type: String,
      default: "",
    },
  },
  data() {
    return {
      position: { x: 0, y: 0 },
      isVisible: true,
    }
  },
  methods: {
    handleMouseMove(e) {
      this.position = {
        x: e.clientX, // Offset from cursor
        y: e.clientY - 30,
      }
    },
    handleMouseEnter() {
      this.isVisible = true
    },
    handleMouseLeave() {
      this.isVisible = false
    },
  },
  mounted() {
    if (this.$refs.tooltipTrigger) {
      this.$refs.tooltipTrigger.addEventListener(
        "mousemove",
        this.handleMouseMove
      )
      this.$refs.tooltipTrigger.addEventListener(
        "mouseenter",
        this.handleMouseEnter
      )
      this.$refs.tooltipTrigger.addEventListener(
        "mouseleave",
        this.handleMouseLeave
      )
    }
  },
  beforeDestroy() {
    if (this.$refs.tooltipTrigger) {
      this.$refs.tooltipTrigger.removeEventListener(
        "mousemove",
        this.handleMouseMove
      )
      this.$refs.tooltipTrigger.removeEventListener(
        "mouseenter",
        this.handleMouseEnter
      )
      this.$refs.tooltipTrigger.removeEventListener(
        "mouseleave",
        this.handleMouseLeave
      )
    }
  },
}
</script>
