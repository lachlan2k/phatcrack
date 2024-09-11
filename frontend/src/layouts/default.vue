<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useConfigStore } from '@/stores/config'
import { onMounted, ref } from 'vue'
import { adminGetVersion } from '@/api/admin'

const authStore = useAuthStore()
const { loggedInUser, isAdmin } = storeToRefs(authStore)

const configStore = useConfigStore()
const { config } = storeToRefs(configStore)
configStore.load()

const router = useRouter()
const route = useRoute()

const pageLinks = [
  { name: 'Project Dashboard', icon: 'fa-folder', to: '/dashboard' },
  { name: 'Listfiles', icon: 'fa-file', to: '/listfiles' },
  { name: 'Hash Search', icon: 'fa-magnifying-glass', to: '/hash-search' },
  { name: 'Utilisation', icon: 'fa-gauge', to: '/utilisation' },
  { name: 'Agents', icon: 'fa-robot', to: '/agents' }
]

const adminPageLinks = [
  { name: 'Configuration', icon: 'fa-gear', to: '/admin/general' },
  { name: 'Manage Users', icon: 'fa-users', to: '/admin/users' },
  { name: 'Manage Agents', icon: 'fa-robot', to: '/admin/agents' }
]

const version = ref('')

onMounted(async () => {
  version.value = await adminGetVersion()
})

async function logout() {
  await authStore.logout()
  router.push('/login')
  authStore.refreshAuth()
}
</script>

<template>
  <div class="drawer lg:drawer-open">
    <input id="my-drawer-2" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content bg-slate-100">
      <router-view />
    </div>
    <div class="drawer-side bg-neutral text-neutral-content">
      <aside class="flex h-full w-72 flex-col p-4">
        <RouterLink to="/dashboard">
          <h2 class="btn btn-ghost w-full text-center text-3xl">Phatcrack</h2>
          <div class="w-full text-center" v-if="version != ''">
            <small class="text-center font-mono text-xs">{{ version }}</small>
          </div>
        </RouterLink>
        <div
          v-if="config?.general?.is_maintenance_mode"
          class="tooltip-white-bg tooltip tooltip-bottom tooltip-warning w-full"
          data-tip="Your administrator has put Phatcrack in maintenance. You can't start new attacks."
        >
          <button class="btn btn-warning btn-sm mb-1 mt-4 w-full py-1">
            <font-awesome-icon icon="fa-solid fa-warning" />
            Maintenance Mode
          </button>
        </div>
        <hr class="mb-8 mt-4 h-px border-0 bg-gray-700" />
        <RouterLink to="/wizard">
          <a class="btn btn-neutral w-full gap-2 bg-slate-600">
            <span>
              <font-awesome-icon icon="fa-solid fa-pencil" />
            </span>
            Get Cracking
          </a>
        </RouterLink>

        <ul class="menu mt-4">
          <li v-for="link in pageLinks" :key="link.name" :class="route.path == link.to ? 'bordered' : 'hover-bordered'" class="mt-2">
            <RouterLink :to="link.to" :class="route.path == link.to ? 'active' : ''">
              <span class="w-6 text-center"><font-awesome-icon :icon="'fa-solid ' + link.icon" /></span>

              {{ link.name }}
            </RouterLink>
          </li>
        </ul>

        <div class="flex flex-grow"></div>

        <ul class="menu justify-self-end">
          <li class="hover-bordered" :class="route.path.startsWith('/admin') ? 'bordered' : 'hover-bordered'" v-if="isAdmin">
            <div class="text-content-neutral dropdown dropdown-top">
              <label tabindex="0" class="col-span-2 w-full cursor-pointer">
                <span class="w-6 text-center"><font-awesome-icon icon="fa-solid fa-lock" /></span>
                <span>Admin Tools</span>
              </label>

              <ul tabindex="0" class="menu dropdown-content rounded-box mb-2 w-52 bg-base-100 p-2 pb-4 text-black shadow">
                <li v-for="link in adminPageLinks" :key="link.name">
                  <RouterLink :to="link.to">
                    <span class="min-w-[20px] text-center"><font-awesome-icon :icon="'fa-solid ' + link.icon" /></span>
                    <span>{{ link.name }}</span>
                  </RouterLink>
                </li>
              </ul>
            </div>
          </li>

          <li class="hover-bordered" :class="route.path == '/account' ? 'bordered' : 'hover-bordered'">
            <div class="text-content-neutral dropdown dropdown-top">
              <label tabindex="0" class="col-span-2 w-full cursor-pointer">
                <span class="w-6 text-center"><font-awesome-icon icon="fa-solid fa-user" /></span>
                <span
                  >Welcome, <strong>{{ loggedInUser?.username }}</strong></span
                >
              </label>

              <ul tabindex="0" class="menu dropdown-content rounded-box w-52 bg-base-100 p-2 pb-4 text-black shadow">
                <li>
                  <RouterLink to="/account">
                    <span><font-awesome-icon icon="fa-solid fa-user" /></span>
                    <span>My account</span>
                  </RouterLink>
                </li>
                <li>
                  <a @click="logout()">
                    <span><font-awesome-icon icon="fa-solid fa-sign-out" /></span>
                    <span>Sign out</span>
                  </a>
                </li>
              </ul>
            </div>
          </li>
        </ul>
      </aside>
    </div>
  </div>
</template>

<style scoped>
/* Backported from Daisy UI v2 */
.menu li {
  margin-top: 0.65rem;
  transition: border 0.125s ease;
}

.menu li a {
  padding: 0.6rem 1.25rem;
}

.menu {
  font-size: 1rem;
}

.menu li.hover-bordered {
  @apply border-l-4 border-transparent hover:border-primary;
}

.menu li.hover-bordered:hover {
  background: rgba(255, 255, 255, 0.1);
}

.menu li.bordered {
  @apply border-l-4 border-primary;
}

.dropdown,
.dropdown label,
.dropdown label:hover,
.dropdown label:active {
  color: inherit !important;
}

.menu .dropdown {
  padding: 0.6rem 1.25rem;
}

/* Same as a menu item */
.menu .dropdown label {
  display: grid;
  grid-auto-flow: column;
  align-content: flex-start;
  align-items: center;
  gap: 0.5rem;
  grid-auto-columns: max-content auto max-content;
}

.menu li:hover a,
.menu li a:active {
  color: inherit;
}

.tooltip-white-bg::before,
.tooltip-white-bg::after {
  --tooltip-color: white;
}
</style>
