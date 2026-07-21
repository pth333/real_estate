<template>
  <div class="flex flex-col min-h-screen">
    <AppHeader />
    <main class="flex-1 p-3">
      <slot />
    </main>
    <AppFooter />
  </div>
</template>
