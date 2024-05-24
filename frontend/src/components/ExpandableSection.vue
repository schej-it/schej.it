<template>
  <div>
    <v-btn
      class="tw-flex tw-items-end tw-justify-start tw-p-1"
      block
      text
      @click="toggle"
    >
      <span class="tw-mr-1" :class="labelClass">
        {{ label }}
      </span>
      <v-spacer />
      <v-icon
        :class="`tw-rotate-${expanded ? '180' : '0'} ${iconClass}`"
        :size="30"
        >mdi-chevron-down</v-icon
      ></v-btn
    >
    <v-expand-transition>
      <div v-show="expanded">
        <slot></slot>
      </div>
    </v-expand-transition>
  </div>
</template>

<script>
export default {
  name: "ExpandableSection",

  props: {
    label: { type: String, default: "" },
    labelClass: { type: String, default: "tw-text-base" },
    iconClass: { type: String, default: "" },
    startExpanded: { type: Boolean, default: false },
  },

  data() {
    return {
      expanded: this.startExpanded,
    }
  },

  methods: {
    toggle() {
      this.expanded = !this.expanded
      this.$emit("toggle", this.expanded)
    },
  },
}
</script>
