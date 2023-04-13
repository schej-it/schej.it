<!-- Displays a number within a green circle -->
<template>
  <div class="
    tw-grid tw-place-content-center 
    tw-rounded-full 
    tw-shrink-0 
    tw-bg-green tw-text-white 
    md:tw-w-10 md:tw-h-10 
    tw-w-7 tw-h-7 
    md:tw-text-xl 
    tw-text-base
    tw-font-medium
  ">
    <slot></slot>
  </div>
</template>

<script>
export default {
  name: 'NumberBullet',
}
</script>