<script setup lang="ts">
import { ref, computed } from 'vue'
import { useSandboxes } from '@/composables/useSandboxes'
import { useImages } from '@/composables/useImages'
import { getApiErrorMessage } from '@/utils/error'
import { toast } from 'vue-sonner'
import { formatDateTime } from '@/utils/formatters'
import type { Sandbox, SandboxStatus } from '@/types'
import PageHeader from '@/components/shared/PageHeader.vue'
import StatusBadge from '@/components/shared/StatusBadge.vue'
import ExtendTtlDialog from '@/components/modals/ExtendTtlDialog.vue'
import ConfirmDeleteDialog from '@/components/modals/ConfirmDeleteDialog.vue'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Clock, Square } from 'lucide-vue-next'

const { sandboxes, deleteSandbox, loading } = useSandboxes('all')
const { images } = useImages('all')

const statusFilter = ref<string>('all')

const filteredSandboxes = computed(() => {
  if (statusFilter.value === 'all') return sandboxes.value
  const activeStatuses: SandboxStatus[] = ['running', 'starting']
  const inactiveStatuses: SandboxStatus[] = ['stopped', 'expired', 'deleted', 'failed']
  if (statusFilter.value === 'active') return sandboxes.value.filter((s) => activeStatuses.includes(s.status))
  if (statusFilter.value === 'inactive') return sandboxes.value.filter((s) => inactiveStatuses.includes(s.status))
  return sandboxes.value
})

const showExtend = ref(false)
const showConfirmDelete = ref(false)
const selectedSandbox = ref<Sandbox | null>(null)

function getImageName(imageId: string): string {
  const image = images.value.find((i) => i.id === imageId)
  return image?.title || image?.name || '—'
}

function handleExtend(sandbox: Sandbox) {
  selectedSandbox.value = sandbox
  showExtend.value = true
}

function handleDelete(sandbox: Sandbox) {
  selectedSandbox.value = sandbox
  showConfirmDelete.value = true
}

async function handleConfirmDelete() {
  if (!selectedSandbox.value) return
  try {
    await deleteSandbox(selectedSandbox.value.id)
    toast.success('Sandbox wurde beendet')
  } catch (e) {
    toast.error(getApiErrorMessage(e, 'Fehler beim Beenden'))
  }
}
</script>

<template>
  <div>
    <PageHeader title="Instanzen" subtitle="Alle Sandbox-Instanzen verwalten." />

    <div class="flex items-center gap-3 mb-4">
      <Select v-model="statusFilter">
        <SelectTrigger class="w-[160px]">
          <SelectValue placeholder="Alle Status" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">Alle Status</SelectItem>
          <SelectItem value="active">Aktiv</SelectItem>
          <SelectItem value="inactive">Abgelaufen</SelectItem>
        </SelectContent>
      </Select>
    </div>

    <div class="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Status</TableHead>
            <TableHead>Vorlage</TableHead>
            <TableHead>Gestartet</TableHead>
            <TableHead>Läuft ab</TableHead>
            <TableHead class="text-right">Aktionen</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableEmpty v-if="loading" :colspan="5">
            Wird geladen...
          </TableEmpty>
          <TableEmpty v-else-if="filteredSandboxes.length === 0" :colspan="5">
            Keine Instanzen gefunden
          </TableEmpty>
          <TableRow v-for="sandbox in filteredSandboxes" :key="sandbox.id">
            <TableCell>
              <StatusBadge :status="sandbox.status" />
            </TableCell>
            <TableCell class="font-medium">{{ getImageName(sandbox.imageId) }}</TableCell>
            <TableCell class="text-muted-foreground">{{ formatDateTime(sandbox.createdAt) }}</TableCell>
            <TableCell class="text-muted-foreground">
              {{ sandbox.expiresAt ? formatDateTime(sandbox.expiresAt) : '—' }}
            </TableCell>
            <TableCell class="text-right">
              <TooltipProvider>
                <div class="flex items-center justify-end gap-1">
                  <Tooltip v-if="sandbox.status === 'running' || sandbox.status === 'starting'">
                    <TooltipTrigger as-child>
                      <Button variant="ghost" size="icon" @click="handleExtend(sandbox)">
                        <Clock class="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Verlängern</TooltipContent>
                  </Tooltip>
                  <Tooltip>
                    <TooltipTrigger as-child>
                      <Button
                        variant="ghost"
                        size="icon"
                        class="text-destructive hover:text-destructive"
                        @click="handleDelete(sandbox)"
                      >
                        <Square class="h-4 w-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Beenden</TooltipContent>
                  </Tooltip>
                </div>
              </TooltipProvider>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <ExtendTtlDialog
      v-model:open="showExtend"
      :sandbox-name="selectedSandbox?.containerName ?? ''"
    />

    <ConfirmDeleteDialog
      v-model:open="showConfirmDelete"
      :sandbox-name="selectedSandbox?.containerName ?? ''"
      @confirm="handleConfirmDelete"
    />
  </div>
</template>
