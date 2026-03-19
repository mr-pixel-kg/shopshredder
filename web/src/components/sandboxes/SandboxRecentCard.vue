<script setup lang="ts">
import type { Sandbox, Image } from '@/types'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import SandboxInstanceItem from './SandboxInstanceItem.vue'

defineProps<{
  sandboxes: Sandbox[]
  images: Image[]
}>()

const emit = defineEmits<{
  delete: [sandbox: Sandbox]
}>()
</script>

<template>
  <Card v-if="sandboxes.length > 0">
    <CardHeader class="pb-3">
      <CardTitle class="text-base">Zuletzt beendet</CardTitle>
    </CardHeader>
    <CardContent class="p-0">
      <div class="divide-y">
        <SandboxInstanceItem
          v-for="sandbox in sandboxes"
          :key="sandbox.id"
          :sandbox="sandbox"
          :image="images.find((i) => i.id === sandbox.imageId)"
          @delete="emit('delete', $event)"
        />
      </div>
    </CardContent>
  </Card>
</template>
