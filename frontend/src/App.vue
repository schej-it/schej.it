<template>
  <v-app>
    <v-app-bar
      app
      color="var(--green)"
      dark
    >
      <div class="d-flex align-center">
        <v-img
          alt="Vuetify Logo"
          class="shrink mr-2"
          contain
          src="@/assets/logo2_shadow_dark.svg"
          transition="scale-transition"
          width="150"
        />
      </div>
    </v-app-bar>

    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data: () => ({
    //
  }),
};
</script>

<style>
:root {
  --light-green: #6FCF97;
  --green: #27AE60;
  --dark-green: #219653;
  --white: #F2F2F2;
  --black: #4F4F4F;
}
</style>