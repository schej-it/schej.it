<template>
  <div>
    <div class="tw-flex tw-gap-1">
      <div class="tw-text-very-dark-gray">Set up email reminders</div>

      <v-tooltip top>
        <template v-slot:activator="{ on, attrs }">
          <v-icon small v-bind="attrs" v-on="on"
            >mdi-information-outline
          </v-icon>
        </template>
        <div>
          Reminder emails will be sent the day of event creation,<br />one day
          after, and three days after.
        </div>
      </v-tooltip>
    </div>
    <div class="tw-mt-1 tw-text-xs tw-text-dark-gray" v-if="true">
      <a class="tw-underline" @click="requestContactsAccess"
        >Enable contacts access</a
      >
      for email auto-suggestions.
    </div>

    <!-- {{ this.searchedContacts }} -->

    <v-combobox
      v-model="remindees"
      :search-input.sync="query"
      :items="searchedContacts"
      item-text="email"
      item-value="email"
      class="tw-mt-2 tw-text-sm"
      placeholder="Enter an email address..."
      multiple
      append-icon=""
      solo
      return-object
      :rules="[rules.validEmails]"
    >
      <template v-slot:selection="data, parent">
        <v-chip
          :key="isContact(data.item) ? data.item.email : data.item"
          size="x-small"
          class="tw-flex tw-items-center tw-bg-light-gray tw-text-very-dark-gray"
        >
            <v-avatar class="bg-accent text-uppercase tw-mr-2" start
              ><img
                v-if="typeof data.item === 'object' && data.item.picture.length > 0"
                :src="data.item.picture"
                referrerpolicy="no-referrer"
                width="10px"
              />
              <v-icon v-else>mdi-account</v-icon></v-avatar
            >
          {{ isContact(data.item) ? data.item.email : data.item }}

          <v-icon small @click="() => removeEmail(data.item)" class="tw-ml-1"
            >mdi-close</v-icon
          >
        </v-chip>
      </template>
      <template v-slot:item="{ item }">
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
      </template>
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
    remindees: [], // Currently displayed emails
    searchedContacts: [], // Contacts that match the search query
    timeout: null, // Timeout for search debouncing

    hasContactsAccess: true,
    query: "",

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

    this.remindees = this.addedEmails
    console.log(this.remindees)
  },

  methods: {
    requestContactsAccess() {
      this.$emit("requestContactsAccess", {
        emails: this.remindees,
      })
    },
    searchContacts() {
      // Searches the user's contacts using the google contacts API
      if (this.hasContactsAccess) {
        if (this.timeout) clearTimeout(this.timeout)
        this.timeout = setTimeout(() => {
          get(`/user/searchContacts?query=${this.query}`).then((results) => {
            this.searchedContacts = results
            console.log(this.searchedContacts)
          })
        }, 500)
      }
    },
    removeEmail(email) {
      this.remindees.splice(this.remindees.indexOf(email), 1)
    },
    isContact(contact) {
      return typeof contact === 'object'
    },
  },

  watch: {
    remindees() {
      this.$emit("update:emails", this.remindees)
    },
    query() {
      if (this.query && this.query.length > 0) {
        this.searchContacts()
      } else {
        this.searchedContacts = []
      }
    },
  },

  components: { UserAvatarContent },
}
</script>
