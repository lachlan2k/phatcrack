<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="containerAddJob height100 mx-auto">
    <v-scale-transition>
      <v-sheet
        v-show="keyspaceKnown && keyspace || estimatedTime"
        class="infobar"
        color="info"
        dark
        elevation="8"
      >
        <div class="d-flex">
          <div
            v-show="keyspace"
            class="text-end"
          >
            <div class="overline">
              <v-icon
                x-small
                class="mr-1"
              >
                mdi-key
              </v-icon>
              <span>
                Keyspace
              </span>
            </div>
            <div>{{ numberFormat(keyspace) }}</div>
          </div>
          <div
            v-show="estimatedTime"
            class="ml-4 text-end"
          >
            <div class="overline">
              <v-icon
                x-small
                class="mr-1"
              >
                mdi-timelapse
              </v-icon>
              <span>
                Est. cracking time
              </span>
            </div>
            <div>{{ estimatedTime }}</div>
          </div>
        </div>
      </v-sheet>
    </v-scale-transition>
    <div
      v-show="loading"
      class="loadingProgressCont"
    >
      <v-progress-circular
        indeterminate
        color="primary"
        size="65"
        class="loadingProgress"
        :width="4"
      />
    </div>

    <v-row>
      <v-col>
        <v-text-field
          v-model="name"
          outlined
          autofocus
          required
          label="Name"
          hint="Give this job a descriptive name"
          persistent-hint
        />
      </v-col>
      <v-col>
        <v-autocomplete
          :items="templates"
          item-text="name"
          item-value="id"
          :value="selectedTemplate"
          auto-select-first
          outlined
          label="Template"
          hint="Prefill the form with a saved template (use empty to quickly reset)"
          persistent-hint
          @input="fetchAndApplyTemplate"
        />
      </v-col>
    </v-row>

    <div>
      <v-col>
        <v-stepper
          id="job-stepper"
          v-model="step"
          vertical
          non-linear
          class="mb-4"
        >
          <v-stepper-step
            id="job-step-1"
            editable
            step="1"
          >
            Input settings
          </v-stepper-step>
          <v-stepper-content step="1">
            <v-container>
              <v-row class="mb-4">
                <v-btn-toggle
                  v-model="inputMethod"
                  mandatory
                  color="primary"
                  class="mr-2"
                >
                  <v-btn id="job-input-mode-manual" value="multipleHashes">
                    Manual entry
                  </v-btn>
                  <v-btn id="job-input-mode-hashlist" value="hashFile">
                    From hash file
                  </v-btn>
                  <v-btn id="job-input-mode-extract" value="extractFromFile">
                    Extract from file
                  </v-btn>
                </v-btn-toggle>
                <v-autocomplete
                  id="hash-type-select"
                  v-model="hashType"
                  editable
                  validate-on-blur
                  clearable
                  label="Select hash type"
                  :items="hashTypes"
                  item-text="name"
                  :filter="hashTypeFilter"
                  return-object
                  required
                  hide-details
                  single-line
                  flat
                  solo-inverted
                  no-data-text="No matching hash type"
                  @change="validateHashes(null)"
                >
                  <template #item="{ item }">
                    <v-list-item-content>
                      <v-list-item-title><b>{{ item.code }}</b> - {{ item.name }}</v-list-item-title>
                    </v-list-item-content>
                  </template>
                </v-autocomplete>
              </v-row>
              <v-row>
                <v-col
                  v-if="inputMethod === 'extractFromFile'"
                  cols="12"
                >
                  <v-alert
                    type="info"
                    text
                  >
                    Supported formats:
                    <v-tooltip top>
                      <template v-slot:activator="{ on }">
                        <span v-on="on"><a href="#" class="filetype-link">MS_OFFICE</a>,</span>
                      </template>
                      <span>Hash types: 9400, 9500, 9600, 9700, 9800</span>
                    </v-tooltip>
                    <v-tooltip top>
                      <template v-slot:activator="{ on }">
                        <span v-on="on"><a href="#" class="filetype-link">PDF</a>,</span>
                      </template>
                      <span>Hash types: 10400, 10500, 10600, 10700</span>
                    </v-tooltip>
                    <v-tooltip top>
                      <template v-slot:activator="{ on }">
                        <span v-on="on"><a href="#" class="filetype-link">7Z</a>,</span>
                      </template>
                      <span>Hash type: 11600</span>
                    </v-tooltip>
                    <v-tooltip top>
                      <template v-slot:activator="{ on }">
                        <span v-on="on"><a href="#" class="filetype-link">RAR</a> and </span>
                      </template>
                      <span>Hash types: 12500, 13000</span>
                    </v-tooltip>
                    <v-tooltip top>
                      <template v-slot:activator="{ on }">
                        <span v-on="on"><a href="#" class="filetype-link">ZIP</a>.</span>
                      </template>
                      <span>Hash types: 13600</span>
                    </v-tooltip>
                  </v-alert>
                  <file-uploader
                    ref="encryptedFileUploader"
                    :url="this.$serverAddr + '/protectedFiles/add'"
                    @uploadComplete="uploadComplete"
                  />
                </v-col>
                <v-col
                  v-if="inputMethod === 'hashFile'"
                  cols="12"
                >
                  <v-alert
                    type="info"
                    text
                  >
                    Input to hashcat. Select a binary hash (WPA/WPA2), or plain text hashlist.
                  </v-alert>
                  <file-uploader
                    ref="hashFileUploader"
                    :multiple="false"
                    no-upload
                    label="Select a file to read"
                    @filesChanged="hashFileSelected"
                  />
                </v-col>
                <v-col cols="12">
                  <v-alert
                    :value="gotBinaryHash"
                    text
                    type="info"
                    color="primary"
                    class="mt-0 mb-0"
                  >
                    You can select only one binary hash.
                  </v-alert>
                  <fc-textarea
                    id="hashes-input"
                    v-if="inputMethod !== null"
                    ref="textarea"
                    v-model="hashList"
                    :class="{'hasherror': hashListError}"
                    class="textarea"
                    max-height="500"
                    :readonly="!(inputMethod === 'multipleHashes' && !gotBinaryHash) "
                    :can-remove-line="true "
                    @blur="validateHashes"
                    @focus="unvalidateHashes"
                  >
                    <div
                      slot="after"
                      class="hashCeckContainer pl-1 pt-2"
                    >
                      <div
                        v-for="hashObj in validatedHashes"
                        :key="hashObj.id"
                      >
                        <v-icon
                          v-if="hashObj.result === 'OK'"
                          small
                          color="success"
                        >
                          check_circle_outlined
                        </v-icon>
                        <v-tooltip
                          v-else
                          left
                        >
                          <template v-slot:activator="{ on }">
                            <v-icon
                              small
                              color="error"
                              class="clickable"
                              v-on="on"
                            >
                              error_circle_outlined
                            </v-icon>
                          </template>
                          <span>{{ hashObj.result }}</span>
                        </v-tooltip>
                        <v-tooltip
                          v-if="hashObj.isInCache"
                          left
                        >
                          <template v-slot:activator="{ on }">
                            <v-icon
                              small
                              color="warning"
                              class="clickable"
                              v-on="on"
                            >
                              error_circle_outlined
                            </v-icon>
                          </template>
                          <span>hash already in hashcache</span>
                        </v-tooltip>
                      </div>
                    </div>
                  </fc-textarea>
                </v-col>
              </v-row>
              <v-row>
                <v-btn
                  v-if="dev"
                  text
                  color="success"
                  @click="getRandomHash"
                >
                  <v-icon left>
                    mdi-auto-fix
                  </v-icon>
                  Random SHA1
                </v-btn>
                <v-checkbox
                  v-show="invalidHashes.length > 0"
                  v-model="ignoreHashes"
                  label="Ignore invalid hashes"
                  hide-details
                  :height="15"
                  color="error"
                  class="ml-2 mt-1"
                />
                <v-spacer />
                <v-btn
                  v-show="hashList !== ''"
                  color="error"
                  class="mr-2"
                  text
                  @click="clearInput"
                >
                  Reset
                </v-btn>
                <v-btn
                  color="primary"
                  @click="step = 2"
                >
                  Next
                </v-btn>
              </v-row>
            </v-container>
          </v-stepper-content>

          <v-stepper-step
            id="job-step-2"
            editable
            step="2"
          >
            Attack settings
          </v-stepper-step>
          <v-stepper-content step="2">
            <v-container>
              <v-card-title>Attack mode</v-card-title>
              <v-item-group
                v-model="attackSettingsTab"
                class="d-flex flex-wrap justify-space-between"
              >
                <v-item
                  v-for="(type, i) in attacks"
                  :key="i"
                  :value="type.handler"
                  #default="{ active, toggle }"
                >
                  <v-btn
                    :id="`attack-mode-${type.handler}`"
                    class="mode-btn"
                    text
                    small
                    :color="active ? 'primary' : ''"
                    @click="toggle"
                  >
                    <div class="my-3">
                      <v-icon class="mb-2">
                        {{ attackIcon(type.handler) }}
                      </v-icon>
                      <div>{{ type.name }}</div>
                    </div>
                  </v-btn>
                </v-item>
              </v-item-group>

              <v-component :id="attackSettingsTab" :is="attackSettingsTab" :keyspace="keyspace" />

              <v-row>
                <v-spacer />
                <v-btn
                  class="mr-6 mt-4"
                  color="primary"
                  @click="step = 3"
                >
                  Next
                </v-btn>
              </v-row>
            </v-container>
          </v-stepper-content>
          <v-stepper-step
            id="job-step-3"
            editable
            step="3"
          >
            Host assignment
          </v-stepper-step>
          <v-stepper-content step="3">
            <v-subheader>Select which hosts to distribute workunits to</v-subheader>
            <div class="scroller">
              <host-selector
                v-model="hosts"
                select-all
                auto-refresh
              />
            </div>
            <v-subheader><b>({{ hosts.length }} selected)</b></v-subheader>
            <v-row>
              <v-spacer />
              <v-btn
                class="mr-6 mt-4"
                color="primary"
                @click="step = 4"
              >
                Next
              </v-btn>
            </v-row>
          </v-stepper-content>
          <v-stepper-step
            id="job-step-4"
            editable
            step="4"
          >
            Additional settings
          </v-stepper-step>
          <v-stepper-content step="4">
            <v-container>
              <v-row>
                <v-textarea
                  id="comment-input"
                  v-model="comment"
                  label="Comment"
                  outlined
                  auto-grow
                  :rows="1"
                />
              </v-row>
              <v-row>
                <v-col>
                  <div class="title mb-2">Planned start</div>
                  <dt-picker
                    v-model="startDate"
                    outlined
                    hide-details
                    :disabled="startNow"
                    single-line
                    label=""
                  />
                  <v-checkbox
                    v-model="startNow"
                    label="Immediately"
                  />
                </v-col>
                <v-col>
                  <div class="title mb-2">
                    Desired time per workunit
                  </div>
                  <v-text-field
                    v-model="timeForJob"
                    outlined
                    hide-details
                    single-line
                    label="Time for job"
                    required
                    type="number"
                    suffix="seconds"
                    :min="60"
                  />
                </v-col>
                <v-col>
                  <div class="title mb-2">Planned end</div>
                  <dt-picker
                    v-model="endDate"
                    outlined
                    hide-details
                    :disabled="endNever"
                    single-line
                  />
                  <v-checkbox
                    v-model="endNever"
                    label="No time limit"
                  />
                </v-col>
              </v-row>
            </v-container>
          </v-stepper-content>
        </v-stepper>

        <v-row
          v-if="hosts.length == 0"
          justify="center"
        >
          <v-alert
            outlined
            type="warning"
          >
            This job has no hosts assigned to it and won't be able to run until edited!
          </v-alert>
        </v-row>

        <v-row
          v-if="timeForJob < 180"
          justify="center"
        >
          <v-alert
            outlined
            type="warning"
          >
            Setting extremely low time per workunit leads to high overhead. The cracking may take much longer than estimated. Consider choosing a higher value.
          </v-alert>
        </v-row>

        <v-row
          justify="center"
          class="mb-5 mt-2"
        >
          <template-modal
            :inherited-name="selectedTemplateName"
            @templatesUpdated="fetchTemplates"
          />
          <v-btn
            large
            color="primary"
            class="ml-2"
            @click="submit"
          >
            <v-icon left>
              {{ hosts.length > 0 ? 'mdi-check' : 'mdi-content-save' }}
            </v-icon>
            {{ hosts.length > 0 ? 'Create' : 'Save for later' }}
          </v-btn>
        </v-row>
      </v-col>
    </div>
  </div>
</template>

<script>
  import sha1 from 'sha1'
  import numberFormat from '@/assets/scripts/numberFormat'
  import { attackIcon } from '@/assets/scripts/iconMaps'

  import combinator from '@/components/job/attacks/combinator.vue'
  import mask from '@/components/job/attacks/mask.vue'
  import dictionary from '@/components/job/attacks/dictionary.vue'
  import hybridMaskWordlist from '@/components/job/attacks/hybridMaskWordlist.vue'
  import hybridWordlistMask from '@/components/job/attacks/hybridWordlistMask.vue'
  import pcfgAttack from '@/components/job/attacks/pcfg.vue'
  import princeAttack from '@/components/job/attacks/prince.vue'
  import FileUploader from "@/components/fileUploader/fileUploader.vue";
  import fcTextarea from '@/components/textarea/fc_textarea.vue'
  import hostSelector from '@/components/selector/hostSelector.vue'
  import templateModal from '@/components/jobTemplate/templateModal.vue'
  import dtPicker from '@/components/picker/datetime.vue'

  import {mapState, mapGetters, mapMutations} from 'vuex'
  import {mapTwoWayState} from 'spyfu-vuex-helpers'
  import {twoWayMap} from '@/store'

  import { attacks } from '@/store/job-form'

  export default {
    name: 'AddJob',
    components: {
      FileUploader,
      'combinator': combinator,
      'maskattack': mask,
      'dictionary': dictionary,
      'hybridMaskWordlist': hybridMaskWordlist,
      'hybridWordlistMask': hybridWordlistMask,
      'pcfgAttack': pcfgAttack,
      'princeAttack': princeAttack,
      'fc-textarea': fcTextarea,
      'host-selector': hostSelector,
      'template-modal': templateModal,
      dtPicker
    },
    data: function () {
      return {
        loading: false,
        hashTypes: [],
        showEstimatedTime: false,
        estimatedTime: null,
        keyspace: null,
        gotBinaryHash: false,
        hashListError: false,
        selectedTemplateName: '',
        attacks,
        templates: [
          {
            name: 'Empty',
            id: 0
          }
        ]
      }
    },
    computed: {
      ...mapState('jobForm', ['selectedTemplate']),
      ...mapTwoWayState('jobForm', twoWayMap([
        'step', 'attackSettingsTab', 'validatedHashes', 'name', 'inputMethod', 'hashList', 'hashType', 'ignoreHashes', 'startDate', 'endDate', 'template', 'comment', 'hosts', 'startNow', 'endNever', 'timeForJob'
      ])),
      ...mapGetters('jobForm', ['jobSettings', 'valid', 'validAttackSpecificSettings', 'keyspaceKnown']),
      templateItems () {
        return this.templates.map((t, i) => ({text: t.template, value: i}))
      },
      invalidHashes () {
        return this.validatedHashes.filter(h => h.result !== 'OK')
      },
      dev () {
        return localStorage.getItem('testmode') == 'true'
      }
    },
    watch: {
      jobSettings (val) {
        if (val.attack_settings != false && this.validAttackSpecificSettings) {
          var boincIds = []
          for (let i = 0; i < this.hosts.length; i++) {
            boincIds.push(this.hosts[i].id)
          }
          /* -1 means no hash entered */
          var hash_code = this.hashType == null ? -1 : this.hashType.code
          this.axios.get(this.$serverAddr + '/job/crackingTime', {
            params: {
              'hash_type_code': hash_code,
              'boinc_host_ids': boincIds.join(","),
              'attack_settings': val.attack_settings
            }
          }).then((response) => {
            if (response['data']) {
              this.estimatedTime = response.data.display_time
              this.keyspace = response.data.keyspace
            }
          })
        }
      }
    },
    mounted: function () {
      this.loadSettings()
      this.getHashTypes()
      this.startDate = this.$moment().format('YYYY-MM-DDTHH:mm')
      this.endDate = this.$moment().format('YYYY-MM-DDTHH:mm')
      if (this.hashList.length > 0) this.validateHashes()
      this.fetchTemplates()
      if (this.name === '') {
        this.generateJobName()
      }
    },
    methods: {
      ...mapMutations('jobForm', ['applyTemplate']),
      numberFormat,
      attackIcon,
      async loadSettings () {
        if (!this.timeForJob) {
          const settings = await this.axios.get(this.$serverAddr + '/settings').then(r => r.data)
          this.timeForJob = settings.default_seconds_per_workunit
        }
      },
      fetchTemplates () {
        this.axios.get(this.$serverAddr + '/template')
        .then((response) => {
          if (response.data && response.data.items) {
            this.templates = [
              { name: 'Empty', id: 0 },
              ...response.data.items
            ]
          }
        })
        .catch(console.error)
      },
      fetchAndApplyTemplate (id) {
        if (id == 0) {
          this.applyTemplate()
          this.selectedTemplateName = ''
          this.$store.commit('jobForm/selectedTemplateMut', 0)
          this.loadSettings()
          return
        }
        this.axios.get(this.$serverAddr + `/template/${id}`)
        .then((response) => {
          if (response.data && response.data.template) {
            const data = JSON.parse(response.data.template)
            this.applyTemplate(data)
            this.selectedTemplateName = data.template
            this.$store.commit('jobForm/selectedTemplateMut', id)
          }
        })
        .catch(console.error)
      },
      hashTypeFilter ({name, code}, query) {
        const q = query.toLowerCase()
        return name.toLowerCase().includes(q) || code.toLowerCase().includes(q)
      },
      subHashtypeChanged: function (key, val) {
        this.hashType.code = this.hashType.code.replace(key, val.code)
        this.validateHashes(null)
      },
      focusTextarea: function () {
        this.$refs.textarea.focus()
      },
      validateHashes: function (data = null) {
        if (data === null) {
          data = this.hashList
        }
        var hashesList = data.split('\n')
        if (data.startsWith("BASE64:")) {
          this.gotBinaryHash = true
        } else {
          this.gotBinaryHash = false
        }
        if (this.hashType === null || isNaN(this.hashType.code)) {
          return
        }
        if (data === '') {
          return
        }

        this.axios.post(this.$serverAddr + '/job/verifyHash', {
          'hashtype': this.hashType.code,
          'hashes': data
        }).then((response) => {
          this.hashListError = response.data.error
          this.validatedHashes = response.data.items
        })
      },
      unvalidateHashes: function (data) {
        this.validatedHashes = []
      },
      getHashTypes: function () {
        this.axios.get(this.$serverAddr + '/hashcat/hashTypes').then((response) => {
          this.hashTypes = response.data.hashtypes
        })
      },
      addHash: function (hash) {
        var parsedHashlist = this.hashList.split('\n')
        var lastHash = parsedHashlist[parsedHashlist.length-1]
        if (lastHash === '') {
          this.hashList = this.hashList + hash
        } else {
          this.hashList = this.hashList + '\n' + hash
        }
      },
      uploadComplete: function (data) {
        this.$success("Successfully extracted hash form file.")
        this.hashType = this.hashTypes.find(h => h.code == data['hash_type'])
        this.addHash(data['hash'])
        this.validateHashes(null)
      },
      isBinaryFile: function(content) {
        for (var i = 0; i < 24; i++) {
          var charCode = content.charCodeAt(i);
          if (charCode === 65533 || charCode <= 8) {
            return true
          }
        }
        return false
      },
      hashFileSelected: function (file) {
        var reader = new FileReader()
        reader.onloadend = function(evt) {
          if (evt.target.readyState == FileReader.DONE) {

            if (this.isBinaryFile(evt.target.result)) {
              // we got binary hash
              var binReader = new FileReader()
              binReader.onloadend = function(evt) {
                if (evt.target.readyState == FileReader.DONE) {
                  this.hashList = 'BASE64:' + evt.target.result.substr(evt.target.result.indexOf(',') + 1)
                  this.validateHashes(null)
                  this.gotBinaryHash = true
                }
              }.bind(this)
              binReader.readAsDataURL(file)

            } else {
              // we got hashlists
              var parsedHashlist = this.hashList.split('\n')
              var lastHash = parsedHashlist[parsedHashlist.length-1]
              if (lastHash === '') {
                this.hashList += evt.target.result
              } else {
                this.hashList += '\n' + evt.target.result
              }
              this.validateHashes(null)
            }
          }
        }.bind(this)
        reader.readAsText(file, 'utf-8')
      },
      clearInput () {
        this.hashType = null
        this.hashList = ''
        this.unvalidateHashes()
      },
      submit () {
        // TODO: maybe delete this condition
        if (this.name === '') {
          this.$error('Job name can not be empty.')
          return
        }

        if (this.inputMethod === 'encryptedFile' && !this.$refs.encryptedFileUploader.fileUploaded ) {
          this.$error('No file uploaded.')
          this.step = 1
          return
        }

        if (this.hashType === null) {
          this.$error('No hash type selected.')
          this.step = 1
          return
        }

        if (this.invalidHashes.length > 0 && !this.ignoreHashes) {
          this.$error('Some hashes are invalid.')
          this.step = 1
          return
        }

        if (this.validatedHashes.length == 0) {
          this.$error('List of validated hashes is empty.')
          this.step = 1
          return
        }

        if (!this.attackSettingsTab) {
          this.$error('No attack selected.')
          this.step = 2
          return
        }

        if (this.jobSettings.attack_settings === false) {
          this.$error('Error in attack settings.')
          this.step = 2
          return
        }

        if (this.keyspace > 1.8446744e+19 /* 2^64 */) {
          this.$error('Job keyspace is higher than maximal allowed value 2^64.')
          this.step = 2
          return
        }

        if (this.timeForJob < 10) {
          this.$error('Time per workunit must be higher or equal to 10 seconds.')
          this.step = 4
          return
        }

        // Check if all job settings are valid
        if (!this.valid) {
          // If all checks above passed and this one did not, it means that attack specific settings are incorrect.
          this.$error('Error in attack settings.')
          this.step = 2
          return
        }

        this.loading = true
        console.log(this.jobSettings.endNever)
        const finalStartTime = this.startNow ? 
          this.jobSettings['time_start'] : 
          this.$moment(this.jobSettings['time_start']).utc().toISOString(true).slice(0, 16)
        const finalEndTime = this.endNever ? 
          this.jobSettings['time_end'] : 
          this.$moment(this.jobSettings['time_end']).utc().toISOString(true).slice(0, 16)
        const finalSettings = {
          ...this.jobSettings,
          'time_start': finalStartTime,
          'time_end': finalEndTime
        }
        this.axios.post(this.$serverAddr + '/job', finalSettings).then((response) => {
          this.$router.push({name: 'jobDetail', params: {id: response.data.job_id}})
          this.applyTemplate() // Clear all
        }).catch((error) => {
          this.loading = false
        })
      },
      getRandomHash () {
        const randomPass = Math.random().toString(36).substring(2,6)
        this.hashList += `${sha1(randomPass)}\n`
        this.hashType = { code: '100', name: 'SHA1' },
        this.comment += `(${randomPass}) `
        this.validateHashes()
      },
      generateJobName () {
        this.name = this.$store.state.project + ' Job – ' + this.$moment().format('DD.MM.YYYY HH:mm')
      }
    }
  }
</script>

<style scoped>
  .containerAddJob {
    padding: 2em;
    padding-top: 54px;
    position: relative;
    max-width: 1300px;
  }

  .addJobContent {
    width: 100%;
  }

  .max500 {
    max-width: 500px;
    width: 100%;
  }

  .max800 {
    max-width: 800px;
    width: 100%;
  }

  .max1000 {
    max-width: 1000px;
  }

  .infobar {
    position: fixed;
    z-index: 5;
    bottom: 1.2em;
    right: 1.2em;
    padding: 0.5em 1.5em;
    border-radius: 2em;
  }

  .filetype-link {
    color: inherit
  }

  .hashCeckContainer {
    display: block;
    max-width: 35px;
    overflow: hidden;
  }

  .mode-btn {
    height: initial !important;
    margin: 1em;
  }

  .scroller {
    max-height: 400px;
    overflow-y: auto;
  }
</style>

<style>
  .hasherror .scrollCont {
    border-color: #e01 !important;
  }
</style>
