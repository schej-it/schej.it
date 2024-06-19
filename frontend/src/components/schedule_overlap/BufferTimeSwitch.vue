<template>
  <div>
    <div class="tw-mb-1 tw-text-sm tw-text-black">Buffer time</div>
    <div class="tw-mb-2 tw-text-xs tw-text-dark-gray">
      Add time around calendar events
    </div>
    <v-switch
      id="buffer-time-switch"
      :input-value="value"
      @change="emitItem"
      inset
      class="tw-flex tw-items-center"
      hide-details
    >
      <template v-slot:label>
        <div
          class="tw-flex tw-items-center tw-justify-center tw-gap-2 tw-text-sm tw-text-black"
        >
          <v-select
            dense
            hide-details
            :items="bufferTimes"
            class="-tw-mt-0.5 tw-w-20 tw-text-xs"
            :value="bufferTime"
            @input="(val) => $emit('update:bufferTime', val)"
            @click="
              (e) => {
                e.preventDefault()
                e.stopPropagation()
              }
            "
          ></v-select>
        </div>
      </template>
    </v-switch>
  </div>
</template>

<script>
export default {
  name: "BufferTimeToggle",

  props: {
    value: { type: Boolean, required: true },
    bufferTime: { type: Number, required: true },
  },

  components: {},

  data() {
    return {
      bufferTimes: [
        { text: "15 min", value: 15 },
        { text: "30 min", value: 30 },
        { text: "45 min", value: 45 },
      ],
    }
  },

  methods: {
    emitItem(e) {
      this.$emit("input", e)
    },
  },

  watch: {
    bufferTime() {
      this.$emit("update:bufferTime", this.bufferTime)
    },
  },
}
</script>
