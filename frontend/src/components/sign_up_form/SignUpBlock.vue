<template>
  <v-container>
    <div
      class="tw-flex tw-flex-col tw-rounded-md tw-border-[1px] tw-p-4"
      :class="unsaved ? 'tw-border-light-green' : 'tw-border-light-gray-stroke'"
    >
      <div class="tw-flex tw-items-start tw-justify-between tw-h-7">
        <div
          v-if="!isEditingName"
          class="tw-flex tw-items-center tw-gap-2 tw-font-medium"
        >
          {{ isEditing ? newName : signUpBlock.name }}
          <v-btn v-if="isEditing" icon x-small @click="isEditingName = true">
            <v-icon x-small>mdi-pencil</v-icon>
          </v-btn>
        </div>
        <div v-else class="tw-flex tw-w-full tw-items-center tw-gap-2 -tw-mt-[6px]">
          <v-text-field
            v-model="newName"
            dense
            hide-details
            autofocus
            @keyup.enter="saveName"
          ></v-text-field>
          <v-btn icon x-small @click="cancelEditName">
            <v-icon x-small>mdi-close</v-icon>
          </v-btn>
          <v-btn icon x-small color="primary" @click="saveName">
            <v-icon x-small>mdi-check</v-icon>
          </v-btn>
        </div>
      </div>
      <div class="tw-text-xs tw-italic tw-text-dark-gray">
        {{ timeRangeString }}
      </div>
      <div class="tw-mt-4 tw-flex tw-items-center tw-gap-4">
        <div class="tw-text-xs">People per slot</div>
        <div class="tw-flex tw-h-4 tw-items-center">
          <div v-if="isEditing" class="-tw-mt-2 tw-w-20">
            <v-select
              :value="signUpBlock.capacity"
              @input="
                $emit('update:signUpBlock', {
                  ...signUpBlock,
                  capacity: $event,
                })
              "
              class="tw-text-xs"
              menu-props="auto"
              :items="capacityOptions"
              hide-details
              dense
            ></v-select>
          </div>
          <div v-else class="tw-text-xs">{{ signUpBlock.capacity }}</div>
        </div>
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

      <div v-if="isEditing" class="tw-mt-2">
        <a
          class="tw-text-xs tw-text-red"
          text
          @click="$emit('delete:signUpBlock', signUpBlock._id)"
          >Delete slot</a
        >
      </div>

      <div v-if="!isOwner && hasCapacity" class="tw-mt-2">
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
    unsaved: { type: Boolean, default: false },
  },

  data: () => ({
    capacityOptions: [...Array(100).keys()].map((i) => i + 1),
    isEditingName: false,
    newName: "",
  }),

  computed: {
    timeRangeString() {
      return getStartEndDateString(
        this.signUpBlock.startDate,
        this.signUpBlock.endDate
      )
    },
    hasCapacity() {
      return (
        !this.signUpBlock.responses ||
        this.signUpBlock.capacity > this.signUpBlock.responses.length
      )
    },
  },

  methods: {
    saveName() {
      console.log(this.newName)
      this.$emit("update:signUpBlock", {
        ...this.signUpBlock,
        name: this.newName,
      })
      this.isEditingName = false
    },
    cancelEditName() {
      this.newName = this.signUpBlock.name
      this.isEditingName = false
    },
  },

  watch: {
    signUpBlock: {
      immediate: true,
      handler(newVal) {
        this.newName = newVal.name
      },
    },
  },
}
</script>
