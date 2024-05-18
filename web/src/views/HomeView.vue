<script setup lang="ts">
import { onMounted} from "vue";
import loginService from "@/services/login-service";

onMounted(async () => {
  await loginService.requestLogin();
});

</script>

<template>
  <h1>Logging in ....</h1>
</template>
