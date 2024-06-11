<template>
  <v-switch
    :input-value="value"
    @change="emitItem"
    inset
    class="tw-flex tw-items-center"
    hide-details
  >
    <template v-slot:label>
      <div class="tw-flex tw-flex-col tw-text-xs">
        <div class="tw-flex tw-items-center tw-justify-center tw-text-black">
          Buffer time
          <v-select
            dense
            :items="bufferTimes"
            class="-tw-mb-[0.7rem] tw-w-20 tw-scale-90 tw-text-xs"
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
      </div>
    </template>
  </v-switch>
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
