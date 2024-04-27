<script setup lang="ts">
import Config from '@/components/Admin/Config.vue'
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">General Settings</h1>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card w-full bg-base-100 shadow-xl">
        <div class="card-body justify-between">
          <Config />
        </div>
      </div>
    </div>
  </main>
</template>
