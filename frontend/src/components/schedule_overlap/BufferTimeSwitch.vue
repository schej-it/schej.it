<template>
  <v-switch
    :value="value"
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
          v-model="bufferTime"
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
  },

  components: {},

  data() {
    return {
      bufferTime: 15,
      bufferTimes: [15, 30, 45],
    }
  },

  methods: {
    emitItem(e) {
      this.$emit("input", e !== null)
    },
  },

  watch: {
    bufferTime() {
      this.$emit("update:bufferTime", this.bufferTime)
    },
  },
}
</script>
