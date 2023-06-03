<script setup lang="ts">
import { maskCharsets } from '@/util/hashcat'

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

function update(val: string) {
  return emit('update:modelValue', val)
}
</script>

<template>
  <label class="label font-bold">Select Mask</label>
  <input
    type="text"
    placeholder="Mask"
    :value="props.modelValue"
    @input="update($event.target.value)"
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
        @click="update(props.modelValue + maskCharset.mask)"
        class="btn-outline btn btn-xs mr-1 normal-case"
      >
        {{ maskCharset.mask }}
      </button>
    </span>
  </div>
</template>
