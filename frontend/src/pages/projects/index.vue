<script setup lang="ts">
import { getAllProjects } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { computed } from 'vue'
import { timeSince } from '@/util/units'

const { data, isLoading } = useApi(getAllProjects)
const projects = computed(() => data.value?.projects)
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Project Folders</h1>
    </div>
    <p v-if="isLoading">Loading</p>
    <div class="mt-6 flex flex-wrap gap-6" v-else>
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">My Projects & Projects Shared with Me</h2>
          <table class="table w-full">
            <!-- head -->
            <thead>
              <tr>
                <th>Name</th>
                <th>Created</th>
                <th>Hashlists</th>
                <th>Attacks</th>
              </tr>
            </thead>
            <tbody class="first-col-bold">
              <RouterLink
                custom
                v-slot="{ navigate }"
                v-for="project in projects"
                :key="project.id"
                :to="`/project/${project.id}`"
              >
                <tr class="hover cursor-pointer" @click="navigate">
                  <td>{{ project.name }}</td>
                  <td>{{ timeSince(project.time_created * 1000) }}</td>
                  <td>foo</td>
                  <td>bar</td>
                  <!-- <td>8</td> -->
                  <!-- <td><div class="badge badge-info">4 attacks running</div></td> -->
                </tr>
              </RouterLink>
            </tbody>
          </table>
          <p class="mt-2">
            <span class="text-sm text-slate-500">
              <font-awesome-icon icon="fa-solid fa-link" /> = Shared
            </span>
          </p>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
thead > tr > th {
  background: none !important;
}

.first-col-bold > tr td:first-of-type {
  font-weight: bold;
}
</style>
