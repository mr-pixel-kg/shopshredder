<script setup>
import Toast from "primevue/toast";
import Header from "@/components/layout/Header.vue";
</script>

<template>
  <Toast />
  <Header></Header>

  <main>
    <RouterView />
  </main>
</template>

<style scoped>

</style>
