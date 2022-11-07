<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-app>
    <Connection 
      :visible="connectionAssistant" 
      @connect="handleReconnection"
    />
    <v-snackbar
      v-model="alert"
      :timeout="6000"
      bottom
      right
      color="transparent"
      class="errorSnackbar pa-0"
    >
      <v-alert
        v-model="alert"
        :type="alertType"
        class="height100 ma-0 width100"
        :value="true"
        dismissible
        transition="none"
      >
        {{ alertText }}
      </v-alert>
    </v-snackbar>
    <router-view 
      :is-dark="isDark" 
      @alert="setAlert"
    />
    <confirm ref="confirm" />
  </v-app>
</template>

<script>
  import confirm from "@/components/confirmDialog/confirm.vue";
  import Connection from "./components/errorPages/connection.vue";
  export default {
    name: 'App',
    components: {
    "confirm": confirm,
    Connection
},
    data: function () {
      return {
        alert: false,
        alertText: '',
        alertType: 'error',
        isDark: false,
        connectionAssistant: false
      }
    },
    created: function () {
      /* ...
      this.$store.state.loggedInLink = this.$route;
      if (this.$needLogin) {
        this.$router.push({
          name: 'login'
        })
      }
      */
     
      this.axios.interceptors.request.use(
        function (request) {
          const token = localStorage.getItem('jwt')
          if (token && token !== '') {
            request.headers.common.Authorization = `Bearer ${token}`
          }
          return request
        }.bind(this),
        function (error) {
          return Promise.reject(error)
        }.bind(this)
      )
       
      this.axios.interceptors.response.use(
        function (response) {
          if (response.data.hasOwnProperty('hideAlert') && response.data.hideAlert) {
            this.alert = false;
          }
          if (response.data.hasOwnProperty('status') && response.data.hasOwnProperty('message') && response.data.status) {
            this.alertType = 'success';
            this.alert = true;
            this.alertText = response.data.message
          }
          return response
        }.bind(this),
        function (error) {
          this.alertType = 'error';
          this.alert = true;
          if (error.response !== undefined) {
            if (error.response.data.message === 'Input payload validation failed') {
              let property = Object.keys(error.response.data.errors)[0]
              this.alertText = 'Wrong settings. ' +  property + ' (' + error.response.data.errors[property] + ')'
            } else {
              this.alertText = error.response.data.message
            }
          } else if (error.message !== undefined) {
            if (error.message === 'Network Error') {
              this.alert = false
              this.connectionAssistant = true
            } else {
              this.alertText = error.message
            }
          } else {
            this.alertText = 'Can not connect to ' + this.$serverAddr
            //this.alertText = error.message
          }
          return Promise.reject(error)
        }.bind(this)
      )

      const mql = window.matchMedia('(prefers-color-scheme: dark)')
      mql.addListener(e => {
        const pref = localStorage.getItem('appearance')
        if (!pref || pref == 'auto') {
          this.$vuetify.theme.dark = e.matches
        }
      })
      const init = localStorage.getItem('appearance')
      if (init && init == 'dark') {
        this.$vuetify.theme.dark = true
      } else if (!init || init == 'auto') {
        this.$vuetify.theme.dark = mql.matches
      }

      window.addEventListener('keypress', e => {
        if (localStorage.getItem('testmode') == 'true' && e.ctrlKey && e.shiftKey && e.code == 'KeyL') {
          this.$vuetify.theme.dark = !this.$vuetify.theme.dark
        }
      })
    },
    mounted () {
      this.$root.$confirm = this.$refs.confirm.open
    },
    methods: {
      setAlert ({type, text}) {
        this.alertType = type
        this.alertText = text
        this.alert = true
      },
      handleReconnection () {
        this.connectionAssistant = false
        this.$router.go() // .replace({ name: 'login' })
      }
    }
  }
</script>
<style scoped>

  .errorSnackbar {
    overflow: visible;
    padding: 0;
    z-index: 6;
  }

</style>
<style>
  .theme--light.v-application {
    background-color: #f5f5f5 !important;
  }

  h1, h2 {
    font-weight: normal;
  }

  .errorSnackbar .v-snack__content {
    padding: 0;
  }

  a {
    text-decoration: none;
  }

  .width100 {
    width: 100%;
  }

  .table-link {
    display: inline-block;
    line-height: 3em;
    font-weight: bold;
  }

  .drag-handle {
    cursor: grab;
  }
  .drag-handle:active {
    cursor: grabbing;
  }

  /**/

  .list-enter-active,
  .list-leave-active,
  .list-move {
    transition: 200ms cubic-bezier(0.59, 0.12, 0.34, 0.95);
    transition-property: opacity, transform;
  }

  .list-enter,
  .list-leave-to {
    opacity: 0;
    transform: translateX(-10px);
  }

  .list-enter-to {
    opacity: 1;
    transform: translateX(0);
  }

  .list-leave-active {
    position: absolute;
  }

</style>

