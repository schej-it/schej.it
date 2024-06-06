<template>
  <v-switch
    :input-value="value"
    @change="emitItem"
    inset
    class="tw-flex tw-items-center"
  >
    <template v-slot:label>
      <div
        class="tw-flex tw-items-center tw-justify-center tw-text-xs tw-text-black"
      >
        Buffer
        <v-select
          dense
          :items="bufferTimes"
          class="-tw-mb-2 tw-w-[3.1rem] tw-scale-75 tw-text-xs"
          :value="bufferTime"
          @input="(val) => $emit('update:bufferTime', val)"
          @click="
            (e) => {
              e.preventDefault()
              e.stopPropagation()
            }
          "
        ></v-select>
        minutes
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
      bufferTimes: [15, 30, 45],
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
