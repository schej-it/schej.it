<!-- Allows user to change timezone -->
<template>
  <div
    class="tw-flex tw-justify-center tw-items-center"
    id="timezone-select-container"
  >
    <div class="tw-mt-px tw-mr-2">Shown in</div>
    <v-select
      id="timezone-select"
      :value="value"
      @input="$emit('input', $event)"
      class="tw-text-sm -tw-mt-px tw-flex-none tw-min-w-min tw-max-w-xl"
      :items="timezones"
      dense
      color="#219653"
      item-color="green"
      hide-details
    ></v-select>
  </div>
</template>

<script>
export default {
  name: "TimezoneSelector",

  props: {
    value: { type: String, required: true },
    timezones: { type: Array, default: () => [] },
  },
};
</script>