
<template>
  <div class="dash-container">
    <v-sheet class="dash-row pa-0">
      <status
        v-if="hostsInfo"
        :jobs="jobsInfo"
        :hosts="hostsInfo"
        :server="serverInfo"
      >
        <server-gauges
          v-bind="actualUsage"
          :services="serverInfo"
        />
      </status>
    </v-sheet>
    <v-divider />
    <v-card-title class="px-8">
      Last jobs
    </v-card-title>
    <div
      v-if="lastJobs"
      class="minijobs overflowing dash-row pt-0 pb-7 px-8"
    >
      <minijob
        detailed
        :data="lastJobs[0]"
      />
      <minijob
        v-for="j in lastJobs.slice(1)"
        :key="j.id"
        :data="j"
      />
      <div class="pr-6">
        <v-btn
          to="/jobs"
          large
          color="primary"
        >
          See all
          <v-icon right>
            mdi-briefcase-search
          </v-icon>
        </v-btn>
      </div>
    </div>
    <v-divider />
    <v-sheet class="dash-row px-4">
      <div class="half">
        <v-card-title>Progress of all jobs</v-card-title>
        <jobProgress
          :from="chartsFrom"
          :to="chartsTo"
        />
      </div>
      <div class="half">
        <v-card-title>Global workunits distribution</v-card-title>
        <jobWorkunits
          logarithmic
          :from="chartsFrom"
          :to="chartsTo"
        />
      </div>
    </v-sheet>
    <v-sheet class="dash-row px-8 justify-end">
      <div class="d-flex">
        <dt-picker
          v-model="chartsFrom"
          class="mr-4"
          label="From"
        />
        <dt-picker
          v-model="chartsTo"
          label="To"
        />
      </div>
    </v-sheet>
  </div>
</template>

<script>
  import status from '@/components/dashboard/status.vue'
  import minijob from '@/components/dashboard/miniJob.vue'
  import serverGauges from '@/components/dashboard/serverGauges.vue'
  
  import jobProgress from '@/components/chart/jobProgress.vue'
  import jobWorkunits from '@/components/chart/jobWorkunits.vue'

  import { attackIcon } from '@/assets/scripts/iconMaps'
  import dtPicker from '@/components/picker/datetime.vue'

  export default {
    name: 'Dashboard',
    components: {
      status,
      minijob,
      serverGauges,
      jobProgress,
      jobWorkunits,
      dtPicker
    },
    data () {
      return {
        chartsFrom: this.$moment().subtract(24, 'hours').format('YYYY-MM-DDTHH:mm'),
        chartsTo: this.$moment().format('YYYY-MM-DDTHH:mm'),
        hostsInfo: null,
        jobsInfo: [],
        serverInfo: [],
        lastJobs: null,
        actualUsage: {
          cpu: 0,
          ram: 0,
          hdd_read: 0,
          hdd_write: 0,
          net_recv: 0,
          net_sent: 0
        }
      }
    },
    mounted () {
      this.loadData()
      this.interval = setInterval(this.loadData, 5000)
    },
    beforeDestroy () {
      clearInterval(this.interval)
    },
    methods: {
      loadData () {
        this.axios.get(this.$serverAddr + '/job/lastJobs').then((response) => {
          this.lastJobs = response.data.items;
        })

        this.axios.get(this.$serverAddr + '/serverInfo/actualUsageData').then((response) => {
          this.actualUsage = response.data
        })

        this.axios.get(this.$serverAddr + '/hosts/info').then((response) => {
          this.hostsInfo = response.data
        })
        this.axios.get(this.$serverAddr + '/job/info').then((response) => {
          this.jobsInfo = response.data
        })
        this.axios.get(this.$serverAddr + '/serverInfo/info').then((response) => {
          this.serverInfo = response.data.subsystems
        })
      }
    }
  }
</script>

<style scoped>
.dash-row {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
}
.overflowing {
  flex-wrap: nowrap;
  overflow-x: auto;
}
.minijobs {
  align-items: stretch;
}

.half {
  width: 50%;
}
@media screen and (max-width: 600px) {
  .half {
    width: 100%;
  }
}
</style>