<!--
  * Author : see AUTHORS
  * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile title="Hybrid attack details">
    <v-card-title>
      Mask
    </v-card-title>
    <v-data-table
      :headers="maskHeaders"
      :items="data.masks"
      hide-default-footer
    />

    <v-divider />
    <v-card-title>
      Dictionaries
    </v-card-title>
    <v-data-table
      :headers="headers"
      :items="dicts"
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
    <template v-if="hasRule">
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Rule:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
            {{ rule }}
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
        maskHeaders: [
          {
            text: 'Mask',
            align: 'left',
            value: 'mask'
          }
        ]
      }
    },
    computed: {
      dicts () { return this.data.attack_mode === '7' ? this.data.right_dictionaries : this.data.left_dictionaries },
      hasRule () { return this.data.rule_left !== '' || this.data.rule_right !== '' },
      rule () { return this.data.attack_mode === '7' ? this.data.rule_right : this.data.rule_left },
      headers () {
        const base = [
          {
            text: 'Name',
            align: 'left',
            value: 'dictionary.name'
          },
          {text: 'Keyspace', value: 'dictionary.keyspace', align: 'right'}
        ]
        return this.data.attack_mode === '6' ? base : [
          ...base,
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
