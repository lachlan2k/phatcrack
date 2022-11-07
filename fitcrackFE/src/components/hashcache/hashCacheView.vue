<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div class="cont">
    <v-card-title class="py-2">
      <v-text-field
        v-model="search"
        clearable
        flat
        outlined
        prepend-inner-icon="mdi-database-search"
        label="Search by password or hash"
        single-line
        hide-details
      />
    </v-card-title>
    <v-divider />
    <v-data-table
      ref="table"
      :headers="headers"
      :items="hashes"
      :search="search"
      :options.sync="pagination"
      :server-items-length="totalItems"
      :loading="loading"
      :footer-props="{itemsPerPageOptions: [25,50,100], itemsPerPageText: 'Passwords per page'}"
    >
      <template v-slot:item.password="{ item }">
        {{ item.password || '–' }}
      </template>
      <template v-slot:item.hash="{ item }">
        <div class="d-flex align-center mw-50">
          <v-dialog max-width="600">
            <template v-slot:activator="{ on }">
              <span 
                class="text-truncate d-inline-block"
                v-on="on"
              >
                {{ item.hash }}
              </span>
            </template>
            <v-card>
              <v-card-title primary-title>
                {{ item.hash_type_name }}
              </v-card-title>
              <v-card-text>
                {{ item.hash }}
              </v-card-text>
            </v-card>
          </v-dialog>
          <!-- <svg 
            xmlns="http://www.w3.org/2000/svg"
            height="18"
            width="100%"
          >
            <rect
              v-for="(char, i) in item.hash"
              :key="i"
              :x="i"
              :y="18 + 5 - char.charCodeAt(0) / 5"
              width="1"
              :height="char.charCodeAt(0) / 5 - 5"
              :fill="'hsla(' + (char.charCodeAt(0) * 5 - 100) + ',100%,50%,50%)'"
            />
          </svg> -->
        </div>
      </template>
      <template v-slot:item.added="{ item }">
        {{ $moment.utc(item.added).local().format('D.M.YYYY H:mm:ss') }}
      </template>
      <template v-slot:item.job="{ item }">
      <router-link
        :to="{ name: 'jobDetail', params: { id: item.job.id } }"
        class="middle"
      >
        {{ item.job.name }}
      </router-link>
      </template>
    </v-data-table>
  </div>
</template>

<script>
  export default {
    name: "HashesView",
    data: function () {
      return {
        interval: null,
        status: 'active',
        search: '',
        totalItems: 0,
        pagination: {},
        loading: true,
        headers: [
          {text: 'Password', value: 'password', align: 'start', sortable: true},
          {text: 'Hash type', value: 'hash_type_name', align: 'start', sortable: true},
          {text: 'Hash', value: 'hash', align: 'start', sortable: true},
          {text: 'Added', value: 'added', align: 'start', sortable: true},
          {text: 'Job', value: 'job', align: 'end', sortable: true}
        ],
        hashes:
          []
      }
    },
    watch: {
      pagination: {
        handler() {
          this.loading = true;
          this.loadHashes()
        },
        deep: true
      }
    },
    mounted() {
      this.interval = setInterval(function () {
        this.loadHashes()
      }.bind(this), 5000)
    },
    beforeDestroy: function () {
      clearInterval(this.interval)
    },
    methods: {
      loadHashes() {
        this.axios.get(this.$serverAddr + '/hashes', {
          params: {
            'page': this.pagination.page,
            'per_page': this.pagination.rowsPerPage,
            'order_by': this.pagination.sortBy,
            'descending': this.pagination.descending,
            'search': this.search
          }
        }).then((response) => {
          this.hashes = response.data.items;
          this.totalItems = response.data.total;
          this.loading = false
        })
      },
      updateStatus: function (e) {
        this.status = e;
        this.$refs.table.updatePagination({page: 1, totalItems: this.totalItems})
      }
    }
  }
</script>

<style scoped>

  .mw-50 {
    max-width: 50vw;
  }

  .cont {
    height: 100%;

  }

</style>
