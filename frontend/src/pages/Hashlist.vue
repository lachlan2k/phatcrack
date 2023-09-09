<script setup lang="ts">
import Modal from '@/components/Modal.vue'
import HashlistEditor from '@/components/HashlistEditor.vue'
import PaginationControls from '@/components/PaginationControls.vue'

import {
  JobStatusAwaitingStart,
  JobStatusCreated,
  JobStatusExited,
  JobStatusStarted,
  JobStopReasonFinished,
  JobStopReasonUserStopped,
  getHashlist
} from '@/api/project'
import { exportResults, ExportFormat } from '@/util/exportHashlist'
import { useApi } from '@/composables/useApi'
import { useResourcesStore } from '@/stores/resources'
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import decodeHex from '@/util/decodeHex'
import { usePagination } from '@/composables/usePagination'
import { storeToRefs } from 'pinia'
import { getAttackModeName, hashrateStr } from '@/util/hashcat'
import { timeDurationToReadable } from '@/util/units'
import type { AttackWithJobsDTO } from '@/api/types'
import { getAttacksWithJobsForHashlist } from '@/api/project'
import JobWizard from '@/components/Wizard/JobWizard.vue'
import AttackDetailsModal from '@/components/AttackDetailsModal/index.vue'

const hashlistId = useRoute().params.id as string
const { data: hashlistData, isLoading: isLoadingHashlist, silentlyRefresh: refreshHashlist } = useApi(() => getHashlist(hashlistId))

const {
  data: attacksData,
  isLoading: isLoadingAttacksData,
  silentlyRefresh: refreshAttack
} = useApi(() => getAttacksWithJobsForHashlist(hashlistId))

const intervalId = ref(0)

async function intervalLoop() {
  await Promise.all([refreshAttack(), refreshHashlist()])
  intervalId.value = setTimeout(intervalLoop, 5 * 1000)
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

const allHashes = computed(() => hashlistData.value?.hashes ?? [])
const crackedHashes = computed(() => allHashes.value.filter((x) => x.is_cracked))

const filteredHashes = computed(() => {
  const arr = onlyShowCracked.value ? crackedHashes.value : allHashes.value

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

const { next: nextPage, prev: prevPage, totalPages, currentItems: currentHashes, activePage } = usePagination(filteredHashes, 20)

const numberOfHashesCracked = computed(() => {
  return crackedHashes.value?.length ?? 0
})

// TODO: this will almost certainl perform terribly, and the code isn't super tidy?
// When maturing this page, it might be worthwhile pulling this out to a weakmap or something computed
const numJobs = (attack: AttackWithJobsDTO) => attack.jobs.length
const numJobsRunning = (attack: AttackWithJobsDTO) => attack.jobs.filter((x) => x.runtime_data.status == JobStatusStarted).length
const numJobsFinished = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter((x) => x.runtime_data.status == JobStatusExited && x.runtime_data.stop_reason == JobStopReasonFinished).length
const numJobsStopped = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter((x) => x.runtime_data.status == JobStatusExited && x.runtime_data.stop_reason == JobStopReasonUserStopped).length
const numJobsFailed = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(
    (x) =>
      x.runtime_data.status == JobStatusExited &&
      x.runtime_data.stop_reason != JobStopReasonFinished &&
      x.runtime_data.stop_reason != JobStopReasonUserStopped
  ).length
const numJobsQueued = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter((x) => x.runtime_data.status == JobStatusAwaitingStart || x.runtime_data.status == JobStatusCreated).length
const hashrateSum = (attack: AttackWithJobsDTO) => attack.jobs.map((x) => x.runtime_summary.hashrate).reduce((a, b) => a + b, 0)

const hashTypeStr = computed(() => {
  if (isLoading.value) {
    return ''
  }
  return getHashTypeName.value(hashlistData.value!.hash_type)
})

const quantityStr = (num: number, str: string) => {
  if (num == 1) {
    return `${num} ${str}`
  }
  return `${num} ${str}s`
}

const isAttackModalOpen = ref(false)
const attackModalAttackIndex = ref(-1)

const selectedAttack = computed(() => {
  return attacksData.value?.attacks[attackModalAttackIndex.value] ?? null
})

function openAttackModal(attackIndex: number) {
  attackModalAttackIndex.value = attackIndex
  isAttackModalOpen.value = true
}
</script>

<template>
  <AttackDetailsModal v-if="selectedAttack != null" :attack="selectedAttack" v-model:isOpen="isAttackModalOpen"></AttackDetailsModal>

  <main class="w-full p-4 h-full">
    <div v-if="isLoading" class="w-full h-full flex justify-center">
      <span class="loading loading-spinner loading-lg"></span>
    </div>
    <div v-else>
      <h1 class="text-4xl font-bold">{{ hashlistData?.name }} {{ hashTypeStr }}</h1>
      <div class="flex flex-wrap gap-4">
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

                <div class="dropdown">
                  <label tabindex="0" class="btn btn-ghost btn-sm m-1">...</label>
                  <ul tabindex="0" class="menu dropdown-content rounded-box z-[1] bg-base-100 p-2 shadow">
                    <li>
                      <button
                        class="btn btn-ghost btn-sm"
                        @click="() => exportResults(hashlistId, ExportFormat.ColonSeparated, onlyShowCracked)"
                      >
                        Export
                      </button>
                    </li>
                    <li>
                      <button class="btn btn-ghost btn-sm" @click="() => (isHashlistEditorOpen = true)">Edit</button>
                    </li>
                  </ul>
                </div>
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
                  <input type="text" class="input input-bordered input-sm" placeholder="Hash or plaintext..." v-model="filterText" />
                </label>
              </div>

              <table class="compact-table compact-table hashlist-table table table-sm w-full">
                <thead>
                  <tr>
                    <th>Original Hash</th>
                    <th>Cracked Plaintext</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="hash in currentHashes" :key="hash.normalized_hash">
                    <td class="font-mono text-ellipsis overflow-hidden whitespace-nowrap" style="max-width: 500px;">{{ hash.input_hash }}</td>
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
                  <JobWizard :firstStep="2" :existingHashlistId="hashlistId" :existingProjectId="hashlistData?.project_id" />
                </Modal>
                <h2 class="card-title">Attacks</h2>
                <button class="btn btn-primary btn-sm" @click="() => (isAttackWizardOpen = true)">New Attack</button>
              </div>
              <table class="compact-table table w-full">
                <thead>
                  <tr>
                    <th>Attack Mode</th>
                    <th>Status</th>
                    <th>Total Hashrate</th>
                    <th>Time Remaining</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    class="cursor-pointer"
                    @click="() => openAttackModal(attackIndex)"
                    v-for="(attack, attackIndex) in attacksData?.attacks"
                    :key="attack.id"
                  >
                    <td>
                      <strong>{{ getAttackModeName(attack.hashcat_params.attack_mode) }}</strong>
                    </td>
                    <td v-if="attack.progress_string != ''">
                      <div class="badge whitespace-nowrap badge-neutral my-1 mr-1">{{ attack.progress_string }}</div>
                    </td>
                    <td v-else-if="numJobs(attack)" style="min-width: 130px">
                      <div class="badge whitespace-nowrap badge-success my-1 mr-1" v-if="numJobsFinished(attack) > 0">
                        {{ quantityStr(numJobsFinished(attack), 'job') }} finished
                      </div>
                      <div class="badge whitespace-nowrap badge-info my-1 mr-1" v-if="numJobsRunning(attack) > 0">
                        {{ quantityStr(numJobsRunning(attack), 'job') }} running
                      </div>
                      <div class="badge whitespace-nowrap badge-secondary my-1 mr-1" v-if="numJobsQueued(attack) > 0">
                        {{ quantityStr(numJobsQueued(attack), 'job') }} pending
                      </div>
                      <div class="badge whitespace-nowrap badge-warning my-1 mr-1" v-if="numJobsStopped(attack)">
                        {{ quantityStr(numJobsStopped(attack), 'job') }} stopped
                      </div>
                      <div class="badge whitespace-nowrap badge-error my-1 mr-1" v-if="numJobsFailed(attack)">
                        {{ quantityStr(numJobsFailed(attack), 'job') }} failed
                      </div>
                    </td>
                    <td style="min-width: 130px" v-else>
                      <div class="badge whitespace-nowrap badge-ghost">No jobs</div>
                    </td>
                    <td>{{ hashrateStr(hashrateSum(attack)) }}</td>
                    <td v-if="attack.jobs.some((x) => x.runtime_summary.estimated_time_remaining > 0)">
                      {{ timeDurationToReadable(Math.max(...attack.jobs.map((x) => x.runtime_summary.estimated_time_remaining), 0)) }}
                      left
                    </td>
                    <td v-else>-</td>
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
