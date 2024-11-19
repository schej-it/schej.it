<template>
  <div
    class="tw-cursor-pointer tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-border-2 tw-border-solid tw-bg-white tw-p-1 tw-text-xs"
    :class="unsaved ? 'tw-border-light-green' : 'tw-border-gray'"
  >
  <div v-if="!titleOnly">
    <div class="ph-no-capture tw-font-medium tw-text-dark-gray">
      {{ signUpBlock.name }}
    </div>
    <div class="ph-no-capture tw-font-medium tw-text-dark-gray">
      ({{ signUpBlock.responses ? signUpBlock.responses.length : 0 }}/{{ signUpBlock.capacity }})
    </div>
  </div>
  <div v-else>
    <div class="tw-text-dark-gray tw-text-xs tw-italic">
      {{ title }}
    </div>
  </div>
  </div>
</template>

<script>
export default {
  name: "SignUpCalendarBlock",

  props: {
    signUpBlock: { type: Object, default: null },
    unsaved: { type: Boolean, default: false },
    titleOnly: { type: Boolean, default: false },
    title: { type: String, default: "" },
  },

  data: () => ({

  }),

  computed: {},

  methods: {
  },
}
</script>
