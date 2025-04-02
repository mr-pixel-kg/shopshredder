<script>
import Card from "primevue/card";
import Button from "primevue/button";
import ActiveSandboxImageCard from "@/components/home/SandboxImageCard.vue";
import ActiveSandboxCard from "@/components/home/SandboxCard.vue";
import { SandboxImageModel } from "@/models/SandboxImageModel.js";

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
            "mpXvoucher",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "mpXvoucher für Shopware 6: Gutscheine erstellen, verwalten & verkaufen – individuell, sicher & einfach.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "mpXcolorTabs",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Ermöglicht das Zuordnen einer individuellen Tab-Hintergrundfarbe und einer Tab-Textfarbe für Kategorien.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "mpXcheckout",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "mpXcheckout erweitert Shopware für Gastronomie: Lieferservice, Vorbestellungen & automatische Zeitberechnung.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Gastro Berlin Theme",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Berlin Theme für mpXcheckout. Es wurde speziell an die Anforderungen eines Shops der Gastronomie erstellt.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Gastro Stockholm Theme",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Stockholm Theme für mpXcheckout. Es wurde speziell an die Anforderungen eines Shops der Gastronomie erstellt.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Gastro Rom Theme",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Rom Theme für mpXcheckout. Es wurde speziell an die Anforderungen eines Shops der Gastronomie erstellt.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Fester Link für Medien",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Das Plugin erstellt statische Medien-Links und leitet diese automatisch auf die aktuelle Datei weiter.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "mpXFAQPro",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "mpXfaqpro zeigt FAQs direkt am Artikel & ermöglicht individuelle Kundenfragen mit Backend-Management.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "mpXDirectOrder",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "B2B-Schnellbestellungen per CSV-Upload – wiederkehrende Bestellungen mit einem Klick erledigen.",
            "https://www.mr-pixel.de",
        ),

        /*new SandboxImageModel(
        "Shopware 6.6 Play",
        "shopware-banner.jpg",
        "dockware/play:6.6.10.0",
        "Eine sofort einsatzbereite, saubere Shopware 6.6.10.0 Installation zum Testen, Experimentieren und Erkunden.",
        "https://www.mr-pixel.de",
      ),*/

        new SandboxImageModel(
            "Shopware 6.6",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.10.0",
            "Eine sofort einsatzbereite, saubere Shopware 6.6.10.0 Installation zum Testen, Experimentieren und Erkunden.",
            "https://www.mr-pixel.de",
        ),
        new SandboxImageModel(
            "Shopware 6.7 RC-1",
            "shopware-banner.jpg",
            "mr-pixel/sw-sandbox:6.6.7.0-rc1",
            "Shopware 6.7.0.0 RC-1 welcher es dir erlaubt Plugins und Themes auf der neuen Shopware Version zu testen.",
            "https://www.mr-pixel.de",
        ),
      ],
    };
  },
};
</script>

<template>
  <div class="max-w-7xl mx-auto p-5">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-gray-200">Shopware Umgebungen</h1>
    </div>
    <div
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      id="theme-grid"
    >
      <ActiveSandboxImageCard
        v-for="sandboxImage in sandboxImages"
        :sandboxImage="sandboxImage"
      />
    </div>
  </div>
</template>

<style scoped></style>
