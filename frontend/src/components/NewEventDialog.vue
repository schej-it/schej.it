<template>
  <v-dialog
    :value="value"
    fullscreen
    hide-overlay
    transition="dialog-bottom-transition"
  >
    <v-card tile class="tw-flex tw-flex-col">
      <v-card-title class="tw-flex">
        <div>New Event</div>
        <v-spacer />
        <v-btn 
          icon 
          @click="$emit('input', false)"
        ><v-icon>mdi-close</v-icon></v-btn>
      </v-card-title>
      <v-card-text class="tw-space-y-4 tw-flex tw-flex-col tw-flex-1">
        <v-text-field 
          ref="name-field"
          v-model="name"
          class="tw-text-white tw-flex-initial"
          placeholder="Name of event..."
          hide-details
          @keyup.enter="blurNameField"
        />

        <div>
          <div class="tw-flex tw-space-x-2 tw-items-baseline tw-justify-center">
            <v-select
              v-model="startTime"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              outlined
              hide-details
            ></v-select>
            <div>to</div>
            <v-select
              v-model="endTime"
              class="tw-flex-initial /*tw-w-28*/"
              menu-props="auto"
              :items="times"
              outlined
              hide-details
            ></v-select>
          </div>
        </div>
        
        <div>
          <div class="tw-text-lg tw-text-black tw-text-center tw-font-medium tw-mt-6 tw-mb-2">What dates would you like to meet?</div>
          <div class="tw-flex tw-flex-col tw-justify-center tw-items-center">
            <vc-date-picker v-model="dateRange" is-range class="tw-min-w-full sm:tw-min-w-0 " />
          </div>
        </div>

        <v-spacer />

        <v-btn 
          :dark="formComplete" 
          class="tw-bg-blue"
          :disabled="!formComplete"
          @click="submit"
        >Create</v-btn>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import { getDateWithTime, post } from '@/utils'
export default {
  name: 'NewEventDialog',

  emits: ['input'],

  props: {
    value: { type: Boolean, required: true },
  },

  data: () => ({
    name: '',
    startTime: 9,
    endTime: 17,
    dateRange: {},
  }),

  computed: {
    formComplete() {
      return this.name.length > 0 && this.dateRange && 'start' in this.dateRange && 'end' in this.dateRange
    },
    times() {
      const times = []

      for (let h = 1; h < 12; ++h) {
        times.push({ text: `${h} am`, value: h })
      }
      for (let h = 0; h < 12; ++h) {
        times.push({ text: `${h == 0 ? 12 : h} pm`, value: h+12 })
      }
      times.push({ text: '12 am', value: 0 })

      return times
    },
  },

  methods: {
    blurNameField() {
      this.$refs['name-field'].blur()
    },
    reset() {
      this.name = ''
      this.startTime = 9
      this.endTime = 17
      this.dateRange = {}
    },
    submit() {
      /* TODO: make sure to strip the time from the date so dates aren't localized to what the client's current time */
      const startDate = getDateWithTime(this.dateRange.start, '00:00')
      const endDate = getDateWithTime(this.dateRange.end, '11:59')
      post('/events', {
        name: this.name,
        startDate,
        endDate,
        startTime: this.startTime,
        endTime: this.endTime,
      }).then(data => {
        console.log(data)
        this.reset()
        this.$emit('input', false)
      })
    },
  },
}
</script>