<script>
import Card from "primevue/card";
import Button from "primevue/button";
import {SandboxEnvironment} from "@/models/SandboxEnvironment.js";

export default {
  name: "SandboxCard",

  components: {
    Card,
    Button,
  },

  props: {
    sandboxEnvironment: SandboxEnvironment
  },

  computed: {
    sandbox() {
      return this.sandboxEnvironment;
    }
  },

  methods: {
    openUrl() {
      window.open("https://" + this.sandbox.url, "_blank").focus();
    }
  }

}
</script>

<template>
  <Card style="width: 25rem; overflow: hidden">
    <template #title>Shopware Sandbox</template>
    <template #subtitle>{{ sandbox.image }}</template>
    <template #footer>
      <div class="flex w-full gap-2 mt-1">
        <Button icon="pi pi-trash" rounded severity="danger" aria-label="Cancel" class="w-1/3" />
        <Button severity="primary" aria-label="Cancel" class="w-full" label="Öffnen" @click="openUrl" />
      </div>
    </template>
  </Card>
</template>

<style scoped>

</style>