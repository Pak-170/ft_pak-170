<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  import Home from './GameChat.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';
  import MonitorFrame from './MonitorFrame.svelte';

  import { loginHasError, passwordMismatch } from './stores/loginStatus.ts';

  let displayedRoute = $state('/');
  let isGameMode = $state(false);
  let targetGameMode = $state(false);
  let isTransitioning = $state(false);
  let pendingRoute = '/';
  let transitionTimeout: ReturnType<typeof setTimeout> | undefined;

  function getRoute() {
    return window.location.hash.slice(1).split('?')[0] || '/';
  }

  function updateRoute() {
    const nextRoute = getRoute();
    const nextGameMode = nextRoute === '/';
    if (nextRoute === displayedRoute && !isTransitioning) {
      return;
    }
    if (nextGameMode === isGameMode && !isTransitioning) {
      displayedRoute = nextRoute;
      return;
    }
    pendingRoute = nextRoute;
    targetGameMode = nextGameMode;
    isTransitioning = true;
    if (transitionTimeout !== undefined) {
      clearTimeout(transitionTimeout);
    }
    transitionTimeout = setTimeout(() => {
      displayedRoute = pendingRoute;
      isGameMode = pendingRoute === '/';
      targetGameMode = isGameMode;
      isTransitioning = false;
      transitionTimeout = undefined;
    }, 3000);
  }

  onMount(() => {
    const initialRoute = getRoute();
    displayedRoute = initialRoute;
    isGameMode = initialRoute === '/';
    targetGameMode = isGameMode;
    window.addEventListener('hashchange', updateRoute);
    return () => {
      window.removeEventListener('hashchange', updateRoute);
    };
  });

  onDestroy(() => {
    if (transitionTimeout !== undefined) {
      clearTimeout(transitionTimeout);
    }
  });
</script>

<div class="monitor-container" class:game-mode={isGameMode}>
  <div
    class="web-content"
    class:game-mode={isGameMode}
    class:transition-hidden={isTransitioning}
  >
    {#if displayedRoute === '/'}
      <Home />
    {:else if displayedRoute === '/login'}
      <Login />
    {:else if displayedRoute === '/chat'}
      <Chat />
    {/if}
  </div>

  <MonitorFrame
    hasError={$loginHasError || $passwordMismatch}
    vertical={targetGameMode}
    transitioning={isTransitioning}
  />
</div>

<style>
  :global(html),
  :global(body) {
    margin: 0;
    width: 100%;
    height: 100%;
  }

  :global(body) {
    min-height: 100vh;
    display: grid;
    place-items: center;
    overflow: hidden;
    background-image: url('/Bg_PAK.jpg');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
  }

  .monitor-container {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: min(100vw, calc(100dvh * 16 / 9));
    aspect-ratio: 16 / 9;
    overflow: visible;
  }

  .web-content {
    position: absolute;
    left: 5%;
    top: 7.5%;
    width: 80%;
    height: 84%;
    overflow: hidden;
    z-index: 1;
  }

  .web-content.transition-hidden {
    visibility: hidden;
  }

  .monitor-container.game-mode {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100dvh;
    aspect-ratio: auto;
    transform: none;
  }

  .web-content.game-mode {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
  }
</style>
