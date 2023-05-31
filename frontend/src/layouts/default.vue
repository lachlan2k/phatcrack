<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const route = useRoute()

const pageLinks = [
  { name: 'Dashboard', icon: 'fa-gauge', to: '/dashboard' },
  { name: 'Project Folders', icon: 'fa-folder', to: '/projects' },
  { name: 'All Running Jobs', icon: 'fa-bars-progress', to: '/jobs' },
  { name: 'Potfile', icon: 'fa-trophy', to: '/potfile' },
  { name: 'Agents', icon: 'fa-robot', to: '/agents' }
]
</script>

<template>
  <div class="drawer-mobile drawer">
    <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content bg-slate-100">
      <router-view />
    </div>
    <div class="drawer-side bg-neutral text-neutral-content">
      <label for="my-drawer-2" class="drawer-overlay"></label>

      <aside class="flex w-80 flex-col p-4">
        <h2 class="btn-ghost btn w-full text-center text-3xl">Phatcrack</h2>
        <hr class="mt-4 mb-8 h-px border-0 bg-gray-200 dark:bg-gray-700" />
        <RouterLink to="/wizard">
          <a class="btn w-full gap-2 bg-slate-600">
            <span>
              <font-awesome-icon icon="fa-solid fa-pencil" />
            </span>
            Get Cracking
          </a>
        </RouterLink>
        <div class="h-8"></div>
        <ul class="menu text-lg">
          <li class="menu-title pb-2"><span>Pages</span></li>

          <li
            v-for="link in pageLinks"
            :key="link.name"
            :class="route.path == link.to ? 'bordered' : 'hover-bordered'"
          >
            <RouterLink :to="link.to">
              <span><font-awesome-icon :icon="'fa-solid ' + link.icon" /></span>

              {{ link.name }}
            </RouterLink>
          </li>
        </ul>

        <div></div>

        <ul class="menu mt-auto text-lg">
          <li
            class="hover-bordered"
            v-if="authStore.isAdmin"
            :class="route.path == 'admin' ? 'bordered' : 'hover-bordered'"
          >
            <a>
              <span><font-awesome-icon icon="fa-solid fa-lock" /></span>
              Admin
            </a>
          </li>
          <li class="hover-bordered">
            <a>
              <span><font-awesome-icon icon="fa-solid fa-user" /></span>
              <span
                >Welcome, <strong>{{ authStore.loggedInUser?.username }}</strong></span
              >
            </a>
          </li>
        </ul>
      </aside>
    </div>
  </div>
</template>
