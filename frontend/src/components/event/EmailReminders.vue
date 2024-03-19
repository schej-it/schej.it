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
          after, and three days after. You will also receive <br />an email when
          everybody has filled out the Schej.
        </div>
      </v-tooltip>
    </div>
    <div class="tw-mt-1 tw-text-xs tw-text-dark-gray" v-if="!hasContactsAccess">
      <a class="tw-underline" @click="requestContactsAccess"
        >Enable contacts access</a
      >
      for email auto-suggestions.
    </div>

    <v-combobox
      v-model="remindees"
      :search-input.sync="query"
      :items="searchedContacts"
      item-text="queryString"
      item-value="queryString"
      class="tw-mt-2 tw-text-sm"
      placeholder="Enter an email address..."
      multiple
      append-icon=""
      solo
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
              v-if="isContact(data.item) && data.item.picture.length > 0"
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
    searchDebounceTime: 250,

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
      if (err.error?.code === 403 || err.error?.code === 401) {
        this.hasContactsAccess = false
      }
    })

    this.remindees = this.addedEmails
  },

  methods: {
    /**
     * Requests access to contacts.
     */
    requestContactsAccess() {
      this.$emit("requestContactsAccess", {
        emails: this.remindees,
      })
    },
    /**
     * Searches contacts based on the query string if the user has access to contacts.
     */
    searchContacts() {
      if (this.hasContactsAccess) {
        if (this.timeout) clearTimeout(this.timeout)
        this.timeout = setTimeout(() => {
          get(`/user/searchContacts?query=${this.query}`).then((results) => {
            this.searchedContacts = results
            this.searchedContacts.map((contact) => {
              contact["queryString"] = this.contactToQueryString(contact)
            })
          })
        }, this.searchDebounceTime)
      }
    },
    /**
     * Removes the specified email from the remindees list.
     */
    removeEmail(email) {
      this.remindees.splice(this.remindees.indexOf(email), 1)
    },
    /**
     * Check if the contact is an object and not a user inputed string.
     */
    isContact(contact) {
      return typeof contact === "object"
    },
    /**
     * Takes a contact object and converts it to a query string.
     */
    contactToQueryString(contact) {
      // Need to split first name to get rid of middle name
      return `${contact["firstName"].split(" ")[0]} ${contact["lastName"]} ${
        contact["email"]
      }`
    },
  },

  watch: {
    remindees() {
      this.$emit(
        "update:emails",
        this.remindees.map((r) => (this.isContact(r) ? r.email : r))
      )
    },
    query() {
      if (this.query && this.query.length > 0) {
        if (this.query[this.query.length - 1] == " ") {
          const pureEmail = this.query.substring(0, this.query.length - 1)
          if (!this.remindees.includes(pureEmail)) this.remindees.push(pureEmail)
          this.query = ""
        } else {
          this.searchContacts()
        }
      } else {
        clearTimeout(this.timeout)
        this.searchedContacts = []
      }
    },
  },

  components: { UserAvatarContent },
}
</script>
