<template>
  <a-config-provider :theme="{ token: { colorPrimary: '#3b82f6', borderRadius: 8 } }">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
</script>

<style>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  color: #1e293b;
}
</style>
