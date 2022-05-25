<template>
  <div class="tw-p-4">

    <UserItem :user="userItem" />
    <TestCalendar :noEventNames="false" ref="calendar" />

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
import UserItem from '@/components/UserItem'
import TestCalendar from '@/components/TestCalendar'
import { mapState } from 'vuex'
import html2canvas from 'html2canvas'

export default {
  name: 'Schedule',

  components: {
    UserItem,
    TestCalendar,
  },

  data() {
    return {
      output: '',
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
    async share() {
      let el = this.$refs.calendar.$el
      this.output = (await html2canvas(el)).toDataURL('image/png')

      let filesArray = [new File([this.dataURItoBlob(this.output)], 'schedule.png', { type: 'image/png' })]
      console.log(filesArray)
      if (navigator.canShare && navigator.canShare({ files: filesArray })) {
        navigator
          .share({
            files: filesArray,
            title: 'Pictures',
            text: 'Our Pictures.',
          })
          .then(() => console.log('Share was successful.'))
          .catch((error) => console.log('Sharing failed', error))
      } else {
        console.log(`Your system doesn't support sharing files.`)
      }
    },
    dataURItoBlob(dataURI) {
      // convert base64 to raw binary data held in a string
      // doesn't handle URLEncoded DataURIs - see SO answer #6850276 for code that does this
      var byteString = atob(dataURI.split(',')[1])

      // separate out the mime component
      var mimeString = dataURI.split(',')[0].split(':')[1].split(';')[0]

      // write the bytes of the string to an ArrayBuffer
      var ab = new ArrayBuffer(byteString.length)
      var ia = new Uint8Array(ab)
      for (var i = 0; i < byteString.length; i++) {
        ia[i] = byteString.charCodeAt(i)
      }
      
      return new Blob([ab], { type: mimeString })
    },
  },
}
</script>
