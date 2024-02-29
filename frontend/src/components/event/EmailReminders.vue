<template>
  <div>
    <div class="tw-flex tw-gap-1">
      <div class="tw-text-very-dark-gray">
      Set up email reminders</div>

      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <v-icon small v-bind="attrs" v-on="on"
            >mdi-information-outline
          </v-icon>
        </template>
        <div>Reminder emails will be sent the day of event creation,<br>one day after, and three days after.</div>
      </v-tooltip>
    </div>
    <div class="tw-mt-1 tw-text-dark-gray tw-text-xs" v-if="!hasContactsAccess">
        <a class="tw-underline" @click="requestContactsAccess"
          >Enable contacts access</a
        >
        for email auto-suggestions.
    </div>

    <v-combobox
      v-model="emails"
      class="tw-text-sm tw-mt-2"
      placeholder="Enter an email address.."
      dense
      multiple
      append-icon=""
      solo
      :rules="[rules.validEmails]"
    >
      <template v-slot:selection="data">
        <v-chip
          :key="data.item"
          size="x-small"
          class="tw-flex tw-items-center tw-bg-light-gray tw-text-very-dark-gray"
        >
          <!-- <template v-slot:prepend>
            <v-avatar class="bg-accent text-uppercase" start
              ><img
                v-if="data.item.picture.length > 0"
                :src="data.item.picture"
                referrerpolicy="no-referrer"
              />
              <v-icon v-else>mdi-account</v-icon></v-avatar
            >
          </template> -->
          {{ data.item }}

          <v-icon
            small
            @click="() => removeEmail(data.item)"
            class="tw-ml-1"
            >mdi-close</v-icon
          >
        </v-chip>
      </template>
      <!-- <template v-slot:item="{ item }">
        <v-list-item-avatar>
          <img
            v-if="item.picture.length > 0"
            :src="item.picture"
            referrerpolicy="no-referrer"
          />
          <v-icon v-else>mdi-account</v-icon>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title
            v-text="`${item.firstName} ${item.lastName}`"
          ></v-list-item-title>
          <v-list-item-subtitle v-text="item.email"></v-list-item-subtitle>
        </v-list-item-content>
      </template> -->
    </v-combobox>
  </div>
</template>

<script>
import UserAvatarContent from "@/components/UserAvatarContent.vue"
import { validateEmail, get, post } from "@/utils"

export default {
  name: "EmailReminders",

  props: {
    addedEmails: {
      type: Array,
      default: () => [],
    },
  },

  data: () => ({
    emails: [], // Currently displayed emails
    searchedContacts: [], // Contacts that match the search query

    hasContactsAccess: true,
    rules: {
      validEmails: (emails) => {
        for (const email of emails) {
          if (email?.length > 0 && !validateEmail(email)) {
            return "Please enter a valid email."
          }
        }
        return true
      },
    },
  }),

  mounted() {
    // Send a warmup request to update cache and check if contacts permissions are enabled
    get(`/user/searchContacts?query=`).catch((err) => {
      // User has not granted contacts permissions
      console.log("ENCOUNTERED ERROR")
      if (err.error?.code === 403) {
        this.hasContactsAccess = false
      }
    })

    this.emails = this.addedEmails;
  },

  methods: {
    requestContactsAccess() {
      this.$emit("requestContactsAccess", {
        // https://github.com/schej-it/schej.it/commit/8ff5aef2d2301478e44022651831b95015160d86
        emails: this.emails,
      })
    },
    searchContacts(emailsIndex, query) {
      // Searches the user's contacts using the google contacts API
      if (this.hasContactsAccess) {
        get(`/user/searchContacts?query=${query}`).then((results) => {
          this.searchedContacts = results
        })
      }
    },
    removeEmail(email) {
      this.emails.splice(this.emails.indexOf(email), 1)
      _delete(`/events/${this.eventId}/attendee`, { email })
    }
  },

  watch: {
    emails() {
      if (this.rules.validEmails(this.emails)) {

      }
    }
  },  

  components: { UserAvatarContent },
}
</script>
