<script>
import Button from "primevue/button";
import Select from "primevue/select";
import LoadingScreen from "@/components/home/LoadingScreen.vue";
import { GeneralStore } from "@/stores/generalStore.js";
import { SandboxImage } from "@/models/SandboxImage.js";
import SandboxImageGallery from "@/components/home/SandboxImageGallery.vue";
import SandboxGallery from "@/components/home/SandboxGallery.vue";

export default {
  components: {
    SandboxGallery,
    SandboxImageGallery,
    LoadingScreen,
    Button,
    Select,
  },

  data() {
    return {
      sandboxImages: [
        new SandboxImage(
          "Shopware 6.6",
          "dockware/dev:6.6.10.0",
          "/shopware-banner.jpg",
        ),
        new SandboxImage(
          "Shopware Play",
          "dockware/dev:6.6.10.0",
          "/shopware-banner.jpg",
        ),
      ],
    };
  },

  setup() {
    const store = GeneralStore();
    return {
      generalStore: store,
    };
  },
};
</script>

<template>
  <div class="max-w-7xl mx-auto" v-if="generalStore.sandboxes.length">
    <SandboxGallery />
  </div>

  <div class="max-w-7xl mx-auto">
    <SandboxImageGallery />
  </div>

  <LoadingScreen></LoadingScreen>
</template>

<style scoped></style>
