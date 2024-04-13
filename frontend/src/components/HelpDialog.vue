<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title>Help</v-card-title>
      <v-card-text>{{ text }}</v-card-text>
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
    text: { type: String, default: "" },
  },
}
</script>
