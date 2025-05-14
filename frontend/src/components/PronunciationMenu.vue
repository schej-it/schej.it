<template>
  <v-menu
    :nudge-bottom="10"
    offset-y
    :close-on-content-click="false"
    @input="handleMenuStateChange"
  >
    <template v-slot:activator="{ on, attrs }">
      <span
        class="tw-cursor-pointer"
        :class="{
          'tw-underline': isMenuOpen,
          'hover:tw-underline': !isMenuOpen,
        }"
        v-bind="attrs"
        v-on="on"
      >
        how to pronounce "schej"?
      </span>
    </template>
    <v-card class="tw-p-3 tw-text-center">
      <div class="tw-text-left tw-text-sm">
        Pronounced like "schedule" but shorter - "skej"
      </div>
      <div
        class="pronunciation-image-container tw-mt-3 tw-flex tw-items-center tw-justify-center"
      >
        <img
          :src="currentImageSrc"
          alt="Schej pronunciation animation"
          class="pronunciation-animation-image"
        />
      </div>
      <audio ref="pronunciationAudio" controls class="tw-hidden" autoplay>
        <source
          src="@/assets/audio/schej_pronunciation.mp3"
          type="audio/mpeg"
        />
        Your browser does not support the audio element.
      </audio>
    </v-card>
  </v-menu>
</template>

<script>
export default {
  name: "PronunciationMenu",
  data() {
    return {
      isMenuOpen: false,
      currentImageIndex: 0,
      images: [
        require("@/assets/doodles/pronunciation/0.jpg"),
        require("@/assets/doodles/pronunciation/1.jpg"),
        require("@/assets/doodles/pronunciation/2.jpg"),
        require("@/assets/doodles/pronunciation/3.jpg"),
        require("@/assets/doodles/pronunciation/4.jpg"),
      ],
      animationInterval: null,
      animationSpeedMs: 200, // Speed of animation frame change
    }
  },
  computed: {
    currentImageSrc() {
      return this.images[this.currentImageIndex]
    },
  },
  methods: {
    handleMenuStateChange(isOpen) {
      this.isMenuOpen = isOpen
      if (isOpen) {
        this.startAnimationAndAudio()
      } else {
        this.stopAnimationAndAudio()
      }
    },
    startAnimationAndAudio() {
      this.currentImageIndex = 1 // Always start from the first image

      if (this.animationInterval) {
        clearInterval(this.animationInterval)
      }

      this.animationInterval = setInterval(() => {
        if (this.currentImageIndex < this.images.length - 1) {
          this.currentImageIndex++
        } else {
          // Last image reached, stop interval and reset to first image (0.jpg)
          clearInterval(this.animationInterval)
          this.animationInterval = null
          // Set a timeout to show 0.jpg for a bit after the animation completes
          // before truly resetting if the menu is still open. Or just set to 0.
          setTimeout(() => {
            // Check if menu is still open before resetting to 0 if desired
            // This ensures it ends on 0.jpg visibly after the animation sequence
            this.currentImageIndex = 0
          }, this.animationSpeedMs)
        }
      }, this.animationSpeedMs)

      if (this.$refs.pronunciationAudio) {
        this.$refs.pronunciationAudio.currentTime = 0
        this.$refs.pronunciationAudio.play().catch((error) => {
          console.warn("Audio play prevented: ", error)
        })
      }
    },
    stopAnimationAndAudio() {
      if (this.animationInterval) {
        clearInterval(this.animationInterval)
        this.animationInterval = null
      }
      this.currentImageIndex = 0 // Reset to the first image on close

      if (this.$refs.pronunciationAudio) {
        this.$refs.pronunciationAudio.pause()
        this.$refs.pronunciationAudio.currentTime = 0
      }
    },
  },
  beforeDestroy() {
    if (this.animationInterval) {
      clearInterval(this.animationInterval)
    }
    // Ensure audio is stopped and cleaned up if component is destroyed while menu is open
    if (this.$refs.pronunciationAudio) {
      this.$refs.pronunciationAudio.pause()
      this.$refs.pronunciationAudio.currentTime = 0
    }
  },
}
</script>

<style scoped>
.pronunciation-image-container {
  min-height: 100px; /* Example: Adjust to your image's aspect ratio or desired space */
}
.pronunciation-animation-image {
  max-width: 100%;
  max-height: 100px; /* Example: Adjust as needed */
  object-fit: contain;
}
</style>
