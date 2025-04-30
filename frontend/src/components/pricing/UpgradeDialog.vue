<template>
  <v-dialog
    :value="value"
    @input="(e) => $emit('input', e)"
    width="400"
    content-class="tw-m-0"
  >
    <v-card class="tw-p-4 sm:tw-p-6">
      <div>Upgrade to Pro</div>
      <div>One-time payment of {{ formattedPrice }}</div>
    </v-card>
  </v-dialog>
</template>

<script>
import { get } from "@/utils"
import { mapState } from "vuex"

export default {
  name: "UpgradeDialog",

  props: {
    value: { type: Boolean, required: true },
  },

  data() {
    return {
      loaded: false,
      price: null,
    }
  },

  computed: {
    ...mapState(["featureFlagsLoaded", "pricingPageConversion"]),
    formattedPrice() {
      if (!this.price) return ""
      return "$" + (this.price.unit_amount / 100).toFixed(2)
    },
  },

  methods: {
    async fetchPrice() {
      const res = await get("/stripe/price?exp=" + this.pricingPageConversion)
      this.price = res.price
      this.loaded = true
    },
  },

  watch: {
    featureFlagsLoaded: {
      handler(newVal) {
        if (newVal) {
          this.fetchPrice()
        }
      },
      immediate: true,
    },
  },
}
</script>
