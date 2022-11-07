<!--
  * Author : see AUTHORS
  * Licence: MIT, see LICENSE
-->

<template>
  <div class="height100">
    <v-navigation-drawer
      id="navigation-drawer"
      v-model="drawer"
      persistent
      class="navigationDrawer"
      enable-resize-watcher
      fixed
      app
      width="290"
    >
      <router-link :to="{ name: 'home'}">
        <div class="logo">
          <img
            v-if="$vuetify.theme.dark"
            :src="require('@/assets/logo-dark.svg')"
            alt="logo"
          >
          <img
            v-else
            :src="require('@/assets/logo-light.svg')"
            alt="logo"
          >
        </div>
      </router-link>

      <div class="dash-link">
        <v-btn
          id="dash-link"
          to="/"
          text
        >
          <v-icon left>
            mdi-view-dashboard
          </v-icon>
          Dashboard
        </v-btn>
      </div>
        

      <v-divider />

      <v-tabs
        v-model="navtab"
        icons-and-text
        grow
        class="notrans"
      >
        <v-tab id="nav-jobs-tab">
          Jobs
          <v-icon>mdi-briefcase-outline</v-icon>
        </v-tab>
        <v-tab id="nav-library-tab">
          Library
          <v-icon>mdi-folder-outline</v-icon>
        </v-tab>
        <v-tab id="nav-system-tab">
          System
          <v-icon>mdi-settings-outline</v-icon>
        </v-tab>
      </v-tabs>

      <v-tabs-items 
        v-model="navtab"
        class="notrans"
      >
        <v-tab-item eager>
          <v-list
            id="nav-jobs-list"
            expand
            nav
            dense
          >
            <v-list-item
              v-if="$userCan('ADD_NEW_JOB')"
              :to="{ name: 'addJob'}"
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('addJob') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Add Job</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item 
              :to="{ name: 'batches'}"
              exact
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('batches') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Batches</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item 
              :to="{ name: 'jobs'}"
              exact
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('jobs') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>{{ $userCan('VIEW_ALL_JOBS') ? 'All' : 'My' }} Jobs</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-divider class="mb-1" />

            <bins id="job-bin-list" />

            <v-divider class="my-1" />

            <v-list-item 
              to="/bins/trash"
              exact
            >
              <v-list-item-action>
                <v-icon>mdi-delete</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Trash</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-tab-item>
        <v-tab-item>
          <v-list
            id="nav-library-list"
            expand
            nav
            dense
          >
            <v-list-item :to="{ name: 'hashes'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('hashes') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Hashes</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'dictionaries'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('dictionaries') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Dictionaries</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'pcfg'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('pcfg') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>PCFG</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'rules'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('rules') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Rules</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'charsets'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('charsets') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Charsets</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'masks'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('masks') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Masks</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'markovChains'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('markovChains') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Markov Chains</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'files'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('files') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Encrypted Files</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-tab-item>
        <v-tab-item>
          <v-list
            id="nav-system-list"
            expand
            nav
            dense
          >
            <v-list-item :to="{ name: 'hosts'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('hosts') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Hosts</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'templates'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('templates') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Job Templates</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item
              v-if="$userCan('MANAGE_USERS')"
              :to="{ name: 'manageUsers'}"
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('manageUsers') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Manage Users</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item
              v-if="$userCan('MANAGE_USERS')"
              :to="{ name: 'server'}"
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('server') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Server Monitor</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item :to="{ name: 'settings'}">
              <v-list-item-action>
                <v-icon>{{ routeIcon('settings') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Settings</v-list-item-title>
              </v-list-item-content>
            </v-list-item>

            <v-list-item
              :to="{ name: 'transfer'}"
            >
              <v-list-item-action>
                <v-icon>{{ routeIcon('transfer') }}</v-icon>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>Data Transfer</v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </v-list>
        </v-tab-item>
      </v-tabs-items>
      <!-- <v-btn
        icon
        absolute
        bottom
        left
      >
        <v-icon>mdi-lifebuoy</v-icon>
      </v-btn> -->
    </v-navigation-drawer>

    <!-- TOP BAR -->

    <v-app-bar
      app
      height="64px"
      :class="['mainToolbar', {'very-important-must-not-miss': daemonWarning}]"
      :color="daemonWarning ? 'error' : ''"
      :dark="daemonWarning"
    >
      <v-app-bar-nav-icon id="nav-toggle" @click.stop="drawer = !drawer" />
      <v-toolbar-title>{{ $store.state.project }}</v-toolbar-title>
      <v-spacer />
      <v-row
        v-if="daemonWarning"
        class="align-center"
      >
        <v-col class="text-center">
          <v-icon
            large
            left
            class="important-must-not-miss"
          >
            mdi-alert-octagon
          </v-icon>
          <span class="mr-4">
            Some Daemons are not running!
          </span>
          <v-btn
            light
            @click="daemonDetails = true"
          >
            Details
          </v-btn>
          <v-btn
            icon
            @click="daemonWarning = false"
          >
            <v-icon>
              mdi-close-circle
            </v-icon>
          </v-btn>
        </v-col>
      </v-row>
      <v-spacer />
      <v-badge
        :value="notificationsCount > 0"
        :content="notificationsCount"
        color="red"
        overlap
        left
        top
        offset-y="18"
      >
        <v-btn
          id="notifications-button"
          large
          rounded
          icon
          @click.stop="toggleNotifications"
        >
          <v-icon>
            {{ notificationsCount > 0 ? 'mdi-bell' : 'mdi-bell-outline' }}
          </v-icon>
        </v-btn>
      </v-badge>
      <v-menu
        v-model="userFlyout"
        :close-on-content-click="false"
        transition="slide-y-transition"
      >
        <template #activator="{ on }">
          <v-btn
            id="user-button"
            large
            rounded
            text
            class="pr-2 pl-3"
            v-on="on"
          >
            <div class="text-none text-end subtitle-2 mr-2">
              {{ user.username }}
            </div>
            <v-avatar
              :color="daemonWarning ? 'secondary' : 'primary'"
              size="32"
            >
              <span class="white--text">
                {{ user.username.split(' ').map(p => p[0]).join('') }}
              </span>
            </v-avatar>
          </v-btn>
        </template>
        <v-card id="user-flyout" min-width="250">
          <v-card-text class="text--primary d-flex align-center">
            <div>
              <div class="subtitle-2">
                {{ user.username }}
              </div>
              <div class="caption">
                {{ user.mail }}
              </div>
            </div>
            <v-spacer />
            <v-btn
              text
              class="ml-4"
              @click="roles = !roles"
            >
              {{ user.role.name }}
              <v-icon right>
                {{ roles ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
              </v-icon>
            </v-btn>
          </v-card-text>
          <v-expand-transition>
            <div v-show="roles">
              <v-card-text class="py-0">
                <div class="subtitle-1">
                  Roles
                </div>
                <v-list dense>
                  <v-list-item
                    v-for="role in roleList"
                    :key="role"
                  >
                    <v-list-item-title class="text-capitalize">
                      {{ role }}
                    </v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-card-text>
            </div>
          </v-expand-transition>
          <v-card-actions>
            <v-btn
              text
              :to="{ name: 'myAccount' }"
              class="mr-2"
              @click="userFlyout = false"
            >
              <v-icon left>
                {{ routeIcon('myAccount') }}
              </v-icon>
              <span>My Account</span>
            </v-btn>
            <v-spacer />
            <v-btn
              text
              class="ml-2"
              @click.stop="$store.dispatch('signOut')"
            >
              <span>Sign out</span>
              <v-icon right>
                mdi-logout-variant
              </v-icon>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-menu>
    </v-app-bar>
    <!-- CONTENT -->
    <v-main class="height100 main">
      <transition
        name="route"
        mode="out-in"
      >
        <router-view />
      </transition>
    </v-main>
    <!-- NOTIFICATION DRAWER -->
    <v-navigation-drawer
      id="notifications-drawer"
      v-model="rightDrawer"
      class="pa-0"
      temporary
      right
      fixed
      app
      width="450"
    >
      <notifications-wrapper
        ref="notifWrapper"
        :count="notificationsCount"
        @close="rightDrawer = false"
      />
    </v-navigation-drawer>
    <!-- D A E M O N (。>︿<) A L E R T -->
    <v-dialog
      v-model="daemonDetails"
      overlay-color="red darken-4"
      :overlay-opacity=".8"
      max-width="550"
    >
      <v-card>
        <v-card-title>
          Daemons not running
        </v-card-title>
        <v-card-text class="pb-0">
          The server daemons are services that carry out regular routines and <strong>are vital</strong>
          for the inner workings of Fitcrack. It seems that, unfortunately, some have crashed. Here's a list of inactive daemons:
        </v-card-text>
        <v-card-text class="py-0">
          <v-list dense>
            <v-list-item
              v-for="daemon in fallenDaemons"
              :key="daemon"
            >              
              <v-list-item-title>
                {{ daemon }}
              </v-list-item-title>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-text>
          <v-sheet
            dark
            color="error"
            class="pa-4 mb-2"
          >
            Follow these steps to summon them again:
            <ol>
              <li>Access your server's shell (ex. via ssh)</li>
              <li>
                Run <code>/home/<em>&lt;BOINC user&gt;</em>/projects/<em>&lt;project name&gt;</em>/bin/start</code>
                as BOINC. BOINC user is <code>boincadm</code> by default.
              </li>
            </ol>
          </v-sheet>
          <span>
            After starting successfully, you can dismiss the warning.
          </span>
        </v-card-text>
      </v-card>
    </v-dialog>
    <!-- TRANSFER SHORTCUT -->
    <v-fab-transition>
      <v-btn
        v-if="transferCount > 0 && $route.name !== 'transfer'"
        fab
        fixed
        bottom
        left
        class="z-top"
        color="primary"
        :to="{ name: 'transfer' }"
      >
        <v-badge
          :content="transferCount"
          color="secondary"
        >
          <v-icon>
            mdi-dolly
          </v-icon>
        </v-badge>
      </v-btn>
    </v-fab-transition>
  </div>
</template>

<script>
  import notifications from '@/components/notification/fc_notifications_wrapper.vue'
  import bins from '@/components/job/bins.vue'
  import { routeIcon } from '@/router'
  import store from '@/store'

  export default {
    components: {
      'notifications-wrapper': notifications,
      bins
    },
    data() {
      return {
        drawer: true,
        rightDrawer: false,
        userFlyout: false,
        roles: false,
        daemonWarning: false,
        daemonDetails: false,
        fallenDaemons: [],
        notificationsCount: 0,
        navtab: parseInt(localStorage.getItem('navtab')) || 0
      }
    },
    computed: {
      user () {
        return this.$store.state.user.userData
      },
      roleList () {
        return Object.entries(this.user.role)
          .filter(e => e[1] === true)
          .map(e => e[0])
          .map(r => r.split('_').join(' ').toLowerCase())
      },
      transferCount () {
        return this.$store.state.transfer.jobs.length
      }
    },
    watch: {
      $route (to) {
        if (to.meta.navtab !== undefined) {
          this.navtab = to.meta.navtab
        }
      },
      navtab (val) {
        localStorage.setItem('navtab', val)
      }
    },
    async beforeRouteEnter (to, from, next) {
      try {
        const loggedIn = await store.dispatch('resume')
        if (!loggedIn) {
          if (from.path === '/login') {
            next()
          }
          sessionStorage.setItem('loginRedirect', to.path)
          next('/login')
        } else {
          next()
        }
      } catch (e) {
        store.commit('_setUser', null)
        next('/login')
      }
    },
    mounted: function () {
      this.getNotificationsCount()
      this.checkDaemons()
      this.ninterval = setInterval(this.getNotificationsCount, 10000)
      this.dinterval = setInterval(this.checkDaemons, 60000)
    },
    beforeDestroy: function () {
      clearInterval(this.ninterval)
      clearInterval(this.dinterval)
    },
    methods: {
      routeIcon,
      toggleNotifications: function () {
        this.rightDrawer = !this.rightDrawer;
        if (this.rightDrawer) {
          this.$refs.notifWrapper.loadNotifications()
        }
      },
      getNotificationsCount: function () {
        this.axios.get(this.$serverAddr + '/notifications/count', {
          withCredentials: true
        }).then((response) => {
          this.notificationsCount = response.data.count
        })
      },
      async checkDaemons () {
        const statuses = await this.axios.get(`${this.$serverAddr}/serverInfo/info`).then(r => r.data.subsystems)
        this.fallenDaemons = statuses.filter(s => s.status !== 'running').map(s => s.name)
        this.daemonWarning = this.fallenDaemons.length > 0
        if (!this.daemonWarning) this.daemonDetails = false
      }
    }
  }
</script>

<style scoped>
  .dash-link {
    display: flex;
    justify-content: center;
    padding: 1em;
    padding-top: 0;
  }

  .notifHeader {
    z-index: 2;
  }

  .logo {
    display: flex;
    width: 100%;
    min-height: 10em;
    margin: 1em auto;
    justify-content: center;
    align-items: center;
  }
  .logo img {
    max-width: 150px;
  }

  .logoText {
    color: #4c4c4c;
    text-align: center;
    padding-bottom: 15px;
    font-size: 26px;
    font-weight: 300;
    text-transform: uppercase;
  }

  .textLogoText {
    font-size: 18px;
    font-weight: 300;
  }

  .height100 {
    height: 100%;
  }

  .mainToolbar {
    z-index: 3;
  }

  .navigationDrawer {
    z-index: 6;
  }

  .z-top {
    z-index: 9001;
  }

  .theme--dark .navigationDrawer {
    background-color: #1e1e1e;
  }

  .main {
    padding-top: 64px !important;
  }

  .logoSmall {
    width: 100%;
    max-width: 40px;
    vertical-align: middle;
  }

  .logoSmallText {
    font-size: 24px;
    text-transform: uppercase;
    font-weight: 300;
    vertical-align: middle;
  }

  .important-must-not-miss {
    animation: notice-me-disappear 1s infinite;
  }
  .very-important-must-not-miss {
    animation: notice-me-bleed 2s infinite;
  }
  
  @keyframes notice-me-disappear {
    50% {
      opacity: 0;
    }
  }
  @keyframes notice-me-bleed {
    65% {
      box-shadow: 0 0 2em red
    }
  }

  /**/

  .route-enter-active {
    transition: all .4s cubic-bezier(0.16, 1, 0.3, 1)
  }
  .route-leave-active {
    transition: all .1s ease-in;
  }
  .route-enter {
    transform: translateY(10px);
    opacity: 0;
  }
  .route-leave-to {
    transform: scale(1.02) translateY(-10px);
    opacity: 0;
  }
</style>


<style>
  .v-list__group__items .v-list__tile {
    padding-left: 35px;
  }

  .notrans,
  .notrans .v-tabs-bar {
    transition: none;
  }
</style>
