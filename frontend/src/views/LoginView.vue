<script setup lang="ts">
import { ref } from 'vue'

const version = ref('v1.0.0')

const loading = ref(false)
const username = ref('')
const password = ref('')

function onSubmit() {
  loading.value = true
  console.log({ username, password })
}
</script>

<template>
    <v-container class="fluid fill-height" fluid>
      <v-layout class="align-center justify-center fill-height">
        <v-col xs="12" sm="8" md="4" lg="3" class="login-col">
          <v-card>

            <v-card-title>
              <p class="text-h4 text--primary text-center mt-4 mb-4">Phatcrack</p>
            </v-card-title>

            <v-card-text>
              <v-form>
                <v-text-field v-model="username" solo label="Username" required />
                <v-text-field v-model="password" solo label="Password" required type="password" />
                <v-btn class="d-block mx-auto mb-4" type="submit" color="primary" :disabled="!username || !password"
                  @click.prevent="onSubmit">
                  Sign In
                </v-btn>
              </v-form>
            </v-card-text>


          </v-card>
        </v-col>
      </v-layout>
    </v-container>
</template>

<style scoped>
.login-col {
  max-width: 500px;
}
</style>