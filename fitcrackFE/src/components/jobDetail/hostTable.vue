<template>
  <v-data-table
    :headers="headers"
    :items="hosts"
    hide-default-footer
    no-data-text="None assigned"
  >
    <template v-slot:item.name="{ item }">
      <div class="hostname">
        <div
          v-if="showChartPatterns"
          class="pattern"
        >
          <canvas
            :ref="`pattern-${item.id}`"
            width="15"
            height="15"
          />
        </div>
        <router-link
          :to="{ name: 'hostDetail', params: { id: item.id}}"
          class="middle"
        >
          {{ item.domain_name + ' (' + fixUserNameEncoding(item.user.name) + ')' }}
        </router-link>
      </div>
    </template>
    <template v-slot:item.last_active="{ item }">
      <span v-if="item.last_active.seconds_delta > 61">{{ parseTimeDelta(item.last_active.last_seen) }}</span>
      <v-icon
        v-else
        color="success"
      >
        mdi-power
      </v-icon>
    </template>
  </v-data-table>
</template>

<script>
import fixUserNameEncoding from '@/assets/scripts/unswedishify'
import { getColors } from '@/components/chart/helpers.js'

export default {
  props: {
    hosts: {
      type: Array,
      default: () => []
    },
    showChartPatterns: Boolean
  },
  data () {
    return {
      headers: [
        {
          text: 'Name',
          align: 'start',
          value: 'name'
        },
        {text: 'IP address', value: 'ip_address', align: 'end', sortable: false},
        {text: 'Online', value: 'last_active', align: 'end', sortable: false}
      ]
    }
  },
  mounted () {
    // draw patterns if shown
    if (this.showChartPatterns) {
      const canvases = Object.values(this.$refs)
      const patterns = getColors(canvases.length, true)
      canvases.forEach((cvs, i) => {
        const tc = cvs.getContext("2d")
        tc.rect(0, 0, cvs.width, cvs.height)
        tc.fillStyle = patterns[i]
        tc.fill()
      })
    }
  },
  methods: {
    fixUserNameEncoding,
    //
    parseTimeDelta (delta) {
      if (delta !== null && delta !== undefined) {
        return this.$moment.utc(delta).fromNow()
      } else {
        return 'Unknown'
      }
    }
  }
}
</script>

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