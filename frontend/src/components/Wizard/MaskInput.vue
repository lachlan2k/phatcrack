<script setup lang="ts">
const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <label class="label font-bold">Select Mask</label>
  <input
    type="text"
    placeholder="Mask"
    :value="props.modelValue"
    @input="emit('update:modelValue', $event.target.value)"
    class="input-bordered input w-full max-w-xs"
  />
</template>
