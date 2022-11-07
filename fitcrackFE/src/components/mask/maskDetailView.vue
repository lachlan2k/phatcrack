<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <v-breadcrumbs
      v-if="data != null"
      :items="
        [
          { text: 'Masks', to: { name: 'masks' }, exact: true },
          { text: data.name }
        ]"
      divider="/"
    />

    <v-container>
        <fc-tile
          title="Mask set"
          :loading="data==null"
          class="mx-2 dictContentContainer mb-4"
        >
          <v-list
            v-if="data != null"
          >
            <v-list-item>
              <v-list-item-action>
                Name:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ data.name }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
            <v-divider />
            <v-list-item>
              <v-list-item-action>
                Added:
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title class="text-right">
                  {{ $moment.utc(data.time).local().format('DD.MM.YYYY HH:mm') }}
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
          <v-row class="mx-2 mb-2">
            <v-col>
              <v-btn
                color="primary"
                :href="$serverAddr + '/masks/' + data.id + '/download'"
                target="_blank"
                class="ml-2"
              >
                Download <v-icon right>mdi-download</v-icon>
              </v-btn>
            </v-col>
          </v-row>
          <v-divider />
          <div
            v-if="data != null"
            class="dictContent pa-2"
            :class="{editing: edit}"
          >
            <code
              class="code elevation-0"
            >{{ data.data }}</code>
          </div>
        </fc-tile>
    </v-container>
  </div>
</template>

<script>
  import tile from '@/components/tile/fc_tile.vue'
  export default {
    name: "MaskDetailView",
    components: {
      'fc-tile': tile
    },
    data: function () {
      return {
        data: null,
        newData: ''
      }
    },
    mounted: function () {
      this.loadData()
    },
    methods: {
      loadData: function ($state) {
        this.axios.get(this.$serverAddr + '/masks/' + this.$route.params.id).then((response) => {
          this.data = response.data
        });
      }
    }
  }
</script>

<style scoped>
  .dictContentContainer {
    max-width: 800px;
  }
  .dictContent {
    overflow: auto;
  }

  .dictContent.editing {
    border: 1px solid;
  }

  .code::before {
    display: none;
  }
  .code {
    width: 100%;
    background: transparent;
    color: #000;
    white-space: pre-wrap;
  }
  .width100 {
    width: 100%;
  }

  .margintop5 {
    margin-top: 10px;
  }
</style>
