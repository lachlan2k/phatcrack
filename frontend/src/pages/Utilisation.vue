<script setup lang="ts">
import { useApi } from '@/composables/useApi'
import { getJobCountPerUser } from '@/api/project';
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';
import { useUsersStore } from '@/stores/users';
import type { RunningJobCountForUserDTO } from '@/api/types';

const { data, silentlyRefresh } = useApi(getJobCountPerUser)

const userStore = useUsersStore()
userStore.load()

let intervalId = ref<number | null>(null)

onMounted(() => {
    intervalId.value = setInterval(() => silentlyRefresh, 30*1000)
})

onBeforeUnmount(() => {
    if (intervalId.value != null) {
        clearInterval(intervalId.value)
    }
})

const enrichRowWithUsername = (x: RunningJobCountForUserDTO) => 
    ({ username: userStore.byId(x.user_id)?.username ?? 'Unknown User', job_count: x.job_count })

const sortedData = computed(
    () => (data.value?.result.sort((a, b) => b.job_count - a.job_count) ?? []).map(enrichRowWithUsername)
)

</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Utilisation</h1>
    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card min-w-[800px] bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex flex-row justify-between">
            <h2 class="card-title">Job Count per User</h2>
          </div>

          <table class="table w-full">
            <thead>
              <tr>
                <th>User</th>
                <th>Job Count</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in sortedData" :key="row.username">
                <td>{{ row.username }}</td>
                <td>{{ row.job_count }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </main>
</template>
