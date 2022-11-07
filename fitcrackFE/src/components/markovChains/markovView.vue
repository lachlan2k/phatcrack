<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max500">
    <fc-tile
      title="Markov chains"
      class="ma-2"
      :icon="$route.meta.icon"
    >
      <v-alert
        tile
        text
        type="warning"
        class="mb-0"
      >
        Markov files must have a .hcstat2 extension.
      </v-alert>
      <v-data-table
        :headers="headers"
        :items="markovChains.items"
        :loading="loading"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Hcstats per page'}"
      >
        <template v-slot:item.time="{ item }">
          {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <a
                :href="$serverAddr + '/markovChains/' + item.id"
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
                @click="deleteMarkov(item)"
                v-on="on"
              >
                <v-icon color="error">
                  mdi-delete-outline
                </v-icon>
              </v-btn>
            </template>
            <span>Delete</span>
          </v-tooltip>
        </template>
      </v-data-table>
      <div class="text-right pa-3">
        <v-btn
          class="d-inline-block"
          color="primary"
          text
          outlined
          @click.native.stop="dialog = true"
        >
          Add new
        </v-btn>
      </div>
    </fc-tile>

    <file-creator
      v-model="dialog"
      title="Add new HcStat"
      :upload-url="this.$serverAddr + '/markovChains/add'"
      :working="working"
      @fileUploaded="loadMarkovChains"
      @dictionarySelected="makeHcStatFromDictionary"
    />
  </v-container>
</template>

<script>
  import tile from '@/components/tile/fc_tile.vue'
  import fileCreator from "@/components/fileUploader/fileCreator.vue";
  export default {
    name: "MarkovView",
    components: {
      fileCreator,
      'fc-tile': tile,
    },
    data: function () {
      return {
        working: false,
        dialog: false,
        loading: false,
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Added', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ],
        markovChains: [],
      }
    },
    mounted: function () {
      this.loadMarkovChains()
    },
    methods: {
      loadMarkovChains: function () {
        this.dialog= false
        this.loading = true;
        this.axios.get(this.$serverAddr + '/markovChains', {}).then((response) => {
          this.markovChains = response.data;
          this.loading = false
        })
      },
      makeHcStatFromDictionary: function (id, name) {
        this.working = true
        this.axios.post(this.$serverAddr + '/markovChains/makeFromDictionary', {
          "dictionary_id": id,
          "name": name
        }).then((response) => {
          this.working = false
          this.dialog = false
          this.loadMarkovChains()
        }).catch((error) => {
          this.working = false
        })
      },
      deleteMarkov: function (item) {
        this.$root.$confirm('Delete', `This will remove ${item.name} from your Markov chains. Are you sure?`).then((confirm) => {
          this.loading = true;
          this.axios.delete(this.$serverAddr + '/markovChains/' + item.id).then((response) => {
            this.loadMarkovChains()
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
    text-align: center;
  }

  .max500 {
    max-width: 600px;
  }

</style>


<style>
  .selectedDict {
    background: rgba(37, 157, 173, 0.85) !important;
    color: white;
  }

  .selectedDict a {
    color: white;
  }

  .clickable {
    cursor: pointer;
  }
</style>
