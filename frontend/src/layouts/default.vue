<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { logout as apiLogout } from '@/api/auth'

const authStore = useAuthStore()

const router = useRouter()
const route = useRoute()

const pageLinks = [
  { name: 'Project Dashboard', icon: 'fa-gauge', to: '/dashboard' },
  { name: 'Wordlists & Rules', icon: 'fa-file', to: '/listfiles' },
  // { name: 'Potfile', icon: 'fa-trophy', to: '/potfile' },
  { name: 'Agents', icon: 'fa-robot', to: '/agents' }
]

async function logout() {
  await apiLogout()
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
        <h2 class="btn btn-ghost w-full text-center text-3xl">Phatcrack</h2>
        <hr class="mb-8 mt-4 h-px border-0 bg-gray-200 dark:bg-gray-700" />
        <RouterLink to="/wizard">
          <a class="btn btn-neutral w-full gap-2 bg-slate-600">
            <span>
              <font-awesome-icon icon="fa-solid fa-pencil" />
            </span>
            Get Cracking
          </a>
        </RouterLink>

        <ul class="menu mt-4">
          <li
            v-for="link in pageLinks"
            :key="link.name"
            :class="route.path == link.to ? 'bordered' : 'hover-bordered'"
            class="mt-2"
          >
            <RouterLink :to="link.to" :class="route.path == link.to ? 'active' : ''">
              <span class="w-6 text-center"
                ><font-awesome-icon :icon="'fa-solid ' + link.icon"
              /></span>

              {{ link.name }}
            </RouterLink>
          </li>
        </ul>

        <div class="flex flex-grow"></div>

        <ul class="menu justify-self-end">
          <li
            class="hover-bordered"
            v-if="authStore.isAdmin"
            :class="route.path == 'admin' ? 'bordered' : 'hover-bordered'"
          >
            <RouterLink to="/admin" :class="route.path == '/admin' ? 'active' : ''">
              <span class="w-6 text-center"><font-awesome-icon icon="fa-solid fa-lock" /></span>
              Admin
            </RouterLink>
          </li>
          <li class="hover-bordered">
            <div class="text-content-neutral dropdown dropdown-top">
              <label tabindex="0" class="w-full cursor-pointer">
                <span class="w-6 text-center"><font-awesome-icon icon="fa-solid fa-user" /></span>
                <span
                  >Welcome, <strong>{{ authStore.loggedInUser?.username }}</strong></span
                >
              </label>

              <ul
                tabindex="0"
                class="dropdown-content menu rounded-box w-52 bg-base-100 p-2 pb-4 text-black shadow"
              >
                <li>
                  <a>
                    <span><font-awesome-icon icon="fa-solid fa-user" /></span>
                    <span>My account</span>
                  </a>
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
</style>
