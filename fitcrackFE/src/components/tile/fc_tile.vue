<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-card flat class="mx-auto">
    <v-card-title>
      <v-icon 
        v-if="icon"
        left
      >
        {{ icon }}
      </v-icon>
      {{ title }}
    </v-card-title>
    <v-card-text class="contentFcTile">
      <v-skeleton-loader
        v-if="loading"
        class="mx-auto"
        type="article"
      />
      <slot v-else />
    </v-card-text>
  </v-card>
</template>

<script>
  export default {
    name: "FcTile",
    props: {
      title: String,
      loading: Boolean,
      icon: String
    }
  }
</script>
