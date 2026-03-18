<script setup lang="ts">
import { computed } from "vue";
import { cn } from "@/lib/utils";

const props = defineProps<{
  tone?: "neutral" | "success" | "danger" | "accent";
  class?: string;
}>();

const classes = computed(() =>
  cn(
    "inline-flex items-center rounded-full px-2.5 py-1 text-xs font-semibold",
    props.tone === "success" && "bg-success/15 text-success",
    props.tone === "danger" && "bg-danger/15 text-danger",
    props.tone === "accent" && "bg-accent/15 text-accent",
    (!props.tone || props.tone === "neutral") && "bg-secondary text-secondary-foreground",
    props.class
  )
);
</script>

<template>
  <span :class="classes"><slot /></span>
</template>
