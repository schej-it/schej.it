<template>
  <div class="tw-flex tw-align-center">
    <v-switch @input="(e) => $emit('input', e)"  inset class="tw-flex tw-items-center"> </v-switch>
    <div
      class="tw-flex tw-items-center tw-justify-center tw-text-xs tw-text-black"
    >
      Buffer
      <v-select
        dense
        :items="bufferTimes"
        class="-tw-mb-2 tw-w-[3.1rem] tw-scale-75 tw-text-xs"
        v-model="bufferTime"
      ></v-select>
      minutes
    </div>
  </div>
</template>

<script>
export default {
  name: "BufferTimeToggle",

  props: {
    value: { type: String, required: true },
  },

  components: {},

  data() {
    return {
      bufferTime: 15,
      bufferTimes: [15, 30, 45],
    }
  },
}
</script>
