<script setup lang="ts">
import { getAllAgents } from '@/api/agent'
import { useApi } from '@/composables/useApi'
import { formatDeviceName } from '@/util/formatDeviceName'

const AgentStatusHealthy = 'AgentStatusHealthy'
const AgentStatusUnhealthyButConnected = 'AgentStatusUnhealthyButConnected'
const AgentStatusUnhealthyAndDisconnected = 'AgentStatusUnhealthyAndDisconnected'
// const AgentStatusDead = 'AgentStatusDead'

const { data: allAgents } = useApi(getAllAgents)
</script>

<template>
  <main class="w-full p-4">
    <div class="prose">
      <h1>Agents</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="grow">
        <div class="stats shadow">
          <div class="stat">
            <div class="stat-title">Agents Online & Healthy</div>
            <div class="stat-value flex justify-between">
              <span
                >{{
                  allAgents?.agents.filter((x) => x.agent_info.status == AgentStatusHealthy)
                    .length ?? '?'
                }}/{{ allAgents?.agents.length ?? '?' }}</span
              >
              <span class="mt-1 text-2xl text-primary">
                <font-awesome-icon icon="fa-solid fa-robot" />
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex basis-full"></div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Agent List</h2>
          <table class="compact-table table w-full">
            <thead>
              <tr>
                <th>Hostname</th>
                <th>Devices</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody class="first-col-bold">
              <tr class="hover" v-for="agent in allAgents?.agents" :key="agent.id">
                <td>{{ agent.name }}</td>
                <td>
                  <span
                    v-for="device in agent.agent_devices"
                    :key="device.device_id + device.device_name"
                  >
                    <font-awesome-icon
                      icon="fa-solid fa-memory"
                      v-if="device.device_type == 'GPU'"
                    />
                    <font-awesome-icon icon="fa-solid fa-microchip" v-else />
                    {{ formatDeviceName(device.device_name) }} ({{ device.temp }} °c)
                    <br />
                  </span>
                </td>

                <td class="text-center">
                  <div
                    class="badge badge-accent badge-sm"
                    v-if="agent.agent_info.status == AgentStatusHealthy"
                    title="Healthy"
                  ></div>
                  <div
                    class="badge badge-warning badge-sm"
                    title="Unhealthy"
                    v-else-if="
                      agent.agent_info.status == AgentStatusUnhealthyAndDisconnected ||
                      agent.agent_info.status == AgentStatusUnhealthyButConnected
                    "
                  ></div>
                  <div class="badge badge-ghost badge-sm" title="Dead" v-else></div>
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
thead > tr > th {
  background: none !important;
  border-bottom-width: 1px;
  /* border-bottom: 1px solid black; */
}

.first-col-bold > tr td:first-of-type {
  font-weight: bold;
}
</style>
