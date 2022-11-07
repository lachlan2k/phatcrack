<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="pt-0 mt-0">
    <v-alert
      v-if="overSizeLimit"
      type="error"
      tile
    >
      Maximum size of uploaded file is 2 GB.
    </v-alert>
    <form
      enctype="multipart/form-data"
      class="form pa-2"
    >
      <v-file-input
        outlined
        chips
        show-size
        class="mr-2"
        :label="label"
        name="file"
        :multiple="multiple"
        @change="fileChange"
      />
      <v-btn
        v-if="!noUpload"
        color="primary"
        outlined
        :disabled="files === null || overSizeLimit"
        @click="upload()"
      >
        Upload
        <v-icon right>
          mdi-upload
        </v-icon>
      </v-btn>
    </form>
    <!--
    <ul v-if="multiple">
      <li v-for="file in selectedFiles">
        <strong>{{ file.name }}</strong> - {{ file.size }} bytes, last modified: {{ file.modified }}
      </li>
    </ul>
    -->
    <v-progress-linear
      v-model="progress"
      :query="true"
      :active="showProgress"
      color="primary"
    />
  </div>
</template>

<script>
  export default {
    name: "FileUploader",
    props:{
      multiple: Boolean,
      url: {
        type: String,
        default: window.$serverAddr
      },
      noUpload: Boolean,
      label: {
        type: String,
        default: 'Select files'
      },
      args: {
        type: Object,
        default: () => {}
      }
    },
    data: function () {
      return {
        selectedFiles: [],
        files: null,
        progress: 0,
        showProgress: false,
        fileUploaded: false
      }
    },
    computed: {
      overSizeLimit () {
        return this.selectedFiles.some(file => file.size > 2000000000)
      }
    },
    methods:{
      uploadChange(progressEvent) {
        this.progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
        console.log(this.progress)
      },
      fileChange(fileOrList) {
        this.files = new FormData()
        this.selectedFiles = []
        if (this.multiple) {
          for (var i = 0; i < fileOrList.length; i++) {
            this.selectedFiles.push({
              'name': fileOrList[i].name,
              'type': fileOrList[i].type || 'n/a',
              'size': fileOrList[i].size
            })
            this.files.append("file", fileOrList[i], fileOrList[i].name);
          }
        } else {
          this.selectedFiles.push({
              'name': fileOrList.name,
              'type': fileOrList.type || 'n/a',
              'size': fileOrList.size
            })
            this.files.append("file", fileOrList, fileOrList.name);
        }
        if (this.args) {
          for (let [name, val] of Object.entries(this.args)) {
            this.files.append(name, val)
          }
        }
        this.fileUploaded = false
        this.progress = 0
        this.$emit('filesChanged', fileOrList)
      },
      upload() {
        this.showProgress = true
        var config = {
          withCredentials: true,
          headers: {
            'Content-type': 'multipart/form-data'
          },
          onUploadProgress: this.uploadChange
        }

        this.axios.post(this.url, this.files, config
        ).then(response => {
          this.fileUploaded = true
          this.$emit('uploadComplete', response.data)
          this.files = null
        }).catch(error => {
          console.log(error)
        })
      },
    }
  }
</script>

<style scoped>
.uploadInput {
  max-width: 250px;
}

.form {
  max-width: 100%;
  text-align: center;
}

.uploadButton {
  margin-top: 15px;
}

.spacer {
  width: 10px;
}
</style>
