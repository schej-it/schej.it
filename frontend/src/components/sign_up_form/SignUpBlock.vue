<template>
  <v-container>
    <div
      class="tw-flex tw-flex-col tw-rounded-md tw-border-[1px] tw-border-light-gray-stroke tw-p-4"
    >
      <div class="tw-font-medium">{{ signUpBlock.name }}</div>
      <div class="tw-text-xs tw-italic tw-text-dark-gray">
        {{ timeRangeString }}
      </div>
      <div class="tw-mt-4 tw-flex tw-items-center tw-gap-4">
        <div class="tw-text-xs">People per slot</div>
        <div v-if="isEditing" class="-tw-mt-2 tw-w-20">
          <v-select
            :value="signUpBlock.capacity"
            @input="
              $emit('update:signUpBlock', { ...signUpBlock, capacity: $event })
            "
            class="tw-text-sm"
            menu-props="auto"
            :items="capacityOptions"
            hide-details
            dense
          ></v-select>
        </div>
        <div v-else class="tw-text-xs">{{ signUpBlock.capacity }}</div>
      </div>
    </div>
  </v-container>
</template>

<script>
import { getStartEndDateString } from "@/utils"

export default {
  name: "SignUpBlock",

  props: {
    signUpBlock: { type: Object, required: true },
    isEditing: { type: Boolean, default: false },
  },

  data: () => ({
    capacityOptions: [...Array(100).keys()].map((i) => i + 1),
  }),

  computed: {
    timeRangeString() {
      return getStartEndDateString(
        this.signUpBlock.startDate,
        this.signUpBlock.endDate
      )
    },
  },

  methods: {},
}
</script>
