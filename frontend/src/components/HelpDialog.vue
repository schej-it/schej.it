<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="425"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title><slot name="header">Help</slot></v-card-title>
      <v-card-text><slot></slot></v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="$emit('input', false)">Ok</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "HelpDialog",

  props: {
    value: { type: Boolean, required: true },
  },
}
</script>
