<script setup lang="ts">
import Badge from "@/components/ui/Badge.vue";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import { formatDateTime, relativeRemaining } from "@/lib/utils";
import type { SandboxRecord } from "@/types/api";

defineProps<{
  sandboxes: SandboxRecord[];
  deletingId?: string | null;
}>();

defineEmits<{
  remove: [id: string];
}>();
</script>

<template>
  <Card class="space-y-5">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h2 class="section-title text-xl">Your demo sandboxes</h2>
        <p class="text-sm text-muted-foreground">Tracked automatically via guest session cookie.</p>
      </div>
      <Badge tone="accent">{{ sandboxes.length }} active</Badge>
    </div>

    <div v-if="sandboxes.length === 0" class="rounded-2xl border border-dashed border-border bg-secondary/40 px-4 py-8 text-sm text-muted-foreground">
      You do not have any guest demo sandboxes yet.
    </div>

    <div v-else class="grid gap-4">
      <article v-for="sandbox in sandboxes" :key="sandbox.id" class="rounded-2xl border border-border/70 bg-white/80 p-4">
        <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
          <div class="space-y-2">
            <div class="flex flex-wrap items-center gap-2">
              <h3 class="font-semibold">{{ sandbox.containerName }}</h3>
              <Badge :tone="sandbox.status === 'running' ? 'success' : 'neutral'">{{ sandbox.status }}</Badge>
            </div>
            <a :href="sandbox.url" target="_blank" rel="noreferrer" class="text-sm font-medium text-primary hover:underline">
              {{ sandbox.url }}
            </a>
            <div class="text-sm text-muted-foreground">
              Expires {{ formatDateTime(sandbox.expiresAt) }} · {{ relativeRemaining(sandbox.expiresAt) }}
            </div>
          </div>

          <Button variant="outline" :disabled="deletingId === sandbox.id" @click="$emit('remove', sandbox.id)">
            Delete
          </Button>
        </div>
      </article>
    </div>
  </Card>
</template>
