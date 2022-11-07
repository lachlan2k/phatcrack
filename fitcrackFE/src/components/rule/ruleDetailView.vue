<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-breadcrumbs
      v-if="data != null"
      :items="
        [
          { text: 'Rules', to: { name: 'rules' }, exact: true },
          { text: data.name }
        ]"
      divider="/"
    />

    <v-container>
        <fc-tile
          title="Rule file"
          :loading="data==null"
          class="mx-2 dictContentContainer mb-4"
        >
          <v-list
            v-if="data != null"
          >
            <v-list-item>
              <v-list-item-action>
                Name:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ data.name }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider />
            <v-list-item>
              <v-list-item-action>
                Rule count:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ fmt(data.count) }}
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
                  {{ $moment.utc(data.time).local().format('DD.MM.YYYY HH:mm') }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-row class="mx-2 mb-2">
            <v-col>
              <v-btn
                color="primary"
                :href="$serverAddr + '/rule/' + data.id + '/download'"
                target="_blank"
              >
                Download <v-icon right>mdi-download</v-icon>
              </v-btn>
            </v-col>
          </v-row>
          <v-divider />
          <div
            v-if="data != null"
            class="dictContent"
          >
            <fc-textarea
              :readonly="true"
              :search-enable="true"
              :url="this.$serverAddr + '/rule/' + this.$route.params.id + '/data'"
              max-height="500"
            />
          </div>
        </fc-tile>
    </v-container>
  </div>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  import fcTextarea from '@/components/textarea/fc_textarea.vue'
  export default {
    name: "RuleDetailView",
    components: {
      'fc-tile': tile,
      'fc-textarea': fcTextarea
    },
    data: function () {
      return {
        data: null
      }
    },
    mounted: function () {
      this.loadData()
    },
    methods: {
      fmt,
      loadData: function ($state) {
        this.axios.get(this.$serverAddr + '/rule/' + this.$route.params.id).then((response) => {
          this.data = response.data
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
    overflow: auto;
  }

  .width100 {
    width: 100%;
  }

  .margintop5 {
    margin-top: 10px;
  }
</style>
