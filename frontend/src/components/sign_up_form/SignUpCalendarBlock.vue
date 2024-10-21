<template>
  <div
    class="tw-h-full tw-w-full tw-overflow-hidden tw-text-ellipsis tw-border-2 tw-border-solid tw-bg-white tw-p-1 tw-text-xs"
    :class="unsaved ? 'tw-border-light-green' : 'tw-border-gray'"
  >
    <div class="ph-no-capture tw-font-medium tw-text-dark-gray">
      {{ title }}
    </div>
  </div>
</template>

<script>
export default {
  name: "SignUpCalendarBlock",

  props: {
    title: { type: String, required: true },
    unsaved: { type: Boolean, default: false },
  },

  data: () => ({}),

  computed: {},

  methods: {},
}
</script>
