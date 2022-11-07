<template>
  <div class="job-info">
    <div class="sticky-wrap">
      <!-- TOP SHEET -->
      <v-sheet
        class="info-sheet top-sheet pa-4"
        dark
        tile
        :color="statusColor"
      >
        <v-icon
          class="attack-icon"
          aria-hidden
        >
          {{ attackIcon(data.attack) }}
        </v-icon>
        <div class="top-sheet-bar mb-2">
          <div>
            <div class="name title mb-2">
              {{ data.name }}
            </div>
            <div class="subtitle-2 ">
              <b class="text-capitalize">{{ data.attack }}</b> Attack
              <br>
              <b>{{ data.hash_type_name }}</b>
            </div>
          </div>
        </div>
        <div class="status">
          <v-progress-circular
            class="fu"
            :value="data.progress"
            :indeterminate="indtProgress"
            size="26"
            rotate="-90"
          />
          <span class="text-capitalize pt-1 ml-2">
            <span v-if="data.batch && data.batch.id && data.status == 0">
              Enqueued
            </span>
            <span v-else>
              <b v-show="data.progress > 0 && data.progress < 100">{{ data.progress }} %</b>
              {{ data.status_text }}
            </span>
          </span>
        </div>
      </v-sheet>
      <!-- ACTION SHEET -->
      <v-sheet
        class="info-sheet cta-sheet pa-2 d-flex"
        dark
        tile
        :color="statusColor"
      >
        <template v-if="!editing">
          <!-- action button -->
          <template v-if="data.permissions.operate || $userCan('OPERATE_ALL_JOBS')">
            <v-btn
              v-if="data.batch && data.batch.id && data.status == 0"
              key="batch-link"
              text
              :to="`/batches/${data.batch.id}`"
            >
              <v-icon left>
                mdi-tray-full
              </v-icon>
              <span>Go to Batch</span>
            </v-btn>
            <v-btn
              v-else-if="data.hosts.length > 0"
              key="operate"
              text
              :disabled="data.status === '12'"
              @click="$emit('operate', operation.text)"
            >
              <v-icon left>
                {{ operation.icon }}
              </v-icon>
              <span>{{ operation.text }}</span>
            </v-btn>
            <v-btn
              v-else
              key="assign"
              text
              @click="$emit('edit-hosts')"
            >
              <v-icon left>
                mdi-desktop-classic
              </v-icon>
              <span>Add Hosts</span>
            </v-btn>
          </template>
          <v-spacer />
          <!-- edit button -->
          <v-btn
            id="job-editor-open"
            v-if="data.permissions.edit || $userCan('EDIT_ALL_JOBS')"
            key="edit-enter"
            text
            @click="editing = true"
          >
            <span>Edit</span>
            <v-icon right>
              mdi-pencil
            </v-icon>
          </v-btn>
        </template>
        <template v-else>
          <v-btn
            id="job-permissions-open"
            v-if="data.permissions.owner"
            key="edit-perms"
            text
            @click="$emit('edit-permissions')"
          >
            <v-icon left>
              mdi-account
            </v-icon>
            Permissions
          </v-btn>
          <v-spacer />
          <!-- cancel edit button -->
          <v-btn
            key="edit-cancel"
            text
            @click="editing = false"
          >
            <span>Cancel</span>
            <v-icon right>
              mdi-undo-variant
            </v-icon>
          </v-btn>
        </template>
      </v-sheet>
      <!-- INFO LIST -->
      <v-list
        class="info-list"
        inactive
        two-line
      >
        <transition name="fade" mode="out-in">
          <transition-group
            v-if="!editing"
            tag="div"
            name="list"
          >
            <v-list-item
              v-for="({title, icon, value}) in details"
              :key="title"
            >
              <v-list-item-content>
                <v-list-item-subtitle class="overline mb-1">
                  <v-icon
                    class="mr-1"
                    small
                  >
                    {{ icon }}
                  </v-icon>
                  {{ title }}
                </v-list-item-subtitle>
                <v-list-item-title class="pre">{{ value }}</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </transition-group>
          <!-- EDITOR -->
          <template v-else>
            <job-editor
              :data="data"
              @save="saveEdit"
            />
          </template>
        </transition>
      </v-list>
    </div>
  </div>
</template>

<script>
import JobEditor from './jobEditor.vue'

import { jobIcon, attackIcon } from '@/assets/scripts/iconMaps'
import fmt from '@/assets/scripts/numberFormat'

export default {
  components: {
    JobEditor
  },
  props: {
    data: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      editing: false,
      testmode: localStorage.getItem('testmode') == 'true' || false
    }
  },
  computed: {
    statusColor () {
      if (this.editing || this.data.hosts.length == 0) {
        return 'secondary'
      } else if (this.data.batch && this.data.batch.id && this.data.status == 0) {
        return 'primary'
      } 
      else {
        return this.data.status_type
      }
    },
    indtProgress () {
      return ((this.data.progress == 0 || this.data.progress >= 100) && this.data.status === '10') || parseInt(this.data.status) > 10
    },
    operation () {
      if (this.data.status === '0') {
        return {text: 'start', icon: 'mdi-play'}
      } else if (this.data.status < 10) {
        return {text: 'restart', icon: 'mdi-restart'}
      } else {
        return {text: 'stop', icon: 'mdi-stop'}
      }
    },
    // Dynamically create the info list items
    details () {
      const d = this.data
      const items = [
        {title: 'Part of batch', icon: 'mdi-tray-full', value: d.batch.name},
        {title: 'Comment', icon: 'mdi-format-quote-close', value: d.comment},
        {title: 'Keyspace', icon: 'mdi-key', value: d.keyspace, format: fmt},
        {title: 'Date Added', icon: 'mdi-calendar-import', value: d.time, format: v => this.$moment.utc(v).local().calendar()},
        {title: 'Workunit Sum Time', icon: 'mdi-sigma', value: d.workunit_sum_time_str},
        {title: 'Cracking Time', icon: 'mdi-timer', value: this.crackingTime},
        {title: 'Estimated Cracking Time', icon: 'mdi-bell-ring-outline', value: d.estimated_cracking_time_str},
        {title: 'Start', icon: 'mdi-ray-start', value: this.startTime},
        {title: 'End', icon: 'mdi-ray-end', value: this.endTime},
        {title: 'Efficiency', hide: () => !this.testmode, icon: 'mdi-chart-timeline-variant', value: d.efficiency, format: v => `${v} %`},
        {title: 'Cracked Hashes', icon: 'mdi-textbox-password', value: d.cracked_hashes_str},
        {title: 'Time per workunit', icon: 'mdi-timeline-clock', value: d.seconds_per_job, format: v => `${v} seconds`},
      ]
      // Leave out missing stuff and format fields
      return items
        .filter(i => !i.hide || !i.hide()) // only leave fields without hide condition or when it resolves false
        .filter(i => i.value != null && i.value != undefined & i.value !== '')
        .map(i => {
          if (i.format) {
            i.value = i.format(i.value)
          }
          return i
        })
    },
    crackingTime () {
      const d = this.data
      if (d.total_time) {
        return this.$moment.duration(d.total_time, 'seconds').format()
      } else {
        return 'Not started yet'
      }
    },
    startTime () {
      const d = this.data
      if (d.time_start) {
        return this.$moment.utc(d.time_start).local().format('DD.MM.YYYY HH:mm')
      } else {
        return 'Not started yet'
      }
    },
    endTime () {
      const d = this.data
      if (d.time_end) {
        return this.$moment.utc(d.time_end).local().format('DD.MM.YYYY HH:mm')
      } else {
        return 'Not planned'
      }
    }
  },
  methods: {
    jobIcon,
    attackIcon,
    fmt,
    saveEdit () {
      this.editing = false
      this.$emit('reload')
    }
  }
}
</script>

<style scoped>
.job-info {
  position: relative;
  display: flex;
  flex-direction: column;
}

.sticky-wrap {
  position: sticky;
  top: 64px;
}

.info-sheet {
  position: relative;
  z-index: 1;
  transition: background-color .5s;
}

.theme--dark .top-sheet::before,
.theme--light .cta-sheet::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0,0,0,.25);
  z-index: -1;
}

.wrap { white-space: normal; }
.pre { white-space: pre; }

.top-sheet-bar {
  display: flex;
  justify-content: space-between;
}

.status {
  display: inline-flex;
  align-items: center;
}

/**/

.attack-icon {
  position: absolute;
  font-size: 8em;
  z-index: 0;
  opacity: .3;
  right: -.23em;
  top: .1em;
}

.z { z-index: 1 }

/**/

  .fade-enter-active {
    transition: transform .5s cubic-bezier(0.16, 1, 0.3, 1);
  }

  .fade-leave-active {
    transition: opacity .2s;
  }

  .fade-leave-to {
    opacity: 0;
  }
  .fade-enter {
    opacity: 0;
    transform: translateY(-15px)
  }

</style>

<style>
.v-progress-circular__underlay {
  stroke: rgba(255,255,255,.3);
}
</style>