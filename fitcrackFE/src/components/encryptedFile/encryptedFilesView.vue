<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max1000">
    <fc-tile
      title="Encrypted files"
      class="ma-2"
      :loading="encryptedFiles === null"
    >
      <v-data-table
        v-if="encryptedFiles !== null"
        :headers="headers"
        :items="encryptedFiles.items"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Files per page'}"
      >
        <template v-slot:item.name="{ item }">
          {{ item.name }}
        </template>
        <template v-slot:item.time="{ item }">
          {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.actions="{ item }">
           <v-tooltip top>
            <template v-slot:activator="{ on }">
              <a
                :href="$serverAddr + '/protectedFiles/' + item.id"
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
        </template>
      </v-data-table>
    </fc-tile>
  </v-container>
</template>

<script>
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "EncryptedFilesView",
    components: {
      'fc-tile': tile,
    },
    data: function () {
      return {
        uploadOptions: {
          url: this.$serverAddr + '/protectedFiles/',
          paramName: 'file'
        },
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Hash', value: 'hash', align: 'end'},
          {text: 'Time', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ],
        encryptedFiles: [],
        uploadProgress: 0
      }
    },
    mounted: function () {
      this.loadFiles()
    },
    methods: {
      loadFiles: function () {
        this.axios.get(this.$serverAddr + '/protectedFiles/', {}).then((response) => {
          this.encryptedFiles = response.data;
        })
      },
      uploadProgressChanged: function (progress, totalBytes, bytesSent) {
        this.uploadProgress = progress
      },

      uploadComplete: function (file, status, xhr) {
        if (status === 'error') {
          this.axios.interceptors.response.handlers["0"].rejected(JSON.parse(xhr.response))
        } else {
          this.loadFiles()
        }
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

  .max1000 {
    max-width: 1000px;
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
