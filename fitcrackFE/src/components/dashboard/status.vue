<template>
  <div class="dashboard-status">
    <div class="py-4 px-8 top-panel d-flex justify-space-between align-center">
      <div>
        <h1>Hello, {{ $store.state.user.userData.username }}!</h1>
        <h2>
          There {{ pluralJobs ? 'are' : 'is' }} currently 
          {{ activeJobs }} active job{{ pluralJobs ? 's' : '' }} 
          with {{ hosts.activeHosts }} host{{ pluralHosts ? 's' : '' }} connected.
        </h2>
      </div>
      <div class="shortcuts">
        <v-btn
          large
          text
          :to="{name: 'hashes'}"
          class="mr-2"
        >
          <v-icon left>
            mdi-book-lock
          </v-icon>
          <span>Hashes</span>
        </v-btn>
        <v-btn
          large
          color="primary"
          elevation="0"
          :to="{name: 'addJob'}"
        >
          <v-icon left>
            mdi-briefcase-plus
          </v-icon>
          <span>New job</span>
        </v-btn>
      </div>
    </div>
    <v-divider />
    <div class="d-flex flex-wrap">
      <div class="status-tile jobs">
        <v-card-title class="py-0">
          Jobs
        </v-card-title>
        <div class="jobs-wrapper">
          <div
            v-for="jg in jobs"
            :key="jg.status"
            class="job-tile"
          >
            <v-icon
              left
              :color="jobColor(jg.status)"
            >
              {{ jobIcon(jg.status) }}
            </v-icon>
            <span><strong>{{ jg.count }}</strong> {{ jg.status }}</span>
          </div>
        </div>
      </div>
      <div class="status-tile hosts">
        <v-card-title class="py-0">
          Hosts
        </v-card-title>
        <div class="host-summary mx-4">
          <span class="mono">{{ hosts.activeHosts }}</span> of <span class="mono">{{ hosts.totalHosts }}</span> online
        </div>
        <host-grid
          class="ma-3"
          :active="hosts.activeHosts"
          :inactive="hosts.inactiveHosts"
        />
      </div>
      <div class="status-tile server">
        <v-card-title class="py-0 d-flex justify-space-between">
          Server
          <v-btn
            x-small
            color="info"
            elevation="0"
            :to="{name: 'server'}"
          >
            <v-icon
              small
              left
            >
              mdi-chart-multiline
            </v-icon>
            <span>Monitor</span>
          </v-btn>
        </v-card-title>
        <div class="server-wrapper">
          <slot />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import hostGrid from '@/components/dashboard/hostGrid.vue'
import { jobIcon, attackIcon } from '@/assets/scripts/iconMaps'

export default {
  components: {
    hostGrid
  },
  props: {
    jobs: {
      type: Array,
      default: () => []
    },
    hosts: {
      type: Object,
      default: () => {}
    },
    server: {
      type: Array,
      default: () => []
    }
  },
  computed: {
    activeJobs () {
      return this.jobs
        .filter(j => j.status === 'running' || j.status === 'finishing')
        .reduce((total, j) => total += j.count, 0)
    },
    pluralJobs () { return this.activeJobs !== 1 },
    pluralHosts () { return this.hosts.activeHosts !== 1 },
  },
  methods: {
    jobIcon,
    jobColor (status) {
      const colors = {
        ready: 'info',
        finished: 'success',
        exhausted: 'error',
        malformed: 'error',
        finishing: 'info'
      }
      return colors[status]
    }
  }
}
</script>

<style scoped>
.dashboard-status {
  width: 100%;
}

.top-panel {
  margin: .95em 0;
}

h1 {
  font-size: 1.7em;
}
h2 {
  font-size: 1.2em;
  font-weight: bold;
}

.status-tile {
  padding: 1em;
  flex: 1 1 25ch;
}
.status-tile:not(:last-child) {
  border-right: 1px solid rgba(0,0,0,.12);
}
.theme--dark .status-tile:not(:last-child) {
  border-right: 1px solid hsla(0,0%,100%,.12);
}

.host-summary {
  font-weight: 500;
}
.host-summary .mono {
  font-family: Roboto Mono;
}

.status-tile.server {
  flex-grow: 3000;
}

/* Jobs tile */
.jobs-wrapper {
  margin: 1em .5em;
  columns: 2;
}

.job-tile {
  position: relative;
  display: flex;
  align-items: center;
  margin: 0 .5em;
  margin-bottom: 1em;
  margin-right: .5em;
  white-space: nowrap;
}

/* Server tile */
.server-wrapper {
  display: flex;
  padding: 0 1em;
}
</style>