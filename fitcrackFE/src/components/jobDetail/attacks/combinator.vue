<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile title="Dictionaries and rules">
    <v-card-title>
      Left
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="data.left_dictionaries"
      hide-default-footer
    >
      <template v-slot:item.dictionary.name="{ item: { dictionary } }">
        <router-link :to="{name: 'dictionaryDetail', params: { id: dictionary.id}}">
          {{ dictionary.name }}
          <v-icon 
            small
            color="primary"
          >
            mdi-open-in-new
          </v-icon>
        </router-link>
      </template>
      <template v-slot:item.dictionary.keyspace="{ item: { dictionary } }">
        {{ fmt(dictionary.keyspace) }}
      </template>
    </v-data-table>
    <template v-if="data.rule_left !== ''">
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Left rule:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
            {{ data.rule_left }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>

    <v-divider />
    <v-card-title>
      Right
    </v-card-title>
    <v-data-table
      :headers="rightHeaders"
      :items="data.right_dictionaries"
      hide-default-footer
    >
      <template v-slot:item.dictionary.name="{ item: { dictionary } }">
        <router-link :to="{name: 'dictionaryDetail', params: { id: dictionary.id}}">
          {{ dictionary.name }}
          <v-icon 
            small
            color="primary"
          >
            mdi-open-in-new
          </v-icon>
        </router-link>
      </template>
      <template v-slot:item.dictionary.keyspace="{ item: { dictionary } }">
        {{ fmt(dictionary.keyspace) }}
      </template>
      <template v-slot:item.progress="{ item: { current_index, dictionary: { keyspace } } }">
        <v-progress-circular
          size="16"
          :width="3"
          :rotate="270"
          color="primary"
          :value="(100 / keyspace) * current_index"
        />
        <span>{{ ((100 / keyspace) * current_index).toFixed() }} %</span>
      </template>
    </v-data-table>
    <template v-if="data.rule_right !== ''">
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Right rule:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
            {{ data.rule_right }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
  </fc-tile>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "CombinatorDetail",
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
            value: 'dictionary.name'
          },
          {text: 'Keyspace', value: 'dictionary.keyspace', align: 'right'}
        ]
      }
    },
    computed: {
      rightHeaders () {
        return [
          ...this.headers,
          {text: 'Progress', value: 'progress', align: 'right'}
        ]
      }
    },
    methods: { fmt }
  }
</script>

<style scoped>
  .borderLeft {
    border-left: 2px solid #eee;
  }

</style>
