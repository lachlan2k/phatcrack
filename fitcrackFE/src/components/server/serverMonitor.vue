<!--
  * Author : see AUTHORS
  * Licence: MIT, see LICENSE
-->

<template>
  <div class="monitor-container">
    <v-sheet class="monitor-sheet services-detail">
      <v-card-title>
        <v-icon left>
          mdi-toolbox-outline
        </v-icon>
        Services
      </v-card-title>
      <div class="services-grid mx-4">
        <div
          v-for="s in services"
          :key="s.name"
          class="service mb-2"
        >
          <div class="status">
            <v-sheet
              :color="statusColor(s.status)"
              class="status-indicator"
            />
            <span :class="['status-text', 'overline', `${statusColor(s.status)}--text`]">{{ s.status }}</span>
          </div>
          {{ s.name }}
        </div>
      </div>
    </v-sheet>
    <v-divider />
    <v-row>
      <v-col class="d-flex px-10 mt-2">
        <v-switch
          v-model="fullDate"
          label="Custom range"
          class="mr-4"
          @input="loadData"
        />
        <div 
          v-if="fullDate"
          class="d-flex"
        >
          <dt-picker
            v-model="chartsFrom"
            class="mr-4"
            label="From"
            @input="loadData"
          />
          <dt-picker
            v-model="chartsTo"
            label="To"
            @input="loadData"
          />
        </div>
        <div v-else>
          <v-radio-group
            v-model="hoursBefore"
            row
            hide-details
            class="dateSelect"
            @change="loadData"
          >
            <v-radio
              value="1"
              label="last hour"
            />
            <v-radio
              value="3"
              label="last 3 hours"
            />
            <v-radio
              value="6"
              label="last 6 hours"
            />
          </v-radio-group>
        </div>
      </v-col>
    </v-row>
    <v-progress-linear
      :style="{opacity: loading ? 1 : 0}"
      indeterminate
    />
    <v-divider />
    <div class="monitor-charts">
      <v-sheet class="monitor-sheet chart">
        <v-card-title>
          <v-icon left>
            mdi-memory
          </v-icon>
          CPU
        </v-card-title>
        <res-mon
          :usage-data="usageData"
          :metrics="['cpu']"
          :max="100"
        />
      </v-sheet>
      <v-sheet class="monitor-sheet chart">
        <v-card-title>
          <v-icon left>
            mdi-chip
          </v-icon>
          Memory
        </v-card-title>
        <res-mon
          :usage-data="usageData"
          :metrics="['ram']"
          :max="100"
        />
      </v-sheet>
      <v-sheet class="monitor-sheet chart">
        <v-card-title>
          <v-icon left>
            mdi-harddisk
          </v-icon>
          Disk
        </v-card-title>
        <res-mon
          :usage-data="usageData"
          :metrics="['hdd_read', 'hdd_write']"
        />
      </v-sheet>
      <v-sheet class="monitor-sheet chart">
        <v-card-title>
          <v-icon left>
            mdi-ethernet
          </v-icon>
          Network
        </v-card-title>
        <res-mon
          :usage-data="usageData"
          :metrics="['net_recv', 'net_sent']"
        />
      </v-sheet>
    </div>
  </div>
</template>

<script>
import ResMon from '@/components/chart/resourceMonitor.vue'
import DtPicker from '@/components/picker/datetime.vue'

export default {
  components: { ResMon, DtPicker },
  data () {
    return {
      infoInterval: null,
      dataInterval: null,
      //
      fullDate: false,
      hoursBefore: '3',
      chartsFrom: this.$moment().subtract(1, 'hours').format('YYYY-MM-DDTHH:mm'),
      chartsTo: this.$moment().format('YYYY-MM-DDTHH:mm'),
      //
      services: [],
      usageData: []
    }
  },
  mounted () {
    this.loadInfo()
    this.infoInterval = setInterval(this.loadInfo, 60000)
    this.loadData(true)
    this.dataInterval = setInterval(this.loadData, 15000)
  },
  beforeDestroy () {
    clearInterval(this.infoInterval)
    clearInterval(this.dataInterval)
  },
  methods: {
    async loadInfo () {
      this.services = await this.axios
        .get(`${this.$serverAddr}/serverInfo/info`)
        .then(r => r.data.subsystems)
    },
    async loadData (manual = false) {
      if (manual) this.loading = true

      let from_date, to_date = ''
      if (this.fullDate) {
        if (this.chartsFrom.length === 16) {
          from_date = this.$moment(this.chartsFrom, 'YYYY-MM-DDTHH:mm').utc().format('YYYY-M-D H:mm:ss')
        } else {
          return
        }
        if (this.chartsTo.length === 16) {
          to_date = this.$moment(this.chartsTo, 'YYYY-MM-DDTHH:mm').utc().format('YYYY-M-D H:mm:ss')
        } else {
          return
        }
      } else {
        from_date = this.$moment().subtract(this.hoursBefore, 'hours').utc().format('YYYY-M-D H:mm:ss')
      }

      this.usageData = await this.axios
        .get(`${this.$serverAddr}/serverInfo/getUsageData`, {
          params: {from_date, to_date}
        })
        .then(r => r.data.items)
      this.loading = false
    },
    statusColor (status) {
      switch (status) {
        case 'running':
          return 'success'
        default:
          return 'error'
      }
    }
  }
}
</script>

<style scoped>
.monitor-sheet {
  padding: 1em;
  padding-top: .5em;
}

.monitor-charts {
  display: flex;
  flex-wrap: wrap;
}

.monitor-charts .chart {
  flex: 1 0 50%;
  min-width: 400px;
}

/**/

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
}

.service {
  font-weight: 500;
}

.status-indicator {
  display: inline-block;
  width: .5em;
  height: .5em;
  border-radius: .5em;
  margin-right: .1em;
}
</style>