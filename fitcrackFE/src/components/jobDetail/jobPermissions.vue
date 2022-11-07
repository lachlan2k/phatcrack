<template>
  <v-dialog 
    max-width="800"
    :value="value"
    @input="$emit('input')"
  >
    <v-card>
      <v-card-title>
        <span>Permissions Manager</span>
        <v-spacer />
        <v-btn
          icon
          @click="$emit('input')"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>
      <v-card-text>
        <div>
          <strong>Assign permissions to specific users for this job.</strong>
          <ul>
            <li>Only you can change these settings as the job owner.</li>
            <li>This does not affect users with global permission roles.</li>
          </ul>
        </div>
      </v-card-text>
      <v-card-text>
        <v-data-table
          :loading="loading"
          :headers="headers"
          :items="perms"
        >
          <template v-slot:item.username="{ item }">
            <div class="py-2">
              <div class="subtitle-2">
                {{ item.username }}
                <span v-if="item.owner">
                  (owner)
                </span>
              </div>
              <div class="caption">
                {{ item.mail }}
              </div>
            </div>
          </template>
          <template v-slot:item.view="{ item }">
            <v-checkbox
              v-model="item.view"
              class="cc d-inline-block ma-0"
              on-icon="mdi-eye"
              off-icon="mdi-eye-off"
              hide-details
              :disabled="item.owner"
              @change="e => update(item.id, 'view', e)"
            />
          </template>
          <template v-slot:item.modify="{ item }">
            <v-checkbox
              v-model="item.modify"
              class="cc d-inline-block ma-0"
              on-icon="mdi-pencil"
              off-icon="mdi-pencil-off"
              hide-details
              :disabled="item.owner"
              @change="e => update(item.id, 'modify', e)"
            />
          </template>
          <template v-slot:item.operate="{ item }">
            <v-checkbox
              v-model="item.operate"
              class="cc d-inline-block ma-0"
              on-icon="mdi-flash"
              off-icon="mdi-flash-off"
              hide-details
              :disabled="item.owner"
              @change="e => update(item.id, 'operate', e)"
            />
          </template>
        </v-data-table>
        <v-menu
          v-model="addUserMenu"
          :close-on-content-click="false"
          top
        >
          <template #activator="{ on }">
            <v-btn
              v-on="on"
              @click="loadUsers"
            >
              Add User
            </v-btn>
          </template>
          <v-card>
            <v-card-text>
              <v-autocomplete
                :loading="fetchingUsers"
                :items="users"
                autofocus
                item-text="username"
                item-value="id"
                label="Type username"
                @input="addUser"
              />
            </v-card-text>
          </v-card>
        </v-menu>
      </v-card-text>
      <v-card-actions>
        <v-progress-linear
          color="secondary"
          :indeterminate="saving"
          rounded
        />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  components: {
  },
  props: {
    value: Boolean,
    jobId: {
      type: Number,
      required: true
    }
  },
  data () {
    return {
      perms: [],
      users: [],
      addUserMenu: false,
      loading: true,
      fetchingUsers: false,
      saving: false,
      //
      headers: [
        {
          text: 'User',
          align: 'start',
          value: 'username'
        },
        {text: 'View', value: 'view', align: 'center', sortable: false},
        {text: 'Edit', value: 'modify', align: 'center', sortable: false},
        {text: 'Operate', value: 'operate', align: 'center', sortable: false}
      ]
    }
  },
  mounted () {
    this.load()
  },
  methods: {
    async load () {
      this.loading = true
      this.perms = await this.axios.get(`${this.$serverAddr}/job/${this.jobId}/permissions`).then(r => r.data.items)
      this.loading = false
    },
    async loadUsers () {
      if (this.users.length == 0) {
        this.fetchingUsers = true
        this.users = await this.axios.get(`${this.$serverAddr}/user/`).then(r => r.data.items)
        this.fetchingUsers = false
      }
    },
    async addUser (id) {
      this.addUserMenu = false
      await this.update(id)
      this.load()
    },
    async update (user_id, permission = 'view', value = false) {
      this.saving = true
      try {
        await this.axios.put(`${this.$serverAddr}/job/${this.jobId}/permissions`, {
          user_id,
          [permission]: value
        })
      } catch (e) {
        console.error(e)
      } finally {
        this.saving = false
      }
    }
  }
}
</script>

<style>
  .cc .v-input--selection-controls__input {
    margin: 0
  }
</style>