<template>
  <v-dialog
      v-model="visible"
      max-width="550"
      persistent
    >
      <v-card>
        <v-card-title>
          API Server {{retries > 5 ? 'Really ' : ''}}Unreachable
        </v-card-title>
        <v-card-subtitle>
          {{ addr }}
        </v-card-subtitle>
        <v-card-text>
          <span v-if="retries <= 5">
            We'll keep trying to connect in case it hasn't started yet.
          </span>
          <span v-else>
            Tried {{retries}} times and still nothing.
            <span v-if="retries >= 15"> Sorry, but at this point, it's pretty clear you'll need to fix it manually.</span>
            <span v-if="retries >= 20"> Seriously, it won't just magically start working.</span>
            <span v-if="retries >= 30"> We're getting nowhere like this. Just follow the <strong>Help button</strong> below, please?</span>
            <span v-if="retries >= 50"> <br><br>You're enjoying this, aren't you?</span>
            <span v-if="retries >= 60"> It's kind of awkward, not gonna lie.</span>
            <span v-if="retries >= 70"> You do realize that this is just an error message, right?</span>
            <span v-if="retries >= 100"> <br><br>OK, it's your free time, not ours. Let's humor you some more, then.</span>
          </span>
          <v-progress-linear
            class="mt-3"
            indeterminate
            :color="statusColor"
          />
        </v-card-text>
        <v-card-text
          class="py-0"
        >
          <v-alert 
            v-if="retries <= 0"
            type="info"
          >
            Hang tight! Reconnecting in a few.
          </v-alert>
          <v-alert 
            v-else-if="retries <= 5"
            type="error"
            :color="statusColor"
          >
            {{retries}} failed reconnection attempt{{retries > 1 ? 's' : ''}}
          </v-alert>
          <v-alert 
            v-else-if="retries < 100"
            type="error"
            icon="mdi-emoticon-dead"
          >
            Yeah, still no connection. Try <strong>Help</strong> below.
          </v-alert>
          <v-alert 
            v-else
            type="success"
            icon="mdi-trophy"
          >
            Yay, <strong>Achievement Unlocked!</strong> Retried {{Math.floor(retries / 10) * 10}} times.
          </v-alert>
        </v-card-text>
        <v-card-actions>
          <v-spacer/>
          <v-btn 
            text
            :color="retries >= 3 ? statusColor : ''"
            @click="showSuggestions"
          >
            Help
          </v-btn>
          <v-btn
            @click="retry"
          >
            Retry Now
          </v-btn>
        </v-card-actions>
      </v-card>

      <v-bottom-sheet 
        v-model="help"
        overlay-color="secondary"
        overlay-opacity="0.3"
        dark
      >
        <v-card tile min-width="100vw">
          <v-card-title>
            <v-icon left>mdi-bandage</v-icon>
            Here's what you can try
            <v-spacer/>
            <v-card-actions>
              <v-btn 
                text
                to="//github.com/nesfit/fitcrack/issues"
                target="_blank"
              >
                <v-icon left>mdi-github</v-icon>
                Ask on Github
              </v-btn>
              <v-btn
                @click="hideSuggestions"
                color="success"
              >
                <v-icon left>mdi-check</v-icon>
                Done
              </v-btn>
            </v-card-actions>
          </v-card-title>
          <v-card-text>
            <v-expansion-panels 
              v-model="suggestionsShown"
              multiple
              light
            >
              <!--            -->
              <v-expansion-panel>
                <v-expansion-panel-header>
                  <span>
                    <v-icon left small class="pb-1">mdi-ethernet</v-icon>
                    <span class="text-subtitle-1">Check network connection</span>
                  </span>
                </v-expansion-panel-header>
                <v-expansion-panel-content class="limit-width">
                  Are you connected to the right network? Your Fitcrack API server might be exposed to a different network, or your internet connection maybe isn't working right and the app loaded out of cache.
                  <ul>
                    <li>Check your internet connection</li>
                    <li>Check where you are – do you need a VPN to connect to Fitcrack?</li>
                    <li>If you are on the same network as the server, try to ping it</li>
                  </ul>
                </v-expansion-panel-content>
              </v-expansion-panel>
              <!--            -->
              <v-expansion-panel>
                <v-expansion-panel-header>
                  <span>
                    <v-icon left small class="pb-1">mdi-wrench</v-icon>
                    <span class="text-subtitle-1">Check client configuration</span>
                  </span>
                </v-expansion-panel-header>
                <v-expansion-panel-content class="limit-width">
                  The webadmin app is configured to connect to the API on a specific address. Check if it hasn't changed server-side, or if it was configured correctly in the first place.
                  <ul>
                    <li>Check if <code>window.serverAddress</code> in the <code>/public/static/configuration.js</code> file of FitcrackFE has the correct address</li>
                    <li>Double-check the server's current address!</li>
                  </ul>
                </v-expansion-panel-content>
              </v-expansion-panel>
              <!--            -->
              <v-expansion-panel>
                <v-expansion-panel-header>
                  <span>
                    <v-icon left small class="pb-1">mdi-server</v-icon>
                    <span class="text-subtitle-1">Check server status</span>
                  </span>
                </v-expansion-panel-header>
                <v-expansion-panel-content class="limit-width">
                  The API server could be offline or malfunctoining. You could try to restart the Fitcrack service.
                  <ul>
                    <li>Access the server shell and check if the service is enabled with <code>service fitcrack status</code></li>
                    <li>Try to start it up if needed with <code>service fitcrack start</code></li>
                  </ul>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-card-text>
        </v-card>
      </v-bottom-sheet>
    </v-dialog>
</template>

<script>
import store from '@/store'
export default {
  props: {
    visible: Boolean
  },
  data () {
    return {
      interval: null,
      retries: 0,
      help: false,
      suggestionsShown: [0]
    }
  },
  computed: {
    addr () {
      return window.serverAddress
    },
    statusColor () {
      if (this.retries == 0) {
        return 'blue'
      } else if (this.retries <= 5) {
        return 'orange'
      } else if (this.retries < 100) {
        return 'red'
      } else {
        return 'green'
      }
    }
  },
  methods: {
    showSuggestions () {
      this.help = true
    },
    hideSuggestions () {
      this.help = false
    },
    async retry () {
      this.retries += 1
      // debug
      // if (this.retries > 5) this.$emit('connect')
      //
      const user = await store.dispatch('refreshUser')
      if (user) {
        this.$emit('connect')
      }
    }
  },
  watch: {
    visible (value) {
      if (value == true) {
        this.retries = 0
        setTimeout(this.retry, 2000)
        this.interval = setInterval(this.retry, 5000)
      } else {
        clearInterval(this.interval)
      }
    }
  }
}
</script>

<style scoped>
.limit-width {
  max-width: 700px;
}
</style>