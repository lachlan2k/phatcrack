<script setup lang="ts">
  //
</script>

<template>
    <div class="dash-container">
        <v-sheet class="dash-row pa-0">
            <status v-if="hostsInfo" :jobs="jobsInfo" :hosts="hostsInfo" :server="serverInfo">
                <server-gauges v-bind="actualUsage" :services="serverInfo" />
            </status>
        </v-sheet>
        <v-divider />
        <v-card-title class="px-8">
            Last jobs
        </v-card-title>
        <div v-if="lastJobs" class="minijobs overflowing dash-row pt-0 pb-7 px-8">
            <minijob detailed :data="lastJobs[0]" />
            <minijob v-for="j in lastJobs.slice(1)" :key="j.id" :data="j" />
            <div class="pr-6">
                <v-btn to="/jobs" large color="primary">
                    See all
                    <v-icon right>
                        mdi-briefcase-search
                    </v-icon>
                </v-btn>
            </div>
        </div>
        <v-divider />
        <v-sheet class="dash-row px-4">
            <div class="half">
                <v-card-title>Progress of all jobs</v-card-title>
                <jobProgress :from="chartsFrom" :to="chartsTo" />
            </div>
            <div class="half">
                <v-card-title>Global workunits distribution</v-card-title>
                <jobWorkunits logarithmic :from="chartsFrom" :to="chartsTo" />
            </div>
        </v-sheet>
        <v-sheet class="dash-row px-8 justify-end">
            <div class="d-flex">
                <dt-picker v-model="chartsFrom" class="mr-4" label="From" />
                <dt-picker v-model="chartsTo" label="To" />
            </div>
        </v-sheet>
    </div>
</template>

<style scoped>
.dash-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
}

.overflowing {
    flex-wrap: nowrap;
    overflow-x: auto;
}

.minijobs {
    align-items: stretch;
}

.half {
    width: 50%;
}

@media screen and (max-width: 600px) {
    .half {
        width: 100%;
    }
}
</style>