<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max500">
    <fc-tile
      title="Charsets"
      class="ma-2"
      :icon="$route.meta.icon"
    >
      <v-alert
        tile
        text
        type="warning"
        class="mb-0"
      >
        Charsets must have a .txt, .hcchr or .charset extension.
      </v-alert>
      <v-data-table
        :headers="headers"
        :items="charsets.items"
        :loading="loading"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Charsets per page'}"
      >
        <template v-slot:item.name="{ item }">
          <router-link :to="`charsets/${item.id}`">
            {{ item.name }}
          </router-link>
        </template>
        <template v-slot:item.keyspace="{ item }">
          {{ fmt(item.keyspace) }}
        </template>
        <template v-slot:item.time="{ item }">
          {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <a
                :href="$serverAddr + '/charset/' + item.id + '/download'"
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
                @click="deleteCharset(item)"
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
      <v-divider />
      <file-uploader
        :url="this.$serverAddr + '/charset/add'"
        @uploadComplete="loadCharsets"
      />
    </fc-tile>
  </v-container>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  import FileUploader from "@/components/fileUploader/fileUploader.vue";
  export default {
    name: "CharsetView",
    components: {
      FileUploader,
      'fc-tile': tile,
    },
    data: function () {
      return {
        loading: false,
        charsets: [],
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Keyspace', value: 'keyspace', align: 'end'},
          {text: 'Added', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ]
      }
    },
    mounted: function () {
      this.loadCharsets()
    },
    methods: {
      fmt,
      loadCharsets: function () {
        this.loading = true;
        this.axios.get(this.$serverAddr + '/charset', {}).then((response) => {
          this.charsets = response.data;
          this.loading = false
        })
      },
      deleteCharset: function (item) {
        this.$root.$confirm('Delete', `This will remove ${item.name} from your charsets. Are you sure?`).then((confirm) => {
          this.loading = true;
          this.axios.delete(this.$serverAddr + '/charset/' + item.id).then((response) => {
            this.loadCharsets()
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
