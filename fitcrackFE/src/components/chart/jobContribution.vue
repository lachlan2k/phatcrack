<template>
  <div class="container">
    <doughnut
      v-if="loaded"
      class="chart chart-wrapper"
      :chart-data="chartdata"
      :overrides="options"
    />
  </div>
</template>

<script>
import Doughnut from './types/doughnut'
import { getColors } from './helpers'
import autoload from './autoupdateMixin'

export default {
  components: { Doughnut },
  mixins: [autoload],
  props: {
    id: {
      type: Number,
      default: undefined
    },
    batch: Boolean
  },
  data: () => ({
    loaded: false,
    chartdata: null,
    options: {
      plugins: {
        legend: {
          display: false
        }
      },
    }
  }),
  methods: {
    loadData () {
      const endpoint = this.batch ? 'batch' : 'job'
      let target = `${this.$serverAddr}/chart/${endpoint}Contribution/${this.id}`
      this.axios.get(target).then(r => {

        const datasetCount = r.data.labels.length
        const patterns = getColors(datasetCount, true)

        this.chartdata = {
          labels: r.data.labels,
          datasets: [
            {
            backgroundColor: patterns,
            data: r.data.data
            }
          ]
        }
        this.loaded = true
      })
    }
  }
}
</script>

<style scoped>
.chart {
  height: 250px;
}
</style>