<script setup lang="ts">
import { computed } from 'vue'
import { maskCharsets } from '@/util/hashcat'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

const value = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})
</script>

<template>
  <label class="label font-bold">Select Mask</label>
  <input
    type="text"
    placeholder="Mask"
    v-model="value"
    class="input-bordered input w-full max-w-xs"
  />
  <div class="mt-4">
    <span
      class="tooltip tooltip-bottom"
      :data-tip="`${maskCharset.description}`"
      v-for="maskCharset in maskCharsets"
      :key="maskCharset.mask"
    >
      <button
        text="foo"
        @click="value += maskCharset.mask"
        class="btn-outline btn-xs btn mr-1 normal-case"
      >
        {{ maskCharset.mask }}
      </button>
    </span>
  </div>
</template>
