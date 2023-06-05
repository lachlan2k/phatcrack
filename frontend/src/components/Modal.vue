<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits(['update:isOpen'])

const isOpen = computed({
  get: () => props.isOpen,
  set: (value: boolean) => emit('update:isOpen', value)
})
</script>

<template>
  <dialog :class="isOpen ? 'modal modal-open' : 'custom-modal modal'">
    <form method="dialog" class="remove-card-backgrounds modal-box">
      <button
        @click="() => (isOpen = false)"
        class="btn-ghost btn-sm btn-circle btn absolute right-2 top-2"
      >
        ✕
      </button>
      <slot></slot>
    </form>
    <form method="dialog" class="modal-backdrop">
      <button @click="() => (isOpen = false)">close</button>
    </form>
  </dialog>
</template>

<style scoped>
.modal-box {
  max-width: 90vw;
  width: auto;
}
</style>

<style>
.modal-box.remove-card-backgrounds .card {
  box-shadow: none !important;
}
</style>
