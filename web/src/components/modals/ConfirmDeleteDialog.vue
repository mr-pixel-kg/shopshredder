<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'

defineProps<{
  open: boolean
  sandboxName: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  confirm: []
}>()

function handleConfirm() {
  emit('confirm')
  emit('update:open', false)
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-[400px]">
      <DialogHeader>
        <DialogTitle>Sandbox beenden</DialogTitle>
        <DialogDescription>
          Bist du sicher, dass du <strong>{{ sandboxName }}</strong> beenden möchtest?
          Diese Aktion kann nicht rückgängig gemacht werden.
        </DialogDescription>
      </DialogHeader>
      <DialogFooter>
        <Button variant="outline" @click="emit('update:open', false)">Abbrechen</Button>
        <Button variant="destructive" @click="handleConfirm">Beenden</Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
