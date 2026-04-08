<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'
import { ComboboxInput as RekaComboboxInput } from 'reka-ui'
import { ref, watch } from 'vue'

import {
  Combobox,
  ComboboxAnchor,
  ComboboxEmpty,
  ComboboxItem,
  ComboboxList,
  ComboboxViewport,
} from '@/components/ui/combobox'
import { Input } from '@/components/ui/input'

export interface Suggestion {
  value: string
  label: string
  description?: string
}

const props = withDefaults(
  defineProps<{
    modelValue: string
    suggestions: Suggestion[]
    loading?: boolean
    placeholder?: string
    disabled?: boolean
    minChars?: number
    id?: string
  }>(),
  {
    loading: false,
    placeholder: '',
    disabled: false,
    minChars: 2,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const open = ref(false)
const inputValue = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    if (val !== inputValue.value) {
      inputValue.value = val
    }
  },
)

watch(inputValue, (val) => {
  emit('update:modelValue', val)
  open.value = val.length >= props.minChars
})

function onSelect(item: Suggestion) {
  inputValue.value = item.value
  emit('update:modelValue', item.value)
  open.value = false
}
</script>

<template>
  <Combobox
    v-model:open="open"
    :ignore-filter="true"
    :reset-search-term-on-blur="false"
    :reset-search-term-on-select="false"
    :open-on-focus="modelValue.length >= minChars"
  >
    <ComboboxAnchor class="w-full">
      <RekaComboboxInput as-child>
        <Input
          :id="id"
          v-model="inputValue"
          :placeholder="placeholder"
          :disabled="disabled"
          autocomplete="off"
        />
      </RekaComboboxInput>
    </ComboboxAnchor>

    <ComboboxList
      v-if="open && (suggestions.length > 0 || loading)"
      class="w-[var(--reka-combobox-trigger-width)]"
    >
      <ComboboxViewport class="max-h-[240px]">
        <div
          v-if="loading && suggestions.length === 0"
          class="text-muted-foreground flex items-center justify-center gap-2 py-4 text-sm"
        >
          <Loader2 class="h-4 w-4 animate-spin" />
          Suche...
        </div>

        <ComboboxEmpty v-if="!loading"> Keine Ergebnisse </ComboboxEmpty>

        <ComboboxItem
          v-for="item in suggestions"
          :key="item.value"
          :value="item"
          class="flex-col items-start"
          @select.prevent="onSelect(item)"
        >
          <span>{{ item.label }}</span>
          <span v-if="item.description" class="text-muted-foreground line-clamp-1 text-xs">{{
            item.description
          }}</span>
        </ComboboxItem>
      </ComboboxViewport>
    </ComboboxList>
  </Combobox>
</template>
