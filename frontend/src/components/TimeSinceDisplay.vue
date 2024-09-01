<script setup lang="ts">
import { timeSince } from '@/util/units'
import { computed } from 'vue'

const props = defineProps<{
  timestamp: number
}>()

const timestampDisplayed = computed(() => new Date(props.timestamp).toLocaleDateString(undefined, { hour: '2-digit', minute: '2-digit' }))
</script>

<template>
  <div class="tooltip" :data-tip="timestampDisplayed">
    <span>{{ timeSince(props.timestamp) }}</span>
  </div>
</template>
