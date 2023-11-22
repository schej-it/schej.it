<!-- Used to indicate that a schedule overlap component has more to scroll -->
<template>
  <div ref="container" class="tw-overflow-hidden">
    <div v-if="left" class="gradient left-gradient"></div>
    <div v-else-if="right" class="gradient right-gradient"></div>
  </div>
</template>

<style scoped>
.gradient {
  position: absolute;
  width: 100%;
  height: 100%;
}
.left-gradient {
  background: linear-gradient(90deg, rgba(0, 0, 0, 0.1), transparent);
}
.right-gradient {
  background: linear-gradient(-90deg, rgba(0, 0, 0, 0.15), transparent);
}
</style>

<script>
export default {
  name: "OverflowIndicator",

  props: {
    left: { type: Boolean, default: false },
    right: { type: Boolean, default: false },
  },
}
</script>
