<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import Modal from '@/components/Modal.vue'
import FileUpload from '@/components/FileUpload.vue'

import { onBeforeUnmount, onMounted, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { bytesToReadable } from '@/util/units'
import { useListfilesStore } from '@/stores/listfiles'
import ConfirmModal from '@/components/ConfirmModal.vue'
import type { ListfileDTO } from '@/api/types'
import { useAuthStore } from '@/stores/auth'
import { useToastError } from '@/composables/useToastError'
import { deleteListfile } from '@/api/listfiles'
import { useToast } from 'vue-toastification'

const listfilesStore = useListfilesStore()
const { load: loadListfiles } = listfilesStore
const { wordlists, rulefiles } = storeToRefs(useListfilesStore())

const isWordlistUploadOpen = ref(false)
const isRulefileUploadOpen = ref(false)

loadListfiles(true)

const refreshTimer = ref(0)

onMounted(() => {
  refreshTimer.value = setInterval(() => {
    loadListfiles(true)
  }, 1000 * 60) // every 60 seconds
})

onBeforeUnmount(() => {
  clearInterval(refreshTimer.value)
})

const authStore = useAuthStore()
const { loggedInUser, isAdmin } = storeToRefs(authStore)

function canDelete(listfile: ListfileDTO) {
  if (listfile.pending_delete) {
    return false
  }

  if (isAdmin) {
    return true
  }
  return listfile.created_by_user_id == loggedInUser.value?.id
}

function isGreyed(listfile: ListfileDTO) {
  return listfile.pending_delete || !listfile.available_for_use
}

const toast = useToast()
const { catcher } = useToastError()

async function onDeleteListfile(listfile: ListfileDTO) {
  try {
    await deleteListfile(listfile.id)
    toast.info(`Marked ${listfile.name} for deletion`)
  } catch (e: any) {
    catcher(e)
  } finally {
    loadListfiles(true)
  }
}
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Listfiles</h1>
    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <Modal v-model:isOpen="isWordlistUploadOpen">
              <FileUpload fileType="Wordlist" />
            </Modal>
            <h2 class="card-title">Wordlists</h2>
            <button class="btn btn-primary btn-sm" @click="() => (isWordlistUploadOpen = true)">Upload Wordlist</button>
          </div>

          <table class="table w-full">
            <thead>
              <tr>
                <th>Name</th>
                <th>Size</th>
                <th>Lines</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                :class="isGreyed(wordlist) ? 'greyed-out-row hover text-gray-500' : 'hover'"
                v-for="wordlist in wordlists"
                :key="wordlist.id"
              >
                <td>
                  <strong>{{ wordlist.name }}</strong>
                  <span class="pl-2 text-sm text-gray-500" v-if="wordlist.pending_delete">
                    <div class="tooltip" data-tip="Marked for death">
                      <font-awesome-icon icon="fa-solid fa-skull-crossbones" title="" />
                    </div>
                  </span>
                </td>

                <td>{{ bytesToReadable(wordlist.size_in_bytes) }}</td>
                <td>{{ wordlist.lines }}</td>
                <td class="text-center">
                  <ConfirmModal @on-confirm="() => onDeleteListfile(wordlist)" v-if="canDelete(wordlist)">
                    <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                  </ConfirmModal>
                  <div v-else class="tooltip cursor-not-allowed text-gray-300" :data-tip="'You can\'t delete this'">
                    <button class="btn btn-ghost btn-xs cursor-not-allowed">
                      <font-awesome-icon icon="fa-solid fa-lock" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <Modal v-model:isOpen="isRulefileUploadOpen">
              <FileUpload fileType="Rulefile" />
            </Modal>
            <h2 class="card-title">Rulefiles</h2>
            <button class="btn btn-primary btn-sm" @click="() => (isRulefileUploadOpen = true)">Upload Rulefile</button>
          </div>

          <table class="table w-full">
            <thead>
              <tr>
                <th>Name</th>
                <th>Size</th>
                <th>Lines</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr
                :class="isGreyed(rulefile) ? 'greyed-out-row hover text-gray-500' : 'hover'"
                v-for="rulefile in rulefiles"
                :key="rulefile.id"
              >
                <td>
                  <strong>{{ rulefile.name }}</strong>
                  <span class="pl-2 text-sm text-gray-500" v-if="rulefile.pending_delete">
                    <div class="tooltip" data-tip="Marked for death">
                      <font-awesome-icon icon="fa-solid fa-skull-crossbones" />
                    </div>
                  </span>
                </td>
                <td>{{ bytesToReadable(rulefile.size_in_bytes) }}</td>
                <td>{{ rulefile.lines }}</td>
                <td class="text-center">
                  <ConfirmModal @on-confirm="() => onDeleteListfile(rulefile)" v-if="canDelete(rulefile)">
                    <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                  </ConfirmModal>
                  <div v-else class="tooltip cursor-not-allowed text-gray-300" :data-tip="'You can\'t delete this'">
                    <button class="btn btn-ghost btn-xs cursor-not-allowed">
                      <font-awesome-icon icon="fa-solid fa-lock" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>

<style scoped>
.greyed-out-row strong {
  font-weight: normal;
}
</style>
