<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { ref } from 'vue'

import Modal from '@/components/Modal.vue'
import EmptyTable from '@/components/EmptyTable.vue'

import { useAttackTemplatesStore } from '@/stores/attackTemplates'

import { Icons } from '@/util/icons'

const attackTemplatesStore = useAttackTemplatesStore()

const { templates } = storeToRefs(attackTemplatesStore)

const isCreateModalOpen = ref(false)
const isEditModalOpen = ref(false)
</script>

<template>
  <Modal v-model:isOpen="isCreateModalOpen">
    <h1>foo</h1>
  </Modal>

  <Modal v-model:isOpen="isEditModalOpen">
    <h1>bar</h1>
  </Modal>

  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Attack Templates</h1>
    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card min-w-[800px] bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <h2 class="card-title">Attack Templates</h2>

            <div>
              <button class="btn btn-primary btn-sm ml-1" @click="() => (isCreateModalOpen = true)">New Template</button>
              <button class="btn btn-primary btn-sm ml-1" @click="() => (isCreateModalOpen = true)">New Template Set</button>
            </div>
          </div>
          <table class="table w-full">
            <thead>
              <tr></tr>
            </thead>
            <tbody>
              <tr></tr>
            </tbody>
          </table>
          <EmptyTable v-if="templates.length == 0" text="No Attack Templates Yet" :icon="Icons.AttackTemplate" />
        </div>
      </div>
    </div>
  </main>
</template>
