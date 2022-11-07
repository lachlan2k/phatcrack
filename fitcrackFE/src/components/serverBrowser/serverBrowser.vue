<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <fc-tile
    :title="data ? 'Import from ' + data.path || 'Server browser' : 'Connecting...'"
    icon="mdi-server-network"
    class="mx-auto dictContentContainer"
  >
    <v-list
      v-if="data!==null"
      tile
      nav
      class="filesCont"
    >
      <v-list-item
        v-show="data.parent !== ''"
        @click="loadDirectory(data.parent)"
      >
        <v-list-item-action>
          <v-icon>keyboard_backspace</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>..</v-list-item-title>
          <v-list-item-subtitle>{{ data.parent }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        v-show="data.folders.length == 0 && data.files.length == 0"
      >
        <v-list-item-action>
          <v-icon>mdi-folder-alert</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>Empty folder</v-list-item-title>
          <v-list-item-subtitle>Go up or add files on server</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        v-for="item in data.folders"
        :key="item.name"
        @click="loadDirectory(item.path)"
      >
        <v-list-item-action>
          <v-icon>folder</v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>{{ item.name }}</v-list-item-title>
          <v-list-item-subtitle>{{ item.path }}</v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
      <v-list-item
        v-for="item in data.files"
        :key="item.name"
        :class="{selected: (selectedFiles[item.path] !== undefined && selectedFiles[item.path] !== null)}"
        @click="fileSelected(item.name, item.path)"
      >
        <v-list-item-action>
          <v-icon class="whiteAfterSelected">
            insert_drive_file
          </v-icon>
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title class="whiteAfterSelected">
            {{ item.name }}
          </v-list-item-title>
          <v-list-item-subtitle class="whiteAfterSelected">
            {{ item.path }}
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>
    <div v-if="Object.keys(selectedFiles).length > 0">
      <v-divider />
      <v-card-title>Selected files:</v-card-title>
      <v-card-text
        v-for="item in selectedFiles"
        :key="item.name"
        class="my-0 py-0"
      >
        {{ item.path }}
      </v-card-text>
      <v-btn
        class="selectBtn d-block mx-auto"
        color="primary"
        @click="upload"
      >
        Add
      </v-btn>
      <v-divider />
      <v-progress-linear
        v-model="progress"
        :query="true"
        :active="adding"
        color="primary"
      />
      <v-card-text v-if="adding"
        class="my-0 py-0"
      >
      It may take a few minutes to complete this action. Please check back later.
      </v-card-text>
    </div>
  </fc-tile>
</template>

<script>
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "ServerBrowser",
    components: {
      'fc-tile': tile
    },
    props: {
      sort: Boolean,
      hex_wordlist: Boolean
    },
    data: function () {
      return {
        data: null,
        selectedFiles:{},
        adding: false,
        progress: 0
      }
    },
    mounted: function () {
      this.loadDirectory('/')
    },
    methods: {
      uploadChange(progressEvent) {
        this.progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      },
      loadDirectory: function (path) {
        // this.data = null
        this.axios.get(this.$serverAddr + '/directory', {
          params: {
            path: path
          }
        })
        .then((response) => {
          this.data = response.data
        })
      },
      fileSelected: function (name, path) {
        if (this.selectedFiles[path] === undefined || this.selectedFiles[path] === null) {
          this.$set(this.selectedFiles, path, {
            name: name,
            path: path
          })
        } else {
          this.$delete(this.selectedFiles,path)
        }
      },
      upload: function () {
        this.progress = 0
        this.adding = true
        var files = Object.keys(this.selectedFiles).map(function(key) {
          return {
            "name": this.selectedFiles[key]["name"],
            "path": this.selectedFiles[key]["path"]
          };
        }.bind(this));

        var config = {
          onUploadProgress: this.uploadChange
        }
        this.axios.post(this.$serverAddr + '/dictionary/fromFile', {
          files,
          sort: this.sort,
          hex_wordlist: this.hex_wordlist
        }, config).then((response) => {
          this.adding = false
          this.$emit('filesuploaded', true)
        }).catch((error) => {
          this.adding = false
          this.$emit('filesuploaded', false)
        })
      }
    }
  }
</script>

<style scoped>
  .selected {
    background: rgba(37, 157, 173, 0.85) !important;
  }

  .selected .whiteAfterSelected {
    color: white !important;
  }

  .dictContentContainer {
    max-width: 500px;
  }


  .filesCont {
    height: max-content;
    max-height: 500px;
    overflow: auto;
  }
</style>
