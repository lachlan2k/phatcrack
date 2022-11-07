<!--
   * Author : see AUTHORS
   * Licence: MIT, see LICENSE
-->

<template>
  <v-jumbotron
    :gradient="gradient"
    class="jumbo"
    height="100%"
  >
    <v-container
      class="fill-height"
      fluid
    >
      <v-row
        class="flex"
        align="center"
        justify="center"
      >
        <div class="jumboCard px-4">
          <div class="contentJumbo">
            <h1 class="d-block">
              404 - Not Found
            </h1>
            <h2 class="d-block">
              We are sorry. Could not find what you are looking for.
            </h2>
          </div>
        </div>
      </v-row>
    </v-container>
  </v-jumbotron>
</template>

<script>
    export default {
      name: "PageNotFound",
      computed: {
        gradient: function () {
          return window.theme.gradient
        }
      },
    }
</script>

<style scoped>
  .contentJumbo {
    position: relative;
  }

  .contentJumbo h1 {
    font-size: 60px;
    color:white;
    font-weight: 400;
  }
  .contentJumbo h2 {
    font-size: 34px;
    color:white;
    font-weight: 100;
  }
  .jumbo {
    position: fixed;
    z-index: 5;
  }
</style>
