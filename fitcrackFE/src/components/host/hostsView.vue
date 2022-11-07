<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="cont">
    <div class="d-flex justify-space-between align-center px-4 pt-2">
      <v-switch
        v-model="viewHidden"
        label="View hidden hosts"
        :prepend-icon="viewHidden ? 'mdi-eye-off' : 'mdi-eye'"
        class="mr-4"
      />
      <v-text-field
        v-model="search"
        class="px-2 pt-3"
        clearable
        prepend-inner-icon="mdi-table-search"
        label="Search"
        single-line
        @input="updateList"
      />
    </div>
    <v-divider />
    <v-data-table
      ref="table"
      :headers="headers"
      :items="hosts"
      :search="search"
      :options.sync="pagination"
      :server-items-length="totalItems"
      :loading="loading"
      :footer-props="{itemsPerPageOptions: [25,50,100], itemsPerPageText: 'Hosts per page'}"
      fixed-header
    >
      <template v-slot:item.domain_name="{ item }">
        <router-link
          :to="{ name: 'hostDetail', params: {id: item.id} }"
          class="middle font-weight-bold"
        >
          {{ item.domain_name + ' (' + fixUserNameEncoding(item.user.name) + ')' }}
        </router-link>
      </template>
      <template v-slot:item.jobs="{ item }">
        {{ item.jobs.map(j => j.status === 10 ? 1 : 0).reduce((a, b) => a + b, 0) }}
      </template>
      <template v-slot:item.p_model="{ item }">
        <span class="oneline">{{ item.p_model.replace(/(?:\(R\)|\(TM\)|Intel|AMD)/g, '') }}</span>
      </template>
      <template v-slot:item.last_active="{ item }">
        <span v-if="item.last_active.seconds_delta > 61">{{ parseTimeDelta(item.last_active.last_seen) }}</span>
        <v-icon
          v-else
          color="success"
        >
          mdi-power
        </v-icon>
      </template>
      <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                class="mx-0"
                v-on="on"
                @click="unassignAllJobs(item.id)"
              >
                <v-icon>
                  {{ 'mdi-lan-disconnect' }}
                </v-icon>
              </v-btn>
            </template>
          <span>Unassign from all jobs</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                class="mx-0"
                v-on="on"
                @click="hideHost(item.id)"
              >
                <v-icon>
                  {{ item.deleted ? 'mdi-eye' : 'mdi-eye-off' }}
                </v-icon>
              </v-btn>
            </template>
          <span>{{ item.deleted ? 'Show' : 'Hide' }}</span>
          </v-tooltip>
      </template>
    </v-data-table>
  </div>
</template>

<script>
  import iconv from 'iconv-lite'
  export default {
    name: "HostsView",
    data: function () {
      return {
        interval: null,
        search: '',
        viewHidden: false,
        totalItems: 0,
        pagination: {},
        loading: true,
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'domain_name'
          },
          {text: 'IP address', value: 'ip_address', align: 'end', sortable: false},
          {text: 'Operating system', value: 'os_name', align: 'end', sortable: false},
          {text: 'Processor', value: 'p_model', align: 'end', width: '200', sortable: false},
          {text: 'Active jobs', value: 'jobs', align: 'center', sortable: false},
          {text: 'Online', value: 'last_active', align: 'center', sortable: false},
          {text: 'Actions', value: 'actions', align: 'center', sortable: false}
        ],
        hosts_statuses: [],
        hosts:
          []
      }
    },
    watch: {
      pagination: {
        handler: 'updateList',
        deep: true
      },
      '$route.name': 'updateList',
      viewHidden: 'updateList'
    },
    mounted() {
      this.interval = setInterval(function () {
        this.loadHosts()
      }.bind(this), 5000)
    },
    beforeDestroy: function () {
      clearInterval(this.interval)
    },
    methods: {
      parseTimeDelta: function (delta) {
        if (delta !== null && delta !== undefined) {
          return this.$moment.utc(delta).fromNow()
        } else {
          return 'Unknown'
        }
      },
      loadHosts() {
        this.axios.get(this.$serverAddr + '/hosts', {
          params: {
            'page': this.pagination.page,
            'per_page': this.pagination.itemsPerPage,
            'order_by': this.pagination.sortBy[0],
            'descending': this.pagination.sortDesc[0],
            'name': this.search,
            'showDeleted': this.viewHidden
          }
        }).then((response) => {
          this.hosts = response.data.items;
          this.totalItems = response.data.total;
          this.loading = false
        })
      },
      updateList () {
        this.loading = true
        this.loadHosts()
      },
      hideHost: function (id) {
        this.loading = true
        this.axios.delete(this.$serverAddr + '/hosts/' + id)
          .then((response) => {
            this.loadHosts()
          })
      },
      unassignAllJobs: function (id) {
        this.axios.put(this.$serverAddr + '/hosts/' + id + "/unassignAllJobs")
        .then((response) => {
            this.loadHosts()
          })
      },
      fixUserNameEncoding : function(username) {
          /* Boinc DB uses latin1_swedish encoding, which breaks names with special characters,
          which are not supported in this encoding. Fix it by converting name to utf8. */
          return iconv.decode(iconv.encode(username, 'latin1'), 'utf-8')
      }
    }
  }
</script>

<style scoped>

  .middle {
    vertical-align: middle;
  }

  .cont {
    height: 100%;

  }

  .oneline {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: block;
    width: 200px;
    line-height: 50px;
    height: 50px;
  }

  .inheritColor {
    color: inherit !important;
  }

</style>
