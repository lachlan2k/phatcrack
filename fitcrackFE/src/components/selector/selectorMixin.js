export default {
  props: {
    selectAll: {
      type: Boolean,
      default: false
    },
    value: {
      type: Array,
      default: function () {
        return []
      }
    }
  },
  data() {
    return {
      items: [],
      loading: false,
      search: '',
      selected: [],
      footerProps: {
        itemsPerPageOptions: [3, 5, 10, 20, -1]
      }
    }
  },
  mounted() {
    this.getData()
    if (this.value) this.selected = this.value
  },
  computed: {
    itemsPerPage() {
      // return value of perpage in localStorage or default value
      return parseInt(localStorage.getItem('defaultItemsPerPage') || 10)
    }
  },
  methods: {
    getData() {
      console.error('Selector component missing getData method')
    },
    updateSelected() {
      this.$emit('input', this.selected)
    }
  },
  watch: {
    value (val) {
      this.selected = val
    }
  }
}