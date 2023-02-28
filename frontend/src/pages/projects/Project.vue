<script setup lang="ts">
import { useProjectsStore } from '~/store/projects';

const route = useRoute()
const projectsStore = useProjectsStore()

</script>

<template>
  <main class="w-full p-6">
    <div class="prose">
      <h1>Project {{ route.params.id }}</h1>
    </div>
    <pre>
        {{ JSON.stringify(projectsStore.projects) }}
    </pre>
    <button class="btn btn-primary" @click="projectsStore.loadProjects">Load</button>
  </main>
</template>
