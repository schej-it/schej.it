<template>
  <div>
    <h1>Test</h1>
    <UpgradeDialog :value="true" />
  </div>
</template>

<script>
import UpgradeDialog from "@/components/pricing/UpgradeDialog"

export default {
  name: "Test",
  components: {
    UpgradeDialog,
  },
  methods: {
    createCheckoutSession() {
      post("/stripe/create-checkout-session").then((res) => {
        console.log(res)
      })
    },
  },
}
</script>
