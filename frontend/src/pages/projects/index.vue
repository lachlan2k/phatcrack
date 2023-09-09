<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import ConfirmModal from '@/components/ConfirmModal.vue'

import { timeSince } from '@/util/units'
import { useProjectsStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'
import { useToastError } from '@/composables/useToastError'
import { deleteProject } from '@/api/project'

const projectsStore = useProjectsStore()
projectsStore.load()

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
    projectsStore.load()
  }
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
                  <td class="cursor-pointer" @click="navigate">{{ project.name }}</td>
                  <td>{{ timeSince(project.time_created * 1000) }}</td>
                  <td>
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
