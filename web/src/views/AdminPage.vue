<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import AppHeader from "@/components/layout/AppHeader.vue";
import Button from "@/components/ui/Button.vue";
import Card from "@/components/ui/Card.vue";
import { api } from "@/lib/api";
import { useAuthStore } from "@/stores/auth";
import type { AuditLogRecord, CreateImagePayload, CreateSandboxPayload, ImageRecord, SandboxRecord } from "@/types/api";
import ImagesManager from "@/features/admin/ImagesManager.vue";
import SandboxesManager from "@/features/admin/SandboxesManager.vue";
import AuditLogManager from "@/features/admin/AuditLogManager.vue";

type AdminTab = "images" | "sandboxes" | "audit";

const auth = useAuthStore();
const router = useRouter();

const tab = ref<AdminTab>("images");
const images = ref<ImageRecord[]>([]);
const sandboxes = ref<SandboxRecord[]>([]);
const logs = ref<AuditLogRecord[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const creatingImage = ref(false);
const creatingSandbox = ref(false);
const deletingImageId = ref<string | null>(null);
const deletingSandboxId = ref<string | null>(null);
const snapshottingSandboxId = ref<string | null>(null);

const token = computed(() => auth.token);

async function load() {
  if (!token.value) return;

  loading.value = true;
  error.value = null;

  try {
    const [imageResult, sandboxResult, auditResult] = await Promise.all([
      api.getImages(token.value),
      api.getSandboxes(token.value),
      api.getAuditLogs(token.value),
    ]);

    images.value = imageResult;
    sandboxes.value = sandboxResult;
    logs.value = auditResult;
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not load admin data";
  } finally {
    loading.value = false;
  }
}

async function createImage(payload: CreateImagePayload) {
  if (!token.value) return;
  creatingImage.value = true;
  try {
    await api.createImage(token.value, payload);
    await load();
    tab.value = "images";
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not create image";
  } finally {
    creatingImage.value = false;
  }
}

async function deleteImage(id: string) {
  if (!token.value) return;
  deletingImageId.value = id;
  try {
    await api.deleteImage(token.value, id);
    await load();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not delete image";
  } finally {
    deletingImageId.value = null;
  }
}

async function createSandbox(payload: CreateSandboxPayload) {
  if (!token.value) return;
  creatingSandbox.value = true;
  try {
    await api.createSandbox(token.value, payload);
    await load();
    tab.value = "sandboxes";
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not create sandbox";
  } finally {
    creatingSandbox.value = false;
  }
}

async function deleteSandbox(id: string) {
  if (!token.value) return;
  deletingSandboxId.value = id;
  try {
    await api.deleteSandbox(token.value, id);
    await load();
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not delete sandbox";
  } finally {
    deletingSandboxId.value = null;
  }
}

async function snapshotSandbox(id: string, payload: CreateImagePayload) {
  if (!token.value) return;
  snapshottingSandboxId.value = id;
  try {
    await api.snapshotSandbox(token.value, id, payload);
    await load();
    tab.value = "images";
  } catch (err) {
    error.value = err instanceof Error ? err.message : "Could not create snapshot";
  } finally {
    snapshottingSandboxId.value = null;
  }
}

function logout() {
  auth.logout();
  router.push("/login");
}

onMounted(load);
</script>

<template>
  <div class="app-shell">
    <AppHeader />

    <section class="hero-panel mb-8 px-6 py-8 md:px-8">
      <div class="flex flex-col gap-5 md:flex-row md:items-end md:justify-between">
        <div>
          <span class="inline-flex rounded-full bg-primary/15 px-3 py-1 text-sm font-semibold text-primary">Admin workspace</span>
          <h1 class="mt-4 text-3xl font-extrabold tracking-tight md:text-5xl">Control images, sandboxes and audit trails.</h1>
          <p class="mt-3 max-w-3xl text-sm text-muted-foreground md:text-base">
            This area is designed for internal users who need full visibility into public demos and editable internal sandbox environments.
          </p>
        </div>
        <Button variant="outline" @click="logout">Logout</Button>
      </div>
    </section>

    <div v-if="error" class="mb-6 rounded-2xl border border-danger/20 bg-danger/10 px-4 py-3 text-sm text-danger">
      {{ error }}
    </div>

    <div class="mb-6 flex flex-wrap gap-2">
      <Button :variant="tab === 'images' ? 'primary' : 'outline'" @click="tab = 'images'">Images</Button>
      <Button :variant="tab === 'sandboxes' ? 'primary' : 'outline'" @click="tab = 'sandboxes'">Sandboxes</Button>
      <Button :variant="tab === 'audit' ? 'primary' : 'outline'" @click="tab = 'audit'">Audit log</Button>
    </div>

    <Card v-if="loading" class="text-sm text-muted-foreground">Loading admin data...</Card>

    <div v-else class="space-y-6">
      <ImagesManager
        v-if="tab === 'images'"
        :images="images"
        :creating="creatingImage"
        :deleting-id="deletingImageId"
        @create="createImage"
        @remove="deleteImage"
      />

      <SandboxesManager
        v-else-if="tab === 'sandboxes'"
        :images="images"
        :sandboxes="sandboxes"
        :creating="creatingSandbox"
        :deleting-id="deletingSandboxId"
        :snapshotting-id="snapshottingSandboxId"
        @create="createSandbox"
        @remove="deleteSandbox"
        @snapshot="snapshotSandbox"
      />

      <AuditLogManager v-else :logs="logs" />
    </div>
  </div>
</template>
