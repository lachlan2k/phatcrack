<script setup lang="ts">
import JobWizard from '@/components/Wizard/JobWizard.vue'
import { useRouter } from 'vue-router'

const router = useRouter()

function navToHashlist(hashlistId: string) {
  router.push(`/hashlist/${hashlistId}`)
}
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Cracking Wizard</h1>
    </div>
    <JobWizard :first-step="0" @successful-start="(details) => navToHashlist(details.hashlistId)" />
  </main>
</template>
