<script setup lang="ts">
import { computed } from "vue";
import { RouterLink, useRouter } from "vue-router";
import { Boxes, ShieldCheck } from "lucide-vue-next";
import { useAuthStore } from "@/stores/auth";
import Button from "@/components/ui/Button.vue";

const auth = useAuthStore();
const router = useRouter();

const userLabel = computed(() => auth.user?.email ?? "Guest");

function logout() {
  auth.logout();
  router.push("/");
}
</script>

<template>
  <header class="hero-panel mb-8 flex flex-col gap-5 px-5 py-5 md:flex-row md:items-center md:justify-between md:px-8">
    <div class="flex items-center gap-4">
      <div class="flex h-14 w-14 items-center justify-center rounded-2xl bg-primary text-primary-foreground">
        <Boxes class="h-7 w-7" />
      </div>
      <div>
        <RouterLink to="/" class="text-lg font-bold tracking-tight">Shopshredder Sandbox Platform</RouterLink>
        <p class="text-sm text-muted-foreground">Public demos for guests, full control for your internal team.</p>
      </div>
    </div>

    <div class="flex flex-wrap items-center gap-3">
      <RouterLink to="/">
        <Button variant="ghost">Storefront</Button>
      </RouterLink>
      <template v-if="auth.isAuthenticated">
        <RouterLink to="/admin">
          <Button variant="secondary">
            <ShieldCheck class="mr-2 h-4 w-4" />
            Admin
          </Button>
        </RouterLink>
        <span class="rounded-full bg-secondary px-3 py-2 text-sm font-medium text-secondary-foreground">
          {{ userLabel }}
        </span>
        <Button variant="outline" @click="logout">Logout</Button>
      </template>
      <template v-else>
        <RouterLink to="/login">
          <Button>Login</Button>
        </RouterLink>
      </template>
    </div>
  </header>
</template>
