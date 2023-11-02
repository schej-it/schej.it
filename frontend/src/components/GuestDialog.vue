<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title class="tw-flex">
        <div>Continue as guest</div>
        <v-spacer />
        <v-btn icon @click="$emit('input', false)">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="name"
          @keyup.enter="submit"
          @beforeinput="restrictInput"
          :error="alreadyTaken"
          class="-tw-mt-1"
          placeholder="Enter your name..."
          :hint="alreadyTaken ? 'Name already taken' : ''"
          persistent-hint
          autofocus
        ></v-text-field>
        <div class="tw-flex">
          <v-spacer />
          <v-btn
            @click="submit"
            class="tw-bg-green"
            :disabled="!formComplete || alreadyTaken"
            :dark="formComplete && !alreadyTaken"
          >
            Continue
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from "@/utils"

export default {
  name: "GuestDialog",

  emits: ["input", "submit"],

  props: {
    value: { type: Boolean, required: true },
    respondents: { type: Array, required: true },
  },

  data() {
    return {
      name: "",
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    formComplete() {
      return this.name.length > 0
    },
    alreadyTaken() {
      return this.respondents.includes(this.name)
    },
  },

  methods: {
    submit() {
      if (!this.alreadyTaken && this.formComplete)
        this.$emit("submit", this.name)
    },
    /** Restricts input to only letters, numbers, and spaces */
    restrictInput(e) {
      // if (e.data && /[^\w\s]/.test(e.data)) {
      //   e.preventDefault()
      // }
    },
  },

  watch: {
    value() {
      if (this.value) {
        this.name = ""
      }
    },
  },
}
</script>
