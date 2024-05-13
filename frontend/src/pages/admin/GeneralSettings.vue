<script setup lang="ts">
import AgentConfig from '@/components/Admin/AgentConfig.vue'
import GeneralConfig from '@/components/Admin/GeneralConfig.vue'
import { computed, ref } from 'vue'
import type { Component } from 'vue'

const tab = (name: string, icon: string, component: Component) => ({ name, icon, component })

const tabs = [tab('General', 'fa-gear', GeneralConfig), tab('Auth', 'fa-passport', AgentConfig), tab('Agent', 'fa-robot', AgentConfig)]

const activeTab = ref(0)

const ComponentToRender = computed(() => tabs[activeTab.value].component)
</script>

<template>
  <main class="w-full p-4">
    <h1 class="text-4xl font-bold">Configuration</h1>

    <div class="mt-6 flex flex-wrap gap-6">
      <div class="card bg-base-100 shadow-xl">
        <div class="card-body justify-between">
          <div class="tabs">
            <a
              v-for="(tab, i) in tabs"
              v-key="i"
              class="tab tab-bordered"
              :class="activeTab == i ? 'tab-active' : ''"
              @click="activeTab = i"
            >
              <font-awesome-icon :icon="'fa-solid ' + tab.icon" class="mr-2" />{{ tab.name }}</a
            >
          </div>

          <ComponentToRender />
        </div>
      </div>
    </div>
  </main>
</template>
