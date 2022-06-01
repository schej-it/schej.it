<template>
  <div class="tw-max-w-6xl tw-mx-auto tw-p-4">

    <UserItem :user="userItem" @showEventNames="showEventNames"/>
    <TestCalendar :noEventNames="hideEventNames" ref="calendar" />

    <v-scale-transition appear origin="center">
      <v-btn
        fab
        fixed
        dark
        class="tw-bg-blue tw-mx-auto tw-left-0 tw-right-0 tw-bottom-16"
        @click="share"
      >
        <v-icon>mdi-share</v-icon>
      </v-btn>
    </v-scale-transition>
  </div>
</template>

<script>
import { mapState } from 'vuex'
import html2canvas from 'html2canvas'

import { dataURItoBlob } from '@/utils'
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
      hideEventNames: false
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

  methods: {
    showEventNames(option) {
      this.hideEventNames = !option
    },
    async share() {
      let el = this.$refs.calendar.$el
      this.output = (await html2canvas(el)).toDataURL('image/png')

      let filesArray = [new File([dataURItoBlob(this.output)], 'schedule.png', { type: 'image/png' })]
      console.log(filesArray)
      if (navigator.canShare && navigator.canShare({ files: filesArray })) {
        navigator
          .share({
            files: filesArray,
            title: 'My Schedule',
            text: `Check out my Schej! My timezone: ${Intl.DateTimeFormat().resolvedOptions().timeZone}`,
          })
          .then(() => console.log('Share was successful.'))
          .catch((error) => console.log('Sharing failed', error))
      } else {
        console.log(`Your system doesn't support sharing files.`)
      }
    }
  },
}
</script>
