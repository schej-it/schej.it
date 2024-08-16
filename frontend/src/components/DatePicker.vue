<template>
  <div>
    <v-date-picker
      ref="datePicker"
      :pickerDate.sync="pickerDate"
      :value="value"
      @touchstart:date="touchstart"
      @mousedown:date="mousedown"
      @mouseover:date="mouseover"
      readonly
      no-title
      multiple
      color="primary"
      :show-current="false"
      class="tw-min-w-full tw-rounded-md tw-border-0 tw-drop-shadow sm:tw-min-w-0"
      :min="minCalendarDate"
      full-width
      :scrollable="false"
    ></v-date-picker>
    <!-- <div class="tw-mt-2 tw-text-xs tw-text-very-dark-gray">
      Drag to select multiple dates
    </div> -->
  </div>
</template>

<script>
export default {
  name: "DatePicker",

  props: {
    value: { type: Array, required: true },
    minCalendarDate: { type: String, default: "" },
  },

  data() {
    return {
      datePickerEl: null,
      dragging: false,
      dragState: "add",
      dragStates: { ADD: "add", REMOVE: "remove" },
      pickerDate: "",
    }
  },

  methods: {
    /** Start drag */
    mousedown(date) {
      this.dragging = true
      this.setDragState(date)
      this.addRemoveDate(date)
    },
    touchstart(date) {
      this.dragging = true
      this.setDragState(date)
      this.addRemoveDate(date)
    },

    /** Dragging */
    mouseover(date) {
      if (!this.dragging) return

      this.addRemoveDate(date)
    },
    touchmove(e) {
      if (!this.dragging) return

      e.preventDefault()

      // Get the target that we are touching
      var target = document.elementFromPoint(
        e.changedTouches[0].clientX,
        e.changedTouches[0].clientY
      )

      // Only care about targets that are within the date picker and are buttons
      if (
        target &&
        this.datePickerEl.contains(target) &&
        target.classList.contains("v-btn__content")
      ) {
        // Get date num from target
        const dateNum = parseInt(target.innerHTML)
        if (dateNum != NaN) {
          const dateNumString = `${dateNum}`
          const date = `${this.pickerDate}-${dateNumString.padStart(2, "0")}`
          this.addRemoveDate(date)
        }
      }
    },

    /** End drag */
    mouseup(e) {
      if (!this.dragging) return

      // Prevent month switching when tap and drag to left / right
      e.preventDefault()
      e.stopPropagation()

      this.dragging = false
    },

    /** Sets the drag state based on the date */
    setDragState(date) {
      const set = new Set(this.value)
      if (set.has(date)) {
        this.dragState = this.dragStates.REMOVE
      } else {
        this.dragState = this.dragStates.ADD
      }
    },
    addRemoveDate(date) {
      if (this.dragState === this.dragStates.ADD) {
        this.addDate(date)
      } else if (this.dragState === this.dragStates.REMOVE) {
        this.removeDate(date)
      }
    },
    addDate(date) {
      const set = new Set(this.value)
      set.add(date)
      this.$emit("input", [...set])
    },
    removeDate(date) {
      const set = new Set(this.value)
      set.delete(date)
      this.$emit("input", [...set])
    },
  },

  mounted() {
    this.datePickerEl = this.$refs.datePicker.$el
    this.datePickerEl.addEventListener("mouseup", this.mouseup)
    this.datePickerEl.addEventListener("touchmove", this.touchmove)
    this.datePickerEl.addEventListener("touchend", this.mouseup, {
      capture: true,
    })
  },

  beforeDestroy() {
    this.datePickerEl.removeEventListener("mouseup", this.mouseup)
    this.datePickerEl.removeEventListener("touchmove", this.touchmove)
    this.datePickerEl.removeEventListener("touchend", this.mouseup)
  },
}
</script>
