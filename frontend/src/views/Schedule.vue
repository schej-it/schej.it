<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-p-4">
    <div class="tw-pt-5 tw-bg-white-white tw-top-0 tw-left-0 tw-w-full tw-h-32 tw-z-20 tw-relative">

      <UserItem :user="userItem" @showEventNames="showEventNames" />
    </div>
    <div ref="calendarContainer" class="-tw-mt-24 tw-z-0 tw-relative">
      <div class="tw-w-full tw-h-24 tw-flex tw-text-center tw-pt-1 tw-flex-col tw-bg-green">
        <h1 class="tw-text-4xl tw-font-bold tw-text-white">{{ authUser.firstName }}'s Schej</h1>
        <h1 class="tw-text-md tw-text-white mt-2"><span class="tw-font-medium">Timezone:</span> {{ currentTimezone }}
        </h1>
      </div>
      <TestCalendar :noEventNames="hideEventNames" />
    </div>

    <v-scale-transition appear origin="center">
      <v-btn :loading="loading" :disabled="loading" fixed
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-16 white--text" fab @click="share">
        <v-icon dark> mdi-share </v-icon>
      </v-btn>
    </v-scale-transition>

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

  async mounted() { },

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
    currentTimezone() {
      return getCurrentTimezone()
    }
  },

  watch: {},

  methods: {
    ...mapActions(['showError', 'showInfo']),
    showEventNames(option) {
      this.hideEventNames = !option
    },
    async share() {
      this.loading = true

      let el = this.$refs.calendarContainer
      let canvas = await html2canvas(el)

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
              text: ``,
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
