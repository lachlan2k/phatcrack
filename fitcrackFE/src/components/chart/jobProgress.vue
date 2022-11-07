<template>
  <div class="container">
    <v-progress-linear
      v-show="loading"
      indeterminate
      absolute
    />
    <timeseries
      v-if="loaded"
      :chart-data="chartdata"
      :overrides="options"
      class="chart-wrapper"
    />
  </div>
</template>

<script>
import Timeseries from './types/timeseries'
import { prepareLines, localizeTime } from './helpers'
import autoload from './autoupdateMixin'
import watcher from './boundWatcherMixin'

export default {
  components: { Timeseries },
  mixins: [autoload, watcher],
  props: {
    id: {
      type: Number,
      default: undefined
    },
    from: {
      type: String,
      default: undefined
    },
    to: {
      type: String,
      default: undefined
    }
  },
  data: () => ({
    loaded: false,
    loading: false,
    chartdata: null,
    options: {
      elements: {
        line: {
          fill: true
        }
      },
      plugins: {
        legend: {
          display: false
        }
      },
      scales: {
        x: {
          gridLines: {
            display: false
          }
        },
        y: {
          ticks: {
            callback (val) {
              return `${val} %`
            }
          },
          min: 0,
          max: 100
        }
      }
    }
  }),
  computed: {
    fromTime () {
      if (!this.id && !this.from) {
        return this.$moment.utc().subtract(24, 'hours').format('YYYY-M-DTH:mm:ss')
      } else {
        return this.$moment(this.from).utc().format('YYYY-M-DTH:mm:ss')
      }
    },
    toTime () {
      return this.to ? this.$moment(this.to).utc().format('YYYY-M-DTH:mm:ss') : undefined
    }
  },
  methods: {
    loadData (manual = false) {
      if (manual) {
        this.loading = true
      }
      let target = `${this.$serverAddr}/chart/jobProgress`
      if (this.id) target += `/${this.id}`
      this.axios.get(target, {
          params: {
            from: this.fromTime,
            to: this.toTime
          }
        }).then(r => {
        this.chartdata = {
          datasets: localizeTime(prepareLines(r.data.datasets))
        }
        this.loading = false
        this.loaded = true
      })
    }
  }
}
</script>

<style scoped>
.container {
  position: relative;
}
</style>