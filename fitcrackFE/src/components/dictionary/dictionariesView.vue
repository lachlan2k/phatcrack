<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max500">
    <fc-tile
      title="Dictionaries"
      :icon="$route.meta.icon"
      class="ma-2"
    >
      <v-alert
        tile
        text
        type="warning"
        class="mb-0"
      >
        Dictionaries must have a .txt extension.
      </v-alert>
      <v-data-table
        :headers="headers"
        :items="dictionaries.items"
        :loading="loading"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Dictionaries per page'}"
      >
        <template v-slot:item.name="{ item }">
          <router-link :to="`dictionaries/${item.id}`">
            {{ item.name }}
          </router-link>
        </template>
        <template v-slot:item.time="{ item }">
          {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.keyspace="{ item }">
          {{ fmt(item.keyspace) }}
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <a
                :href="$serverAddr + '/dictionary/' + item.id + '/download'"
                target="_blank"
                download
                v-on="on"
              >
                <v-btn icon>
                  <v-icon>mdi-file-download-outline</v-icon>
                </v-btn>
              </a>
            </template>
            <span>Download</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                class="mx-0"
                @click="deleteDictionary(item)"
                v-on="on"
              >
                <v-icon color="error">
                  mdi-delete-outline
                </v-icon>
              </v-btn>
            </template>
            <span>Delete dictionary</span>
          </v-tooltip>
        </template>
      </v-data-table>
      <v-card-actions v-if="$userCan('UPLOAD_DICTIONARIES')">
        <v-checkbox
          v-model="sortUploaded"
          label="Sort on upload"
          hint="Sort by password length"
          persistent-hint
        />
        <v-checkbox
          v-model="hexDict"
          label="HEX dictionary"
          hint="Assume words in the dictionary are given in hex"
          persistent-hint
        />
        <v-spacer />
        <v-btn
          color="primary"
          outlined
          @click.native.stop="browser = true"
        >
          Add from server
        </v-btn>
        <v-btn
          color="primary"
          outlined
          @click.native.stop="uploader = true"
        >
          Upload new
        </v-btn>
      </v-card-actions>
    </fc-tile>

    <v-dialog
      v-model="browser"
      max-width="500"
    >
      <server-browser
        :sort="sortUploaded"
        :hex_dict="hexDict"
        @filesuploaded="browser = false;loadDictionaries()"
      />
    </v-dialog>

    <v-dialog
      v-model="uploader"
      max-width="500"
    >
      <v-card>
        <v-card-title>
          Upload a dictionary file
        </v-card-title>
        <v-card-text>
          <file-uploader
            :url="this.$serverAddr + '/dictionary/add'"
            :args="{sort: sortUploaded, hex_dict: hexDict}"
            @uploadComplete="uploader = false;loadDictionaries()"
          />
        </v-card-text>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  import serverBrowser from '@/components/serverBrowser/serverBrowser.vue'
  import FileUploader from "@/components/fileUploader/fileUploader.vue";

  export default {
    name: "DictionariesView",
    components: {
      FileUploader,
      'fc-tile': tile,
      'server-browser': serverBrowser
    },
    data: function () {
      return {
        loading: false,
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Keyspace', value: 'keyspace', align: 'end'},
          {text: 'Time', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ],
        dictionaries: [],
        browser: false,
        uploader: false,
        sortUploaded: false
      }
    },
    mounted: function () {
      this.loadDictionaries()
    },
    methods: {
      fmt,
      loadDictionaries: function () {
        this.loading = true;
        this.axios.get(this.$serverAddr + '/dictionary', {}).then((response) => {
          this.dictionaries = response.data;
          this.loading = false
        })
      },
      deleteDictionary: function (item) {
        this.$root.$confirm('Delete', `This will remove ${item.name} from your dictionaries. Are you sure?`).then((confirm) => {
          this.loading = true;
          this.axios.delete(this.$serverAddr + '/dictionary/' + item.id).then((response) => {
            this.loadDictionaries()
          })
        })
      }
    }
  }
</script>

<style scoped>
  .noEvent {
    pointer-events: none;
    display: inline-block;
  }

  .dz-message {
    cursor: pointer;
    text-align: right;
  }

  .max500 {
    max-width: 770px;
  }

</style>
