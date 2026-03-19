<script setup lang="ts">
import type { Sandbox, Image } from '@/types'
import StatusDot from '@/components/shared/StatusDot.vue'
import TtlChip from './TtlChip.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ExternalLink, Clock, Square } from 'lucide-vue-next'

defineProps<{
  sandbox: Sandbox
  image?: Image
}>()

const emit = defineEmits<{
  open: [sandbox: Sandbox]
  extend: [sandbox: Sandbox]
  delete: [sandbox: Sandbox]
}>()

function isActive(status: string) {
  return status === 'running' || status === 'starting'
}
</script>

<template>
  <div class="flex items-center gap-3 py-3 px-4">
    <StatusDot :status="sandbox.status" />

    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-2">
        <span class="text-sm font-medium truncate">
          {{ image?.title || image?.name || sandbox.containerName }}
        </span>
        <Badge v-if="image?.tag" variant="secondary" class="text-xs">
          {{ image.tag }}
        </Badge>
      </div>
    </div>

    <TtlChip
      v-if="isActive(sandbox.status)"
      :expires-at="sandbox.expiresAt"
      :created-at="sandbox.createdAt"
    />

    <div class="flex items-center gap-1">
      <template v-if="isActive(sandbox.status)">
        <Button
          variant="ghost"
          size="sm"
          @click="emit('open', sandbox)"
        >
          <ExternalLink class="h-4 w-4 mr-1" />
          Öffnen
        </Button>
        <Button
          variant="ghost"
          size="sm"
          @click="emit('extend', sandbox)"
        >
          <Clock class="h-4 w-4 mr-1" />
          Verlängern
        </Button>
        <Button
          variant="ghost"
          size="sm"
          class="text-destructive hover:text-destructive"
          @click="emit('delete', sandbox)"
        >
          <Square class="h-4 w-4 mr-1" />
          Beenden
        </Button>
      </template>
      <template v-else>
        <Button
          variant="ghost"
          size="sm"
          class="text-destructive hover:text-destructive"
          @click="emit('delete', sandbox)"
        >
          Löschen
        </Button>
      </template>
    </div>
  </div>
</template>
