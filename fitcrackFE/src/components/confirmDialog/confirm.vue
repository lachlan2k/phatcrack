<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-dialog
    v-model="dialog"
    :max-width="options.width"
  >
    <v-card
      flat
    >
      <v-card-title>
        {{ title }}
      </v-card-title>
      <v-card-text v-show="!!message">
        {{ message }}
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          text
          @click.native="cancel()"
        >
          Cancel
        </v-btn>
        <v-btn
          color="primary"
          outlined
          text
          @click.native="agree()"
        >
          <v-icon left>
            mdi-keyboard-return
          </v-icon>
          Yes
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
  /**
   * Vuetify Confirm Dialog component
   *
   * Insert component where you want to use it:
   * <confirm ref="confirm"></confirm>
   *
   * Call it:
   * this.$refs.confirm.open('Delete', 'Are you sure?', { color: 'red' }).then((confirm) => {});
   *
   * Alternatively you can place it in main App component and access it globally via this.$root.$confirm
   * <template>
   *   <v-app>
   *     ...
   *     <confirm ref="confirm"></confirm>
   *   </v-app>
   * </template>
   *
   * mounted() {
   *   this.$root.$confirm = this.$refs.confirm.open;
   * }
   */
  export default {
    data() {
      return {
        dialog: false,
        resolve: null,
        reject: null,
        message: null,
        title: null,
        options: {
          color: 'primary',
          width: 290,
        }
      };
    },
    mounted () {
      window.addEventListener('keydown', this.shortcutHandler)
    },
    beforeDestroy () {
      window.removeEventListener('keydown', this.shortcutHandler)
    },
    methods: {
      open(title, message, options) {
        this.dialog = true;
        this.title = title;
        this.message = message;
        this.options = Object.assign(this.options, options);
        return new Promise((resolve, reject) => {
          this.resolve = resolve;
          this.reject = reject;
        });
      },
      agree() {
        this.resolve(true);
        this.reject(false);
        this.dialog = false;
      },
      cancel() {
        this.dialog = false;
      },
      shortcutHandler (e) {
        switch (e.key) {
          case 'Escape':
            if (this.reject) {
              this.cancel()
            }
            break
          case 'Enter':
            if (this.resolve) {
              this.agree()
            }
            break
          default:
            return
        }
      }
    }
  }
</script>
