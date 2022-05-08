<template>
  <div>
    <div class="tw-text-5xl tw-font-bold">Hello world</div>
    <vc-calendar />
  </div>
</template>

<script>
  import HelloWorld from '../components/HelloWorld'

  export default {
    name: 'Home',

    components: {
      HelloWorld,
    },
  }
</script>
