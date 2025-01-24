<script setup>
import Button from 'primevue/button';
import Card from 'primevue/card';
import InputText from 'primevue/inputtext';
</script>

<template>
  <div class="login-form">
    <Card>
      <template #title>Login</template>
      <template #content>
        <div class="flex justify-center flex-col gap-4">
          <div class="flex flex-col gap-1">
            <InputText name="username" type="text" placeholder="Username" />
          </div>
          <div class="flex flex-col gap-1">
            <InputText name="email" type="text" placeholder="Email" />
          </div>
          <Button label="Submit" />
        </div>
      </template>
    </Card>
  </div>
</template>

<style scoped>
.login-form {
  max-width: 400px;
  margin: 50px auto 0;
}
</style>