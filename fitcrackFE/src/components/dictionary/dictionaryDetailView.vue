<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-breadcrumbs
      v-if="info != null"
      :items="
        [
          { text: 'Dictionaries', to: { name: 'dictionaries' }, exact: true },
          { text: info.name }
        ]"
      divider="/"
    />

    <v-container>
        <fc-tile
          title="Dictionary"
          :loading="info==null"
          class="mx-2 dictContentContainer mb-4"
        >
          <v-list
            v-if="info != null"
          >
            <v-list-item>
              <v-list-item-action>
                Name:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ info.name }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider />
            <v-list-item>
              <v-list-item-action>
                Keyspace:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ fmt(info.keyspace) }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider />
            <v-list-item>
              <v-list-item-action>
                Added:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ $moment.utc(info.time).local().format('DD.MM.YYYY HH:mm') }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-divider />
          <div
            v-if="info != null"
            class="dictContent"
          >
            <fc-textarea
              :search-enable="true"
              :readonly="true"
              :url="this.$serverAddr + '/dictionary/' + this.$route.params.id + '/data'"
              max-height="500"
            />
          </div>
        </fc-tile>
    </v-container>
  </div>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import fcTextarea from '@/components/textarea/fc_textarea.vue'
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "DictionaryView",
    components: {
      'fc-tile': tile,
      'fc-textarea': fcTextarea,
    },
    data: function () {
      return {
        info: null
      }
    },
    mounted: function () {
      this.loadData()
    },
    methods: {
      fmt,
      loadData: function ($state) {
        this.axios.get(this.$serverAddr + '/dictionary/' + this.$route.params.id).then((response) => {
            this.info = response.data
        });
      }
    }
  }
</script>

<style scoped>
  .dictContentContainer {
    max-width: 800px;
  }
  .dictContent {
    max-height: 600px;
  }

</style>
