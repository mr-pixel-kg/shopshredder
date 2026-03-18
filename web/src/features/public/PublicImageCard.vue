<script setup lang="ts">
import { computed } from "vue";
import { ArrowUpRight, Clock3 } from "lucide-vue-next";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import type { ImageRecord } from "@/types/api";

const props = defineProps<{
  image: ImageRecord;
  busy?: boolean;
}>();

const emit = defineEmits<{
  demo: [imageId: string];
}>();

const title = computed(() => props.image.title || `${props.image.name}:${props.image.tag}`);
</script>

<template>
  <Card class="flex h-full flex-col gap-4">
    <div class="overflow-hidden rounded-2xl bg-secondary">
      <img
        v-if="image.thumbnailUrl"
        :src="image.thumbnailUrl"
        :alt="title"
        class="h-44 w-full object-cover"
      />
      <div v-else class="flex h-44 items-center justify-center bg-gradient-to-br from-primary/15 to-accent/10 text-center text-sm text-muted-foreground">
        No thumbnail configured
      </div>
    </div>

    <div class="space-y-2">
      <h3 class="text-lg font-bold">{{ title }}</h3>
      <p class="text-sm text-muted-foreground">
        {{ image.description || "No description available yet for this image." }}
      </p>
    </div>

    <div class="mt-auto flex items-center justify-between text-xs text-muted-foreground">
      <span class="inline-flex items-center gap-1">
        <Clock3 class="h-3.5 w-3.5" />
        Demo lifetime managed by the backend
      </span>
      <span>{{ image.tag }}</span>
    </div>

    <Button class="w-full" :disabled="busy" @click="emit('demo', image.id)">
      <ArrowUpRight class="mr-2 h-4 w-4" />
      Start demo
    </Button>
  </Card>
</template>
