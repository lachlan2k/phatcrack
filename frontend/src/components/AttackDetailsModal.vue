<script setup lang="ts">
import {
  JobStatusAwaitingStart,
  JobStatusCreated,
  JobStatusExited,
  JobStatusStarted,
  JobStopReasonFinished
} from '@/api/project'

import type { AttackWithJobsDTO } from '@/api/types'
import Modal from '@/components/Modal.vue'
import { useAgentsStore } from '@/stores/agents'
import { getAttackModeName, hashrateStr } from '@/util/hashcat'
import { timeBetween, timeDurationToReadable } from '@/util/units'
import { timeSince } from '@/util/units'

import { computed } from 'vue'
import AttackConfigDetails from './AttackConfigDetails.vue'

const props = defineProps<{
  isOpen: boolean
  attack: AttackWithJobsDTO
}>()

const emit = defineEmits(['update:isOpen'])

const isOpen = computed({
  get: () => props.isOpen,
  set: (value: boolean) => emit('update:isOpen', value)
})

const agentStore = useAgentsStore()
agentStore.load()
const getAgentName = (id: string) => agentStore.byId(id)?.name ?? 'Unknown'
</script>

<template>
  <Modal v-model:isOpen="isOpen">
    <h2 class="mb-8 text-center text-xl font-bold">
      {{ getAttackModeName(props.attack.hashcat_params.attack_mode) }} Attack
    </h2>
    <AttackConfigDetails :hashcatParams="attack.hashcat_params"></AttackConfigDetails>
    <div class="my-8"></div>

    <table class="compact-table table w-full" v-if="attack.jobs.length > 0">
      <!-- head -->
      <thead>
        <tr>
          <th>Running Agent</th>
          <th>Status</th>
          <th>Total Hashrate</th>
          <th>Time Started</th>
          <th>Time Remaining</th>
          <th>Time Taken</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="job in attack.jobs" :key="job.id">
          <td>
            <strong>{{ getAgentName(job.assigned_agent_id) }}</strong>
          </td>
          <td>
            <div
              class="badge badge-success mr-1"
              v-if="
                job.runtime_data.status == JobStatusExited &&
                job.runtime_data.stop_reason == JobStopReasonFinished
              "
            >
              Job finished
            </div>
            <div
              class="badge badge-info mr-1"
              v-else-if="job.runtime_data.status == JobStatusStarted"
            >
              Job running
            </div>
            <div
              class="badge badge-secondary mr-1"
              v-else-if="
                job.runtime_data.status == JobStatusAwaitingStart ||
                job.runtime_data.status == JobStatusCreated
              "
            >
              Job pending
            </div>
            <div class="badge badge-error" v-else-if="job.runtime_data.status == JobStatusExited">
              Job failed
            </div>
            <div class="badget badge-ghost" v-else>Unknown state</div>
          </td>
          <td>{{ hashrateStr(job.runtime_summary.hashrate) }}</td>
          <td>{{ timeSince(job.runtime_summary.started_time * 1000) }}</td>

          <td v-if="job.runtime_summary.estimated_time_remaining > 0">
            {{ timeDurationToReadable(job.runtime_summary.estimated_time_remaining) }}
            left
          </td>
          <td v-else>-</td>

          <td v-if="job.runtime_summary.stopped_time > 0">
            {{ timeBetween(job.runtime_summary.started_time * 1000, job.runtime_summary.stopped_time * 1000) }}
          </td>
          <td v-else>-</td>
        </tr>
      </tbody>
    </table>

    <p v-else>No jobs were found for this attack</p>
  </Modal>
</template>
