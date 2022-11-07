<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-data-table
    v-model="selected"
    :headers="headers"
    :items="items"
    :items-per-page="itemsPerPage"
    :footer-props="footerProps"
    :search="search"
    item-key="id"
    show-select
    :single-select="!selectAll"
    @input="updateSelected"
  >
    <template v-slot:item.data-table-select="{ select, isSelected }">
      <v-simple-checkbox
        :value="isSelected"
        :disabled="markovSubmode == 0"
        @input="select"
      />
    </template>
    <template v-slot:item.time="{ item }">
      {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
    </template>
  </v-data-table>
</template>

<script>
  import selector from './selectorMixin'
  export default {
    name: "MarkovSelector",
    mixins: [selector],
    props: {
      value: {
        type: Array,
        default: () => ([])
      },
      markovSubmode: {
        type: Number,
        default: 2
      }
    },
    data() {
      return {
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Added', value: 'time', align: 'end'}
        ],
      }
    },
    methods: {
      getData() {
        this.loading = true
        this.axios.get(this.$serverAddr + '/markovChains').then((response) => {
          this.items = response.data.items
          this.selected = [this.items.find(i => i.name === 'hashcat.hcstat2')]
          this.loading = false
        })
      }
    }
  }
</script>
