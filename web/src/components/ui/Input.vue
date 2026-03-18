<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{
  modelValue?: string | number | null;
  type?: string;
  placeholder?: string;
  class?: string;
}>();

defineEmits<{
  "update:modelValue": [value: string];
}>();
</script>

<template>
  <input
    :type="type ?? 'text'"
    :value="modelValue ?? ''"
    :placeholder="placeholder"
    :class="cn('flex h-11 w-full rounded-xl border border-input bg-white/80 px-3 text-sm outline-none transition focus:border-primary focus:ring-2 focus:ring-primary/20', $props.class)"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>
