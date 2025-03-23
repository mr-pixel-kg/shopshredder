import { defineStore } from "pinia";

export const GeneralStore = defineStore("general", {
  state: () => ({
    showLoadingScreen: false, // Initialer Zustand
    sandboxes: [],
  }),
  getters: {
    isLoading() {
      return this.showLoadingScreen;
    },
    getSandboxes() {
      return this.sandboxes;
    },
  },
  actions: {
    setLoading(value) {
      this.showLoadingScreen = value;
    },
    addSandbox(sandbox) {
      // Überprüfen, ob die Sandbox schon existiert, um Duplikate zu vermeiden
      if (!this.sandboxes.some((env) => env.id === sandbox.id)) {
        this.sandboxes.push(sandbox);
      }
    },
    removeSandbox(sandboxId) {
      this.sandboxes = this.sandboxes.filter((env) => env.id !== sandboxId);
    },
  },
});
