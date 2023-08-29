<script setup lang="ts">
import Config from '@/components/Admin/Config.vue'
import UsersTable from '@/components/Admin/UsersTable.vue'
import AgentsTable from '@/components/Admin/AgentsTable.vue'
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Admin</h1>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <UsersTable />
        </div>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <AgentsTable />
        </div>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body justify-between">
          <h2 class="card-title">Configuration</h2>
          <Config />
        </div>
      </div>
    </div>
  </main>
</template>
