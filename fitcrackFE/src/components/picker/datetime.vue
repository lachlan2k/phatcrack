<template>
  <v-text-field
    v-if="native"
    v-bind="{...$props, ...$attrs}"
    :disabled="disabled"
    type="datetime-local"
    :value="value"
    @input="emitNative"
  />
  <div 
    v-else
    class="polyfill d-flex"
  >
    <v-menu
      ref="dateMenu"
      v-model="datePickerOpen"
      :close-on-content-click="false"
      :return-value.sync="date"
      offset-y
      min-width="290px"
    >
      <template v-slot:activator="{ on }">
        <v-text-field
          v-model="date"
          readonly
          v-bind="{...$props, ...$attrs}"
          :disabled="disabled"
          class="datepicker-field mr-2"
          placeholder="YYYY-MM-DD"
          v-on="on"
        />
      </template>
      <v-date-picker 
        v-model="date" 
        no-title
        scrollable
      >
        <v-spacer />
        <v-btn 
          text
          color="primary"
          @click="datePickerOpen = false"
        >
          Cancel
        </v-btn>
        <v-btn 
          text
          color="primary"
          @click="$refs.dateMenu.save(date)"
        >
          OK
        </v-btn>
      </v-date-picker>
    </v-menu>
    <v-menu
      ref="timeMenu"
      v-model="timePickerOpen"
      :close-on-content-click="false"
      :nudge-right="40"
      :return-value.sync="time"
      offset-y
      max-width="290px"
      min-width="290px"
    >
      <template v-slot:activator="{ on }">
        <v-text-field
          v-model="time"
          readonly
          v-bind="{...$props, ...$attrs}"
          :disabled="disabled"
          class="timepicker-field"
          placeholder="00:00"
          label=""
          v-on="on"
        />
      </template>
      <v-time-picker
        v-if="timePickerOpen"
        v-model="time"
        full-width
        @click:minute="$refs.timeMenu.save(time)"
      />
    </v-menu>
  </div>
</template>

<script>
export default {
  props: {
    disabled: Boolean,
    value: {
      type: String,
      default: () => this.$moment().local().format('YYYY-MM-DDTHH:mm')
    }
  },
  data () {
    return {
      native: true,
      datePickerOpen: false,
      timePickerOpen: false,
      date: null,
      time: null
    }
  },
  watch: {
    date () { this.emitPolyfill() },
    time () { this.emitPolyfill() },
    value (val) {
      [this.date, this.time] = val.split('T')
    }
  },
  created () {
    const testElement = document.createElement('input')
    testElement.setAttribute('type', 'datetime-local')
    this.native = testElement.type !== 'text' // set to false to test in supported browser
    if (this.value) {
      [this.date, this.time] = this.value.split('T')
    }
  },
  methods: {
    emitNative (result) {
      this.$emit('input', result)
    },
    emitPolyfill () {
      if (this.date && this.time) {
        this.$emit('input', `${this.date}T${this.time}`)
      }
    }
  }
}
</script>