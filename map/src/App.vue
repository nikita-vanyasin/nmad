<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
    <img alt="Nmad logo" class="logo" src="@/assets/logo.png" width="60" height="60" />
    <div class="wrapper">
      <nav>
        <RouterLink to="/">Home</RouterLink>
        |
        <RouterLink to="/about">About</RouterLink>
      </nav>
    </div>
  </header>

  <RouterView />
</template>

<style scoped>
.wrapper {
  display: inline-block;
  vertical-align: top;
}

.wrapper nav {
  margin-top: 22px;
  margin-left: 20px;
}

.logo {
  border-radius: 30px;
  display: inline-block;
  margin-left: 5px;
  margin-top: 5px;
}

</style>
