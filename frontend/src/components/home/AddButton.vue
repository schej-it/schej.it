<template>
  <v-menu
    v-model="showAddMenu"
    ref="menu"
    :close-on-content-click="false"
    transition="slide-y-transition"
    bottom
    offset-y
    nudge-left="116"
    nudge-bottom="8"
  >
    <template v-slot:activator="{ on, attrs }">
      <v-btn
        color="white"
        class="tw-rounded-md tw-text-black"
        v-bind="attrs"
        v-on="on"
        @click.prevent
      >
        + Add
      </v-btn>
    </template>

    <div class="tw-py-1 tw-bg-white">
      <div class="tw-flex tw-flex-row tw-w-48 tw-justify-between tw-items-center tw-px-4 tw-py-4 tw-cursor-pointer">
        <div class="tw-flex tw-gap-2">
          <v-icon class="tw-text-very-dark-gray">mdi-folder-plus</v-icon>
          <div>New folder</div>
        </div>
        <div class="tw-flex">
          <div>
            <v-icon class="tw-text-gold" small>mdi-crown</v-icon>
          </div>
        </div>
      </div>
    </div>
  </v-menu>
</template>

<script>
export default {
  name: "AddButton",

  props: {},

  data() {
    return {
      showAddMenu: false,
    }
  },
}
</script>
