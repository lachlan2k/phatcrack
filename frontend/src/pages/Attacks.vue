<script setup lang="ts">
import RecentAttacks from '@/components/RecentAttacks.vue'
</script>

<template>
  <main class="w-full p-6">
    <div class="prose">
      <h1>Attacks</h1>
    </div>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="grow">
        <div class="stats shadow">
          <div class="stat">
            <div class="stat-title">Total Hashrate</div>
            <div class="stat-value flex justify-between">
              <span>3,801MH/s</span>
              <span class="mt-1 ml-4 text-2xl text-primary">
                <font-awesome-icon icon="fa-solid fa-calculator" />
              </span>
            </div>
          </div>

          <div class="stat">
            <div class="stat-title">Total Power Draw</div>
            <div class="stat-value flex justify-between">
              <span>3,801w</span>
              <span class="mt-1 ml-4 text-2xl text-yellow-400">
                <font-awesome-icon icon="fa-solid fa-bolt" />
              </span>
            </div>
          </div>

          <div class="stat">
            <div class="stat-title">Attacks Running</div>

            <div class="stat-value flex justify-between">
              <span>8</span>
              <span class="mt-1 text-2xl text-info">
                <font-awesome-icon icon="fa-solid fa-bars-progress" />
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="flex basis-full"></div>

      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Recent Attacks</h2>
          <RecentAttacks />
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
