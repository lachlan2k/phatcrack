<template>
  <v-dialog 
    max-width="1200"
    :value="value"
    @input="$emit('input')"
  >
    <v-card>
      <v-card-title>
        <span>Host Mapping</span>
        <v-spacer />
        <v-btn
          icon
          @click="$emit('input')"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text class="scroller">
        <host-selector
          v-model="hosts"
          select-all
          auto-refresh
        />
      </v-card-text>
      <v-card-actions class="pa-6">
        <v-btn
          text
          to="/hosts"
        >
          <v-icon left>
            mdi-desktop-classic
          </v-icon>
          Host List
        </v-btn>
        <v-spacer />
        <b class="mr-4">{{ hosts.length }} selected</b>
        <v-btn 
          id="host-mapper-assign"
          :loading="saving"
          color="success"
          @click="save"
        >
          <v-icon left>
            mdi-check
          </v-icon>
          Assign
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import hostSelector from '@/components/selector/hostSelector.vue'

export default {
  components: {
    hostSelector
  },
  props: {
    value: Boolean,
    jobIds: {
      type: Array,
      required: true
    }
  },
  data () {
    return {
      hosts: [],
      saving: false
    }
  },
  methods: {
    async save () {
      this.saving = true
      try {
        await this.axios.post(this.$serverAddr + '/job/host', {
          'job_ids': this.jobIds,
          'newHost_ids': this.hosts.map(host => host.id)
        })
        this.$emit('input') // close
        this.$emit('reload')
      } catch (e) {
        console.error(e)
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style scoped>
.scroller {
  max-height: 60vh;
  overflow-y: auto;
}
</style>