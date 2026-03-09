import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/main.css'
import { AllCommunityModule, ModuleRegistry } from 'ag-grid-community'

// Register all Community features
ModuleRegistry.registerModules([AllCommunityModule])

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
