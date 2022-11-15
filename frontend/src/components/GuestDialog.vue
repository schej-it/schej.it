<template>
  <v-dialog 
    :value="value"
    @input="e => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title class="tw-flex">
        <div>Continue as guest</div>
        <v-spacer />
        <v-btn 
          icon 
          @click="$emit('input', false)"
        ><v-icon>mdi-close</v-icon></v-btn>
      </v-card-title>
      <v-card-text>
        <v-text-field 
          v-model="name"
          autofocus
          class="tw-text-white tw-mb-4"
          placeholder="Enter your name..."
          hide-details
        />
        <div class="tw-flex">
          <v-spacer />
          <v-btn
            @click="submit"
            class="tw-bg-green"
            :disabled="!formComplete"
            :dark="formComplete"
          >
            Continue
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>
import { isPhone } from '@/utils'

export default {
  name: 'GuestDialog',

  emits: ['input', 'submit'],

  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      name: '',
    }
  },

  computed: {
    isPhone() {
      return isPhone(this.$vuetify)
    },
    formComplete() {
      return this.name.length > 0
    },
  },

  methods: {
    submit() {
      this.$emit('submit', this.name)
    },
  },

  watch: {
    value() {
      if (this.value) {
        this.name = ''
      }
    },
  },
}
</script>