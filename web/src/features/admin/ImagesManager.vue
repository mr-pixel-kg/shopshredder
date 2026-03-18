<script setup lang="ts">
import { reactive } from "vue";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import Input from "@/components/ui/Input.vue";
import Textarea from "@/components/ui/Textarea.vue";
import Badge from "@/components/ui/Badge.vue";
import type { CreateImagePayload, ImageRecord } from "@/types/api";

defineProps<{
  images: ImageRecord[];
  creating: boolean;
  deletingId?: string | null;
}>();

const emit = defineEmits<{
  create: [payload: CreateImagePayload];
  remove: [id: string];
}>();

const form = reactive<CreateImagePayload>({
  name: "",
  tag: "latest",
  title: "",
  description: "",
  thumbnailUrl: "",
  isPublic: true,
});

function submit() {
  emit("create", {
    ...form,
    title: form.title || null,
    description: form.description || null,
    thumbnailUrl: form.thumbnailUrl || null,
  });
}
</script>

<template>
  <div class="grid gap-6 xl:grid-cols-[360px_minmax(0,1fr)]">
    <Card class="space-y-4">
      <div>
        <h2 class="section-title text-xl">Register image</h2>
        <p class="text-sm text-muted-foreground">The backend pulls the Docker image if it is not already available locally.</p>
      </div>

      <div class="space-y-3">
        <div class="space-y-1">
          <label class="text-sm font-medium">Name</label>
          <Input v-model="form.name" placeholder="ghcr.io/shopshredder/shopware-demo" />
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium">Tag</label>
          <Input v-model="form.tag" placeholder="latest" />
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium">Title</label>
          <Input v-model="form.title" placeholder="Shopware Demo" />
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium">Description</label>
          <Textarea v-model="form.description" placeholder="Public demo image for storefront previews" />
        </div>
        <div class="space-y-1">
          <label class="text-sm font-medium">Thumbnail URL</label>
          <Input v-model="form.thumbnailUrl" placeholder="https://..." />
        </div>
        <label class="flex items-center gap-3 rounded-xl bg-secondary/60 px-3 py-3 text-sm">
          <input v-model="form.isPublic" type="checkbox" class="h-4 w-4 rounded border-input" />
          Visible on the public storefront
        </label>
      </div>

      <Button class="w-full" :disabled="creating" @click="submit">Create image</Button>
    </Card>

    <Card class="space-y-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="section-title text-xl">Image catalog</h2>
          <p class="text-sm text-muted-foreground">All registered templates managed by the internal team.</p>
        </div>
        <Badge>{{ images.length }} images</Badge>
      </div>

      <div class="overflow-hidden rounded-2xl border border-border/70">
        <table class="min-w-full divide-y divide-border text-sm">
          <thead class="bg-secondary/60 text-left text-muted-foreground">
            <tr>
              <th class="px-4 py-3 font-medium">Image</th>
              <th class="px-4 py-3 font-medium">Title</th>
              <th class="px-4 py-3 font-medium">Public</th>
              <th class="px-4 py-3 font-medium">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-border bg-white/80">
            <tr v-for="image in images" :key="image.id">
              <td class="px-4 py-3 font-medium">{{ image.name }}:{{ image.tag }}</td>
              <td class="px-4 py-3">{{ image.title || "Untitled" }}</td>
              <td class="px-4 py-3">
                <Badge :tone="image.isPublic ? 'success' : 'neutral'">{{ image.isPublic ? "Yes" : "No" }}</Badge>
              </td>
              <td class="px-4 py-3">
                <Button variant="outline" size="sm" :disabled="deletingId === image.id" @click="$emit('remove', image.id)">
                  Delete
                </Button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </Card>
  </div>
</template>
