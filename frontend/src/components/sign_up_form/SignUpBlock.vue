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

      <div v-if="signUpBlock.responses" class="tw-mt-2">
        <div
          v-for="(response, i) in signUpBlock.responses"
          :key="i"
          class="tw-relative tw-flex tw-items-center"
        >
          <div class="tw-ml-1 tw-mr-3">
            <v-avatar v-if="response.user.picture != ''" :size="16">
              <img
                v-if="response.user.picture"
                :src="response.user.picture"
                referrerpolicy="no-referrer"
              />
            </v-avatar>
            <v-avatar v-else :size="16">
              <v-icon small>mdi-account</v-icon>
            </v-avatar>
          </div>
          <div class="tw-mr-1 tw-transition-all">
            {{ response.user.firstName + " " + response.user.lastName }}
          </div>
        </div>
      </div>
      <div v-if="!isOwner" class="tw-mt-2">
        <a
          class="tw-text-xs tw-text-green"
          text
          @click="$emit('signUpForBlock', signUpBlock._id)"
          >+ Join this slot</a
        >
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
    isOwner: { type: Boolean, default: false },
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
