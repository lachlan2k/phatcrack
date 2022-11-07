export default {
  data () {
    return {
      _timeout: null
    }
  },
  watch: {
    from () {this.triggerLoad()},
    to () {this.triggerLoad()}
  },
  methods: {
    triggerLoad () {
      this.loading = true
      clearTimeout(this._timeout)
      this._timeout = setTimeout(() => {
        this.loadData(true)
      }, 100)
    }
  }
}