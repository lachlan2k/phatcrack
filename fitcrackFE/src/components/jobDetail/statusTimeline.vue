<template>
  <div class="d-flex history-wrap">
    <v-tooltip top>
      <template #activator="{ on }">
        <v-btn
          v-show="entries.length > visible.length"
          class="mr-4 mt-3"
          icon
          @click="showAll = true"
          v-on="on"
        >
          <v-icon>mdi-history</v-icon>
        </v-btn>
      </template>
      <span>All Entries</span>
    </v-tooltip>
    <div class="d-flex justify-space-between history-container">
      <div
        v-for="{ time, status_text, status_type } in visible"
        :key="time"
        class="history-entry"
      >
        <div class="overline">{{ status_text }}</div>
        <v-sheet class="mx-3">
          <v-icon :color="status_type">
            {{ jobIcon(status_text) }}
          </v-icon>
        </v-sheet>
        <div class="subtitle-2 mt-2">
          {{ $moment.utc(time).local().format('DD.MM.YYYY') }}
        </div>
        <div>
          {{ $moment.utc(time).local().format('HH:mm') }}
        </div>
      </div>
    </div>
    <v-dialog 
      v-model="showAll"
      max-width="500"
    >
      <v-card>
        <v-card-title>
          <span>Full History</span>
          <v-spacer />
          <v-btn
            icon
            @click="showAll = false"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>
        <v-card-text>
          <v-data-table
            :headers="[
              {text: 'Time', value: 'time', align: 'start'},
              {text: 'Status', value: 'status', align: 'end'},
            ]"
            :items="entries"
            :footer-props="{itemsPerPageOptions: [10,25,50,{text: 'All', value: -1}]}"
          >
            <template v-slot:item.time="{ item }">
              <strong>
                {{ $moment.utc(item.time).local().format('DD.MM.YYYY') }}
              </strong>
              <span>
                {{ $moment.utc(item.time).local().format('HH:mm') }}
              </span>
            </template>
            <template v-slot:item.status="{ item }">
              <v-tooltip top>
                <template v-slot:activator="{ on }">
                  <span 
                    :class="['font-weight-medium', item.status_type + '--text']"
                    v-on="on"
                  >
                    {{ item.status_text }}
                  </span>
                </template>
                <span>{{ item.status_tooltip }}</span>
              </v-tooltip>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script>
import { jobIcon } from '@/assets/scripts/iconMaps'

export default {
  props: {
    entries: {
      type: Array,
      default: () => []
    }
  },
  data () {
    return {
      showAll: false
    }
  },
  computed: {
    visible () {
      let lastReady = 0
      for (let i = this.entries.length - 1; i >= 0; i--) {
        if (this.entries[i].status === '0') {
          lastReady = i
          break
        }
      }
      return this.entries.slice(lastReady)
    }
  },
  methods: {
    jobIcon
  }
}
</script>

<style>
  .history-container {
    position: relative;
    flex-grow: 1;
  }
  .history-container::before {
    content: '';
    position: absolute;
    width: calc(100% - 6em);
    height: 2px;
    background-color: #888;
    margin: 0 2em;
    top: 2.75em;
  }

  .history-entry {
    text-align: center;
    z-index: 2;
  }
</style>