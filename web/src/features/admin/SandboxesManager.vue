<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import Input from "@/components/ui/Input.vue";
import Textarea from "@/components/ui/Textarea.vue";
import Badge from "@/components/ui/Badge.vue";
import { formatDateTime, relativeRemaining } from "@/lib/utils";
import type { CreateImagePayload, CreateSandboxPayload, ImageRecord, SandboxRecord } from "@/types/api";

const props = defineProps<{
  images: ImageRecord[];
  sandboxes: SandboxRecord[];
  creating: boolean;
  deletingId?: string | null;
  snapshottingId?: string | null;
}>();

const emit = defineEmits<{
  create: [payload: CreateSandboxPayload];
  remove: [id: string];
  snapshot: [sandboxId: string, payload: CreateImagePayload];
}>();

const createForm = reactive<CreateSandboxPayload>({
  imageId: "",
  ttlMinutes: 120,
});

const snapshotForm = reactive<CreateImagePayload>({
  name: "",
  tag: "",
  title: "",
  description: "",
  thumbnailUrl: "",
  isPublic: false,
});

const snapshotTarget = ref<string | null>(null);

const selectableImages = computed(() => props.images.map((image) => ({
  value: image.id,
  label: `${image.name}:${image.tag}`,
})));

function createSandbox() {
  emit("create", {
    imageId: createForm.imageId,
    ttlMinutes: createForm.ttlMinutes || null,
  });
}

function openSnapshot(id: string) {
  snapshotTarget.value = id;
}

function submitSnapshot() {
  if (!snapshotTarget.value) return;

  emit("snapshot", snapshotTarget.value, {
    ...snapshotForm,
    title: snapshotForm.title || null,
    description: snapshotForm.description || null,
    thumbnailUrl: snapshotForm.thumbnailUrl || null,
  });
}
</script>

<template>
  <div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
    <Card class="space-y-4">
      <div>
        <h2 class="section-title text-xl">Start internal sandbox</h2>
        <p class="text-sm text-muted-foreground">Create editable employee sandboxes from registered images.</p>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <label class="text-sm font-medium">Image</label>
          <select v-model="createForm.imageId" class="flex h-11 w-full rounded-xl border border-input bg-white/80 px-3 text-sm outline-none">
            <option disabled value="">Select an image</option>
            <option v-for="option in selectableImages" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium">TTL in minutes</label>
          <Input v-model="createForm.ttlMinutes" type="number" placeholder="120" />
        </div>
      </div>

      <Button class="w-full" :disabled="creating || !createForm.imageId" @click="createSandbox">Create sandbox</Button>
    </Card>

    <Card class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="section-title text-xl">Running sandboxes</h2>
          <p class="text-sm text-muted-foreground">Includes guest demos and internal employee environments.</p>
        </div>
        <Badge>{{ sandboxes.length }} active</Badge>
      </div>

      <div class="grid gap-4">
        <article v-for="sandbox in sandboxes" :key="sandbox.id" class="rounded-2xl border border-border/70 bg-white/80 p-4">
          <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
            <div class="space-y-2">
              <div class="flex flex-wrap items-center gap-2">
                <h3 class="font-semibold">{{ sandbox.containerName }}</h3>
                <Badge :tone="sandbox.createdByUserId ? 'accent' : 'neutral'">{{ sandbox.createdByUserId ? "Employee" : "Guest" }}</Badge>
                <Badge :tone="sandbox.status === 'running' ? 'success' : 'neutral'">{{ sandbox.status }}</Badge>
              </div>
              <a :href="sandbox.url" target="_blank" rel="noreferrer" class="text-sm font-medium text-primary hover:underline">
                {{ sandbox.url }}
              </a>
              <p class="text-sm text-muted-foreground">
                Expires {{ formatDateTime(sandbox.expiresAt) }} · {{ relativeRemaining(sandbox.expiresAt) }}
              </p>
            </div>

            <div class="flex flex-wrap gap-2">
              <Button variant="secondary" size="sm" @click="openSnapshot(sandbox.id)">Snapshot</Button>
              <Button variant="outline" size="sm" :disabled="deletingId === sandbox.id" @click="$emit('remove', sandbox.id)">
                Delete
              </Button>
            </div>
          </div>
        </article>
      </div>
    </Card>
  </div>

  <Card v-if="snapshotTarget" class="mt-6 space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="section-title text-xl">Create snapshot image</h2>
        <p class="text-sm text-muted-foreground">Commit the selected sandbox into a reusable image template.</p>
      </div>
      <Button variant="ghost" @click="snapshotTarget = null">Close</Button>
    </div>

    <div class="grid gap-3 md:grid-cols-2">
      <div class="space-y-1">
        <label class="text-sm font-medium">Name</label>
        <Input v-model="snapshotForm.name" placeholder="ghcr.io/shopshredder/shopware-custom" />
      </div>
      <div class="space-y-1">
        <label class="text-sm font-medium">Tag</label>
        <Input v-model="snapshotForm.tag" placeholder="demo-v2" />
      </div>
      <div class="space-y-1 md:col-span-2">
        <label class="text-sm font-medium">Title</label>
        <Input v-model="snapshotForm.title" placeholder="Shopware Custom Demo V2" />
      </div>
      <div class="space-y-1 md:col-span-2">
        <label class="text-sm font-medium">Description</label>
        <Textarea v-model="snapshotForm.description" placeholder="Snapshot taken from a configured employee sandbox." />
      </div>
      <div class="space-y-1 md:col-span-2">
        <label class="text-sm font-medium">Thumbnail URL</label>
        <Input v-model="snapshotForm.thumbnailUrl" placeholder="https://..." />
      </div>
      <label class="flex items-center gap-3 rounded-xl bg-secondary/60 px-3 py-3 text-sm md:col-span-2">
        <input v-model="snapshotForm.isPublic" type="checkbox" class="h-4 w-4 rounded border-input" />
        Publish snapshot to the public storefront
      </label>
    </div>

    <Button :disabled="snapshottingId === snapshotTarget" @click="submitSnapshot">Commit snapshot</Button>
  </Card>
</template>
