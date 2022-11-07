<!--
  * Author : see AUTHORS
  * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max600">
    <fc-tile
      title="Job Templates"
      class="ma-2"
      :icon="$route.meta.icon"
    >
      <v-data-table
        :headers="headers"
        :items="templates.items"
        :loading="loading"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Templates per page'}"
      >
        <template v-slot:item.created="{ item }">
          {{ $moment.utc(item.created).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                @click="deleteTemplate(item)"
                v-on="on"
              >
                <v-icon color="error">
                  mdi-delete-outline
                </v-icon>
              </v-btn>
            </template>
            <span>Delete</span>
          </v-tooltip>
        </template>
      </v-data-table>
    </fc-tile>
  </v-container>
</template>

<script>
  import FcTile from '@/components/tile/fc_tile.vue'
  export default {
    components: {
      FcTile
    },
    data () {
      return {
        loading: false,
        templates: [],
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Added', value: 'created', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ]
      }
    },
    mounted () {
      this.load()
    },
    methods: {
      load () {
        this.loading = true
        this.axios.get(this.$serverAddr + '/template', {}).then((response) => {
          this.templates = response.data
          this.loading = false
        })
      },
      deleteTemplate ({ id }) {
        this.loading = true
        this.axios.delete(this.$serverAddr + '/template/' + id).then((response) => {
          this.load()
        })
      }
    }
  }
</script>

<style scoped>
  .max600 {
    max-width: 600px;
  }
</style>