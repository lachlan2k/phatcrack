<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div
    :class="classes"
    class="logItem"
  >
    <v-row
      align="center"
      
      :class="{new: isNew}"
      class="newOverlay px-2"
    >
      <span class="time pr-2">{{ $moment.utc(time).local().format('D.M.YYYY H:mm:ss') }}: </span>
      <v-col><slot /></v-col>
      <v-btn
        v-if="hostID"
        :to="'hosts/' + hostID"
        icon
        class="btns"
      >
        <v-icon>desktop_windows</v-icon>
      </v-btn>
      <v-btn
        v-if="jobID"
        :to="'jobs/' + jobID"
        icon
        class="btns"
      >
        <v-icon>work</v-icon>
      </v-btn>
    </v-row>
  </div>
</template>

<script>
  export default {
    name: "LogItem",
    props: {
      type: {
        type: String,
        default: 'LOG'
      },
      hostID: {
        type: [Number, null],
        default: null
      },
      jobID: {
        type: [Number, null],
        default: null
      },
      time: {
        type: String,
        default: ''
      },
    },
    data() {
      return {
        isNew: true
      }
    },
    computed: {
      classes () {
        switch (this.type) {
          case 'LOG': return 'infolog'
          case 'ERROR': return 'error'
          case 'WARN': return 'warning'
        }
      }
    },
    mounted () {
      setTimeout(function(){ this.isNew = false }.bind(this), 1000);
    },
  }
</script>

<style scoped>
  .logItem {
    border-width: 4px 0 0 0;
    border-style: solid;
    border-color: rgba(0,0,0,0.12) !important;
    font-family: monospace;
  }
  .error,
  .warning {
    color: white;
  }

  .icon {
    color: rgba(0,0,0,0.3);
    font-size: 24px;
  }

  .time {
    font-weight: bold;
  }

  .btns {

  }

  .newOverlay {
    background-color: rgba(0,0,0,0);
    -webkit-transition: background-color 3000ms linear;
    -ms-transition: background-color 3000ms linear;
    transition: background-color 3000ms linear;
  }

  .new {
    background-color: rgba(255, 246, 10, 0.5);
  }
</style>
