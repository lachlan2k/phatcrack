<script setup lang="ts">
import Config from '@/components/Admin/Config.vue'
import UsersTable from '@/components/Admin/UsersTable.vue'
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Admin</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <UsersTable />
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
