<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'
import TimeSinceDisplay from '@/components/TimeSinceDisplay.vue'

import { useProjectsStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'
import { useToastError } from '@/composables/useToastError'
import { deleteProject } from '@/api/project'
import { useUsersStore } from '@/stores/users'
import { useAuthStore } from '@/stores/auth'
import { useActiveAttacksStore } from '@/stores/activeAttacks'

const projectsStore = useProjectsStore()
projectsStore.load(true)

const activeAttacksStore = useActiveAttacksStore()

const authStore = useAuthStore()
const { loggedInUser } = storeToRefs(authStore)

const usersStore = useUsersStore()

const { projects } = storeToRefs(projectsStore)

const toast = useToast()
const { catcher } = useToastError()

async function onDeleteProject(id: string) {
  try {
    await deleteProject(id)
    toast.info('Deleted project')
  } catch (e: any) {
    catcher(e)
  } finally {
    projectsStore.load(true)
  }
}

const quantityStr = (num: number, str: string) => {
  if (num == 1) {
    return `${num} ${str}`
  }
  return `${num} ${str}s`
}
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Projects</h1>
    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card w-full bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">My Projects & Projects Shared with Me</h2>
          <table class="table w-full">
            <thead>
              <tr>
                <th>Name</th>
                <th>Created</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody class="first-col-bold">
              <RouterLink custom v-slot="{ navigate }" v-for="project in projects" :key="project.id" :to="`/project/${project.id}`">
                <tr class="hover">
                  <td class="cursor-pointer" @click="navigate">
                    {{ project.name }}
                    <div
                      class="badge badge-neutral float-right my-1 mr-1 whitespace-nowrap font-normal"
                      v-if="activeAttacksStore.initialisingAttacksForProject(project.id).length > 0"
                    >
                      {{ quantityStr(activeAttacksStore.initialisingAttacksForProject(project.id).length, 'attack') }} processing
                    </div>
                    <div
                      class="badge badge-info float-right my-1 mr-1 whitespace-nowrap font-normal"
                      v-if="activeAttacksStore.jobsForProject(project.id).length > 0"
                    >
                      {{ quantityStr(activeAttacksStore.jobsForProject(project.id).length, 'job') }} running
                    </div>
                    <span class="ml-1 text-xs font-normal text-slate-500" v-if="project.owner_user_id != loggedInUser?.id">
                      <font-awesome-icon icon="fa-solid fa-link" /> Shared by
                      {{ usersStore.byId(project.owner_user_id)?.username ?? 'Unknown user' }}
                    </span>
                  </td>
                  <td class="cursor-pointer" @click="navigate">
                    <TimeSinceDisplay :timestamp="project.time_created * 1000" />
                  </td>
                  <td class="w-0 text-center">
                    <ConfirmModal @on-confirm="() => onDeleteProject(project.id)">
                      <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                    </ConfirmModal>
                  </td>
                </tr>
              </RouterLink>
            </tbody>
          </table>
          <p class="mt-2">
            <span class="text-sm text-slate-500"> <font-awesome-icon icon="fa-solid fa-link" /> = Shared </span>
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
