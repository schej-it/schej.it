<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    content-class="tw-max-w-xl"
  >
    <v-card>
      <v-card-title class="tw-flex">
        <div>Confirm emails</div>
        <v-spacer />
        <v-btn icon @click="$emit('input', false)">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text>
        <div class="tw-mb-4">
          Google Calendar invites will be sent to people at the following email
          addresses
        </div>
        <div class="tw-max-h-96 tw-overflow-y-auto tw-table-auto">
          <table class="tw-text-black tw-w-full tw-text-left">
            <thead>
              <tr class="tw-font-medium tw-bg-white">
                <th class="tw-pb-4 tw-bg-white tw-sticky tw-top-0 tw-z-10">
                  Name
                </th>
                <th class="tw-pb-4 tw-bg-white tw-sticky tw-top-0 tw-z-10">
                  Email
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(respondent, r) in respondents">
                <td class="tw-pb-4 tw-pr-4">
                  <div class="tw-flex tw-items-center">
                    <UserAvatarContent
                      v-if="respondent.email.length > 0"
                      :user="respondent"
                      class="tw-w-4 tw-h-4 -tw-ml-3 -tw-mr-1"
                    ></UserAvatarContent>
                    <v-icon v-else class="tw-ml-1 tw-mr-3" small>
                      mdi-account
                    </v-icon>

                    {{ respondent.firstName }} {{ respondent.lastName }}
                  </div>
                </td>
                <td class="tw-pr-4">
                  <div class="tw-pb-4" v-if="respondent.email.length > 0">
                    {{ respondent.email }}
                  </div>
                  <v-text-field
                    v-else
                    v-model="emails[r]"
                    class="tw-pt-2"
                    placeholder="Email"
                    outlined
                    dense
                    :rules="[rules.validEmail]"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" @click="confirm" :disabled="!confirmEnabled">
          Confirm
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import UserAvatarContent from "./UserAvatarContent.vue"
import { validateEmail } from "@/utils"

export default {
  name: "ConfirmEmailsDialog",

  props: {
    value: { type: Boolean, required: true },
    respondents: { type: Array, default: () => [] },
  },

  data: () => ({
    emails: [],
    rules: {
      validEmail: (email) => {
        if (email.length > 0 && !validateEmail(email)) {
          return "Please enter a valid email."
        }
        return true
      },
    },
  }),

  mounted() {
    this.emails = this.respondents.map((r) => r.email)
  },

  computed: {
    confirmEnabled() {
      // Only enable confirm button if all emails are valid
      for (const email of this.emails) {
        if (this.rules.validEmail(email) !== true) {
          return false
        }
      }

      return true
    },
  },

  methods: {
    confirm() {
      this.$emit("confirm", this.emails)
    },
  },

  components: { UserAvatarContent },
}
</script>
