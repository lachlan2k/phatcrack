<template>
  <div class="server-gauges-container">
    <gauge
      :value="cpu"
      :value-format="x => `${x || '--'}%`"
      :color="gaugeColor(cpu)"
      class="mr-8"
    >
      <div class="d-flex flex-column g-label-wrapper">
        <v-icon
          class="mb-1"
        >
          mdi-memory
        </v-icon>
        <span class="overline g-label-text">CPU</span>
      </div>
    </gauge>
    <gauge
      :value="ram"
      :value-format="x => `${x || '--'}%`"
      :color="gaugeColor(ram)"
      class="mr-10"
    >
      <div class="d-flex flex-column g-label-wrapper">
        <v-icon
          class="mb-1"
        >
          mdi-chip
        </v-icon>
        <span class="overline g-label-text">Memory</span>
      </div>
    </gauge>
    <div class="text-module disk-usage my-2 mr-10">
      <div class="stat-label mb-1">
        <v-icon>mdi-harddisk</v-icon>
        <span class="overline g-label-text">Disk</span>
      </div>
      <div class="stat">
        <v-icon
          left
          small
        >
          mdi-arrow-top-right
        </v-icon>
        <span class="overline">read</span>
        <div class="stat-value">
          {{ hdd_read || '--' }}<span class="stat-unit"> kb/s</span>
        </div>
      </div>
      <div class="stat">
        <v-icon
          left
          small
        >
          mdi-arrow-bottom-left
        </v-icon>
        <span class="overline">write</span>
        <div class="stat-value">
          {{ hdd_write || '--' }}<span class="stat-unit"> kb/s</span>
        </div>
      </div>
    </div>
    <div class="text-module net-usage my-2 mr-10">
      <div class="stat-label mb-1">
        <v-icon>mdi-ethernet</v-icon>
        <span class="overline g-label-text">Network</span>
      </div>
      <div class="stat">
        <v-icon
          left
          small
        >
          mdi-arrow-bottom-left
        </v-icon>
        <span class="overline">download</span>
        <div class="stat-value">
          {{ net_recv || '--' }}<span class="stat-unit"> kb/s</span>
        </div>
      </div>
      <div class="stat">
        <v-icon
          left
          small
        >
          mdi-arrow-top-right
        </v-icon>
        <span class="overline">upload</span>
        <div class="stat-value">
          {{ net_sent || '--' }}<span class="stat-unit"> kb/s</span>
        </div>
      </div>
    </div>
    <div class="text-module services my-2">
      <div class="stat-label mb-1">
        <v-icon>mdi-toolbox-outline</v-icon>
        <span class="overline g-label-text">Services</span>
      </div>
      <v-sheet
        rounded
        v-if="failedServices.length > 0"
        color="error"
        dark
        class="highlighted-stat px-2 pb-1 mb-1"
      >
        <v-icon
          left
          small
        >
          mdi-heart-broken
        </v-icon>
        <span class="overline">Down</span>
        <ul class="font-weight-bold">
          <li
            v-for="s in failedServices"
            :key="s"
          >
            {{ s }}
          </li>
        </ul>
      </v-sheet>
      <div class="stat">
        <v-icon
          left
          small
        >
          mdi-heart
        </v-icon>
        <span class="overline">running</span>
        <div
          v-if="failedServices.length == 0"
          class="stat-value"
        >
          All good
        </div>
        <div
          v-else
          class="stat-value"
        >
          {{ runningServices.length || '--' }}<span class="stat-unit"> daemons OK</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
/* eslint-disable vue/prop-name-casing */
import Gauge from '@/components/chart/gauge.vue'

export default {
  components: { Gauge },
  props: {
    cpu: {
      type: Number,
      default: 0
    },
    ram: {
      type: Number,
      default: 0
    },
    hdd_read: {
      type: Number,
      default: 0
    },
    hdd_write: {
      type: Number,
      default: 0
    },
    net_recv: {
      type: Number,
      default: 0
    },
    net_sent: {
      type: Number,
      default: 0
    },
    services: {
      type: Array,
      default: () => []
    }
  },
  computed: {
    runningServices () {
      return this.services
        .filter(s => s.status === 'running')
        .map(s => s.name)
    },
    failedServices () {
      return this.services
        .filter(s => s.status !== 'running')
        .map(s => s.name)
    }
  },
  methods: {
    gaugeColor (value) {
      const colors = this.$vuetify.theme.currentTheme
      if (value < 50) {
        return colors.success
      } else if (value < 80) {
        return colors.warning
      } else {
        return colors.error
      }
    }
  }
}
</script>

<style scoped>
.server-gauges-container {
  display: flex;
  flex-wrap: wrap;
}
.server-gauges-container .gauge-container {
  margin-right: 1.5em;
}

.g-label-wrapper {
  position: relative;
  margin-top: -2em;
}
.g-label-text {
  font-weight: bold;
  font-size: .8em !important;
}

.text-module {
  min-width: 128px;
}

.stat {
  padding-left: .15em;
}
.highlighted-stat {
  margin-left: -.35em;
}

.stat-value {
  padding-left: .1em;
  font-family: Roboto Mono;
  font-size: 1.15em;
}

.stat-unit {
  font-family: Roboto;
  font-weight: 500;
  font-size: .75em;
}
</style>