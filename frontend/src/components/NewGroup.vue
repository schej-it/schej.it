<template>
  <v-card
    :flat="dialog"
    :class="{ 'tw-py-4': !dialog }"
    class="tw-overflow-none tw-relative tw-flex tw-max-w-[28rem] tw-flex-col tw-rounded-lg tw-transition-all"
  >
    <v-card-title class="tw-mb-2 tw-flex tw-px-4 sm:tw-px-8">
      <div>
        {{ edit ? "Edit group" : "New group" }}
      </div>
      <v-spacer />
      <template v-if="dialog">
        <v-btn v-if="showHelp" icon @click="helpDialog = true">
          <v-icon>mdi-help-circle</v-icon>
        </v-btn>
        <v-btn v-else @click="$emit('input', false)" icon>
          <v-icon>mdi-close</v-icon>
        </v-btn>
        <HelpDialog
          v-model="helpDialog"
          text="Use groups to see people's calendar availabilities each week"
        />
      </template>
    </v-card-title>
    <v-card-text class="tw-flex-1 tw-overflow-auto tw-px-4 tw-py-1 sm:tw-px-8">
      <div class="tw-flex tw-flex-col tw-space-y-6">
        <v-text-field
          ref="name-field"
          v-model="name"
          placeholder="Name your group..."
          autofocus
          :disabled="loading"
          hide-details
          solo
          @keyup.enter="blurNameField"
        />

        <div>
          <div class="tw-mb-2 tw-text-lg tw-text-black">Time range</div>
          <div class="tw-flex tw-items-baseline tw-justify-center tw-space-x-2">
            <v-select
              v-model="startTime"
              :disabled="loading"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
            <div>to</div>
            <v-select
              v-model="endTime"
              :disabled="loading"
              menu-props="auto"
              :items="times"
              hide-details
              solo
            ></v-select>
          </div>
        </div>

        <div>
          <div class="tw-mb-2 tw-text-lg tw-text-black">Day range</div>
          <div>
            <v-btn-toggle
              v-model="selectedDaysOfWeek"
              multiple
              solo
              color="primary"
            >
              <v-btn> S </v-btn>
              <v-btn> M </v-btn>
              <v-btn> T </v-btn>
              <v-btn> W </v-btn>
              <v-btn> T </v-btn>
              <v-btn> F </v-btn>
              <v-btn> S </v-btn>
            </v-btn-toggle>
          </div>
        </div>

        <div>
          <v-btn
            class="tw-justify-start tw-pl-0"
            block
            text
            @click="showAdvancedOptions = !showAdvancedOptions"
            ><span class="tw-mr-1">Advanced options</span>
            <v-icon :class="`tw-rotate-${showAdvancedOptions ? '180' : '0'}`"
              >mdi-chevron-down</v-icon
            ></v-btn
          >
          <v-expand-transition>
            <div v-show="showAdvancedOptions">
              <div class="tw-my-2">
                <TimezoneSelector v-model="timezone" label="Timezone" />
              </div>
            </div>
          </v-expand-transition>
        </div>
      </div>
    </v-card-text>
    <v-card-actions class="tw-relative tw-px-8">
      <v-btn
        block
        :loading="loading"
        :dark="formComplete"
        class="tw-mt-4 tw-bg-green"
        :disabled="!formComplete"
        @click="submit"
      >
        {{ edit ? "Edit" : "Create" }}
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { validateEmail, isPhone } from "@/utils"
import { mapState, mapActions } from "vuex"
import HelpDialog from "./HelpDialog.vue"
import TimezoneSelector from "./schedule_overlap/TimezoneSelector.vue"

export default {
  name: "NewGroup",

  emits: ["input"],

  props: {
    event: { type: Object },
    edit: { type: Boolean, default: false },
    dialog: { type: Boolean, default: true },
    showHelp: { type: Boolean, default: false },
  },

  components: {
    HelpDialog,
    TimezoneSelector,
  },

  data: () => ({
    name: "",
    startTime: 9,
    endTime: 17,
    loading: false,
    selectedDaysOfWeek: [],
    emails: [],

    showAdvancedOptions: false,
    timezone: {},

    helpDialog: false,
  }),

  computed: {
    ...mapState(["authUser"]),
    formComplete() {
      let emailsValid = true

      for (const email of this.emails) {
        if (!validateEmail(email)) {
          emailsValid = false
          break
        }
      }

      return (
        this.name.length > 0 &&
        this.selectedDaysOfWeek.length > 0 &&
        emailsValid //&&
        // (this.startTime < this.endTime ||
        //   (this.endTime === 0 && this.startTime != 0))
      )
    },
    isPhone() {
      return isPhone(this.$vuetify)
    },
    times() {
      const times = []

      for (let h = 1; h < 12; ++h) {
        times.push({ text: `${h} am`, value: h })
      }
      for (let h = 0; h < 12; ++h) {
        times.push({ text: `${h == 0 ? 12 : h} pm`, value: h + 12 })
      }
      times.push({ text: "12 am", value: 0 })

      return times
    },
  },

  methods: {
    ...mapActions(["showError"]),
    blurNameField() {
      this.$refs["name-field"].blur()
    },
    reset() {
      this.name = ""
      this.startTime = 9
      this.endTime = 17
      this.selectedDaysOfWeek = []
    },
    submit() {},
  },

  watch: {
    event: {
      immediate: true,
      handler() {
        // Populate event fields if this.event exists
        if (this.event) {
          this.name = this.event.name
          this.startTime = Math.floor(dateToTimeNum(this.event.dates[0]))
          this.endTime = (this.startTime + this.event.duration) % 24

          const selectedDaysOfWeek = []
          for (const date of this.event.dates) {
            selectedDaysOfWeek.push(new Date(date).getDay())
          }
          this.selectedDaysOfWeek = selectedDaysOfWeek
        }
      },
    },
  },
}
</script>
