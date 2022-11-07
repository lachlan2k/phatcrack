<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile title="PCFG info">
    <v-list
      single-line
      class="width100"
    >
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Grammar:
        </v-list-item-action>
        <v-list-item-content class="text-right">
          <v-list-item-title>
            <router-link :to="{name: 'pcfg', params: { id: data.grammar_id }}">
              {{ data.grammar_name || `ID ${data.grammar_id}` }}
              <v-icon 
                small
                color="primary"
              >
                mdi-open-in-new
              </v-icon>
            </router-link>
          </v-list-item-title>
          <span>
            Keyspace: {{ fmt(data.grammar_keyspace) || 'not specified' }}
          </span>
        </v-list-item-content>
      </v-list-item>
      <template v-if="data.rulesFile.id !== null">
        <v-divider />
        <v-list-item class="px-2 py-1">
          <v-list-item-action class="pr-3 key">
            Rules:
          </v-list-item-action>
          <v-list-item-content class="text-right">
            <v-list-item-title>
              <router-link :to="{name: 'ruleDetail', params: { id: data.rulesFile.id}}">
                {{ data.rulesFile.name }}
                <v-icon 
                  small
                  color="primary"
                >
                  mdi-open-in-new
                </v-icon>
              </router-link>
            </v-list-item-title>
          <span>
            Count: {{ data.rulesFile.count || 'not specified' }}
          </span>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </fc-tile>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "PcfgDetail",
    components: {
      'fc-tile': tile
    },
    props: ['data'],
    data() {
      return {
        headers: [
          {
            text: 'Name',
            align: 'left',
            value: 'name'
          },
          {text: 'Keyspace', value: 'keyspace', align: 'right'},
          {text: 'Time', value: 'time_added', align: 'right'},
          {text: 'Link to', value: 'name', sortable: false, align: 'right', width: "1"}
        ]
      }
    },
    methods: { fmt }
  }
</script>

<style scoped>

</style>
