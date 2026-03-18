<script setup lang="ts">
import Card from "@/components/ui/Card.vue";
import Badge from "@/components/ui/Badge.vue";
import { formatDateTime } from "@/lib/utils";
import type { AuditLogRecord } from "@/types/api";

defineProps<{
  logs: AuditLogRecord[];
}>();
</script>

<template>
  <Card class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="section-title text-xl">Audit log</h2>
        <p class="text-sm text-muted-foreground">Recent backend actions for troubleshooting and traceability.</p>
      </div>
      <Badge>{{ logs.length }} entries</Badge>
    </div>

    <div class="overflow-hidden rounded-2xl border border-border/70">
      <table class="min-w-full divide-y divide-border text-sm">
        <thead class="bg-secondary/60 text-left text-muted-foreground">
          <tr>
            <th class="px-4 py-3 font-medium">Time</th>
            <th class="px-4 py-3 font-medium">Action</th>
            <th class="px-4 py-3 font-medium">User</th>
            <th class="px-4 py-3 font-medium">IP</th>
            <th class="px-4 py-3 font-medium">Details</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border bg-white/80">
          <tr v-for="entry in logs" :key="entry.id">
            <td class="px-4 py-3">{{ formatDateTime(entry.createdAt) }}</td>
            <td class="px-4 py-3 font-medium">{{ entry.action }}</td>
            <td class="px-4 py-3">{{ entry.userId || "Guest" }}</td>
            <td class="px-4 py-3">{{ entry.ipAddress || "-" }}</td>
            <td class="max-w-lg px-4 py-3 text-xs text-muted-foreground">
              <pre class="whitespace-pre-wrap break-all">{{ JSON.stringify(entry.details ?? {}, null, 2) }}</pre>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </Card>
</template>
