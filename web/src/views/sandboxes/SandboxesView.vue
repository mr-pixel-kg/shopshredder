<script setup lang="ts">
import { ref } from 'vue'
import { useSandboxes } from '@/composables/useSandboxes'
import { useImages } from '@/composables/useImages'
import { getApiErrorMessage } from '@/utils/error'
import { toast } from 'vue-sonner'
import type { Sandbox } from '@/types'
import PageHeader from '@/components/shared/PageHeader.vue'
import SandboxActiveCard from '@/components/sandboxes/SandboxActiveCard.vue'
import SandboxRecentCard from '@/components/sandboxes/SandboxRecentCard.vue'
import NewSandboxDialog from '@/components/modals/NewSandboxDialog.vue'
import ExtendTtlDialog from '@/components/modals/ExtendTtlDialog.vue'
import ConfirmDeleteDialog from '@/components/modals/ConfirmDeleteDialog.vue'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-vue-next'

const { activeSandboxes, recentSandboxes, createSandbox, deleteSandbox, refresh } = useSandboxes()
const { images } = useImages()

const showNewSandbox = ref(false)
const showExtend = ref(false)
const showConfirmDelete = ref(false)
const selectedSandbox = ref<Sandbox | null>(null)

function handleOpen(sandbox: Sandbox) {
  if (sandbox.url) {
    window.open(sandbox.url, '_blank')
  }
}

function handleExtend(sandbox: Sandbox) {
  selectedSandbox.value = sandbox
  showExtend.value = true
}

function handleDelete(sandbox: Sandbox) {
  selectedSandbox.value = sandbox
  showConfirmDelete.value = true
}

async function handleCreateSandbox(
  payload: { imageId: string; ttlMinutes: number },
  done: (success: boolean) => void,
) {
  try {
    await createSandbox(payload)
    toast.success('Sandbox wird gestartet')
    refresh()
    done(true)
  } catch (e) {
    toast.error(getApiErrorMessage(e, 'Fehler beim Starten der Sandbox'))
    done(false)
  }
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
    <PageHeader title="Sandboxes" subtitle="Deine aktiven und kürzlich beendeten Sandboxes.">
      <template #actions>
        <Button @click="showNewSandbox = true">
          <Plus class="h-4 w-4 mr-1" />
          Neue Sandbox
        </Button>
      </template>
    </PageHeader>

    <div class="space-y-6">
      <SandboxActiveCard
        :sandboxes="activeSandboxes"
        :images="images"
        @open="handleOpen"
        @extend="handleExtend"
        @delete="handleDelete"
      />

      <SandboxRecentCard
        :sandboxes="recentSandboxes"
        :images="images"
        @delete="handleDelete"
      />
    </div>

    <NewSandboxDialog
      v-model:open="showNewSandbox"
      :images="images"
      @submit="handleCreateSandbox"
    />

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
