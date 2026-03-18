<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import AppHeader from "@/components/layout/AppHeader.vue";
import Card from "@/components/ui/Card.vue";
import { api } from "@/lib/api";
import type { ImageRecord, SandboxRecord } from "@/types/api";
import PublicImageCard from "@/features/public/PublicImageCard.vue";
import GuestSandboxList from "@/features/public/GuestSandboxList.vue";

const images = ref<ImageRecord[]>([]);
const sandboxes = ref<SandboxRecord[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const creatingId = ref<string | null>(null);
const deletingId = ref<string | null>(null);

const heroCount = computed(() => `${images.value.length} demo templates`);

async function load() {
  loading.value = true;
  error.value = null;

  try {
    const [publicImages, guestSandboxes] = await Promise.all([
      api.getPublicImages(),
      api.getGuestSandboxes(),
    ]);
    images.value = publicImages;
    sandboxes.value = guestSandboxes;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not load storefront";
  } finally {
    loading.value = false;
  }
}

async function startDemo(imageId: string) {
  creatingId.value = imageId;
  try {
    await api.createDemo({ imageId });
    await load();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not create demo";
  } finally {
    creatingId.value = null;
  }
}

async function deleteGuestSandbox(id: string) {
  deletingId.value = id;
  try {
    await api.deleteGuestSandbox(id);
    await load();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not delete demo sandbox";
  } finally {
    deletingId.value = null;
  }
}

onMounted(load);
</script>

<template>
  <div class="app-shell">
    <AppHeader />

    <section class="hero-panel mb-8 px-6 py-10 md:px-10">
      <div class="grid gap-10 lg:grid-cols-[minmax(0,1fr)_320px] lg:items-end">
        <div class="space-y-5">
          <span class="inline-flex rounded-full bg-accent/15 px-3 py-1 text-sm font-semibold text-accent">
            Public storefront
          </span>
          <div class="space-y-3">
            <h1 class="text-4xl font-extrabold tracking-tight text-foreground md:text-6xl">
              Launch guest-ready Docker demos in one click.
            </h1>
            <p class="max-w-2xl text-base text-muted-foreground md:text-lg">
              Browse public templates, create time-limited sandbox instances, and jump straight into live demos behind Traefik-managed URLs.
            </p>
          </div>
        </div>

        <Card class="bg-white/70">
          <div class="space-y-3">
            <p class="text-sm font-semibold uppercase tracking-[0.2em] text-muted-foreground">Live overview</p>
            <div class="text-3xl font-extrabold">{{ heroCount }}</div>
            <p class="text-sm text-muted-foreground">
              {{ sandboxes.length }} active demo sandboxes tracked for this browser session.
            </p>
          </div>
        </Card>
      </div>
    </section>

    <div v-if="error" class="mb-6 rounded-2xl border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      {{ error }}
    </div>

    <div class="grid gap-8">
      <section class="space-y-5">
        <div class="flex items-end justify-between gap-4">
          <div>
            <h2 class="section-title">Public images</h2>
            <p class="text-sm text-muted-foreground">Start a fresh demo sandbox from any published image template.</p>
          </div>
        </div>

        <div v-if="loading" class="rounded-2xl border border-dashed border-border bg-secondary/30 px-4 py-12 text-center text-sm text-muted-foreground">
          Loading public image catalog...
        </div>

        <div v-else class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
          <PublicImageCard
            v-for="image in images"
            :key="image.id"
            :image="image"
            :busy="creatingId === image.id"
            @demo="startDemo"
          />
        </div>
      </section>

      <GuestSandboxList :sandboxes="sandboxes" :deleting-id="deletingId" @remove="deleteGuestSandbox" />
    </div>
  </div>
</template>
