<script>
import Card from "primevue/card";
import Button from "primevue/button";
import {SandboxImageModel} from "@/models/SandboxImageModel.js";

export default {
  name: "ActiveSandboxImageCard",

  props: {
    sandboxImage: {
      type: SandboxImageModel,
      required: true
    }
  },

  components: {
    Card,
    Button,
  },

  methods: {
    createDemo() {
      console.log("Create Demo");
    }
  }

}
</script>

<template>
  <Card style="overflow: hidden" class="w-full">
    <template #header>
      <img alt="Sandbox Image Thumbnail" :src="sandboxImage.thumbnailUrl" />
    </template>
    <template #title>{{ sandboxImage.title }}</template>
    <template #subtitle>{{ sandboxImage.imageName }}</template>
    <template #content>
      <p class="m-0">
        {{ sandboxImage.description }}
      </p>
    </template>
    <template #footer>
      <div class="flex gap-4 mt-1">
        <a :href="sandboxImage.infoLink" class="w-full">
          <Button label="Zum Store" severity="secondary" class="w-full" outlined />
        </a>
        <Button label="Demo" class="w-full" @click="createDemo" />
      </div>
    </template>
  </Card>
</template>

<style scoped>

</style>