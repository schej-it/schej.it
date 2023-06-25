<!-- Displays a number within a green circle -->
<template>
  <div
    class="tw-grid tw-place-content-center tw-rounded-full tw-shrink-0 tw-bg-dark-green tw-text-white tw-w-7 tw-h-7 tw-text-sm tw-font-medium"
  >
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: "NumberBullet",
}
</script>
