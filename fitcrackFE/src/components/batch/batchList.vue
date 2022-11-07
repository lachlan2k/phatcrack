<template>
  <div>
    <v-data-table
      ref="table"
      :loading="loading"
      :headers="headers"
      :items="batches"
      :options.sync="pagination"
      :server-items-length="totalItems"
      :footer-props="{itemsPerPageOptions: [25,10,50,100], itemsPerPageText:'Batches per page'}"
      fixed-header
      disable-sort
    >
      <template v-slot:item.name="{ item }">
        <router-link
          :to="`/batches/${item.id}`"
          class="table-link"
        >
          {{ item.name }}
        </router-link>
      </template>
      <template v-slot:item.status="{ item }">
        <v-tooltip top>
          <template v-slot:activator="{ on }">
            <span
              :class="[`${item.status.color}--text`, 'text-capitalize']"
              v-on="on"
            >
              {{ item.status.text }}
              <v-icon
                right
                :color="item.status.color"
              >
                {{ item.status.icon }}
              </v-icon>
            </span>
          </template>
          <span>{{ item.status.description }}</span>
        </v-tooltip>
      </template>
      <template v-slot:item.total_jobs="{ item }">
        <span><strong>{{ item.total_jobs }}</strong> jobs</span>
      </template>
      <template v-slot:item.waiting_jobs="{ item }">
        <v-progress-circular
          size="18"
          :width="3"
          :rotate="270"
          :value="100 / item.total_jobs * (item.total_jobs - item.waiting_jobs)"
        />
        <span class="ml-2"><strong>{{ item.waiting_jobs }}</strong> in queue</span>
      </template>
      <template v-slot:item.actions="{ item }">
        <div class="actionsBtns">
          <!-- opration buttons -->
          <template v-if="item.current_user_can_operate">
            <v-btn
              v-if="item.status.code < 10"
              :disabled="item.status.code == 1"
              text
              color="success"
              @click="runBatch(item.id)"
            >
              Run
              <v-icon right>
                mdi-play
              </v-icon>
            </v-btn>
            <v-btn
              v-else-if="item.status.code < 20"
              text
              color="error"
              @click="interruptBatch(item.id)"
            >
              Interrupt
              <v-icon right>
                mdi-pause
              </v-icon>
            </v-btn>
            <v-btn
              v-else
              text
              color="error"
              @click="runBatch(item.id)"
            >
              Resume
              <v-icon right>
                mdi-skip-next
              </v-icon>
            </v-btn>
          </template>
          <!-- delete button -->
          <v-btn
            v-if="item.current_user_can_edit"
            text
            @click="deleteBatchConfirm(item.id)"
          >
            Unlink
            <v-icon right>
              mdi-link-variant-off
            </v-icon>
          </v-btn>
        </div>
      </template>
    </v-data-table>
  </div>
</template>

<script>
export default {
  data () {
    return {
      interval: null,
      search: '',
      totalItems: 0,
      pagination: {},
      loading: true,
      batches: [],
      //
      headers: [
        {
          text: 'Name',
          align: 'start',
          value: 'name'
        },
        {text: 'Status', value: 'status', align: 'end'},
        {text: 'Length', value: 'total_jobs', align: 'end'},
        {text: 'Progress', value: 'waiting_jobs', align: 'end'},
        {text: 'Actions', value: 'actions', align: 'end'}
      ]
    }
  },
  mounted () {
    this.load()
    this.interval = setInterval(() => {
      this.load()
    }, 5000)
  },
  beforeDestroy: function () {
    clearInterval(this.interval)
  },
  methods: {
    async load () {
      const data = await this.axios.get(`${this.$serverAddr}/batches`, {
        params: {
          'page': this.pagination.page,
          'per_page': this.pagination.itemsPerPage
        }
      }).then(r => r.data)
      this.batches = data.items
      this.totalItems = data.total
      this.loading = false
    },
    async runBatch (id) {
      this.loading = true
      try {
        await this.axios.post(`${this.$serverAddr}/batches/${id}/run`)
        this.load()
      } catch (e) {
        console.error(e)
      }
    },
    async interruptBatch (id) {
      this.loading = true
      try {
        await this.axios.post(`${this.$serverAddr}/batches/${id}/interrupt`)
        this.load()
      } catch (e) {
        console.error(e)
      }
    },
    deleteBatchConfirm (id) {
      this.$root.$confirm('Unlink Batch', `This will unlink jobs from this batch and remove the batch. Jobs will not be discarded. If a job is currently running, it will finish normally. Are you sure?`)
      .then(async confirm => {
        this.loading = true
        try {
          await this.axios.delete(`${this.$serverAddr}/batches/${id}`)
          this.load()
        } catch (e) {
          console.error(e)
        }
      })
    }
  }
}
</script>