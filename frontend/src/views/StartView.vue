<script>
import Button from 'primevue/button';
import Select from 'primevue/select';

export default {
  // Properties returned from data() become reactive state
  // and will be exposed on `this`.
  components: {
    Button,
    Select
  },
  data() {
    return {
      shopwareVersion: "v6.6.0.0",
      versions: [
        { name: 'v6.6', code: 'v6.6.0.0' },
        { name: 'v6.5', code: 'v6.6.5.0' },
        { name: 'v6.4', code: 'v6.6.4.0' },
      ],
      plugin: null,
      plugins: [
        { name: 'MrpixCloudPrint', repo: 'todo' },
        { name: 'MrpixGastronomy', repo: 'todo' },
        { name: 'MrpixColorTabs', repo: 'todo' },
      ]
    }
  }
}
</script>

<template>
  <h1>Create Shopware Sandbox...</h1>
  <Select v-model="shopwareVersion" :options="versions" optionLabel="name" placeholder="Select Shopware Version" class="w-full md:w-56" />
  <Select v-model="plugin" :options="plugins" optionLabel="name" placeholder="Select Plugin" class="w-full md:w-56" />
  <Button label="Start" />
</template>

<style scoped>

</style>