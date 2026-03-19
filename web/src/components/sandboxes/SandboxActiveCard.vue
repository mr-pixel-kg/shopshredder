<script setup lang="ts">
import type { Sandbox, Image } from '@/types'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import SandboxInstanceItem from './SandboxInstanceItem.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

defineProps<{
  sandboxes: Sandbox[]
  images: Image[]
}>()

const emit = defineEmits<{
  open: [sandbox: Sandbox]
  extend: [sandbox: Sandbox]
  delete: [sandbox: Sandbox]
}>()

function findImage(imageId: string, images: Image[]): Image | undefined {
  return images.find((i) => i.id === imageId)
}
</script>

<template>
  <Card>
    <CardHeader class="flex-row items-center justify-between space-y-0 pb-3">
      <CardTitle class="text-base">Aktiv</CardTitle>
      <Badge v-if="sandboxes.length > 0" variant="secondary">
        {{ sandboxes.length }} {{ sandboxes.length === 1 ? 'läuft' : 'laufen' }}
      </Badge>
    </CardHeader>
    <CardContent class="p-0">
      <EmptyState
        v-if="sandboxes.length === 0"
        title="Keine aktiven Sandboxes"
        description="Starte eine neue Sandbox, um loszulegen."
      />
      <div v-else class="divide-y">
        <SandboxInstanceItem
          v-for="sandbox in sandboxes"
          :key="sandbox.id"
          :sandbox="sandbox"
          :image="findImage(sandbox.imageId, images)"
          @open="emit('open', $event)"
          @extend="emit('extend', $event)"
          @delete="emit('delete', $event)"
        />
      </div>
    </CardContent>
  </Card>
</template>
