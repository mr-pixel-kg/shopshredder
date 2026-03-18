<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  modelValue?: string | null;
  placeholder?: string;
  rows?: number;
  class?: string;
}>();

defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <textarea
    :value="modelValue ?? ''"
    :rows="rows ?? 4"
    :placeholder="placeholder"
    :class="cn('flex min-h-[110px] w-full rounded-xl border border-input bg-white/80 px-3 py-3 text-sm outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20', $props.class)"
    @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
</template>
