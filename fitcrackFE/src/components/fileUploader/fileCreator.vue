<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-dialog
    v-model="open"
    max-width="600"
  >
    <v-card>
      <v-card-title class="headline">
        {{ title }}
      </v-card-title>
      <v-tabs
        v-model="active"
      >
        <v-tab ripple>
          Upload file
        </v-tab>
        <v-tab ripple>
          Make from dictionary
        </v-tab>
        <v-tab-item>
          <v-card text>
            <file-uploader
              :url="uploadUrl"
              @uploadComplete="fileUploaded"
            />
          </v-card>
        </v-tab-item>
        <v-tab-item>
          <div
            v-if="working"
            class="text-center"
          >
            <v-progress-circular
              class="ma-5"
              size="50"
              :width="3"
              indeterminate
              color="primary"
            />
          </div>
          <v-card
            v-else
            flat
          >
            <dict-selector v-model="selectedDict" />
            <v-card-actions>
              <v-text-field
                v-model="newName"
                outlined
                dense
                hide-details
                label="Name"
                class="mr-2"
                @input="customName = true"
              />
              <v-btn
                color="primary"
                outlined
                :disabled="selectedDict.length == 0 || newName === ''"
                @click="dictionarySelected"
              >
                Make from dictionary
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-tab-item>
      </v-tabs>
    </v-card>
  </v-dialog>
</template>

<script>
import dictSelector from '@/components/selector/dictionarySelector.vue'
import fileUploader from '@/components/fileUploader/fileUploader.vue'

export default {
  components: {
    dictSelector,
    fileUploader
  },
  props: {
    value: Boolean,
    title: {
      type: String,
      default: 'Upload or make from dictionary'
    },
    uploadUrl: {
      type: String,
      default: '#'
    },
    working: Boolean
  },
  data () {
    return {
      active: null,
      // Dictionary selection
      newName: '',
      customName: false,
      selectedDict: []
    }
  },
  computed: {
    open: {
      get () {
        return this.value
      },
      set (newValue) {
        this.$emit('input', newValue)
      }
    }
  },
  watch: {
    selectedDict (val) {
      if (this.newName.length == 0) {
        this.customName = false
      }
      if (val.length > 0 && !this.customName) {
        this.newName = val[0].name
      }
    }
  },
  methods: {
    fileUploaded () {
      this.$emit('fileUploaded')
      this.open = false
    },
    dictionarySelected () {
      this.$emit('dictionarySelected', this.selectedDict[0].id, this.newName)
    }
  }
}
</script>