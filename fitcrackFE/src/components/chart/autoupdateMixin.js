export default {
  props: {
    updateInterval: {
      type: Number,
      default: 5000
    }
  },
  data () {
    return {
      _interval: null
    }
  },
  mounted () {
    this.loadData()
    if (this.updateInterval > 0) {
      this._interval = setInterval(this.loadData, this.updateInterval)
    }
  },
  beforeDestroy () {
    clearInterval(this._interval)
  }
}