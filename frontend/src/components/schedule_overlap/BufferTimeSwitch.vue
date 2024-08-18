<template>
  <div>
    <div class="tw-mb-1 tw-text-sm tw-text-black">Buffer time</div>
    <div class="tw-mb-2 tw-text-xs tw-text-dark-gray">
      Add time around calendar events
    </div>
    <v-switch
      id="buffer-time-switch"
      :input-value="bufferTime.enabled"
      @change="(val) => updateBufferTime('enabled', val)"
      inset
      class="tw-flex tw-items-center"
      hide-details
    >
      <template v-slot:label>
        <div
          class="tw-flex tw-items-center tw-justify-center tw-gap-2 tw-text-sm tw-text-black"
        >
          <v-select
            menu-props="auto"
            dense
            hide-details
            :items="bufferTimes"
            class="-tw-mt-0.5 tw-w-20 tw-text-xs"
            :value="bufferTime.time"
            @input="(val) => updateBufferTime('time', val)"
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
import { patch } from "@/utils"

export default {
  name: "BufferTimeToggle",

  props: {
    bufferTime: { type: Object, required: true },
    syncWithBackend: { type: Boolean, default: false },
  },

  components: {},

  data() {
    return {
      bufferTimes: [
        { text: "15 min", value: 15 },
        { text: "30 min", value: 30 },
        { text: "45 min", value: 45 },
        { text: "1 hour", value: 60 },
      ],
    }
  },

  methods: {
    updateBufferTime(key, val) {
      const bufferTime = {
        ...this.bufferTime,
        [key]: val,
      }
      if (this.syncWithBackend) {
        patch(`/user/calendar-options`, {
          bufferTime,
        })
      }
      this.$emit("update:bufferTime", bufferTime)
    },
  },
}
</script>
