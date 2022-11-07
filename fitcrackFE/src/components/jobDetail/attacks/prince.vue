<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile title="PRINCE attack details">
    <v-data-table
      :headers="headers"
      :items="data.left_dictionaries"
      hide-default-footer
    >
      <template v-slot:item.dictionary.name="{ item: { dictionary } }">
        <router-link :to="{name: 'dictionaryDetail', params: { id: dictionary.id}}" target='_blank'>
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
      <template v-slot:item.dictionary.time="{ item: { dictionary } }">
        {{ $moment.utc(dictionary.time).local().format('DD.MM.YYYY HH:mm') }}
      </template>
    </v-data-table>
    <template v-if="data.rulesFile.id !== null">
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Rules:
        </v-list-item-action>
        <v-list-item-content class="text-right">
          <v-list-item-title>
            <router-link :to="{name: 'ruleDetail', params: { id: data.rulesFile.id}}" target='_blank'>
              {{ data.rulesFile.name }}
            </router-link>
          </v-list-item-title>
          <span>
            Keyspace: {{ fmt(data.rulesFile.count) || 'not specified' }}
          </span>
        </v-list-item-content>
      </v-list-item>
    </template>

      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Number of generated random rules:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.generate_random_rules }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Check for password duplicates:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.check_duplicates ? 'Yes' : ' No' }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Case permutation:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.case_permute ? 'Enabled' : ' Disabled' }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Minimal length of passwords:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.min_password_len }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Maximal length of passwords:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.max_password_len }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Minimal number of elements per chain:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.min_elem_in_chain }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-divider />
      <v-list-item class="px-2 py-1">
        <v-list-item-action class="pr-3 key">
          Maximal number of elements per chain:
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="text-right">
          {{ data.max_elem_in_chain }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
  </fc-tile>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "PrinceDetail",
    components: {
      'fc-tile': tile
    },
    props: {
      data: {
        type: Object,
        default: () => {}
      }
    },
    data() {
      return {
        headers: [
          {
            text: 'Dictionary name',
            align: 'start',
            value: 'dictionary.name'
          },
          {text: 'Keyspace', value: 'dictionary.keyspace', align: 'end'},
          {text: 'Time', value: 'dictionary.time', align: 'end'}
        ]
      }
    },
    methods: { fmt }
  }
</script>

<style scoped>

</style>
