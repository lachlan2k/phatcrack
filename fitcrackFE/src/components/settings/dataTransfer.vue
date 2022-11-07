<template>
  <v-container>
    <v-expand-transition>
      <div
        v-if="listing.length == 0"
        class="text-center mt-8 pa-8"
      >
        <h1>Job packages for takeout.</h1>
        <box class="illust" />
        <div>Export your jobs to a file that you can take with you and import it to another server or use it as backup.</div>
      </div>
    </v-expand-transition>
    <v-row>
      <v-col>
        <v-card 
          class="text-start"
          flat
          min-width="300"
        >
          <v-card-title>Export jobs to a file</v-card-title>
          <v-card-subtitle>Perform selective export</v-card-subtitle>
          <v-card-text v-if="count == 0">
            <div class="mb-2 subtitle-1">
              How to select jobs:
            </div>
            <ol>
              <li>
                List <router-link to="/jobs">
                  <i class="mdi mdi-briefcase" /> <em>All Jobs</em>
                </router-link> or a bin
              </li>
              <li>Check <i class="mdi mdi-checkbox-marked-outline" /> the jobs you want to export</li>
              <li>Click <i class="mdi mdi-package-variant" /> <em>Add to Export</em> in the selection toolbar</li>
            </ol>
            <v-alert 
              type="info"
              class="mt-4"
              text
            >
              No specific jobs selected. <strong>All jobs</strong> will be exported.
            </v-alert>
          </v-card-text>
          <v-card-text v-else>
            <v-alert 
              type="warning"
              class="mt-4"
              text
            >
              Note that only hashes and masks are exported,
              but all asset files will need to be uploaded before importing.
            </v-alert>
          </v-card-text>
          <v-card-actions>
            <v-btn
              v-if="count > 0"
              text
              @click="reset"
            >
              <v-icon left>
                mdi-delete-empty
              </v-icon>
              Clear selection
            </v-btn>
            <v-spacer />
            <v-btn
              text
              :href="`${$serverAddr}/serverInfo/transfer${count > 0 ? '?jobs=' + job_ids.join(',') : ''}`"
            >
              Export {{ count > 0 ? `${count} Job${count > 1 ? 's' : ''}` : 'Everything' }}
              <v-icon right>
                mdi-database-export
              </v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      <v-col>
        <v-card 
          class="text-start"
          flat
          min-width="300"
        >
          <v-card-title>Import jobs to this system</v-card-title>
          <v-card-subtitle>Upload a package file</v-card-subtitle>
          <v-card-text>
            <v-file-input
              v-model="packageFile"
              filled
              truncate-length="100"
              label="Select a Fitcrack package"
              hint="Fitcrack packages are binary files with the FCP extension"
              persistent-hint
            />
            <div v-if="packageContents">
              Server <strong>{{ packageContents.packager }}</strong> packaged this export
              on {{ this.$moment.unix(packageContents.timestamp).format('DD.MM.YYYY HH:mm') }}
            </div>
            <v-alert
              v-if="!validFile"
              type="error"
            >
              This does not appear to be a valid package.
            </v-alert>
          </v-card-text>
          <v-card-actions>
            <v-btn
              v-if="dependencyList"
              text
              @click="showDeps = !showDeps"
            >
              Dependencies ({{ dependencyList.length }})
              <v-icon>
                {{ showDeps ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
              </v-icon>
            </v-btn>
            <v-spacer />
            <v-btn
              :text="!canImport"
              :disabled="!canImport"
              :color="canImport ? 'success' : ''"
              :class="{'mr-2': canImport}"
              @click="upload"
            >
              Start Import
              <v-icon right>
                mdi-database-import
              </v-icon>
            </v-btn>
          </v-card-actions>
          <v-card-text>
            <template v-if="dependencyList">
              <v-alert
                text
                :type="depCheckStatusType"
              >
                {{ depCheckStatusText }}
              </v-alert>
              <v-progress-linear
                v-if="!depsValidated"
                indeterminate
                color="info"
              />
              <v-list
                v-if="showDeps"
                dense
              >
                <v-list-item
                  v-for="{dep, missing, type, name} in dependencyList"
                  :key="dep"
                >
                  <v-list-item-action>
                    <v-icon
                      :color="missing ? 'error' : 'success'"
                    >
                      {{ missing ? 'mdi-close-circle' : 'mdi-check-circle' }}
                    </v-icon>
                  </v-list-item-action>
                  <v-list-item-content>
                    <v-list-item-title>
                      {{ name }}
                    </v-list-item-title>
                    <v-list-item-subtitle class="text-capitalize">
                      {{ type }}
                      <span 
                        v-if="missing"
                        class="error--text"
                      >
                        MISSING
                      </span>
                    </v-list-item-subtitle>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </template>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-expand-transition>
      <v-row v-if="listing.length > 0">
        <v-col>
          <v-card
            flat
          >
            <v-card-title>
              Jobs {{ packageFile ? 'in package' : 'selected for export' }}
            </v-card-title>
            <v-card-text>
              <v-data-table
                :items="listing"
                :headers="headers"
              >
                <template v-slot:item.name="{ item }">
                  <strong>
                    <v-component
                      :is="packageFile ? 'div' : 'router-link'"
                      :to="{ name: 'jobDetail', params: { id: item.id }}"
                    >
                      {{ item.name }}
                    </v-component>
                  </strong>
                </template>
                <template v-slot:item.attack="{ item }">
                  <span class="text-capitalize">
                    {{ item.attack }}
                    <v-icon
                      right
                    >
                      {{ attackIcon(item.attack) }}
                    </v-icon>
                  </span>
                </template>
                <template v-slot:item.actions="{ item }">
                  <v-btn
                    v-if="!packageFile"
                    text
                    @click="remove(item.id)"
                  >
                    Remove
                    <v-icon right>
                      mdi-minus
                    </v-icon>
                  </v-btn>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-expand-transition>
  </v-container>
</template>

<script>
import { mapState, mapGetters, mapMutations } from 'vuex'
import { attackIcon } from '@/assets/scripts/iconMaps'
import { decode as msgpackDecode } from 'msgpack-lite'

import box from '@/components/details/boxIllustration.vue'

export default {
  components: {
    box
  },
  data () {
    return {
      packageFile: null,
      validFile: true,
      packageContents: null,
      missingDeps: [],
      depsValidated: false,
      headers: [
        {
          text: 'Name',
          align: 'start',
          value: 'name'
        },
        {text: 'Attack type', value: 'attack', align: 'end'},
        {text: 'Actions', value: 'actions', sortable: false, align: 'end'}
      ],
      showDeps: false,
      progress: 0
    }
  },
  computed: {
    ...mapState('transfer', ['jobs']),
    ...mapGetters('transfer', ['job_ids']),
    count () {
      return this.jobs.length
    },
    listing () {
      return this.packageContents ? this.packageContents.jobs : this.jobs
    },
    dependencyList () {
      if (
        !this.packageContents || 
        (this.packageContents && !this.packageContents.timestamp) ||
        !this.depsValidated
      ) {
        return null
      }
      const md = this.missingDeps
      return this.packageContents.deps.map(dep => {
        const [type, name] = dep.split('/')
        const missing = md.includes(dep)
        return {dep, type, name, missing}
      })
    },
    depCheckStatusType () {
      if (!this.depsValidated) {
        return 'info'
      } else if (this.depsValidated && this.missingDeps.length > 0) {
        return 'error'
      } else {
        return 'success'
      }
    },
    depCheckStatusText () {
      if (!this.depsValidated) {
        return 'Please wait, checking dependencies...'
      } else if (this.depsValidated && this.missingDeps.length > 0) {
        return 'Some dependencies are missing! Please add them first.'
      } else {
        return 'All dependencies in place. Ready to import!'
      }
    },
    canImport () {
      return this.packageFile && this.depsValidated && this.missingDeps.length == 0
    }
  },
  watch: {
    async packageFile (file) {
      if (file) {
        const buffer = await file.arrayBuffer()
        const uarr = new Uint8Array(buffer)
        this.packageContents = msgpackDecode(uarr)
        if (this.packageContents.timestamp) {
          this.validFile = true
          this.checkDeps(this.packageContents.deps)
        } else {
          this.validFile = false
          this.packageContents = null
        }
      } else {
        this.validFile = true
        this.packageContents = null
        this.depsValidated = false
        this.missingDeps = []
      }
    }
  },
  methods: {
    ...mapMutations('transfer', ['remove', 'reset']),
    attackIcon,
    async checkDeps (deps) {
      try {
        this.missingDeps = await this.axios.post(`${this.$serverAddr}/serverInfo/transfer/validate`, {deps}).then(r => r.data.missing)
        this.depsValidated = true
        if (this.missingDeps.length > 0) {
          this.showDeps = true
        }
      } catch (err) {
        console.error(err)
        this.packageFile = null
      }
    },
    async upload () {
      const config = {
        withCredentials: true,
        headers: {
          'Content-type': 'multipart/form-data'
        },
        onUploadProgress: this.uploadChange
      }
      const files = new FormData()
      files.append('file', this.packageFile)
      try {
        await this.axios.post(`${this.$serverAddr}/serverInfo/transfer`, files, config)
        this.$router.push({ name: 'jobs' })
      } catch (err) {
        console.error(err)
      }
    },
    uploadChange(progressEvent) {
      this.progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
    }
  }
}
</script>

<style scoped>
  .detail-container {
    margin: 0 auto;
    max-width: 800px;
  }

  .illust {
    width: 80%;
    max-width: 300px;
    margin: 2em auto;
  }
</style>