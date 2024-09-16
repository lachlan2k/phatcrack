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
    <h1 class="text-4xl font-bold">Cracking Wizard</h1>
    <JobWizard :first-step="0" @successful-start="details => navToHashlist(details.hashlistId)" />
  </main>
</template>
