<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div
    ref="logsCont"
    class="logsCont"
    infinite-wrapper
  >
    <fc-log
      v-for="log, i in logs"
      :key="log.raw"
      :type="log.type"
      :time="log.time"
      :host-i-d="log.host"
      :job-i-d="log.job"
    >
      {{ log.log }}
    </fc-log>
    <infinite-loading @infinite="loadLogs">
      <v-progress-circular
        slot="spinner"
        size="50"
        :width="3"
        indeterminate
        color="primary"
      />
      <span slot="no-more">
        You reached the end of file.
      </span>
      <span slot="no-results">
        You reached the end of file.
      </span>
    </infinite-loading>
  </div>
</template>

<script>
  import logItem from '@/components/logs/logItem.vue'
  import InfiniteLoading from 'vue-infinite-loading';
  export default {
    name: "LogsView",
    components: {
      'fc-log': logItem,
      InfiniteLoading
    },
    data() {
      return {
        height: 0,
        logs: [],
        skipCount: 0,
        lastLog: null,
        interval: null
      }
    },
    mounted () {
      this.interval = setInterval(function () {
        this.loadNewLogs();
      }.bind(this), 2500);
    },
    beforeDestroy: function(){
      clearInterval(this.interval);
    },
    methods: {
      loadLogs($state) {
        this.axios.get(this.$serverAddr + '/log', {
          params: {
            'skip_count': this.skipCount,
            'last_log': this.lastLog,
          }
        }).then((response) => {
          this.logs = this.logs.concat(response.data.items)
          this.lastLog = response.data.items[response.data.items.length - 1].raw
          this.skipCount += response.data.items.length
          $state.loaded()
        }).catch((response) => {
          $state.complete()
        })
      },
      loadNewLogs($state) {
        if (this.logs.length === 0)
          return
        this.axios.get(this.$serverAddr + '/log/new', {
          params: {
            'log': this.logs[0].raw
          }
        }).then((response) => {
          this.logs = response.data.items.concat(this.logs)
          this.skipCount += response.data.items.length
          if (this.$refs.logsCont.scrollTop > 100) {
            this.$nextTick(() => {
              this.$refs.logsCont.scrollTop += response.data.items.length * 52
            })
          }
        })
      }
    }
  }
</script>

<style scoped>
  .logsCont {
    height: 100%;
    overflow: auto;
    max-height: 500px;
  }



</style>
