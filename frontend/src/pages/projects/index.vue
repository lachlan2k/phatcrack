<script setup lang="ts">
import { getAllProjects } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { computed } from 'vue'
import { timeSince } from '@/util/timeSince'

const { data, isLoading } = useApi(getAllProjects)
const projects = computed(() => data.value?.projects)
</script>

<template>
  <main class="w-full p-6">
    <div class="prose">
      <h1>Project Folders</h1>
    </div>
    <p v-if="isLoading">Loading</p>
    <div class="mt-6 flex flex-col flex-wrap gap-6" v-else>
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
              <!-- row 1 -->
              <tr class="hover">
                <td>41234 - Customer A</td>
                <td>2 hours ago</td>
                <td>8</td>
                <td><div class="badge-info badge">4 attacks running</div></td>
              </tr>
              <tr class="hover" v-for="project in projects" :key="project.id">
                <td>{{ project.name }}</td>
                <td>{{ timeSince(project.time_created) }}</td>
                <td>8</td>
                <td><div class="badge-info badge">4 attacks running</div></td>
              </tr>
              <!-- row 2 -->
              <tr class="hover">
                <td>
                  41235 - Customer B
                  <span class="ml-2 text-slate-500">
                    <font-awesome-icon icon="fa-solid fa-link" />
                  </span>
                </td>
                <td>10 days ago</td>
                <td>3</td>
                <td><div class="badge-ghost badge">Idle</div></td>
              </tr>
              <!-- row 3 -->
              <tr class="hover">
                <td>41236 - Customer C</td>
                <td>2 months ago</td>
                <td>2</td>
                <td><div class="badge-ghost badge">Idle</div></td>
              </tr>
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
