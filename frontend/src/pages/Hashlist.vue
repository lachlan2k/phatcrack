<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useToast } from 'vue-toastification'

import Modal from '@/components/Modal.vue'
import PaginationControls from '@/components/PaginationControls.vue'
import JobWizard from '@/components/Wizard/JobWizard.vue'
import AttackDetailsModal from '@/components/AttackDetailsModal/index.vue'
import HashesInput from '@/components/HashesInput.vue'
import PageLoading from '@/components/PageLoading.vue'

import {
  JobStatusAwaitingStart,
  JobStatusCreated,
  JobStatusExited,
  JobStatusStarted,
  JobStopReasonFinished,
  JobStopReasonUserStopped,
  appendToHashlist,
  getHashlist,
  getAttacksWithJobsForHashlist
} from '@/api/project'
import type { AttackWithJobsDTO } from '@/api/types'

import { useApi } from '@/composables/useApi'
import { usePagination } from '@/composables/usePagination'
import { useHashesInput } from '@/composables/useHashesInput'
import { useToastError } from '@/composables/useToastError'

import { useResourcesStore } from '@/stores/resources'
import { useProjectsStore } from '@/stores/projects'

import { exportResults, ExportFormat } from '@/util/exportHashlist'
import decodeHex from '@/util/decodeHex'
import { getAttackModeName, hashrateStr } from '@/util/hashcat'
import { timeDurationToReadable } from '@/util/units'
import { Icons } from '@/util/icons'

const hashlistId = useRoute().params.id as string
const { data: hashlistData, isLoading: isLoadingHashlist, silentlyRefresh: refreshHashlist } = useApi(() => getHashlist(hashlistId))

const {
  data: attacksData,
  isLoading: isLoadingAttacksData,
  silentlyRefresh: refreshAttacks
} = useApi(() => getAttacksWithJobsForHashlist(hashlistId))

const intervalId = ref<number | null>(null)

async function intervalLoop() {
  await Promise.all([refreshAttacks(), refreshHashlist()])
  intervalId.value = setTimeout(intervalLoop, 5 * 1000)
}

onMounted(() => {
  intervalLoop()
})

onBeforeUnmount(() => {
  if (intervalId.value != null) clearInterval(intervalId.value)
})

const resourcesStore = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resourcesStore)
resourcesStore.loadHashTypes()

const projectStore = useProjectsStore()
projectStore.load(true)

const project = computed(() => {
  const id = hashlistData.value?.project_id
  if (id == null) {
    return null
  }
  return projectStore.byId(id)
})

const projectName = computed(() => {
  return project.value?.name ?? ''
})

const projectUrl = computed(() => {
  const proj = project.value
  if (proj == null) {
    return '#'
  }

  return `/project/${proj.id}`
})

const isLoading = computed(() => {
  return isLoadingHashlist.value || !isHashTypesLoaded.value || isLoadingAttacksData.value
})

const isAttackWizardOpen = ref(false)

const filterText = ref('')
const onlyShowCracked = ref(false)

const allHashes = computed(() => hashlistData.value?.hashes ?? [])
const crackedHashes = computed(() => allHashes.value.filter(x => x.is_cracked))

const filteredHashes = computed(() => {
  const arr = onlyShowCracked.value ? crackedHashes.value : allHashes.value

  if (filterText.value == '') {
    return arr
  }

  const searchTerm = filterText.value.trim().toLowerCase()

  return arr.filter(
    x =>
      x.username.toLowerCase().includes(searchTerm) ||
      decodeHex(x.plaintext_hex).toLowerCase().includes(searchTerm) ||
      x.input_hash.toLowerCase().includes(searchTerm) ||
      x.normalized_hash.toLowerCase().includes(searchTerm)
  )
})

const { next: nextPage, prev: prevPage, totalPages, currentItems: currentHashes, activePage } = usePagination(filteredHashes, 20)

const numberOfHashesCracked = computed(() => {
  return crackedHashes.value?.length ?? 0
})

// TODO: this will almost certainl perform terribly, and the code isn't super tidy?
// When maturing this page, it might be worthwhile pulling this out to a weakmap or something computed
const numJobs = (attack: AttackWithJobsDTO) => attack.jobs.length
const numJobsRunning = (attack: AttackWithJobsDTO) => attack.jobs.filter(x => x.runtime_data.status == JobStatusStarted).length
const numJobsFinished = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(x => x.runtime_data.status == JobStatusExited && x.runtime_data.stop_reason == JobStopReasonFinished).length
const numJobsStopped = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(x => x.runtime_data.status == JobStatusExited && x.runtime_data.stop_reason == JobStopReasonUserStopped).length
const numJobsFailed = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(
    x =>
      x.runtime_data.status == JobStatusExited &&
      x.runtime_data.stop_reason != JobStopReasonFinished &&
      x.runtime_data.stop_reason != JobStopReasonUserStopped
  ).length
const numJobsQueued = (attack: AttackWithJobsDTO) =>
  attack.jobs.filter(x => x.runtime_data.status == JobStatusAwaitingStart || x.runtime_data.status == JobStatusCreated).length

const hashrateSum = (attack: AttackWithJobsDTO) =>
  attack.jobs
    .filter(x => x.runtime_data.status == JobStatusStarted)
    .map(x => x.runtime_summary.hashrate)
    .reduce((a, b) => a + b, 0)

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

const isHashAddModalOpen = ref(false)
const isAppendHashesLoading = ref(false)
const { hashesInput: appendHashesInput, hashesArr: appendHashesArr } = useHashesInput()

const toast = useToast()
const { catcher } = useToastError()

async function onAppendHashes() {
  isAppendHashesLoading.value = true

  try {
    const res = await appendToHashlist(hashlistId, appendHashesArr.value)
    toast.success(
      `Added ${res.num_new_hashes} new hashes.` +
        (res.num_populated_from_potfile > 0 ? `${res.num_populated_from_potfile} already cracked.` : '')
    )
    appendHashesInput.value = ''
    isHashAddModalOpen.value = false
    refreshHashlist()
  } catch (e) {
    catcher(e, 'Failed to append hashes: ')
  } finally {
    isAppendHashesLoading.value = false
  }
}
</script>

<template>
  <AttackDetailsModal v-if="selectedAttack != null" :attack="selectedAttack" v-model:isOpen="isAttackModalOpen"></AttackDetailsModal>

  <Modal v-model:isOpen="isHashAddModalOpen">
    <div class="w-screen max-w-[600px]">
      <h3 class="text-lg font-bold">Add new hashes</h3>
      <small class="text-sm">Note: This will not affect your current attacks. You must start new attacks to attack these hashes.</small>

      <HashesInput class="mt-4" v-model="appendHashesInput" />

      <div class="mt-4 flex justify-end">
        <button class="btn btn-primary" :disabled="isAppendHashesLoading" @click="() => onAppendHashes()">Append</button>
      </div>
    </div>
  </Modal>

  <main class="h-full w-full p-4">
    <PageLoading v-if="isLoading" />
    <div v-else>
      <h1 class="text-4xl font-bold">{{ hashlistData?.name }} - {{ hashTypeStr }}</h1>
      <div class="breadcrumbs pl-1 text-sm">
        <ul>
          <li>
            <RouterLink to="/dashboard"> Dashboard </RouterLink>
          </li>
          <li>
            <RouterLink :to="projectUrl">Project {{ projectName }}</RouterLink>
          </li>
          <li>This hashlist</li>
        </ul>
      </div>
      <div class="flex flex-wrap gap-4">
        <div class="mt-3 flex flex-wrap gap-6">
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body" style="min-width: 500px">
              <div class="flex flex-row justify-between">
                <h2 class="card-title">
                  Hashlist ({{ numberOfHashesCracked }}/{{ hashlistData?.hashes.length ?? 0 }}
                  cracked)
                </h2>

                <div>
                  <div class="tooltip" data-tip="Add more hashes">
                    <button class="btn btn-ghost btn-sm" @click="() => (isHashAddModalOpen = true)">
                      <font-awesome-icon :icon="Icons.Add" />
                    </button>
                  </div>

                  <div class="tooltip" data-tip="Export as colon-separated file">
                    <button
                      class="btn btn-ghost btn-sm"
                      @click="() => exportResults(hashlistId, ExportFormat.ColonSeparated, onlyShowCracked)"
                    >
                      <font-awesome-icon :icon="Icons.Download" />
                    </button>
                  </div>
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

              <table class="compact-table hashlist-table table table-sm w-full">
                <thead>
                  <tr>
                    <th v-if="hashlistData?.has_usernames">Username</th>
                    <th>Original Hash</th>
                    <th>Cracked Plaintext</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="hash in currentHashes" :key="hash.id + '|' + hash.normalized_hash">
                    <td class="font-mono" v-if="hashlistData?.has_usernames">
                      {{ hash.username }}
                    </td>
                    <td>
                      <div
                        class="tooltip tooltip-right mr-[5px] h-[15px] w-[15px]"
                        data-tip="This hash was unexpected and automatically appended (issue #22)"
                        v-if="hash.is_unexpected"
                      >
                        <font-awesome-icon :icon="Icons.Warning" class="align-middle" />
                      </div>
                      <span
                        class="inline-block overflow-hidden text-ellipsis whitespace-nowrap align-middle font-mono"
                        :class="hash.is_unexpected ? 'max-w-[480px]' : 'max-w-[500px]'"
                        >{{ hash.input_hash }}</span
                      >
                    </td>
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
        <div class="mt-3 flex flex-wrap gap-6">
          <div class="card bg-base-100 shadow-xl">
            <div class="card-body">
              <div class="flex flex-row justify-between">
                <Modal v-model:isOpen="isAttackWizardOpen">
                  <JobWizard
                    :firstStep="2"
                    :existingHashlistId="hashlistId"
                    :existingProjectId="hashlistData?.project_id"
                    @created-attack="refreshAttacks()"
                  />
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
                      <div class="badge badge-neutral mr-1 whitespace-nowrap">{{ attack.progress_string }}</div>
                    </td>
                    <td v-else-if="numJobs(attack)" style="min-width: 130px">
                      <div class="badge badge-success mr-1 whitespace-nowrap" v-if="numJobsFinished(attack) > 0">
                        {{ quantityStr(numJobsFinished(attack), 'job') }} finished
                      </div>
                      <div class="badge badge-info mr-1 whitespace-nowrap" v-if="numJobsRunning(attack) > 0">
                        {{ quantityStr(numJobsRunning(attack), 'job') }} running
                      </div>
                      <div class="badge badge-secondary mr-1 whitespace-nowrap" v-if="numJobsQueued(attack) > 0">
                        {{ quantityStr(numJobsQueued(attack), 'job') }} pending
                      </div>
                      <div class="badge badge-warning mr-1 whitespace-nowrap" v-if="numJobsStopped(attack)">
                        {{ quantityStr(numJobsStopped(attack), 'job') }} stopped
                      </div>
                      <div class="badge badge-error mr-1 whitespace-nowrap" v-if="numJobsFailed(attack)">
                        {{ quantityStr(numJobsFailed(attack), 'job') }} failed
                      </div>
                    </td>
                    <td style="min-width: 130px" v-else>
                      <div class="badge badge-ghost whitespace-nowrap">No jobs</div>
                    </td>
                    <td>{{ hashrateStr(hashrateSum(attack)) }}</td>
                    <td v-if="attack.jobs.some(x => x.runtime_summary.estimated_time_remaining > 0)">
                      {{ timeDurationToReadable(Math.max(...attack.jobs.map(x => x.runtime_summary.estimated_time_remaining), 0)) }}
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
