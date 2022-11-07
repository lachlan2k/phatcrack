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
      <router-link :to="{name: 'maskDetail', params: { id: item.id}}">
        {{ item.name }}
        <v-icon 
          small
          color="primary"
        >
          mdi-open-in-new
        </v-icon>
      </router-link>
    </template>
    <template v-slot:item.time="{ item }">
      {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
    </template>
    <!-- ???
    <v-checkbox
      :input-value="selected === props.item.id"
      primary
      hide-details
      @click="updateSelected(props.item.id, props.item)"
    />
    -->
  </v-data-table>
</template>

<script>
  import selector from './selectorMixin'
  export default {
    name: "MaskFileSelector",
    mixins: [selector],
    props: {
      value: {
        type: Object,
        default: function () {
          return {}
        }
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
          {text: 'Added', value: 'time', align: 'end'},
        ]
      }
    },
    computed: {
      selectedId () {
        return selected[0].id
      }
    },
    watch: {
      value (val) {
        this.selected = [val]
      }
    },
    methods: {
      getData() {
        this.loading = true
        this.axios.get(this.$serverAddr + '/masks').then((response) => {
          this.items = response.data.items
          this.loading = false
        })
      },
      updateSelected() {
        this.$emit('input', this.selected[0])
      }
    }
  }
</script>

<style scoped>
  .oneline {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
    width: 200px;
    line-height: 50px;
    height: 50px;
  }
</style>
