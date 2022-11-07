// Vuex module for the Data Transfer process

export default {
  namespaced: true,
  state: {
    jobs: []
  },
  getters: {
    job_ids: state => state.jobs.map(j => j.id)
  },
  mutations: {
    add (state, jobs) {
      const adding = jobs.filter(j => !state.jobs.some(job => j.id == job.id))
      state.jobs = state.jobs.concat(adding)
    },
    remove (state, jobId) {
      state.jobs = state.jobs.filter(j => j.id != jobId)
    },
    reset (state) {
      state.jobs = []
    }
  }
}