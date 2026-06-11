<template>
  <div id="app" class="min-h-screen bg-gray-50 text-gray-800">
    <DefaultLayout>
      <router-view />
    </DefaultLayout>
  </div>
</template>

<script setup lang="ts">
import DefaultLayout from './layouts/DefaultLayout.vue'
</script>

<style>
/* CSS global nhẹ */
body {
  margin: 0;
  font-family: 'Inter', system-ui, sans-serif;
}
</style>
