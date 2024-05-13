<script setup lang="ts">
import AgentConfig from '@/components/Admin/AgentConfig.vue'
import AgentsTable from '@/components/Admin/AgentsTable.vue'
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Agent Management</h1>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card w-full bg-base-100 shadow-xl">
        <div class="card-body">
          <AgentsTable />
        </div>
      </div>

      <div class="card w-full bg-base-100 shadow-xl">
        <div class="card-body">
          <AgentConfig />
        </div>
      </div>
    </div>
  </main>
</template>
