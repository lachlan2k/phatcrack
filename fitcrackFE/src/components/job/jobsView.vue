<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-toolbar flat>
      <v-chip>
        <v-fade-transition mode="out-in">
          <v-progress-circular
            v-if="loading"
            class="back-margin"
            size="16"
            width="2"
            indeterminate
          />
          <span v-else>
            {{ totalItems }}
          </span>
        </v-fade-transition>
      </v-chip>
      <h2 class="ml-4">
        {{ binTitle }}
        <span
          v-show="$vuetify.breakpoint.lgAndUp"
          class="text--secondary"
        >
          <span v-show="status || attackType">that </span>
          <span v-show="status">are {{ status }}</span>
          <span v-show="status && attackType"> and </span>
          <span v-show="attackType">use {{ attackType }} attack</span>
        </span>
      </h2>
      <v-spacer />
      <div v-if="isBin && $userCan('EDIT_ALL_JOBS')">
        <transition 
          name="fade"
          mode="out-in"
          @enter="focusRenameInput"
        >
          <div 
            v-if="renamingBin" 
            key="renameBinPanel"
            class="d-flex align-center"
          >
            <v-text-field
              ref="binRenameField"
              v-model="newBinName"
              :placeholder="bin.name"
              hide-details
              class="mr-2"
              @keydown="renameBinKeyHandler"
            />
            <v-btn
              text
              @click="renameBinExec"
            >
              <v-icon left>
                mdi-content-save
              </v-icon>
              Save
            </v-btn>
            <v-btn
              text
              :ripple="false"
              @click="renamingBin = false"
            >
              <v-icon left>
                mdi-undo-variant
              </v-icon>
              Cancel
            </v-btn>
          </div>
          <div 
            v-else
            key="binActionsPanel"
          >
            <v-btn
              text
              :ripple="false"
              @click="renamingBin = true"
            >
              <v-icon left>
                mdi-pencil
              </v-icon>
              Rename Bin
            </v-btn>
            <v-btn
              text
              @click="deleteBinConfirm"
            >
              <v-icon left>
                mdi-delete
              </v-icon>
              Delete Bin
            </v-btn>
          </div>
        </transition>
      </div>
    </v-toolbar>
    <div class="overflow-hidden toolbar-cont">
      <transition name="swap">
        <div
          v-if="selectedJobs.length == 0"
          key="filtersToolbar"
          class="toolbar-module d-flex justify-space-between align-center px-2"
        >
          <v-text-field
            v-model="search"
            class="mx-2 search-field"
            clearable
            prepend-icon="mdi-table-search"
            label="Search"
            single-line
            @input="updateList"
          />
          <v-select
            v-model="status"
            class="mx-2"
            :items="jobs_statuses"
            label="Status"
            single-line
            item-text="text"
            item-value="text"
            prepend-icon="mdi-timelapse"
            clearable
            @change="updateList"
          />
          <v-select
            v-model="attackType"
            class="mx-2"
            :items="jobs_types"
            label="Attack type"
            single-line
            item-text="text"
            item-value="text"
            prepend-icon="mdi-crosshairs-gps"
            clearable
            @change="updateList"
          />
        </div>
        <div
          v-else 
          key="selectionToolbar"
          class="toolbar-module d-flex align-center swap-pad"
        >
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                id="batch-creator-open"
                text
                class="mr-2"
                color="success"
                :disabled="selectedJobs.filter(job => job.status == 0 && !job.batch_id).length < 2"
                @click="batchCreatorOpen = true"
                v-on="on"
              >
                <v-icon left>
                  mdi-tray-plus
                </v-icon>
                Batch Run
              </v-btn>
            </template>
            <span>Only <strong class="info--text">ready</strong> jobs can be added</span>
          </v-tooltip>
          <v-btn
            id="host-mapper-open"
            text
            class="mr-2"
            color="warning"
            @click="hostEditorOpen = true"
          >
            <v-icon left>
              mdi-desktop-classic
            </v-icon>
            Hosts
          </v-btn>
          <v-btn
            text
            class="mr-2"
            @click="exportAdd"
          >
            <v-icon left>
              mdi-package-variant
            </v-icon>
            Add to Export
          </v-btn>
          <v-btn
            text
            class="mr-2"
            @click="multiTrash"
          >
            <v-icon left>
              {{ isTrash ? 'mdi-delete-restore' : 'mdi-delete' }}
            </v-icon>
            {{ isTrash ? 'Restore' : 'Discard' }}
          </v-btn>
          <v-btn
            text
            class="mr-2"
            @click="multiStart"
          >
            <v-icon left>
              mdi-play
            </v-icon>
            Start
          </v-btn>
          <v-btn
            text
            class="mr-2"
            @click="multiStop"
          >
            <v-icon left>
              mdi-stop
            </v-icon>
            Stop
          </v-btn>
          <v-btn
            text
            class="mr-2"
            @click="multiPurge"
          >
            <v-icon left>
              mdi-undo-variant
            </v-icon>
            Purge
          </v-btn>
        </div>
      </transition>
    </div>

    <v-divider />

    <div 
      v-if="!loading && binEmpty"
      class="text-center mt-8 pa-8 detail-container"
    >
      <h1>You had one job. And you binned it.</h1>
      <bin-illust 
        :label="bin.name" 
        class="illust"
      />
      <v-card class="text-start" flat>
        <v-card-title>Here's how:</v-card-title>
        <v-card-text>
          <ol>
            <li>List <router-link to="/jobs"><i class="mdi mdi-briefcase" /> <em>All Jobs</em></router-link> or another bin</li>
            <li>Check <i class="mdi mdi-checkbox-marked-outline" /> the jobs you want to add here</li>
            <li>Select <span class="success--text"><i class="mdi mdi-plus" /> <em>Add to {{bin.name}}</em></span> in the sidebar</li>
          </ol>
        </v-card-text>
      </v-card>
    </div>

    <v-data-table
      id="jobs-table"
      v-show="!binEmpty"
      ref="table"
      v-model="selectedJobs"
      :headers="headers"
      :items="jobs"
      :search="search"
      :options.sync="pagination"
      :server-items-length="totalItems"
      :footer-props="{itemsPerPageOptions: [25,10,50,100], itemsPerPageText:'Jobs per page'}"
      show-select
      fixed-header
    >
      <!-- Job name cell -->
      <template v-slot:item.name="{ item }">
        <router-link
          :to="{ name: 'jobDetail', params: { id: item.id } }"
          class="table-link"
        >
          <span>{{ item.name }}</span>
          <v-icon
            v-if="item.batch_id"
            right
            small
            color="primary"
          >
            mdi-tray-full
          </v-icon>
        </router-link>
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
      <!-- Status text cell -->
      <template v-slot:item.status="{ item }">
        <v-tooltip
          v-if="item.batch_id && item.status == 0"
          top
        >
          <template v-slot:activator="{ on }">
            <router-link :to="`/batches/${item.batch_id}`">
              <span
                class="text-capitalize"
                v-on="on"
              >
                enqueued
                <v-icon
                  right
                  color="primary"
                >
                  mdi-timetable
                </v-icon>
              </span>
            </router-link>
          </template>
          <span>This job is waiting in queue. Click to see the batch.</span>
        </v-tooltip>
        <v-tooltip
          v-else-if="item.host_count > 0"
          top
        >
          <template v-slot:activator="{ on }">
            <span
              :class="[item.status_type + '--text', 'text-capitalize']"
              v-on="on"
            >
              {{ item.status_text }}
              <v-icon
                right
                :color="item.status_type"
              >
                {{ jobIcon(item.status_text) }}
              </v-icon>
            </span>
          </template>
          <span>{{ item.status_tooltip }}</span>
        </v-tooltip>
        <v-tooltip
          v-else
          top
        >
          <template v-slot:activator="{ on }">
            <span
              class="text-capitalize"
              v-on="on"
            >
              no hosts
              <v-icon
                right
              >
                mdi-alert
              </v-icon>
            </span>
          </template>
          <span>No hosts, open job detail and assign them!</span>
        </v-tooltip>
      </template>
      <!-- Progress indicator cell -->
      <template v-slot:item.progress="{ item }">
        <div class="d-flex align-center justify-end">
          <span class="mr-2">{{ progressToPercentage(item.progress) }}</span>
          <v-progress-circular
            size="18"
            :width="3"
            :rotate="270"
            color="primary"
            class="jobProgress"
            :value="item.progress"
          />
        </div>
      </template>
      <!-- Date added cell -->
      <template v-slot:item.time="{ item }">
        {{ $moment.utc(item.time).local().format('D.M.YYYY H:mm:ss') }}
      </template>
      <!-- Action buttons cell -->
      <template v-slot:item.actions="{ item }">
        <div class="actionsBtns">
          <v-btn
            v-if="item.batch_id != null && (item.status == 0 || item.status >= 10)"
            text
            color="primary"
            :to="`/batches/${item.batch_id}`"
          >
            Go to Batch
            <v-icon right>
              mdi-tray-full
            </v-icon>
          </v-btn>
          <template v-else-if="item.permissions.operate || $userCan('OPERATE_ALL_JOBS')">
            <v-btn
              v-if="item.status === '0'"
              :disabled="item.host_count == 0"
              text
              color="success"
              @click="operateJob(item.id, 'start')"
            >
              Start
              <v-icon right>
                mdi-play
              </v-icon>
            </v-btn>
            <v-btn
              v-else-if="item.status < 10"
              :disabled="item.host_count == 0"
              text
              color="info"
              @click="operateJob(item.id, 'restart')"
            >
              Restart
              <v-icon right>
                mdi-restart
              </v-icon>
            </v-btn>
            <v-btn
              v-else
              text
              color="error"
              :disabled="item.status === '12'"
              @click="operateJob(item.id, 'stop')"
            >
              Stop
              <v-icon right>
                mdi-stop
              </v-icon>
            </v-btn>
          </template>

          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                v-if="item.permissions.edit || $userCan('EDIT_ALL_JOBS')"  
                icon
                class="mx-0"
                v-on="on"
                @click="hideJob(item.id)"
              >
                <v-icon>
                  {{ isTrash ? 'mdi-delete-restore' : 'mdi-delete' }}
                </v-icon>
              </v-btn>
            </template>
            <span>{{ isTrash ? 'Restore' : 'Discard' }}</span>
          </v-tooltip>
        </div>
      </template>
      <template #footer.page-text="{ pageStart, pageStop, itemsLength }">
        <strong>Page {{ pagination.page }}</strong> – {{ pageStart }}-{{ pageStop }} of {{ itemsLength }}
      </template>
    </v-data-table>
    <!-- Dialogs -->
    <host-editor
      v-model="hostEditorOpen"
      :job-ids="selectedJobs.map(j => j.id)"
      @reload="updateList"
    />
    <v-dialog 
      v-model="batchCreatorOpen"
      max-width="900"
    >
      <batch-creator 
        :jobs="selectedJobs"
        @close="batchCreatorOpen = false"
      />
    </v-dialog>
  </div>
</template>

<script>
  import { jobIcon, attackIcon } from '@/assets/scripts/iconMaps'
  import { mapTwoWayState } from 'spyfu-vuex-helpers'
  import { twoWayMap } from '@/store'
  import { mapState, mapMutations, mapActions } from 'vuex'

  import BinIllust from '@/components/details/binIllustration.vue'
  import HostEditor from '@/components/jobDetail/hostEditor.vue'
  import BatchCreator from '@/components/batch/batchCreator.vue'

  export default {
    name: 'JobsView',
    components: {
      BinIllust,
      HostEditor,
      BatchCreator
    },
    data () {
      return {
        interval: null,
        renamingBin: false,
        newBinName: null,
        hostEditorOpen: false,
        batchCreatorOpen: false,
        //
        status: null,
        attackType: null,
        search: '',
        //
        totalItems: 0,
        pagination: {},
        loading: true,
        //
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Attack type', value: 'attack', align: 'end', sortable: false},
          {text: 'Status', value: 'status', align: 'end'},
          {text: 'Progress', value: 'progress', align: 'end'},
          {text: 'Added', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', sortable: false, align: 'end'}
        ],
        jobs_statuses: [
          {
            'text': 'ready',
            'code': 0
          },
          {
            'text': 'finished',
            'code': 1
          },
          {
            'text': 'exhausted',
            'code': 2
          },
          {
            'text': 'malformed',
            'code': 3
          },
          {
            'text': 'timeout',
            'code': 4
          },
          {
            'text': 'running',
            'code': 10
          },
          {
            'text': 'finishing',
            'code': 12
          }
        ],
        jobs_types: [
          {
            'text': 'dictionary',
            'code': 0
          },
          {
            'text': 'mask',
            'code': 1
          },
          {
            'text': 'combinator',
            'code': 2
          },
          {
            'text': 'pcfg',
            'code': 3
          },
          {
            'text': 'prince',
            'code': 4
          },
          {
            'text': 'hybrid (wordlist + mask)',
            'code': 5
          },
          {
            'text': 'hybrid (mask + wordlist)',
            'code': 6
          }
        ],
        jobs:
          [
          ]
      }
    },
    computed: {
      ...mapState('binInterface', ['clean']),
      ...mapTwoWayState('binInterface', twoWayMap(['selectedJobs'])),
      isTrash () {
        return this.$route.params['id'] === 'trash'
      },
      isBin () {
        return this.$route.name === 'bins' && !this.isTrash
      },
      binId () {
        return parseInt(this.$route.params['id'])
      },
      bin () {
        return this.$store.state.binInterface.bins.find(bin => bin.id == this.binId)
      },
      binTitle () {
        if (this.bin) {
          return this.bin.name
        } else if (this.isTrash) {
          return 'Discarded Jobs'
        } else {
          return this.$userCan('VIEW_ALL_JOBS') ? 'All Jobs' : 'My Jobs'
        }
      },
      binEmpty () {
        return this.bin && this.bin.job_count == 0
      }
    },
    watch: {
      pagination: {
        deep: true,
        handler ({ page }) {
          this.$paginator(this.$route.path, page)
          this.updateList(false)
        }
      },
      $route ({ path }) {
        this.pagination.page = this.$paginator(path)
        this.status = null
        this.attackType = null
        this.search = ''
        this.updateList(false)
      },
      clean (isClean) {
        if (!isClean) {
          this.updateList(false)
          this.$store.commit('binInterface/clean')
        }
      }
    },
    mounted () {
      // this.loadJobs()
      this.pagination.page = this.$paginator(this.$route.path)
      this.interval = setInterval(() => {
        this.loadJobs()
      }, 5000)
      window.addEventListener('keydown', this.shortcutHandler)
    },
    beforeDestroy: function () {
      this.selectedJobs = []
      clearInterval(this.interval)
      window.removeEventListener('keydown', this.shortcutHandler)
    },
    methods: {
      ...mapMutations('transfer', { addJobsForTransfer: 'add' }),
      ...mapActions('binInterface', { renameBin: 'rename', deleteBin: 'delete' }),
      async loadJobs () {
        const data = await this.axios.get(this.$serverAddr + '/job', {
          params: {
            'page': this.pagination.page,
            'per_page': this.pagination.itemsPerPage,
            'order_by': this.pagination.sortBy[0],
            'descending': this.pagination.sortDesc[0],
            'name': this.search,
            'status': this.status,
            'attack_mode': this.attackType,
            'showDeleted': this.isTrash,
            'bin': this.binId || null
          }
        }).then(r => r.data)
        this.jobs = data.items
        this.totalItems = data.total
        this.loading = false
      },
      updateList (resetPagination = true) {
        if (resetPagination) this.pagination.page = 1
        this.selectedJobs = []
        this.loading = true
        this.loadJobs()
      },
      updateStatus: function (e) {
        this.status = e;
        this.$refs.table.updatePagination({ page: 1, totalItems: this.totalItems })
      },
      progressToPercentage: function (progress) {
        if(progress > 100){
          progress = 100
        }
        return progress.toFixed() + '%'
      },
      operateJob: function (id, operation) {
        this.axios.get(this.$serverAddr + '/job/' + id + '/action', {
          params: {
            'operation': operation
          }
        }).then((response) => {
          console.log(response.data);
          this.loadJobs()
        })
      },
      hideJob: function (id) {
        this.loading = true
        this.axios.delete(this.$serverAddr + '/job/' + id)
        .then((response) => {
          this.loadJobs()
        })
      },
      multiTrash () {
        if (this.selectedJobs.length == 0) return
        //
        this.loading = true
        this.axios.patch(this.$serverAddr + '/job', {
          job_ids: this.selectedJobs.map(j => j.id)
        })
        .then((response) => {
          this.loadJobs()
        })
      },
      multiStop () {
        if (this.selectedJobs.length == 0) return
        //
        this.loading = true
        this.axios.post(this.$serverAddr + '/job/multiJobOperation', {
          job_ids: this.selectedJobs.map(j => j.id),
          operation: 'stop'
        })
        .then((response) => {
          this.loadJobs()
        })
      },
      multiStart () {
        if (this.selectedJobs.length == 0) return
        //
        this.loading = true
        this.axios.post(this.$serverAddr + '/job/multiJobOperation', {
          job_ids: this.selectedJobs.map(j => j.id),
          operation: 'start'
        })
        .then((response) => {
          this.loadJobs()
        })
      },
      multiPurge () {
        if (this.selectedJobs.length == 0) return
        //
        this.loading = true
        this.axios.post(this.$serverAddr + '/job/multiJobOperation', {
          job_ids: this.selectedJobs.map(j => j.id),
          operation: 'kill'
        })
        .then((response) => {
          this.loadJobs()
        })
      },
      focusRenameInput () {
        if (this.renamingBin) {
          this.$refs.binRenameField.focus()
        }
      },
      renameBinExec () {
        this.renameBin({
          id: this.binId,
          newName: this.newBinName
        }).then(() => {
          this.renamingBin = false
        })
      },
      renameBinKeyHandler ({ keyCode }) {
        switch (keyCode) {
          case 13:
            this.renameBinExec()
            break
          case 27:
            this.renamingBin = false
            break
        }
      },
      deleteBinConfirm () {
        this.$root.$confirm('Delete Bin', `This will remove ${this.bin.name}. Jobs will be unassigned from the deleted bin, but will not be discarded. Are you sure?`)
          .then((confirm) => {
            this.deleteBin(this.binId)
          })
      },
      exportAdd () {
        this.addJobsForTransfer(this.selectedJobs)
      },
      jobIcon,
      attackIcon,
      //
      shortcutHandler (e) {
        if (e.key === 'Delete' && this.isBin) {
          this.deleteBinConfirm()
        }
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

  .back-margin {
    margin: 0 -.3em;
  }

  .toolbar-cont {
    position: relative;
    height:4.15em;
  }
  .toolbar-module {
    position: absolute;
    flex-wrap: wrap-reverse;
    width:100%;
  }
  @media screen and (max-width: 60em) {
    .toolbar-cont {
      height: auto;
    }
    .toolbar-module {
      position: static;
    }
  }

  .search-field {
    min-width: 15em;
  }

  .cont {
    height: 100%;
  }

  .actionsBtns{
    text-align: right;
  }

  .swap-pad {
    padding: 1.05em;
  }

  .swap-enter-active, .swap-leave-active {
    transition: transform .2s cubic-bezier(0.785, 0.135, 0.15, 0.86);
  }

  .swap-leave-to {
    transform: translateY(-100%)
  }
  .swap-enter {
    transform: translateY(100%)
  }

  .fade-enter-active, .fade-leave-active {
    transition: opacity .1s, transform .2s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .fade-leave-to {
    opacity: 0;
  }
  .fade-enter {
    transform: translateX(10%)
  }

</style>
