<script>
import Card from "primevue/card";
import Button from "primevue/button";
import ActiveSandboxImageCard from "@/components/home/ActiveSandboxImageCard.vue";
import ActiveSandboxCard from "@/components/home/ActiveSandboxCard.vue";
import {SandboxImageModel} from "@/models/SandboxImageModel.js";

export default {
  name: "SandboxImageGallery",

  components: {
    ActiveSandboxCard,
    ActiveSandboxImageCard,
    Card,
    Button,
  },

  data() {
    return {
      sandboxImages: [
        new SandboxImageModel(
          "Shopware 6.6 Demo",
          "shopware-banner.jpg",
          "dockware/dev:6.6.10.0",
          "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis!",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Shopware 6.6 Play",
            "shopware-banner.jpg",
            "dockware/dev:6.6.10.0",
            "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Inventore sed consequuntur error repudiandae numquam deserunt quisquam repellat libero asperiores earum nam nobis!",
            "https://www.shopware.com",
        )
      ]
    };
  },
}
</script>

<template>
  <div class="max-w-7xl mx-auto bg-red-500 p-5">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Shopware Umgebung</h1>
    </div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" id="theme-grid">

      <ActiveSandboxImageCard
          v-for="sandboxImage in sandboxImages"
          :sandboxImage="sandboxImage"
      />

    </div>
  </div>
</template>

<style scoped>

</style>