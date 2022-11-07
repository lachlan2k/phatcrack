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
    <template v-slot:item.name="{ item }">
      <router-link :to="{name: 'ruleDetail', params: { id: item.id}}">
        {{ item.name }}
        <v-icon 
          small
          color="primary"
        >
          mdi-open-in-new
        </v-icon>
      </router-link>
    </template>
    <template v-slot:item.count="{ item }">
      {{ fmt(item.count) }}
    </template>
    <template v-slot:item.time="{ item }">
      {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
    </template>
  </v-data-table>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import selector from './selectorMixin'
  export default {
    name: "RulesSelector",
    mixins: [selector],
    data() {
      return {
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Count', value: 'count', align: 'end'},
          {text: 'Added', value: 'time', align: 'end'}
        ]
      }
    },
    methods: {
      fmt,
      getData() {
        this.loading = true
        this.axios.get(this.$serverAddr + '/rule').then((response) => {
          this.items = response.data.items
          this.loading = false
        })
      }
      // EMIT OBJ
    },
  }
</script>
