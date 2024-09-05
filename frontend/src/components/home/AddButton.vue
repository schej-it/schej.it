<template>
  <v-menu
    v-model="showAddMenu"
    ref="menu"
    :close-on-content-click="false"
    transition="slide-y-transition"
    bottom
    offset-y
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
      <div class="tw-flex tw-flex-row tw-w-64">
        <div>
          <v-icon>mdi-folder-plus</v-icon>
        </div>
        <div class="tw-flex">
          <div>New folder</div>

          <div>
            <v-icon class="tw-text-yellow">mdi-crown</v-icon>
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
