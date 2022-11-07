<template>
  <v-card>
    <v-card-title>
      Prepare Job Queue
    </v-card-title>
    <v-card-text>
      <v-row class="align-center">
        <v-col cols="8">
          Jobs will run in sequence as listed. Change the order by dragging them or select a sort.
        </v-col>
        <v-col cols="4">
          <v-select
            v-model="sorting"
            label="Sort"
            outlined
            hide-details
            :items="sorts"
            @change="sortQueue"
          />
        </v-col>
      </v-row>
      <draggable
        v-model="queue"
        :animation="150"
        @change="sorting = undefined"
      >
        <v-list-item
          v-for="{ id, name, keyspace, hc_keyspace, attack } in queue"
          :key="id"
          class="hand"
        >
          <v-list-item-action>
            <v-icon>{{ attackIcon(attack) }}</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title 
              class="font-weight-bold"
              :style="{color: identColor(id, keyspace)}"
            >
              {{ name }}
            </v-list-item-title>
            <v-list-item-subtitle>
              <v-icon x-small>
                mdi-key
              </v-icon>
              <span>{{ numberFormat(keyspace) }} possible passwords</span>
              <span>(computing <strong>{{ numberFormat(hc_keyspace) }}</strong>)</span>
            </v-list-item-subtitle>
          </v-list-item-content>
          <v-list-item-action>
            <v-btn
              text
              :disabled="queue.length <= 2"
              @click="deselect(id)"
            >
              <span>Unqueue</span>
              <v-icon right>
                mdi-tray-remove
              </v-icon>
            </v-btn>
          </v-list-item-action>
        </v-list-item>
      </draggable>
    </v-card-text>
    <transition-group
      tag="div"
      class="d-flex mx-6 mb-4"
      name="ks-graph"
    >
      <div
        v-for="{ id, name, keyspace, hc_keyspace } in queue"
        :key="name"
        class="ks subtitle-2"
        :style="{flexGrow: hc_keyspace}"
      >
        <div
          class="ks-number"
          :style="{backgroundColor: identColor(id, keyspace)}"
        >
          {{ numberFormat(hc_keyspace) }}
        </div>
      </div>
    </transition-group>
    <v-card-actions class="align-center">
      <v-btn
        text
        @click="$emit('close')"
      >
        <v-icon left>
          mdi-close
        </v-icon>
        <span>Cancel</span>
      </v-btn>
      <v-spacer />
      <v-text-field
        v-model="batchName"
        label="Batch name"
        hide-details
        outlined
        dense
        class="mr-2"
      />
      <v-btn
        color="success"
        :disabled="batchName === ''"
        :loading="submitting"
        @click="submit"
      >
        <span>Create Batch</span>
        <v-icon right>
          mdi-tray-plus
        </v-icon>
      </v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import numberFormat from '@/assets/scripts/numberFormat'
import { attackIcon } from '@/assets/scripts/iconMaps'

import draggable from 'vuedraggable'

export default {
  components: {
    draggable
  },
  props: {
    jobs: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      queue: [],
      batchName: '',
      sorts: [
        {text: 'Fastest first', value: {by: 'hc_keyspace', asc: true}},
        {text: 'Slowest first', value: {by: 'hc_keyspace', asc: false}},
        {text: 'Less passwords first', value: {by: 'keyspace', asc: true}},
        {text: 'More passwords first', value: {by: 'keyspace', asc: false}},
        {text: 'Newest first', value: {by: 'time', asc: false}},
        {text: 'Oldest first', value: {by: 'time', asc: true}},
      ],
      sorting: undefined,
      submitting: false
    }
  },
  watch: {
    jobs: {
      immediate: true,
      handler (val) {
        this.queue = val.filter(job => job.status == 0 && !job.batch_id)
        if (this.sorting) {
          this.sortQueue(this.sorting)
        }
      }
    }
  },
  methods: {
    numberFormat,
    attackIcon,
    identColor (v1, v2) {
      return `hsl(${v1 * v2}, 65%, 45%)`
    },
    sortQueue ({by, asc}) {
      this.queue.sort((a, b) => {
        let x, y
        if (by === 'time') {
          x = this.$moment(a[by]).unix()
          y = this.$moment(b[by]).unix()
        } else {
          x = parseInt(a[by])
          y = parseInt(b[by])
        }
        return asc ? x - y : y - x
      })
    },
    deselect (id) {
      this.queue = this.queue.filter(job => job.id != id)
    },
    async submit () {
      this.submitting = true
      try {
        const res = await this.axios.put(this.$serverAddr + '/batches', {
          name: this.batchName,
          jobs: this.queue.map(job => job.id)
        }).then(r => r.data)
        this.$store.commit('binInterface/selectedJobsMut', [])
        this.$router.push(`/batches/${res}`)
      } catch (e) {
        console.error(e)
      } finally {
        this.submitting = false
      }
    }
  }
}
</script>

<style scoped>
.hand {
  cursor: grab;
}
.ks {
  text-align: center;
  user-select: none;
}
.ks:not(:last-child) {
  margin-right: 2px;
}
.ks-number {
  color: white;
  line-height: 1.8em;
  padding: 0 .5em;
}
.ks-graph-move {
  transition: transform .4s cubic-bezier(.55,0,.15,1);
}
</style>