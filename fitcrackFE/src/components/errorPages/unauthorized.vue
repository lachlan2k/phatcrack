<template>
  <div class="text-center mt-8 pa-8 container">
    <h1>Sorry {{ $store.state.user.userData.username }}, I'm afraid I can't do that...</h1>
    <svg class="illust" viewBox="0 0 867 468" fill="none" xmlns="http://www.w3.org/2000/svg">
      <g clip-path="url(#clip0)">
        <path d="M45.6 429.28a5 5 0 011.65-.28H818a5 5 0 011.62.27l43.9 15c5.35 1.83 4.03 9.73-1.62 9.73H4.41c-5.63 0-6.96-7.86-1.65-9.72l42.84-15z" fill="#000" fill-opacity=".15"/>
        <rect x="172.5" y="184.5" width="55" height="252" rx="12.5" fill="#D78CA7" stroke="#DC2365" stroke-width="15"/>
        <rect x="7.5" y="-7.5" width="55" height="252" rx="12.5" transform="matrix(-1 0 0 1 701.5 190)" fill="#D78CA7" stroke="#DC2365" stroke-width="15"/>
        <rect x="119.5" y="313.5" width="628" height="57" rx="12.5" fill="#F3D2DE" stroke="#DC2365" stroke-width="15"/>
        <path d="M212.75 184.5h-46.02a12.5 12.5 0 00-11.4 7.35l-113.97 252a12.5 12.5 0 0011.39 17.65h46.02a12.5 12.5 0 0011.4-7.35l113.97-252a12.5 12.5 0 00-11.39-17.65z" fill="#F3D2DE"/>
        <path d="M212.75 184.5h-46.02a12.5 12.5 0 00-11.4 7.35l-113.97 252a12.5 12.5 0 0011.39 17.65h46.02a12.5 12.5 0 0011.4-7.35l113.97-252a12.5 12.5 0 00-11.39-17.65z" fill="url(#paint0_linear)"/>
        <path d="M212.75 184.5h-46.02a12.5 12.5 0 00-11.4 7.35l-113.97 252a12.5 12.5 0 0011.39 17.65h46.02a12.5 12.5 0 0011.4-7.35l113.97-252a12.5 12.5 0 00-11.39-17.65z" stroke="#DC2365" stroke-width="15"/>
        <path d="M653.75 182.5h46.02a12.5 12.5 0 0111.4 7.35l113.97 252a12.5 12.5 0 01-11.39 17.65h-46.02a12.5 12.5 0 01-11.4-7.35l-113.97-252a12.5 12.5 0 0111.39-17.65z" fill="#F3D2DE"/>
        <path d="M653.75 182.5h46.02a12.5 12.5 0 0111.4 7.35l113.97 252a12.5 12.5 0 01-11.39 17.65h-46.02a12.5 12.5 0 01-11.4-7.35l-113.97-252a12.5 12.5 0 0111.39-17.65z" fill="url(#paint1_linear)"/>
        <path d="M653.75 182.5h46.02a12.5 12.5 0 0111.4 7.35l113.97 252a12.5 12.5 0 01-11.39 17.65h-46.02a12.5 12.5 0 01-11.4-7.35l-113.97-252a12.5 12.5 0 0111.39-17.65z" stroke="#DC2365" stroke-width="15"/>
        <rect x="65.5" y="6.5" width="732" height="226" rx="12.5" fill="#F3F3F3" stroke="#DC2365" stroke-width="15"/>
        <path d="M198 14h72l-98 211h-72l98-211zM530 14h72l-98 211h-72l98-211zM364 14h72l-98 211h-72l98-211zM696 14h72l-98 211h-72l98-211z" fill="#FF7168"/>
      </g>
      <defs>
        <linearGradient id="paint0_linear" x1="132.75" y1="192" x2="132.75" y2="454" gradientUnits="userSpaceOnUse">
          <stop stop-color="#D78CA7"/>
          <stop offset=".35" stop-color="#D78CA7"/>
          <stop offset=".35" stop-color="#D78CA7" stop-opacity="0"/>
        </linearGradient>
        <linearGradient id="paint1_linear" x1="733.75" y1="190" x2="733.75" y2="452" gradientUnits="userSpaceOnUse">
          <stop stop-color="#D78CA7"/>
          <stop offset=".35" stop-color="#D78CA7"/>
          <stop offset=".35" stop-color="#D78CA7" stop-opacity="0"/>
        </linearGradient>
        <clipPath id="clip0">
          <path fill="#fff" d="M0 0h867v468H0z"/>
        </clipPath>
      </defs>
    </svg>
    <p>
      Your role doesn't permit you to access <strong>{{ $route.query.accessing || 'that route' }}</strong>.<br>
      If you are not <strong>{{ $store.state.user.userData.username }}</strong>, you can try to sign in with another account.
    </p>
    <v-btn 
      class="mb-4"
      color="primary"
      @click="switchAccount"
    >
      <v-icon 
        left
        small
      >
        mdi-account-switch
      </v-icon>
      Use another account
    </v-btn>
    <p>
      Otherwise, please contact an admin and ask for the {{ $route.query.missing }} permission.
    </p>
  </div>
</template>

<script>
export default {
  methods: {
    switchAccount () {
      sessionStorage.setItem('loginRedirect', this.$route.query.accessing)
      this.$store.dispatch('signOut')
    }
  }
}
</script>

<style scoped>
.container {
  margin: 0 auto;
  max-width: 800px;
}

.illust {
  width: 80%;
  max-width: 300px;
  margin: 2em auto;
}
</style>