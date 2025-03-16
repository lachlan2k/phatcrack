<script setup lang="ts">
import { computed } from 'vue'

import { timeSince } from '@/util/units'

const props = defineProps<{
  timestamp: number
}>()

const timestampDisplayed = computed(() => new Date(props.timestamp).toLocaleDateString(undefined, { hour: '2-digit', minute: '2-digit' }))
</script>

<template>
  <div class="tooltip" :data-tip="timestampDisplayed" v-if="props.timestamp > 0">
    <span>{{ timeSince(props.timestamp) }}</span>
  </div>
  <div v-else>
    <span>Unknown</span>
  </div>
</template>
