<script setup lang="ts">
import RecentAttacks from '@/components/RecentAttacks.vue'
</script>

<template>
  <main class="w-full p-6">
    <div class="prose">
      <h1>Dashboard</h1>
    </div>
    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body">
          <h2 class="card-title">Recent Projects</h2>
          <table class="table w-full">
            <!-- head -->
            <thead>
              <tr>
                <th>Name</th>
                <th>Created</th>
                <th>Hashlists</th>
                <th>Attacks</th>
                <!-- <th>Actions</th> -->
              </tr>
            </thead>
            <tbody class="first-col-bold">
              <!-- row 1 -->
              <tr class="hover">
                <td>41234 - Customer A</td>
                <td>2 hours ago</td>
                <td>8</td>
                <td><div class="badge-info badge">4 attacks running</div></td>
              </tr>
              <!-- row 2 -->
              <tr class="hover">
                <td>41235 - Customer B</td>
                <td>10 days ago</td>
                <td>3</td>
                <td><div class="badge-ghost badge">Idle</div></td>
              </tr>
              <!-- row 3 -->
              <tr class="hover">
                <td>41236 - Customer C</td>
                <td>2 months ago</td>
                <td>2</td>
                <td><div class="badge-ghost badge">Idle</div></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

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
}

.first-col-bold > tr td:first-of-type {
  font-weight: bold;
}
</style>
