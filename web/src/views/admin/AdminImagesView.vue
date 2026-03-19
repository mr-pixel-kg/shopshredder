<script setup lang="ts">
import { ref } from 'vue'
import { useImages } from '@/composables/useImages'
import { getApiErrorMessage } from '@/utils/error'
import { toast } from 'vue-sonner'
import PageHeader from '@/components/shared/PageHeader.vue'
import AddImageDialog from '@/components/modals/AddImageDialog.vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableEmpty,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { Plus, Trash2 } from 'lucide-vue-next'

const { images, loading, createImage, deleteImage } = useImages('all')

const showAddImage = ref(false)

async function handleCreateImage(
  payload: { name: string; tag: string; title: string; description: string; isPublic: boolean },
  done: (success: boolean) => void,
) {
  try {
    await createImage(payload)
    toast.success('Vorlage wurde hinzugefügt')
    done(true)
  } catch (e) {
    toast.error(getApiErrorMessage(e, 'Fehler beim Hinzufügen'))
    done(false)
  }
}

async function handleDeleteImage(id: string) {
  try {
    await deleteImage(id)
    toast.success('Vorlage wurde gelöscht')
  } catch (e) {
    toast.error(getApiErrorMessage(e, 'Fehler beim Löschen'))
  }
}

function handleToggleVisibility() {
  // TODO: Call update API when available
  toast.info('Sichtbarkeit ändern ist noch nicht verfügbar')
}
</script>

<template>
  <div>
    <PageHeader title="Vorlagen" subtitle="Docker-Images als Sandbox-Vorlagen verwalten.">
      <template #actions>
        <Button @click="showAddImage = true">
          <Plus class="h-4 w-4 mr-1" />
          Vorlage hinzufügen
        </Button>
      </template>
    </PageHeader>

    <div class="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Vorlage</TableHead>
            <TableHead>Image</TableHead>
            <TableHead>Öffentlich</TableHead>
            <TableHead class="text-right">Aktionen</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableEmpty v-if="loading" :colspan="4">
            Wird geladen...
          </TableEmpty>
          <TableEmpty v-else-if="images.length === 0" :colspan="4">
            Keine Vorlagen vorhanden
          </TableEmpty>
          <TableRow v-for="image in images" :key="image.id">
            <TableCell>
              <div>
                <div class="font-medium">{{ image.title || image.name }}</div>
                <div v-if="image.description" class="text-xs text-muted-foreground">{{ image.description }}</div>
              </div>
            </TableCell>
            <TableCell>
              <Badge variant="secondary">{{ image.name }}:{{ image.tag }}</Badge>
            </TableCell>
            <TableCell>
              <Switch :checked="image.isPublic" @update:checked="handleToggleVisibility" />
            </TableCell>
            <TableCell class="text-right">
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger as-child>
                    <Button
                      variant="ghost"
                      size="icon"
                      class="text-destructive hover:text-destructive"
                      @click="handleDeleteImage(image.id)"
                    >
                      <Trash2 class="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Löschen</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <AddImageDialog
      v-model:open="showAddImage"
      @submit="handleCreateImage"
    />
  </div>
</template>
