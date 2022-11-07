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
          { text: 'Charsets', to: { name: 'charsets' }, exact: true },
          { text: data.name }
        ]"
      divider="/"
    />

    <v-container>
      <v-row justify="center">
        <fc-tile
          title="Charsets"
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
                Keyspace:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ fmt(data.keyspace) }}
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
          <v-row class="mx-2 mb-2 align-center">
            <v-col>
              <v-combobox
                v-model="as"
                :items="encodings"
                label="View as"
                outlined
                dense
                hide-details
              />
            </v-col>
            <v-col>
              <v-btn
                color="primary"
                :href="$serverAddr + '/charset/' + data.id + '/download'"
                target="_blank"
              >
                Download
                <v-icon right>
                  mdi-download
                </v-icon>
              </v-btn>
            </v-col>
          </v-row>
          <v-divider />
          <div
            v-if="data != null"
            class="dictContent pa-2"
          >
            <code
              class="code elevation-0"
            >{{ repr }}</code>
          </div>
        </fc-tile>
      </v-row>
    </v-container>
  </div>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import iconv from 'iconv-lite'
  import tile from '@/components/tile/fc_tile.vue'
  import { Buffer } from 'buffer'
  export default {
    name: "CharsetDetailView",
    components: {
      'fc-tile': tile
    },
    data: function () {
      return {
        data: null,
        as: 'hex',
        encodings: [
          'hex',
          'utf-8',
          'cp1250',
          'cp1251',
          'cp1252',
          'iso-8859-1',
          'iso-8859-2',
          'iso-8859-5',
          'iso-8859-15',
          'koi8-r',
          'koi8-u'
        ]
      }
    },
    computed: {
      repr () {
        if (iconv.encodingExists(this.as)) {
          return iconv.decode(Buffer.from(this.data.data, 'base64'), this.as)
        } else {
          this.$error('Unrecognized encoding!')
          return '<unrecognized encoding selected>'
        }
      }
    },
    mounted: function () {
      this.loadData()
    },
    methods: {
      fmt,
      loadData: function ($state) {
        this.axios.get(this.$serverAddr + '/charset/' + this.$route.params.id).then((response) => {
          this.data = response.data
          this.as = this.encodings.find(enc => this.data.name.search(new RegExp(enc, 'i')) >= 0) || 'hex'
        })
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

  .code::before {
    display: none;
  }
  .code {
    width: 100%;
    background: transparent;
    color: #222;
  }
  .theme--dark .code {
    color: #dedede;
  }
  .width100 {
    width: 100%;
  }

  .margintop5 {
    margin-top: 10px;
  }
</style>
