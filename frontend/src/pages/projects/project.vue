<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import JobWizard from '@/components/Wizard/JobWizard.vue'

import { ref } from 'vue'
import { getProject, getHashlistsForProject } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { timeSince } from '@/util/units'

const projId = useRoute().params.id as string

const { data: projectData, isLoading: isLoadingProject } = useApi(() => getProject(projId))
const { data: hashlistData, isLoading: isLoadingHashlists } = useApi(() =>
  getHashlistsForProject(projId)
)

const isWizardOpen = ref(false)

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingProject.value || isLoadingHashlists.value || !isHashTypesLoaded.value
})
</script>

<template>
  <main class="w-full p-4">
    <p v-if="isLoading">Loading</p>
    <div v-else>
      <div class="prose">
        <h1>{{ projectData?.name }}</h1>
      </div>
      <div class="mt-6 flex flex-wrap gap-6">
        <div class="card bg-base-100 shadow-xl min-w-[400px]">
          <div class="card-body">
            <div class="flex flex-row justify-between">
              <h2 class="card-title">Hashlists</h2>

              <button class="btn-primary btn-sm btn" @click="() => (isWizardOpen = true)">
                New Hashlist
              </button>

              <Modal v-model:isOpen="isWizardOpen">
                <JobWizard :firstStep="1" :existingProjectId="projId" />
              </Modal>
            </div>
            <table class="table w-full">
              <!-- head -->
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Hash Type</th>
                  <th>Created</th>
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
                    <td>{{ timeSince(hashlist.time_created * 1000) }}</td>
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
