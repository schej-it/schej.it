<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-p-4">
    <UserItem :user="userItem" @showEventNames="showEventNames" />
    <TestCalendar :noEventNames="hideEventNames" ref="calendar" />

    <v-scale-transition appear origin="center">
      <v-btn
        :loading="loading"
        :disabled="loading"
        fixed
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-16 white--text"
        fab
        @click="share"
      >
        <v-icon dark> mdi-share </v-icon>
      </v-btn>
    </v-scale-transition>

    <div id="test"></div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex'
import html2canvas from 'html2canvas'
import { copyImageToClipboard } from 'copy-image-clipboard'

import { dataURItoBlob, isPhone, getCurrentTimezone } from '@/utils'
import UserItem from '@/components/UserItem'
import TestCalendar from '@/components/TestCalendar'

export default {
  name: 'Schedule',

  components: {
    UserItem,
    TestCalendar,
  },

  data() {
    return {
      output: '',
      hideEventNames: false,
      loading: false,
    }
  },

  async mounted() {},

  computed: {
    ...mapState(['authUser']),
    userItem() {
      const { firstName, lastName, picture } = this.authUser
      return {
        name: firstName + ' ' + lastName,
        picture,
        status: 'free',
      }
    },
  },

  watch: {},

  methods: {
    ...mapActions(['showError', 'showInfo']),
    showEventNames(option) {
      this.hideEventNames = !option
    },
    async share() {
      this.loading = true

      let el = this.$refs.calendar.$el
      let canvas = await html2canvas(el)
      let ctx = canvas.getContext("2d")

      console.log(ctx)
      ctx.setFont("30 Arial")
      ctx.fillText("Hello World", 200, 200)
      document.getElementById("test").appendChild(canvas)
      this.output = canvas.toDataURL('image/png')


      if (isPhone(this.$vuetify)) {
        let filesArray = [
          new File([dataURItoBlob(this.output)], 'schedule.png', {
            type: 'image/png',
          }),
        ]

        if (navigator.canShare && navigator.canShare({ files: filesArray })) {
          navigator
            .share({
              files: filesArray,
              title: 'My Schedule',
              text: `Check out my Schej! My timezone: ${
                getCurrentTimezone()
              }`,
            })
            .then(() => {
              console.log('Share was successful.')
            })
            .catch((error) => {
              console.log('Sharing failed', error)
              this.showError('There was an issue when trying to share')
            })
        } else {
          console.log(`Your system doesn't support sharing files.`)
        }
      } else {
        copyImageToClipboard(this.output)
          .then(() => {
            console.log('Image Copied')
            this.showInfo('Copied to clipboard! Try pasting in a chat')
            this.loading = false
          })
          .catch((e) => {
            console.log('Error: ', e.message)
            this.showError(
              'There was an issue when trying to copy your schedule'
            )
            this.loading = false
          })
      }
    },
  },
}
</script>
