<script setup lang="ts">
import { cn } from "@/lib/utils";

defineProps<{ class?: string }>();
</script>

<template>
  <section :class="cn('rounded-3xl border border-border/70 bg-card/95 p-5 shadow-panel', $props.class)">
    <slot />
  </section>
</template>
