<script lang="ts">
  import Router from 'svelte-spa-router';
  import { onMount } from 'svelte';

  import Home from './GameChat.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';
  import MonitorFrame from './MonitorFrame.svelte';

  import { loginHasError, passwordMismatch } from './stores/loginStatus.ts';

  let keyStates = $state([true, false, false, false]);

  let isGameMode = $state(false);

  onMount(() => {
    function updateRoute() {
      const route = window.location.hash.slice(1).split('?')[0] || '/';
      isGameMode = route === '/';
    }
    updateRoute();
    window.addEventListener('hashchange', updateRoute);
    return () => {
      window.removeEventListener('hashchange', updateRoute);
    };
  });
  function activateKey(index: number) {
    keyStates = keyStates.map((_, i) => i === index);
  }

  const routes = {
    '/': Home,
    '/login': Login,
    '/chat': Chat,
  };
</script>

<div class="monitor-container" class:game-mode={isGameMode}>
  <div class="web-content" class:game-mode={isGameMode}>
    <Router {routes} />
  </div>

  {#if !isGameMode}
    <MonitorFrame
      hasError={$loginHasError || $passwordMismatch}
      ledState="all_off"
      {keyStates}
      onKeyClick={activateKey}
    />
  {/if}
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
    position: relative;
    width: min(100vw, calc(100vh * 1280 / 720));
    aspect-ratio: 1280 / 720;
    overflow: hidden;
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

  .monitor-container.game-mode {
    position: fixed;
    inset: 0;
    width: 100vw;
    height: 100dvh;
    aspect-ratio: auto;
  }

  .web-content.game-mode {
    inset: 0;
    width: 100%;
    height: 100%;
  }
</style>
