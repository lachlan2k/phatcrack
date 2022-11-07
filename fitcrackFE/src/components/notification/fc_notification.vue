<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-list-item 
    :to="'/jobs/' + jobId"
    :class="['py-2', 'pr-6', { seen }]"
  >
    <v-list-item-avatar>
      <v-icon :color="seen ? null : type">
        {{ typeIcon(type) }}
      </v-icon>
    </v-list-item-avatar>
    <v-list-item-content>
      <div class="d-flex justify-space-between align-center mb-2">
        <v-list-item-title :class="['font-weight-bold', `${seen ? '' : type}--text`]">
          {{ title }}
        </v-list-item-title>
        <v-list-item-subtitle class="text-end nogrow">
          {{ $moment.utc(time).local().format('DD. MM. YYYY HH:mm') }}
        </v-list-item-subtitle>
      </div>
      <div class="body-2">
        {{ text }}
      </div>
    </v-list-item-content>
  </v-list-item>
</template>

<script>
  import { typeIcon } from '@/assets/scripts/iconMaps'
  export default {
    name: 'Notification',
    props: ['type', 'title', 'text', 'seen', 'time', 'jobId'],
    methods: { typeIcon }
  }
</script>

<style scoped>
.nogrow {
  flex: none;
}

.seen {
  opacity: .7;
}
</style>
