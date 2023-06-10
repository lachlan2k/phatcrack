<script setup lang="ts">
const props = defineProps<{
  currentPage: number
  totalPages?: number
}>()

const emit = defineEmits(['next', 'prev'])
</script>

<template>
  <div class="join">
    <button class="btn-sm join-item btn cursor-pointer" @click="() => emit('prev')">«</button>
    <button class="btn-sm join-item btn">
      Page {{ props.currentPage + 1
      }}{{ props.totalPages != null ? ` of ${props.totalPages}` : '' }}
    </button>
    <button class="btn-sm join-item btn cursor-pointer" @click="() => emit('next')">»</button>
  </div>
</template>
