<template>
  <div>
    <v-progress-linear
      v-if="!tree"
      indeterminate
    />
    <template v-if="tree">
      <v-toolbar flat>
        <h2>
          {{ tree.name }}
        </h2>
      </v-toolbar>
      <v-container>
        <v-row>
          <v-col>
            <v-card min-width="300">
              <v-card-text>
                <v-treeview
                  :items="tree.children"
                  item-key="path"
                  item-disabled="empty"
                  hoverable
                  activatable
                  open-on-click
                  transition
                  @update:active="open"
                />
              </v-card-text>
            </v-card>
          </v-col>
          <v-col class="grow-greedy">
            <v-card min-width="300">
              <v-card-title>
                {{ file || 'Pick a file to open' }}
              </v-card-title>
              <v-card-text>
                <pre class="monospace wrap">
{{ data }}
                </pre>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </template>
  </div>
</template>

<script>
export default {
  data () {
    return {
      loading: true,
      tree: null,
      file: null,
      data: null
    }
  },
  mounted () {
    this.load()
  },
  methods: {
    async load () {
      this.loading = true
      const id = this.$route.params.id
      this.tree = await this.axios.get(`${this.$serverAddr}/pcfg/${id}/tree`).then(r => r.data)
      this.loading = false
    },
    async open ([path]) {
      if (!path) {
        this.file = null
        this.data = null
        return
      }
      this.loading = true
      const id = this.$route.params.id
      const fileInfo = await this.axios.get(`${this.$serverAddr}/pcfg/${id}/file`, {
        params: {path}
      }).then(r => r.data)
      this.file = fileInfo.path
      this.data = fileInfo.data
      this.loading = false
    }
  }
}
</script>

<style scoped>
.monospace {
  font-family: Roboto Mono;
}

.wrap {
  white-space: pre-wrap;
}

.grow-greedy {
  flex-grow: 2;
}
</style>