<script setup lang="ts">
import { getProject, getHashlistsForProject } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

const projId = useRoute().params.id as string

const { data: projectData, isLoading: isLoadingProject } = useApi(() => getProject(projId))
const { data: hashlistData, isLoading: isLoadingHashlists } = useApi(() =>
  getHashlistsForProject(projId)
)

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingProject.value || isLoadingHashlists.value || !isHashTypesLoaded.value
})
</script>

<template>
  <main class="w-full p-6">
    <p v-if="isLoading">Loading</p>
    <div v-else>
      <div class="prose">
        <h1>{{ projectData?.name }}</h1>
      </div>
      <div class="mt-6 flex flex-col flex-wrap gap-6">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Hashlists</h2>
            <table class="table w-full">
              <!-- head -->
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Hash Type</th>
                </tr>
              </thead>

              <tbody class="first-col-bold">
                <RouterLink
                  custom
                  v-slot="{ navigate }"
                  v-for="hashlist in hashlistData?.hashlists"
                  :key="hashlist.id"
                  :to="`/hashlist/${hashlist.id}`"
                >
                  <tr class="hover cursor-pointer" @click="navigate">
                    <td>{{ hashlist.name }}</td>
                    <td>{{ hashlist.hash_type }} - {{ getHashTypeName(hashlist.hash_type) }}</td>
                  </tr>
                </RouterLink>
              </tbody>
            </table>
          </div>
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
