<!--class="tw-flex tw-items-end tw-justify-start tw-p-1"-->
<template>
  <div>
    <v-btn
      class="-tw-ml-2 tw-w-[calc(100%+1rem)] tw-justify-between tw-px-2"
      block
      text
      @click="toggle"
    >
      <span class="-tw-ml-px tw-mr-1" :class="labelClass">
        {{ label }}
      </span>
      <v-spacer />
      <v-icon
        :class="`tw-rotate-${value ? '180' : '0'} ${iconClass}`"
        :size="30"
        >mdi-chevron-down</v-icon
      ></v-btn
    >
    <v-expand-transition>
      <div v-show="value">
        <slot></slot>
      </div>
    </v-expand-transition>
    <div ref="scrollTo"></div>
  </div>
</template>

<script>
export default {
  name: "ExpandableSection",

  props: {
    value: { type: Boolean, required: true },
    label: { type: String, default: "" },
    labelClass: { type: String, default: "tw-text-base" },
    iconClass: { type: String, default: "" },
    autoScroll: { type: Boolean, default: false },
  },

  methods: {
    toggle() {
      this.$emit("input", !this.value)
    },
    scrollToElement(element) {
      if (this.autoScroll && element) {
        setTimeout(() => element.scrollIntoView({ behavior: "smooth" }), 200)
      }
    },
  },

  watch: {
    value() {
      if (this.value) {
        this.scrollToElement(this.$refs.scrollTo)
      }
    },
  },
}
</script>
