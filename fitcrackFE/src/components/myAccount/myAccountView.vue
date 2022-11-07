<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-container>
    <h2> My account </h2>
    <v-row justify="center">
      <v-card
        class="mb-4"
        min-width="300"
        max-width="500"
      >
        <v-card-title>
          {{ username }}
        </v-card-title>
        <v-card-text>
          {{ mail }}
        </v-card-text>
      </v-card>
    </v-row>
    <v-row justify="center">
      <v-card
        min-width="300"
        max-width="500"
      >
        <v-card-title>
          Change password
        </v-card-title>
        <v-card-text>
          <v-form
            ref="form"
            lazy-validation
          >
            <v-text-field
              v-model="oldPassword"
              type="password"
              label="Current password"
              required
            />
            <v-text-field
              v-model="newPassword0"
              type="password"
              label="New password"
              required
            />
            <v-text-field
              v-model="newPassword1"
              type="password"
              label="Confirm new password"
              required
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn
            color="primary"
            text
            :disabled="newPassword0 != newPassword1 || oldPassword == null || newPassword0 == null"
            @click.stop="editPassword"
          >
            Update
          </v-btn>
        </v-card-actions>
      </v-card>
    <!--    <h3> Message: {{ user }} </h3> -->
    </v-row>
  </v-container>
</template>

<script>
    export default {
      name: "MyAccount",
      data: function() {
       return {
          username: this.$store.state.user.userData.username,
          mail: this.$store.state.user.userData.mail,
          oldPassword: null,
          newPassword0: null,
          newPassword1: null
        }
      },
      computed: {
        gradient: function () {
          return window.theme.gradient
        }
      },
      methods: {
        editPassword() {
          this.axios.post(this.$serverAddr + '/user/password/change_my_password', {
            old_password: this.oldPassword,
            new_password: this.newPassword1
          }).then((response) => {
            console.log(response.data);
            console.log('user password changed');
            this.oldPassword = '';
            this.newPassword0 = '';
            this.newPassword1 = ''
          })
        }
      }
  }
</script>



<style scoped>
h2 {
  padding-bottom: 5px;
  text-align: center;
}
.table {
  max-width: 500px;
  margin: auto;
}
.width150 {
  min-width: 100px;
}

</style>
