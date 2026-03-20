<script setup lang="ts">
import type { Sandbox } from '@/types'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import StatusBadge from '@/components/shared/StatusBadge.vue'
import TtlChip from '@/components/sandboxes/TtlChip.vue'
import ActionButton from './ActionButton.vue'
import type { CardAction } from './ActionButton.vue'

// TODO: Replace with dynamic schema from API
export interface SandboxCardMeta {
  label: string
  value: string
}

defineProps<{
  sandbox: Sandbox
  title: string
  actions?: CardAction[]
  metadata?: SandboxCardMeta[]
}>()
</script>

<template>
  <Card class="flex flex-col">
    <CardHeader>
      <div class="flex items-start justify-between gap-2">
        <CardTitle class="text-sm truncate">{{ title }}</CardTitle>
        <StatusBadge :status="sandbox.status" />
      </div>
    </CardHeader>
    <CardContent class="flex-1 space-y-2">
      <TtlChip
        v-if="sandbox.expiresAt"
        :expires-at="sandbox.expiresAt"
        :created-at="sandbox.createdAt"
      />
      <!-- TODO: Replace with dynamic schema from API -->
      <div v-if="metadata?.length" class="space-y-1">
        <div
          v-for="meta in metadata"
          :key="meta.label"
          class="flex items-center justify-between text-xs"
        >
          <span class="text-muted-foreground">{{ meta.label }}</span>
          <span class="font-mono">{{ meta.value }}</span>
        </div>
      </div>
    </CardContent>
    <!-- TODO: Replace with dynamic schema from API -->
    <CardFooter v-if="actions?.length" class="flex gap-2">
      <ActionButton
        v-for="action in actions"
        :key="action.label"
        :action="action"
      />
    </CardFooter>
  </Card>
</template>
