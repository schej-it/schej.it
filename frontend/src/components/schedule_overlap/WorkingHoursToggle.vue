<template>
  <div>
    <div class="tw-mb-2 tw-text-sm tw-text-black">Working hours</div>
    <v-switch
      id="working-hours-toggle"
      inset
      :input-value="value"
      @change="(val) => $emit('input', val)"
      hide-details
    >
      <template v-slot:label>
        <div class="tw-text-sm tw-text-black">
          <div class="tw-flex tw-items-center tw-gap-2">
            <v-select
              dense
              hide-details
              class="-tw-mt-0.5 tw-w-20 tw-text-xs"
              :items="times"
              :value="startTime"
              @input="(val) => $emit('update:startTime', val)"
              @click="
                (e) => {
                  e.preventDefault()
                  e.stopPropagation()
                }
              "
            />
            <div>to</div>
            <v-select
              dense
              hide-details
              class="-tw-mt-0.5 tw-w-20 tw-text-xs"
              :items="times"
              :value="endTime"
              @input="(val) => $emit('update:endTime', val)"
              @click="
                (e) => {
                  e.preventDefault()
                  e.stopPropagation()
                }
              "
            />
          </div>
        </div>
      </template>
    </v-switch>
    <div class="tw-mt-2 tw-text-xs tw-text-dark-gray">
      Only autofill availability between working hours
    </div>
  </div>
</template>

<script>
import { getTimeOptions } from "@/utils"

export default {
  name: "WorkingHoursToggle",

  props: {
    value: { type: Boolean, required: true },
    startTime: { type: Number, required: true },
    endTime: { type: Number, required: true },
    timezone: { type: Object, required: true },
  },

  computed: {
    times() {
      return getTimeOptions()
    },
  },
}
</script>
