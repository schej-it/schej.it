<!-- Allows user to change timezone -->
<template>
  <div
    class="tw-flex tw-items-center tw-justify-center"
    id="timezone-select-container"
  >
    <div class="tw-mr-2 tw-mt-px">Shown in</div>
    <v-select
      id="timezone-select"
      :value="value"
      @input="$emit('input', $event)"
      :items="timezones"
      :menu-props="{ auto: true }"
      class="-tw-mt-px tw-w-52 tw-text-sm"
      dense
      color="#219653"
      item-color="green"
      hide-details
      item-text="label"
      return-object
    ></v-select>
  </div>
</template>

<script>
import { allTimezones } from "@/constants"
import spacetime from "spacetime"

export default {
  name: "TimezoneSelector",

  props: {
    value: { type: Object, required: true },
  },

  created() {
    const localTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone
    let timezoneObject = this.timezones.find((t) => t.value === localTimezone)

    if (!timezoneObject) {
      const offset = spacetime.now(localTimezone).timezone().current.offset * 60
      timezoneObject = this.timezones.find((t) => t.offset === offset)
    }

    this.$emit("input", timezoneObject)
  },

  computed: {
    timezones() {
      // ===============================================================================
      // Source: https://github.com/ndom91/react-timezone-select/blob/main/src/index.tsx
      // ===============================================================================

      const t = Object.entries(allTimezones)
        .map((zone) => {
          try {
            const now = spacetime.now(zone[0])
            const tz = now.timezone()

            const min = tz.current.offset * 60
            const hr = `${(min / 60) ^ 0}:${
              min % 60 === 0 ? "00" : Math.abs(min % 60)
            }`
            const prefix = `(GMT${hr.includes("-") ? hr : `+${hr}`}) ${zone[1]}`

            let label = prefix

            return {
              value: tz.name,
              label: label,
              offset: tz.current.offset * 60,
            }
          } catch (e) {
            console.error(e)
            return null
          }
        })
        .filter(Boolean)
        .sort((a, b) => a.offset - b.offset)
      return t
    },
  },
}
</script>
