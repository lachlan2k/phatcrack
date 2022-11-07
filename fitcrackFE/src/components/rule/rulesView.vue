<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container class="max700">
    <fc-tile
      title="Rules"
      class="ma-2"
      :icon="$route.meta.icon"
    >
      <v-alert
        tile
        text
        type="warning"
        class="mb-0"
      >
        Rule files must have a .txt or .rule extension.
      </v-alert>
      <v-data-table
        :headers="headers"
        :items="rules.items"
        :loading="loading"
        :footer-props="{itemsPerPageOptions: [10,25,50], itemsPerPageText: 'Rules per page'}"
      >
        <template v-slot:item.name="{ item }">
          <router-link :to="`rules/${item.id}`">
            {{ item.name }}
          </router-link>
        </template>
        <template v-slot:item.count="{ item }">
          {{ fmt(item.count) }}
        </template>
        <template v-slot:item.time="{ item }">
          {{ $moment.utc(item.time).local().format('DD.MM.YYYY HH:mm') }}
        </template>
        <template v-slot:item.actions="{ item }">
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <a
                :href="$serverAddr + '/rule/' + item.id + '/download'"
                target="_blank"
                download
                v-on="on"
              >
                <v-btn icon>
                  <v-icon>mdi-file-download-outline</v-icon>
                </v-btn>
              </a>
            </template>
            <span>Download</span>
          </v-tooltip>
          <v-tooltip top>
            <template v-slot:activator="{ on }">
              <v-btn
                icon
                @click="deleteRule(item)"
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
      <v-divider />
      <file-uploader
        :url="this.$serverAddr + '/rule'"
        @uploadComplete="loadRules"
      />
    </fc-tile>
  </v-container>
</template>

<script>
  import fmt from '@/assets/scripts/numberFormat'
  import tile from '@/components/tile/fc_tile.vue'
  import FileUploader from "@/components/fileUploader/fileUploader.vue";
  export default {
    name: "RulesView",
    components: {
      FileUploader,
      'fc-tile': tile,
    },
    data: function () {
      return {
        loading: false,
        rules: [],
        headers: [
          {
            text: 'Name',
            align: 'start',
            value: 'name'
          },
          {text: 'Count', value: 'count', align: 'end'},
          {text: 'Added', value: 'time', align: 'end'},
          {text: 'Actions', value: 'actions', align: 'end', sortable: false}
        ]
      }
    },
    mounted: function () {
      this.loadRules()
    },
    methods: {
      fmt,
      loadRules: function () {
        this.loading = true;
        this.axios.get(this.$serverAddr + '/rule', {}).then((response) => {
          this.rules = response.data;
          this.loading = false
        })
      },
      deleteRule: function (item) {
        this.$root.$confirm('Delete', `This will remove ${item.name} from your rules. Are you sure?`).then((confirm) => {
          this.loading = true;
          this.axios.delete(this.$serverAddr + '/rule/' + item.id).then((response) => {
            this.loadRules()
          })
        })
      }
    }
  }
</script>

<style scoped>
  .noEvent {
    pointer-events: none;
    display: inline-block;
  }

  .dz-message {
    cursor: pointer;
    text-align: center;
  }

  .max700 {
    max-width: 700px;
  }

</style>


<style>
  .selectedDict {
    background: rgba(37, 157, 173, 0.85) !important;
    color: white;
  }

  .selectedDict a {
    color: white;
  }

  .clickable {
    cursor: pointer;
  }
</style>
