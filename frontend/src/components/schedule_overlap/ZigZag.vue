<!-- Used to indicate that a schedule overlap component has more to scroll -->
<template>
  <div ref="container" class="tw-overflow-hidden">
    <div v-if="left" class="left-gradient" :style="gradientStyle"></div>
    <div v-else class="right-gradient" :style="gradientStyle"></div>
    <!-- <div :class="left ? 'line1-left' : 'line1-right'" :style="lineStyle"></div>
    <div :class="left ? 'line2-left' : 'line2-right'" :style="lineStyle"></div> -->
  </div>
</template>

<style scoped>
.left-gradient {
  background: linear-gradient(90deg, rgba(0, 0, 0, 0.1), transparent);
}
.right-gradient {
  background: linear-gradient(-90deg, rgba(0, 0, 0, 0.15), transparent);
}
.line1-left {
  background: linear-gradient(
    45deg,
    white,
    white 49%,
    black 49%,
    transparent 51%
  );
}
.line2-left {
  background: linear-gradient(
    -45deg,
    transparent,
    transparent 49%,
    black 49%,
    white 51%
  );
}

.line1-right {
  background: linear-gradient(
    45deg,
    transparent,
    transparent 49%,
    black 51%,
    white 51%
  );
}
.line2-right {
  background: linear-gradient(
    -45deg,
    white,
    white 49%,
    black 51%,
    transparent 51%
  );
}
</style>

<script>
export default {
  name: "ZigZag",

  props: {
    left: { type: Boolean, default: false },
    right: { type: Boolean, default: false },
  },

  mounted() {
    // Background size is 2 * width of the element
    this.backgroundSize = this.$refs.container.offsetWidth * 2
  },

  data() {
    return {
      backgroundSize: 0,
    }
  },

  computed: {
    gradientStyle() {
      return {
        position: "absolute",
        width: "100%",
        height: "100%",
        // transform: this.left
        //   ? `translate(${-this.backgroundSize / 2}px, 0)`
        //   : "",
      }
    },
    lineStyle() {
      return {
        position: "absolute",
        width: "200%",
        height: "100%",
        backgroundSize: `${this.backgroundSize}px ${this.backgroundSize}px`,
        transform: this.left
          ? `translate(${-this.backgroundSize / 2}px, 0)`
          : "",
      }
    },
  },
}
</script>
