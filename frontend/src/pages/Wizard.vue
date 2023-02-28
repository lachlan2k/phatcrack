<script setup lang="ts">
import JobWizard from '@/components/JobWizard.vue'
</script>

<template>
  <main class="w-full p-6">
    <div class="prose">
      <h1>New Project Wizard</h1>
    </div>
    <JobWizard :first-step="0" />
  </main>
</template>