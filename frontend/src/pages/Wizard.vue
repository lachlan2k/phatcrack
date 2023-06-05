<script setup lang="ts">
import JobWizard from '@/components/Wizard/JobWizard.vue'
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Cracking Wizard</h1>
    </div>
    <JobWizard :first-step="0" />
  </main>
</template>
