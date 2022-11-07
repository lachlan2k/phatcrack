import { Doughnut, mixins } from '@oviovocny/vue-chartjs'
import deepMerge from 'deepmerge'

export default {
  extends: Doughnut,
  mixins: [mixins.reactiveProp],
  props: {
    overrides: {
      type: Object,
      default: null
    }
  },
  computed: {
    options () {
      return deepMerge({
        maintainAspectRatio: false,
        plugins: {
          tooltip: {
            enabled: false
          }
        },
        hover: {
          mode: null
        },
        animation: {
          animateScale: true
        },
        cutout: "65%",
        rotation: -135,
        circumference: 270
      }, this.overrides)
    }
  },
  mounted () {
    this.renderChart(this.chartData, this.options)
  }
}