<template>
  <div class="host-grid">
    <v-sheet
      rounded
      v-for="host in active"
      :key="'ha' + host"
      class="host-icon host-online rounded-xl"
      color="success"
    />
    <v-sheet
      rounded
      v-for="host in inactive"
      :key="'hi' + host"
      class="host-icon host-offline rounded-xl"
    />
  </div>
</template>

<script>
export default {
  props: {
    active: {
      type: Number,
      default: 0
    },
    inactive: {
      type: Number,
      default: 0
    }
  }
}
</script>

<style scoped>
.host-grid {
  display: flex;
  flex-wrap: wrap;
  max-width: 7.5em;
}

.host-icon {
  position: relative;
  width: 1em;
  height: 1em;
  margin: .25em;
}

.host-offline {
  background-color: rgba(128,128,128,.3);
}
</style>