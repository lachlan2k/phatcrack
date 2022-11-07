<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <div>
    <div class="notif-top d-flex align-center pa-3 pb-0">
      <v-badge
        :value="count > 0"
        :content="count"
        inline
        color="secondary"
      >
        <v-card-title>Notifications</v-card-title>
      </v-badge>
      <v-spacer />
      <v-btn
        icon
        @click="$emit('close')"
      >
        <v-icon>
          mdi-close
        </v-icon>
      </v-btn>
    </div>
    <v-progress-linear
      v-if="notifications === null"
      indeterminate
    />
    <div v-else>
      <v-list>
        <notification
          v-for="(item, i) in notifications"
          :key="i"
          :type="item.type"
          :title="item.title"
          :text="item.text"
          :seen="item.seen"
          :time="item.time"
          :job-id="item.job_id"
        />
      </v-list>
      <div class="notif-actions d-flex pa-3">
        <v-btn
          v-show="page < pages"
          text
          color="primary"
          @click="loadNotifications"
        >
          <v-icon left>
            mdi-timeline-text
          </v-icon>
          older
        </v-btn>
        <v-spacer />
        <!-- <v-btn
          text
          color="error"
        >
          clear all
          <v-icon right>
            mdi-notification-clear-all
          </v-icon>
        </v-btn> -->
      </div>
    </div>
  </div>
</template>

<script>
  import notification from '@/components/notification/fc_notification.vue'

  export default {
    components: {
      'notification': notification
    },
    props: {
      count: {
        type: Number,
        default: 0
      }
    },
    data: function () {
      return {
        notifications: [],
        page: 1,
        pages: undefined
      }
    },
    methods: {
      nextPage () {
        if (this.pages && this.page < this.pages) {
          this.pages += 1
          return true
        } else {
          return false
        }
      },
      loadNotifications () {
        if (this.nextPage() || !this.pages) {
          this.axios.get(this.$serverAddr + '/notifications', {
            params: {
              page: this.page
            }
          }).then((response) => {
            this.notifications = [...this.notifications, ...response.data.items]
            this.pages = response.data.pages
          })
        }
      }
    }
  }
</script>
<style scoped>

</style>
