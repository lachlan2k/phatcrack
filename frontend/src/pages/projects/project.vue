<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'

import Modal from '@/components/Modal.vue'
import ProjectShare from '@/components/ProjectShare.vue'
import IconButton from '@/components/IconButton.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import JobWizard from '@/components/Wizard/JobWizard.vue'
import { getProject, getHashlistsForProject, deleteHashlist, deleteProject } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { useToastError } from '@/composables/useToastError'
import { useProjectsStore } from '@/stores/projects'
import { useAuthStore } from '@/stores/auth'
import { useUsersStore } from '@/stores/users'
import { useActiveAttacksStore } from '@/stores/activeAttacks'
import TimeSinceDisplay from '@/components/TimeSinceDisplay.vue'

const projId = useRoute().params.id as string

const router = useRouter()
const usersStore = useUsersStore()

const activeAttacksStore = useActiveAttacksStore()

const authStore = useAuthStore()
const { loggedInUser, isAdmin } = storeToRefs(authStore)

const projectsStore = useProjectsStore()

const { data: projectData, isLoading: isLoadingProject } = useApi(() => getProject(projId))
const { data: hashlistData, isLoading: isLoadingHashlists, fetchData: fetchHashlists } = useApi(() => getHashlistsForProject(projId))

const isShareModalOpen = ref(false)
const isWizardOpen = ref(false)

const resourcesStore = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resourcesStore)
resourcesStore.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingProject.value || isLoadingHashlists.value || !isHashTypesLoaded.value
})

const toast = useToast()
const { catcher } = useToastError()

async function onDeleteHashlist(id: string) {
  try {
    await deleteHashlist(id)
    toast.info('Deleted hashlist')
  } catch (e: any) {
    catcher(e)
  } finally {
    fetchHashlists()
  }
}

async function onDeleteProject(id: string) {
  try {
    await deleteProject(id)
    toast.info('Deleted project')
    router.push('/dashboard')
  } catch (e: any) {
    catcher(e)
  } finally {
    projectsStore.load(true)
  }
}

const hasOwnereshipRights = computed(() => {
  const user = loggedInUser.value
  if (user == null) {
    return false
  }

  return isAdmin || user.id == projectData.value?.owner_user_id
})

const quantityStr = (num: number, str: string) => {
  if (num == 1) {
    return `${num} ${str}`
  }
  return `${num} ${str}s`
}
</script>

<template>
  <main class="h-full w-full p-4">
    <div v-if="isLoading" class="flex h-full w-full justify-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    <div v-else>
      <Modal v-model:isOpen="isShareModalOpen" v-if="hasOwnereshipRights">
        <ProjectShare :projectId="projId" />
      </Modal>

      <div class="flex justify-between">
        <div>
          <h1 class="inline text-4xl font-bold">{{ projectData?.name }}</h1>
          <span class="ml-2 text-xs font-normal text-slate-500" v-if="projectData != null && projectData.owner_user_id != loggedInUser?.id">
            <font-awesome-icon icon="fa-solid fa-link" /> Shared by
            {{ usersStore.byId(projectData?.owner_user_id)?.username ?? 'Unknown user' }}
          </span>
        </div>
        <div class="flex items-center">
          <button v-if="hasOwnereshipRights" class="btn btn-ghost btn-sm" @click="() => (isShareModalOpen = true)">
            <font-awesome-icon icon="fa-solid fa-link" />Share with others
          </button>
          <ConfirmModal v-if="hasOwnereshipRights" @on-confirm="() => onDeleteProject(projId)"
            ><button class="btn btn-ghost btn-sm"><font-awesome-icon icon="fa-solid fa-trash" />Delete</button></ConfirmModal
          >
        </div>
      </div>
      <div class="breadcrumbs pl-1 text-sm">
        <ul>
          <li>
            <RouterLink to="/dashboard"> Dashboard </RouterLink>
          </li>
          <li>This project</li>
        </ul>
      </div>
      <div class="mt-3 flex flex-wrap gap-6">
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
                    <td @click="navigate" class="cursor-pointer">
                      {{ hashlist.name }}
                      <div
                        class="badge badge-neutral float-right mr-1 whitespace-nowrap font-normal"
                        v-if="activeAttacksStore.initialisingAttacksForHashlist(hashlist.id).length > 0"
                      >
                        {{ quantityStr(activeAttacksStore.initialisingAttacksForHashlist(hashlist.id).length, 'attack') }} processing
                      </div>
                      <div
                        class="badge badge-info float-right mr-1 whitespace-nowrap font-normal"
                        v-if="activeAttacksStore.jobsForHashlist(hashlist.id).length > 0"
                      >
                        {{ quantityStr(activeAttacksStore.jobsForHashlist(hashlist.id).length, 'job') }} running
                      </div>
                    </td>
                    <td @click="navigate" class="cursor-pointer">{{ hashlist.hash_type }} - {{ getHashTypeName(hashlist.hash_type) }}</td>
                    <td @click="navigate" class="cursor-pointer">
                      <TimeSinceDisplay :timestamp="hashlist.time_created * 1000" />
                    </td>
                    <td class="w-0 text-center">
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
