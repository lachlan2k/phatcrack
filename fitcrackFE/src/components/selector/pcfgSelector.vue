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
      <router-link :to="{name: 'pcfgDetail', params: { id: item.id}}">
        {{ item.name }}
        <v-icon 
          small
          color="primary"
        >
          mdi-open-in-new
        </v-icon>
      </router-link>
    </template>
    <template v-slot:item.keyspace="{ item }">
      {{ fmt(item.keyspace) }}
    </template>
    <template v-slot:item.time_added="{ item }">
      {{ $moment.utc(item.time_added).local().format('DD.MM.YYYY HH:mm') }}
    </template>
  </v-data-table>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import selector from './selectorMixin'
  export default {
    name: "PcfgSelector",
    mixins: [selector],
    data() {
      return {
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Keyspace', value: 'keyspace', align: 'end'},
          {text: 'Added', value: 'time_added', align: 'end'}
        ]
      }
    },
    watch:{
      value: function(){
        if (this.value) {
          this.selected = this.value
        }
      }
    },
    methods: {
      fmt,
      getData() {
        this.loading = true
        this.axios.get(this.$serverAddr + '/pcfg').then((response) => {
          this.items = response.data.items
          this.loading = false
        })
      }
      // EMIT OBJ
    }
  }
</script>
