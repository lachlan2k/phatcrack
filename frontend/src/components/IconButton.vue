<script setup lang="ts">
const props = defineProps<{
  color: string
  icon: string
  tooltip: string
}>()
</script>

<template>
  <div class="tooltip" :data-tip="props.tooltip">
    <button :class="'hover:btn- btn-ghost btn-xs btn' + props.color">
      <font-awesome-icon :icon="icon" />
    </button>
  </div>
</template>
