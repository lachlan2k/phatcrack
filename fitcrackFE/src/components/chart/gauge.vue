<template>
  <div class="gauge-container">
    <div class="gauge-inner-wrapper">
      <gauge
        class="gauge"
        :style="{maxHeight: `${size}px`, maxWidth: `${size}px`}"
        :chart-data="chartdata"
        :overrides="options"
      />
      <div class="gauge-value">
        {{ valueFormat(value) }}
      </div>
    </div>
    <div class="gauge-label">
      <slot />
    </div>
  </div>
</template>

<script>
import Gauge from './types/gauge'

export default {
  components: { Gauge },
  props: {
    min: {
      type: Number,
      default: 0
    },
    max: {
      type: Number,
      default: 100
    },
    value: {
      type: Number,
      default: 50
    },
    color: {
      type: String,
      default: undefined
    },
    size: {
      type: [String, Number],
      default: 128
    },
    valueFormat: {
      type: Function,
      default: x => x
    }
  },
  data: () => ({
    options: {
      plugins: {
        legend: {
          display: false
        }
      },
    }
  }),
  computed: {
    fillColor () {
      return this.color || this.$vuetify.theme.currentTheme.primary
    },
    chartdata () {
      const scale = (this.max - this.min)
      const filled = this.value
      const rest = scale - filled

      return {
        datasets: [
          {
            backgroundColor: [this.fillColor, '#9994'],
            borderWidth: 0,
            data: [filled, rest]
          }
        ]
      }
    }
  }
}
</script>

<style scoped>
.gauge-container,
.gauge-inner-wrapper {
  position: relative;
}

.gauge-value {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -.4em);
  font-family: Roboto Mono;
  font-size: 1.25em;
  font-weight: bold;
}

.gauge-label {
  text-align: center;
  animation: in .7s
}

@keyframes in {
  from {
    transform: scale(0) translateY(-1em);
    opacity: 0;
  }
}
</style>