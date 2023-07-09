<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import IconButton from '@/components/IconButton.vue'
import HashlistEditor from '@/components/HashlistEditor.vue'
import PaginationControls from '@/components/PaginationControls.vue'

import {
  JobStatusAwaitingStart,
  JobStatusCreated,
  JobStatusExited,
  JobStatusStarted,
  JobStopReasonFinished,
  getHashlist
} from '@/api/project'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import decodeHex from '@/util/decodeHex'
import { usePagination } from '@/composables/usePagination'
import { storeToRefs } from 'pinia'
import { getAttackModeName, hashrateStr } from '@/util/hashcat'
import type { AttackWithJobsDTO } from '@/api/types'
import { getAttacksWithJobsForHashlist } from '@/api/project'
import JobWizard from '@/components/Wizard/JobWizard.vue'

const hashlistId = useRoute().params.id as string
const {
  data: hashlistData,
  isLoading: isLoadingHashlist,
  silentlyRefresh: refreshHashlist
} = useApi(() => getHashlist(hashlistId))

const {
  data: attacksData,
  isLoading: isLoadingAttacksData,
  silentlyRefresh: refreshAttack
} = useApi(() => getAttacksWithJobsForHashlist(hashlistId))

const intervalId = ref(0)

async function intervalLoop() {
  await Promise.all([refreshAttack(), refreshHashlist])
  intervalId.value = setTimeout(intervalLoop, 3 * 1000)
}

onMounted(() => {
  intervalLoop()
})

onBeforeUnmount(() => {
  clearInterval(intervalId.value)
})

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingHashlist.value || !isHashTypesLoaded.value || isLoadingAttacksData.value
})

const isAttackWizardOpen = ref(false)
const isHashlistEditorOpen = ref(false)

const filterText = ref('')
const onlyShowCracked = ref(false)

const crackedHashes = computed(() => {
  return hashlistData.value?.hashes.filter((x) => x.is_cracked) ?? []
})

const filteredHashes = computed(() => {
  if (onlyShowCracked.value) {
    return crackedHashes.value
  }

  const arr = hashlistData.value?.hashes ?? []

  if (filterText.value == '') {
    return arr
  }

  return arr.filter(
    (x) =>
      decodeHex(x.plaintext_hex).toLowerCase().includes(filterText.value) ||
      x.input_hash.toLowerCase().includes(filterText.value) ||
      x.normalized_hash.toLowerCase().includes(filterText.value)
  )
})

const {
  next: nextPage,
  prev: prevPage,
  totalPages,
  currentItems: currentHashes,
  activePage
} = usePagination(filteredHashes, 20)

const numberOfHashesCracked = computed(() => {
  return crackedHashes.value?.length ?? 0
})

// TODO: this will almost certainl perform terribly, and the code isn't super tidy?
// When maturing this page, it might be worthwhile pulling this out to a weakmap or something computed
const numJobs = (attack: AttackWithJobsDTO) => attack.jobs.length
const numJobsRunning = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter((x) => x.runtime_data.status == JobStatusStarted).length
const numJobsFinished = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(
    (x) =>
      x.runtime_data.status == JobStatusExited &&
      x.runtime_data.stop_reason == JobStopReasonFinished
  ).length
const numJobsFailed = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(
    (x) =>
      x.runtime_data.status == JobStatusExited &&
      x.runtime_data.stop_reason != JobStopReasonFinished
  ).length
const numJobsQueued = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(
    (x) =>
      x.runtime_data.status == JobStatusAwaitingStart || x.runtime_data.status == JobStatusCreated
  ).length
const hashrateSum = (attack: AttackWithJobsDTO) =>
  attack.jobs.map((x) => x.runtime_summary.hashrate).reduce((a, b) => a + b)

const hashTypeStr = computed(() => {
  if (isLoading.value) {
    return ''
  }
  return getHashTypeName.value(hashlistData.value!.hash_type)
})
</script>

<template>
  <main class="w-full p-4">
    <p v-if="isLoading">Loading</p>
    <div v-else>
      <div class="prose">
        <h1>{{ hashlistData?.name }} {{ hashTypeStr }}</h1>
      </div>
      <div class="flex gap-4">
        <div class="mt-6 flex flex-wrap gap-6">
          <div class="card bg-base-100 shadow-xl">
            <Modal v-model:isOpen="isHashlistEditorOpen">
              <HashlistEditor v-if="isHashlistEditorOpen" :hashlistId="hashlistData!.id" />
            </Modal>
            <div class="card-body" style="min-width: 500px">
              <div class="flex flex-row justify-between">
                <h2 class="card-title">
                  Hashlist ({{ numberOfHashesCracked }}/{{ hashlistData?.hashes.length ?? 0 }}
                  cracked)
                </h2>
                <button class="btn-ghost btn-sm btn" @click="() => (isHashlistEditorOpen = true)">
                  Edit
                </button>
              </div>
              <div class="form-control">
                <label class="label cursor-pointer">
                  <span class="label-text">Only show cracked</span>
                  <input type="checkbox" class="toggle" v-model="onlyShowCracked" />
                </label>
              </div>
              <div class="form-control">
                <label class="label">
                  <span class="label-text">Filter</span>
                  <input
                    type="text"
                    class="input-bordered input input-sm"
                    placeholder="Hash or plaintext..."
                    v-model="filterText"
                  />
                </label>
              </div>

              <table class="compact-table compact-table hashlist-table table-sm table w-full">
                <thead>
                  <tr>
                    <th>Original Hash</th>
                    <th>Cracked Plaintext</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="hash in currentHashes" :key="hash.normalized_hash">
                    <td class="font-mono">{{ hash.input_hash }}</td>
                    <td class="font-mono">
                      <strong>{{ decodeHex(hash.plaintext_hex) || '-' }}</strong>
                    </td>
                  </tr>
                </tbody>
              </table>

              <div class="mt-2 w-full text-center">
                <PaginationControls
                  @next="() => nextPage()"
                  @prev="() => prevPage()"
                  :current-page="activePage"
                  :total-pages="totalPages"
                />
              </div>
            </div>
          </div>
        </div>
        <div class="mt-6 flex flex-wrap gap-6">
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <div class="flex flex-row justify-between">
                <Modal v-model:isOpen="isAttackWizardOpen">
                  <JobWizard
                    :firstStep="2"
                    :existingHashlistId="hashlistId"
                    :existingProjectId="hashlistData?.project_id"
                  />
                </Modal>
                <h2 class="card-title">Attacks</h2>
                <button class="btn-primary btn-sm btn" @click="() => (isAttackWizardOpen = true)">
                  New Attack
                </button>
              </div>
              <table class="compact-table table w-full">
                <!-- head -->
                <thead>
                  <tr>
                    <th>Attack Mode</th>
                    <th>Status</th>
                    <th>Total Hashrate</th>
                    <th>Controls</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="attack in attacksData?.attacks" :key="attack.id">
                    <td>
                      {{ getAttackModeName(attack.hashcat_params.attack_mode) }}
                    </td>
                    <td v-if="numJobs(attack)">
                      <div class="badge badge-success mr-1" v-if="numJobsFinished(attack) > 0">
                        {{ numJobsFinished(attack) }} jobs finished
                      </div>
                      <div class="badge badge-info mr-1" v-if="numJobsRunning(attack) > 0">
                        {{ numJobsRunning(attack) }} jobs running
                      </div>
                      <div class="badge badge-secondary mr-1" v-if="numJobsQueued(attack) > 0">
                        {{ numJobsQueued(attack) }} jobs pending
                      </div>
                      <div class="badge badge-error" v-if="numJobsFailed(attack)">
                        {{ numJobsFailed(attack) }} jobs failed
                      </div>
                    </td>
                    <td v-else>
                      <div class="badge badge-info">No jobs for attack</div>
                    </td>
                    <td>{{ hashrateStr(hashrateSum(attack)) }}</td>
                    <td class="text-center">
                      <IconButton tooltip="Start" icon="fa-solid fa-play" color="success" />
                      <IconButton tooltip="Stop" icon="fa-solid fa-stop" color="warning" />
                      <IconButton tooltip="Delete" icon="fa-solid fa-trash" color="error" />
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
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

.hashlist-table.table-sm :where(th, td) {
  padding-top: 0.4rem;
  padding-bottom: 0.4rem;
}
</style>
