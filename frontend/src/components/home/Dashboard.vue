<template>
  <div class="tw-rounded-md tw-px-6 tw-py-4 sm:tw-mx-4 sm:tw-bg-[#f3f3f366]">
    <div class="tw-mb-3 tw-flex tw-items-center tw-justify-between">
      <div class="tw-text-xl tw-font-medium tw-text-dark-green sm:tw-text-2xl">
        Dashboard
      </div>
      <v-btn
        text
        @click="createFolder"
        class="tw-hidden tw-text-very-dark-gray sm:tw-block"
      >
        <v-icon class="tw-mr-2 tw-text-lg">mdi-folder-plus</v-icon>
        New folder
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: "Dashboard",

  methods: {
    createFolder() {},
  },
}
</script>
