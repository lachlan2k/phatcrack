import { Line, mixins } from '@oviovocny/vue-chartjs'
import deepMerge from 'deepmerge'
import { gridColor } from '../helpers'

export default {
  extends: Line,
  mixins: [mixins.reactiveProp],
  props: {
    overrides: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    options () {
      return deepMerge({
        maintainAspectRatio: false,
        elements: {
          line: {
            fill: false,
            tension: 0
          }
        },
        scales: {
          x: {
            type: 'time',
            time: {
              tooltipFormat: 'YYYY-MM-DD HH:mm',
              displayFormats: {
                  millisecond: 'HH:mm:ss.SSS',
                  second: 'HH:mm:ss',
                  minute: 'HH:mm',
                  hour: 'HH'
              }
            },
            grid: {
              color: gridColor
            }
          },
          y: {
            grid: {
              color: gridColor
            }
          }
        },
        plugins: {
          tooltip: {
            mode: "nearest",
            intersect: false
          }
        }
      }, this.overrides)
    }
  },
  mounted () {
    this.renderChart(this.chartData, this.options)
  }
}