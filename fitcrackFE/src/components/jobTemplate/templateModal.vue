<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-dialog
    v-model="dialog"
    width="800"
  >
    <template #activator="{ on }">
      <v-btn
        :disabled="!validTemplate"
        large
        color="primary"
        v-on="on"
        @click="nameMe"
      >
        <v-icon left>
          mdi-file-plus
        </v-icon>Manage Template
      </v-btn>
    </template>

    <v-card>
      <v-card-title
        primary-title
      >
        Job template
      </v-card-title>

      <v-card-text>
        <!-- <json-viewer :value="data" ></json-viewer> -->
        <div class="flex mb-3">
          <v-icon left>
            info
          </v-icon>
          <div>
            The choices you made in the job configuration will be saved as a template.<br>
            You can pre-fill any new job configuration with templates you saved.
            <div v-show="name !== '' && name === inheritedName">
              Keep the name to update an existing template.
            </div>
          </div>
        </div>
        <v-text-field
          id="templateName"
          v-model="name"
          name="templateName"
          outlined
          dense
          hide-details
          label="Template name"
        />
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn
          text
          @click="dialog = !dialog"
        >
          Back to editing
        </v-btn>
        <v-btn
          color="primary"
          text
          @click="submit"
          :disabled="name === ''"
        >
          {{ name !== '' && name === inheritedName ? 'Update' : 'Save' }} template
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
  import {mapGetters} from 'vuex'

  export default {
    name: "TemplateModal",
    props: {
      inheritedName: {
        type: String,
        default: ''
      }
    },
    data () {
      return {
        dialog: false,
        name: ''
      }
    },
    computed: mapGetters('jobForm', ['validTemplate']),
    methods: {
      submit() {
        const template = this.$store.getters['jobForm/makeTemplate'](this.name)
        console.log(template)
        this.loading = true
        this.axios.put(this.$serverAddr + '/template', template).then((response) => {
          this.$emit('templatesUpdated')
        }).catch((error) => {
          this.loading = false
        })
        this.dialog = false
      },
      nameMe () {
        this.name = this.inheritedName
      }
    }
  }
</script>

<style scoped>
.flex {
  display: flex;
}
</style>
