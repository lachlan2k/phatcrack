<template>
  <v-card :class="['minijob rounded-lg', {detailed}]">
    <div class="minijob-main">
      <v-sheet
        class="ma-2 overflow-hidden minijob-sheet"
        :color="data.status_type"
        dark
        rounded
      >
        <router-link
          :to="{name: 'jobDetail', params: {id: data.id}}"
          class="transparent-link"
        >
          <v-card-title class="mj-title pb-0">
            {{ data.name }}
          </v-card-title>
        </router-link>
        <div class="mj-progress px-4">
          {{ data.progress }}<span> %</span>
        </div>
        <v-progress-linear
          :value="data.status_text === 'finishing' ? 0 : data.progress"
          color="white"
          :indeterminate="data.status_text === 'finishing'"
        />
      </v-sheet>
      <v-list
        dense
        class="minijob-list"
      >
        <v-list-item>
          <v-list-item-action>
            <v-icon>{{ jobIcon(data.status_text) }}</v-icon>
          </v-list-item-action>
          <v-list-item-content class="text-capitalize">
            <v-list-item-subtitle class="overline">
              Status
            </v-list-item-subtitle>
            <v-list-item-title>
              {{ data.status_text }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action>
            <v-icon>{{ attackIcon(data.attack) }}</v-icon>
          </v-list-item-action>
          <v-list-item-content class="text-capitalize">
            <v-list-item-subtitle class="overline">
              Type
            </v-list-item-subtitle>
            <v-list-item-title>
              {{ data.attack }} attack
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-show="parseInt(data.status) > 0">
          <v-list-item-action>
            <v-icon>mdi-textbox-password</v-icon>
          </v-list-item-action>
          <v-list-item-content class="text-capitalize">
            <v-list-item-subtitle class="overline">
              Cracked Hashes
            </v-list-item-subtitle>
            <v-list-item-title>
              {{ data.cracked_hashes_str }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item v-show="parseInt(data.status) >= 10">
          <v-list-item-action>
            <v-icon>mdi-bell-check-outline</v-icon>
          </v-list-item-action>
          <v-list-item-content class="text-capitalize">
            <v-list-item-subtitle class="overline">
              ETA
            </v-list-item-subtitle>
            <v-list-item-title>
              {{ data.estimated_cracking_time_str }}
            </v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </div>
    <v-divider
      v-if="detailed"
      vertical
      class="ml-2"
    />
    <div
      v-if="detailed"
      class="minijob-detail"
    >
      <wu :id="data.id" />
    </div>
  </v-card>
</template>

<script>
import wu from '@/components/chart/jobWorkunits.vue'
import { jobIcon, attackIcon } from '@/assets/scripts/iconMaps'

export default {
  components: {
    wu
  },
  props: {
    detailed: Boolean,
    data: {
      type: Object,
      default: () => {}
    }
  },
  methods: {
    jobIcon,
    attackIcon
  }
}
</script>

<style>
:root {
  --mj-height: 300px;
  --mj-basis: 28ch;
  --mj-detail: 450px;
}

.minijob-detail > .container > div { /* fixed size for chart */
  height: calc(var(--mj-height) - 1.5em) !important
}
</style>
<style scoped>
.minijob {  
  display: flex;
  margin-right: 1.5em;
  min-width: var(--mj-basis);
  min-height: var(--mj-height);
}

.minijob.detailed {
  min-width: calc(var(--mj-basis) + var(--mj-detail));
}

.minijob-list {
  background: transparent;
  margin-top: .5em;
}

.minijob-sheet {
  box-shadow: 0 .3em 1em -.5em rgba(0,0,0,.3);
}

.transparent-link {
  color: unset;
}

.mj-title {
  word-break: break-word;
}

.mj-progress {
  font-weight: bold;
  font-size: 2.5em;
}
.mj-progress span {
  font-size: 1rem;
}

.minijob-main {
  flex: 0 0 var(--mj-basis);
}

.minijob-detail {
  width: var(--mj-detail);
}

.detailed .bottom-tiled {
  border-top-right-radius: 0;
}
.bottom-tiled {
  border-bottom-left-radius: 0;
  border-bottom-right-radius: 0;
}
</style>
