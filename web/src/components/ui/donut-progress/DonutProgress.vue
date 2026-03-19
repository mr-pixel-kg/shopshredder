<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue?: number
  }>(),
  {
    modelValue: 0,
  },
)

// Circle geometry — matches Lucide icon viewBox (0 0 24 24)
const cx = 12
const cy = 12
const r = 10
const circumference = 2 * Math.PI * r
const dashoffset = computed(() => circumference * ((100 - Math.min(Math.max(props.modelValue, 0), 100)) / 100))
</script>

<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width="24"
    height="24"
    viewBox="0 0 24 24"
    fill="none"
    class="-rotate-90"
  >
    <!-- Track -->
    <circle
      :cx="cx"
      :cy="cy"
      :r="r"
      stroke="currentColor"
      stroke-width="2"
      opacity="0.2"
    />
    <!-- Progress arc -->
    <circle
      :cx="cx"
      :cy="cy"
      :r="r"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      :stroke-dasharray="circumference"
      :stroke-dashoffset="dashoffset"
      class="transition-[stroke-dashoffset] duration-300 ease-out"
    />
  </svg>
</template>
