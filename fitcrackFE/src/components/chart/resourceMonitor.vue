<template>
  <div class="container">
    <timeseries
      :chart-data="chartdata"
      :overrides="options"
      class="chart-wrapper"
    />
  </div>
</template>

<script>
import Timeseries from './types/timeseries'
import { prepareLines, plotUsage } from './helpers'

export default {
  components: { Timeseries },
  props: {
    usageData: {
      type: Array,
      default: () => []
    },
    metrics: {
      type: Array,
      default: () => []
    },
    unitCallback: {
      type: Function,
      default: x => x
    },
    min: {
      type: Number,
      default: 0
    },
    max: {
      type: Number,
      default: undefined
    },
    rightScale: {
      type: Boolean
    },
    unitCallbackRight: {
      type: Function,
      default: x => x
    },
    minRight: {
      type: Number,
      default: 0
    },
    maxRight: {
      type: Number,
      default: undefined
    }
  },
  computed: {
    chartdata () {
      return {
        datasets: prepareLines(
          this.metrics.map(metric => ({
            data: plotUsage(this.usageData, metric),
            label: this.metricLabel(metric)
          }))
        )
      }
    },
    options () {
      return {
        elements: {
          line: {
            fill: true
          },
          point: {
            radius: 0
          }
        },
        plugins: {
          tooltip: {
            mode: "index",
            position: "nearest",
            callbacks: {
              label: ctx => {
                const name = ctx.dataset.label
                return `${ctx.formattedValue} ${this.labelUnit(name)}`
              }
            }
          }
        },
        scales: {
          x: {
            type: 'time',
            time: {
              stepSize: 5,
              tooltipFormat: 'YYYY-MM-DD HH:mm',
              displayFormats: {
                  millisecond: 'HH:mm:ss.SSS',
                  second: 'HH:mm:ss',
                  minute: 'HH:mm',
                  hour: 'HH'
              }
            },
            grid: {
              display: false
            }
          },
          y: {
            ticks: {
              callback: x => {
                const name = this.metricLabel(this.metrics[0])
                return `${x} ${this.labelUnit(name)}`
              }
            },
            min: this.min,
            max: this.max
          },
          y1: {
            position: 'right',
            display: this.rightScale,
            ticks: {
              callback: x => {
                const name = this.metricLabel(this.metrics[1])
                return `${x} ${this.labelUnit(name)}`
              }
            },
            min: this.minRight,
            max: this.maxRight
          }
        }
      }
    }
  },
  methods: {
    metricLabel (metric) {
      const map = {
        cpu: 'CPU',
        hdd_read: 'Disk Read',
        hdd_write: 'Disk Write',
        net_recv: 'Download',
        net_sent: 'Upload',
        ram: 'Memory',
        speed: 'Speed',
        temperature: 'Temperature',
        utilization: 'Utilization'
      }
      return map[metric]
    },
    labelUnit (label) {
      const map = {
        'CPU': '%',
        'Disk Read': 'kb/s',
        'Disk Write': 'kb/s',
        'Download': 'kb/s',
        'Upload': 'kb/s',
        'Memory': '%',
        'Speed': 'H/s',
        'Temperature': '°C',
        'Utilization': '%'
      }
      return map[label]
    }
  }
}
</script>