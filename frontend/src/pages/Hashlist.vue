<script setup lang="ts">
import { getAttacksForHashlist, getHashlist } from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import decodeHex from '@/util/decodeHex'
import { storeToRefs } from 'pinia'

const hashlistId = useRoute().params.id as string
const { data: hashlistData, isLoading: isLoadingHashlist } = useApi(() => getHashlist(hashlistId))

const { data: attacksData, isLoading: isLoadingAttacksData } = useApi(() => getAttacksForHashlist(hashlistId))

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingHashlist.value || !isHashTypesLoaded.value || isLoadingAttacksData.value
})

const hashTypeStr = computed(() => {
  if (isLoading.value) {
    return ''
  }
  return getHashTypeName.value(hashlistData.value!.hash_type)
})
</script>

<template>
  <main class="w-full p-6">
    <p v-if="isLoading">Loading</p>
    <div v-else>
      <div class="prose">
        <h1>{{ hashlistData?.name }} {{ hashTypeStr }}</h1>
      </div>
      <div class="mt-6 flex flex-col flex-wrap gap-6">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Hashes</h2>
            <table class="table w-full">
              <!-- head -->
              <thead>
                <tr>
                  <th>Original Hash</th>
                  <th>Normalized Hash</th>
                  <th>Cracked Plaintext</th>
                </tr>
              </thead>

              <tbody>
                <tr v-for="hash in hashlistData?.hashes" :key="hash.normalized_hash">
                  <td>{{ hash.input_hash }}</td>
                  <td>{{ hash.normalized_hash }}</td>
                  <td>{{ decodeHex(hash.plaintext_hex) || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div class="mt-6 flex flex-col flex-wrap gap-6">
        <div class="card bg-base-100 shadow-xl">
          <div class="card-body">
            <h2 class="card-title">Attacks</h2>
            <table class="table w-full">
              <!-- head -->
              <thead>
                <tr>
                  <th>Original Hash</th>
                  <th>Normalized Hash</th>
                  <th>Cracked Plaintext</th>
                </tr>
              </thead>

              <tbody>
                <tr v-for="attack in attacksData?.attacks" :key="attack.id">
                  <td>
                    {{ attack.id }}
                  </td>
                </tr>
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
</style>
