<template>
  <footer class="tw-bg-green tw-p-8 tw-text-white">
    <div class="tw-mx-auto tw-max-w-7xl">
      <div class="tw-flex tw-flex-col tw-gap-8 md:tw-grid md:tw-grid-cols-4">
        <!-- Logo and Info -->
        <div>
          <div class="tw-mb-3 tw-flex tw-items-center">
            <v-icon class="tw-mr-1 tw-text-white">mdi-calendar-blank</v-icon>
            <span class="tw-text-2xl tw-font-medium">schej</span>
          </div>
          <div class="tw-mb-6 tw-leading-snug">
            Made with 💚 by
            <a
              href="https://x.com/jonathanzliu/"
              target="_blank"
              class="tw-text-white tw-underline"
              >@jonathanzliu</a
            >
          </div>
          <div class="tw-mb-6 tw-flex tw-items-center tw-space-x-3">
            <a href="https://x.com/schej_it" aria-label="X" target="_blank">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 512 512"
                class="tw-mt-px tw-size-5 tw-fill-current tw-text-white/80 hover:tw-text-white"
              >
                <path
                  d="M389.2 48h70.6L305.6 224.2 487 464H345L233.7 318.6 106.5 464H35.8L200.7 275.5 26.8 48H172.4L272.9 180.9 389.2 48zM364.4 421.8h39.1L151.1 88h-42L364.4 421.8z"
                ></path>
              </svg>
            </a>
            <a
              href="https://github.com/schej-it/schej.it"
              aria-label="GitHub"
              target="_blank"
              ><v-icon class="tw-text-white/80 hover:tw-text-white"
                >mdi-github</v-icon
              ></a
            >
            <a
              href="https://www.linkedin.com/company/schejit"
              aria-label="LinkedIn"
              target="_blank"
              ><v-icon class="tw-text-white/80 hover:tw-text-white"
                >mdi-linkedin</v-icon
              ></a
            >
            <a
              href="https://www.reddit.com/r/schej/"
              target="_blank"
              aria-label="Reddit"
              ><v-icon class="tw-text-white/80 hover:tw-text-white"
                >mdi-reddit</v-icon
              ></a
            >
          </div>
          <a href="/privacy-policy" class="tw-text-sm">Privacy Policy</a>
        </div>
        <!-- Links -->
        <div class="tw-flex tw-flex-col tw-gap-2">
          <div class="tw-mb-1 tw-font-bold">Links</div>
          <div>
            <a href="/blog" class="tw-text-sm">Blog</a>
          </div>
          <div>
            <v-menu :nudge-bottom="10" offset-y :close-on-content-click="false">
              <template v-slot:activator="{ on, attrs }">
                <span
                  class="tw-cursor-pointer tw-text-sm tw-text-white/80 hover:tw-text-white"
                  v-bind="attrs"
                  v-on="on"
                >
                  Support
                </span>
              </template>
              <v-card class="tw-p-3">
                <div class="tw-text-sm">
                  Email support at
                  <span class="tw-text-green tw-underline"
                    >schej.team@gmail.com</span
                  >
                </div>
              </v-card>
            </v-menu>
          </div>
          <div>
            <a
              href="https://forms.gle/9AgRy4PQfWfVuBnw8"
              target="_blank"
              class="tw-text-sm"
              >Give Feedback</a
            >
          </div>
        </div>
        <!-- Articles -->
        <div class="tw-flex tw-flex-col tw-gap-2">
          <div class="tw-mb-1 tw-font-bold">Articles</div>
          <div>
            <a href="/blog/about/" class="tw-text-sm"
              >Welcome to the Schej blog!</a
            >
          </div>
          <div>
            <a href="/blog/doodle/" class="tw-text-sm">Schej vs Doodle</a>
          </div>
          <div>
            <a href="/blog/when2meet/" class="tw-text-sm">Schej vs When2meet</a>
          </div>
          <div>
            <a href="/blog/gcalandwhen2meet/" class="tw-text-sm"
              >How to use When2meet with Google Calendar</a
            >
          </div>
        </div>
        <!-- Reddit posts -->
        <div class="tw-flex tw-flex-col tw-gap-2">
          <div class="tw-mb-1 tw-font-bold">Reddit posts</div>
          <div>
            <a
              href="https://www.reddit.com/r/opensource/comments/1dlol7r/i_made_a_better_when2meet/"
              target="_blank"
              class="tw-text-sm"
              >I made a better when2meet</a
            >
          </div>
          <div>
            <a
              href="https://www.reddit.com/r/opensource/comments/1klu471/i_made_a_doodle_alternative/"
              target="_blank"
              class="tw-text-sm"
              >I made a Doodle alternative</a
            >
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
div:has(a) {
  @apply tw-leading-none;
}
a {
  @apply tw-text-white/80 !important;
}
a:hover {
  @apply tw-text-white !important;
}
</style>

<script>
export default {
  name: "Footer",
}
</script>
