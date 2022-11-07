<template>
  <div>
    <v-card-title>
      <v-icon left>
        mdi-package-variant-closed
      </v-icon>
      <b class="title">Workunits</b>
      <v-spacer />
      <v-chip class="mr-2">
        <v-avatar left>
          <v-icon small>
            mdi-chart-donut
          </v-icon>
        </v-avatar>
        {{ title.benchmarks }} benchmark
      </v-chip>
      <v-chip class="mr-2">
        <v-avatar left>
          <v-icon small>
            mdi-hammer
          </v-icon>
        </v-avatar>
        {{ title.valid }} work
      </v-chip>
      <v-chip>
        <v-avatar left>
          <v-icon small>
            mdi-key
          </v-icon>
        </v-avatar>
        {{ fmt(title.avgKeyspace) }} average keyspace
      </v-chip>
    </v-card-title>

    <div class="d-flex">
      <div
        v-for="workunit in graphical"
        :key="workunit.id"
        :style="{ 'flex-grow': workunit.keyspace }"
        :class="['workunit-child', workunit.state]"
      >
        <v-tooltip top>
          <template v-slot:activator="{ on }">
            <div v-on="on">
&nbsp;
            </div>
          </template>
          <span>{{ workunit.text }}</span>
        </v-tooltip>
      </div>
    </div>

    <v-data-table
      :footer-props="footer"
      :headers="headers"
      :items="data.workunits"
      show-expand
    >
      <template v-slot:item.host_id="{ item }">
        <div class="hostname">
          <div
            v-if="showChartPatterns"
            class="pattern"
          >
            <canvas
              :ref="`pattern-${item.id}`"
              :data-host-id="item.host_id"
              width="15"
              height="15"
            />
          </div>
          <v-tooltip top>
          <template #activator="{ on }">
            <router-link
              :to="{ name: 'hostDetail', params: { id: item.host.id}}"
              class="middle"
            >
              <span v-on="on">{{ item.host.domain_name }}</span>
            </router-link>
          </template>
          <span>User: {{ fixUserNameEncoding(item.host.user.name) }}</span>
        </v-tooltip>
        </div>
      </template>
      <template v-slot:item.progress="{ item }">
        <div class="d-flex align-center justify-end">
          <span class="mr-2">{{ item.progress }}&nbsp;%</span>
          <v-progress-circular
            size="18"
            :width="3"
            :rotate="270"
            color="primary"
            class="jobProgress"
            :value="item.progress"
          />
        </div>
      </template>
      <template v-slot:item.speed="{ item }">
        <span v-if="item.speed >= 1000000">
          {{ fmt(item.speed / 1000, {maximumFractionDigits: 0}) }}&nbsp;kH/s
        </span>
        <span v-else-if="item.speed > 0">
          {{ fmt(item.speed) }}&nbsp;H/s
        </span>
        <span v-else>–</span>
      </template>
      <template v-slot:item.time="{ item }">
        {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
      </template>
      <template v-slot:item.start_index_real="{ item }">
        {{ fmt(item.start_index_real) }}
      </template>
      <template v-slot:item.keyspace="{ item }">
        {{ fmt(item.keyspace) }}
      </template>
      <template v-slot:item.retry="{ item }">
        <v-chip
          small
          :color="item.retry ? 'error' : 'success'"
        >
          {{ item.retry ? 'Yes' : 'No' }}
        </v-chip>
      </template>
      <template v-slot:item.finished="{ item }">
        <v-chip
          small
          :color="item.finished ? 'success' : 'error'"
        >
          {{ item.finished ? 'Yes' : 'No' }}
        </v-chip>
      </template>
      <template
        v-slot:expanded-item="{ headers, item }"
      >
        <td :colspan="headers.length">
          <fc-textarea
            max-height="500"
            :readonly="true"
            :value="item.result.stderr_out_text"
          />
        </td>
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { getColors } from '@/components/chart/helpers.js'
import fixUserNameEncoding from '@/assets/scripts/unswedishify'
import fmt from '@/assets/scripts/numberFormat'
import fcTextarea from '@/components/textarea/fc_textarea.vue'

export default {
  components: {
    fcTextarea
  },
  props: {
    data: {
      type: Object,
      default: () => {}
    },
    showChartPatterns: Boolean
  },
  mounted () {
    // draw patterns if shown
    if (this.showChartPatterns) {
      const canvases = Object.values(this.$refs)
      const patterns = getColors(canvases.length, true)
      canvases.forEach(cvs => {
        console.log(cvs)
        const tc = cvs.getContext("2d")
        tc.rect(0, 0, cvs.width, cvs.height)
        tc.fillStyle = patterns[cvs.dataset.hostId - 1]
        tc.fill()
      })
    }
  },
  data () {
    return {
      title: {
        valid: 0,
        benchmarks: 0,
        avgKeyspace: 0,
        efectivity: 0
      },
      headers: [
        {
          text: 'Host',
          value: 'host_id',
        },
        {text: 'Progress', align: 'end', value: 'progress'},
        {text: 'Speed', align: 'end', value: 'speed'},
        {text: 'Cracking time', align: 'end', value: 'cracking_time_str'},
        {text: 'Generated', align: 'end', value: 'time'},
        {text: 'Start index', align: 'end', value: 'start_index_real'},
        {text: 'Keyspace', align: 'end', value: 'keyspace'},
        {text: 'Retry', align: 'center', value: 'retry'},
        {text: 'Finished', align: 'center', value: 'finished'},
        {text: 'Log', align: 'center', value: 'data-table-expand'}
      ],
      footer: {
        itemsPerPageOptions: [15,30,60,{text: 'All', value: -1}], 
        itemsPerPageText: 'Workunits per page'
      },
      graphical: []
    }
  },
  watch: {
    data: {
      immediate: true,
      deep: true,
      handler (job) {
              
        // Computing of counts and avg keyspace
        if (job.workunits.length > 0) {
          const validWorkunits = job.workunits.filter(workunit => workunit.hc_keyspace !== 0)
          this.title.valid = validWorkunits.length
          this.title.benchmarks = job.workunits.length - validWorkunits.length
          if (validWorkunits.length > 0) {
            this.title.avgKeyspace = (validWorkunits.map(workunit => workunit.keyspace).reduce((a, b) => a + b) / validWorkunits.length)
          }
          

          // Prepare objects for workunits
          this.graphical = validWorkunits.map(workunit => {

            // retry && finished
            let state = ""

            if (!workunit.retry && workunit.finished) {
              state = 'success'
            } else if (!workunit.retry && !workunit.finished){
              state = 'info'
            } else if (workunit.retry && !workunit.finished) {
              state = 'error'
            } else {
              state = 'warning'
            }

            return {
              id: workunit.id,
              keyspace: workunit.keyspace,
              state,
              text: `ID ${workunit.id} (keyspace ${workunit.keyspace}, hc_keyspace ${workunit.hc_keyspace}) ${workunit.retry ? 'retried and' : ''} ${workunit.finished ? 'finished' : 'not finished'}`,
            }
          })
        }
      }
    }
  },
  methods: {
    fixUserNameEncoding,
    fmt
  }
}
</script>

<style>
.workunit-child:not(:last-child) {
  margin-right: 2px;
  user-select: none;
}
</style>

<style scoped>
.hostname {
  display: flex;
  align-items: center;
}
.pattern {
  margin-right: 1ch;
  display: flex;
  align-items: center;
}
.pattern canvas {
  border: 2px solid #fffd;
  border-radius: 5px;
}
</style>