<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import JobWizard from '@/components/Wizard/JobWizard.vue'

import { ref } from 'vue'
import { getProject, getHashlistsForProject, deleteHashlist } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { timeSince } from '@/util/units'
import { useToast } from 'vue-toastification'
import { useToastError } from '@/composables/useToastError'

const projId = useRoute().params.id as string

const { data: projectData, isLoading: isLoadingProject } = useApi(() => getProject(projId))
const { data: hashlistData, isLoading: isLoadingHashlists, fetchData: fetchHashlists } = useApi(() => getHashlistsForProject(projId))

const isWizardOpen = ref(false)

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingProject.value || isLoadingHashlists.value || !isHashTypesLoaded.value
})

const toast = useToast()
const { catcher } = useToastError()

async function onDeleteHashlist(id: string) {
  try {
    await deleteHashlist(id)
    toast.info('Deleted hashlist')
  } catch(e: any) {
    catcher(e)
  } finally {
    fetchHashlists()
  }
}
</script>

<template>
  <main class="w-full p-4">
    <p v-if="isLoading">Loading</p>
    <div v-else>
      <h1 class="text-4xl font-bold">{{ projectData?.name }}</h1>
      <div class="mt-6 flex flex-wrap gap-6">
        <div class="card w-full bg-base-100 shadow-xl">
          <div class="card-body">
            <div class="flex flex-row justify-between">
              <h2 class="card-title">Hashlists</h2>

              <button class="btn btn-primary btn-sm" @click="() => (isWizardOpen = true)">New Hashlist</button>

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
                  <th>Actions</th>
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
                  <tr class="hover">
                    <td @click="navigate" class="cursor-pointer">{{ hashlist.name }}</td>
                    <td>{{ hashlist.hash_type }} - {{ getHashTypeName(hashlist.hash_type) }}</td>
                    <td>{{ timeSince(hashlist.time_created * 1000) }}</td>
                    <td>
                      <ConfirmModal @on-confirm="() => onDeleteHashlist(hashlist.id)">
                        <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                      </ConfirmModal>
                    </td>
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
