<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  rows?: number
  modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

const rows = props.rows ?? 8

const hashes = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})
</script>

<template>
  <textarea
    placeholder="Hashes"
    class="hashes-input textarea textarea-bordered w-full font-mono focus:outline-none"
    :rows="rows"
    v-model="hashes"
  ></textarea>
</template>

<style scoped>
textarea.hashes-input {
  background-image: linear-gradient(to bottom, rgba(87, 87, 87, 0.05) 50%, transparent 50%);
  background-repeat: repeat-y;

  background-size: 100% 50px;

  line-height: 25px;
  font-size: 15px;

  padding: 0;
  padding-left: 3px;

  white-space: pre;
  resize: none;

  background-attachment: local;
}
</style>
