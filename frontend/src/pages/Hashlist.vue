<script setup lang="ts">
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
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import decodeHex from '@/util/decodeHex'
import { storeToRefs } from 'pinia'
import { getAttackModeName } from '@/util/hashcat'
import type { AttackWithJobsDTO } from '@/api/types'
import { getAttacksWithJobsForHashlist } from '@/api/project'

const hashlistId = useRoute().params.id as string
const { data: hashlistData, isLoading: isLoadingHashlist } = useApi(() => getHashlist(hashlistId))

const { data: attacksData, isLoading: isLoadingAttacksData } = useApi(() =>
  getAttacksWithJobsForHashlist(hashlistId)
)

const resources = useResourcesStore()

const { getHashTypeName, isHashTypesLoaded } = storeToRefs(resources)
resources.loadHashTypes()

const isLoading = computed(() => {
  return isLoadingHashlist.value || !isHashTypesLoaded.value || isLoadingAttacksData.value
})

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
                  <th>Attack Mode</th>
                  <th>Status of Jobs</th>
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

                  <td>Stop button/restart/start button? delete button?</td>
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
