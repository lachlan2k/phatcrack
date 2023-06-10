<script setup lang="ts">
import IconButton from '@/components/IconButton.vue'
import Modal from '@/components/Modal.vue'
import FileUpload from '@/components/FileUpload.vue'

import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import { bytesToReadable } from '@/util/units'
import { useListfilesStore } from '@/stores/listfiles'

const listfilesStore = useListfilesStore()
const { loadListfiles } = listfilesStore
const { wordlists, rulefiles } = storeToRefs(useListfilesStore())

const isWordlistUploadOpen = ref(false)
const isRulefileUploadOpen = ref(false)

loadListfiles()
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Listfiles</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <Modal v-model:isOpen="isWordlistUploadOpen">
              <FileUpload fileType="Wordlist" />
            </Modal>
            <h2 class="card-title">Wordlists</h2>
            <button class="btn-primary btn-sm btn" @click="() => (isWordlistUploadOpen = true)">
              Upload Wordlist
            </button>
          </div>

          <table class="table w-full">
            <!-- head -->
            <thead>
              <tr>
                <th>Name</th>
                <th>Size</th>
                <th>Lines</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <!-- row 1 -->
              <tr class="hover" v-for="wordlist in wordlists" :key="wordlist.id">
                <td>
                  <strong>{{ wordlist.name }}</strong>
                </td>
                <td>{{ bytesToReadable(wordlist.size_in_bytes) }}</td>
                <td>{{ wordlist.lines }}</td>
                <td>
                  <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
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
              <button class="btn-primary btn-sm btn" @click="() => (isRulefileUploadOpen = true)">
                Upload Rulefile
              </button>
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
              <tr class="hover" v-for="rulefile in rulefiles" :key="rulefile.id">
                <td>
                  <strong>{{ rulefile.name }}</strong>
                </td>
                <td>{{ bytesToReadable(rulefile.size_in_bytes) }}</td>
                <td>{{ rulefile.lines }}</td>
                <td>
                  <IconButton icon="fa-solid fa-trash" color="error" tooltip="Delete" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>
