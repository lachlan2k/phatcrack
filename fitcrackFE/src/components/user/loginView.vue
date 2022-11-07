<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <transition name="fade-transition">
    <v-container
      class="fill-height"
      fluid
    >
      <img 
        :src="require('@/assets/motive.svg')"
        class="motive"
      >
      <div class="main mx-auto">
        <img
          v-if="$vuetify.theme.dark"
          :src="require(`@/assets/logo-dark.svg`)"
          class="mx-auto pt-3 pb-3 d-block logo"
          alt="logo"
        >
        <img
          v-else
          :src="require(`@/assets/logo-light.svg`)"
          class="mx-auto pt-3 pb-3 d-block logo"
          alt="logo"
        >
        <div class="content">
          <v-progress-circular
            size="50"
            :width="3"
            class="progress trans"
            indeterminate
            color="primary"
            :class="{hide: !loading}"
          />
          <v-form
            ref="form"
            lazy-validation
            :class="{hide: loading}"
            class="trans"
          >
            <v-text-field
              v-model="username"
              solo
              label="Username"
              required
            />
            <v-text-field
              v-model="password"
              solo
              label="Password"
              required
              type="password"
            />
            <v-btn
              class="d-block mx-auto mb-4"
              type="submit"
              color="primary"
              :disabled="!username || !password"
              @click.prevent="submit"
            >
              sign in
            </v-btn>
          </v-form>
        </div>
        <div class="version">
          {{ version }}
        </div>
      </div>
    </v-container>
  </transition>
</template>

<script>
  export default {
    name: 'LoginView',
    data () {
      return {
        username: null,
        password: null,
        loading: true,
        version: __APP_VERSION__
      }
    },
    created () {
      if (this.$store.state.user.loggedIn) {
        this.$router.replace({name: 'home'})
      } else {
        this.loading = false
      }
    },
    methods: {
      async submit() {
        this.loading = true
        const succeeded = await this.$store.dispatch('signIn', {
          username: this.username,
          password: this.password
        })
        if (!succeeded) {
          this.loading = false
        }
      }
    }
  }
</script>

<style scoped>
  .hide {
    opacity: 0;
    visibility: hidden;
  }

  .trans {
    -webkit-transition: opacity .3s ease-in-out;
    -moz-transition: opacity .3s ease-in-out;
    -ms-transition: opacity .3s ease-in-out;
    -o-transition: opacity .3s ease-in-out;
    transition: opacity .3s ease-in-out;
  }

  .main {
    z-index: 1;
    min-width: 260px;
  }

  .content {
    position: relative;
    max-width: 600px;
  }

  .version {
    position: fixed;
    bottom: 1.5em;
    left: 0;
    width: 100%;
    font-size: 0.8em;
    font-weight: bold;
    color: #888;
    text-align: center;
  }

  .progress {
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    margin: auto;
    position: absolute;
  }

  .motive {
    position: fixed;
    z-index: 0;
    opacity: .2;
    left: 5%;
    width: 90%;
    transform: rotate(-10deg)
  }

  .logo {
    width: 100%;
    max-width: 260px;
  }
</style>

