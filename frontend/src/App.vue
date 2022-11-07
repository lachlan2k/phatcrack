<script setup lang="ts">
import { RouterView } from 'vue-router'
import Footer from '@/components/Footer.vue'
</script>

<template>
  <v-app theme="light">
    <v-main>
      <RouterView />
      <Footer />
    </v-main>
  </v-app>
</template>