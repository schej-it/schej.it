<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="425"
    content-class="tw-m-0"
  >
    <v-card>
      <v-card-title>Unsaved changes</v-card-title>
      <v-card-text>
        Are you sure you want to leave this page? 
Changes you made will not be saved.
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="$emit('input', false)" class="tw-text-green">Cancel</v-btn>
        <v-btn text @click="$emit('input', false)">Leave page</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: "SaveProgressDialog",

  props: {
    value: { type: Boolean, required: true },
  },
}
</script>
