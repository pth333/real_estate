import { AllCommunityModule, ModuleRegistry } from 'ag-grid-community'

ModuleRegistry.registerModules([AllCommunityModule])

export default defineNuxtPlugin(() => {
  // AG Grid modules are registered at module level
})
