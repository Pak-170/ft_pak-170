<script>
  import Router from 'svelte-spa-router';

  import Home from './Home.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';
  import MonitorFrame from './MonitorFrame.svelte';

  import { loginHasError, passwordMismatch } from './stores/loginStatus.ts';

  let keyStates = $state([true, false, false, false]);

  function activateKey(index) {
    keyStates = keyStates.map((_, i) => i === index);
  }

  const routes = {
    '/': Home,
    '/login': Login,
    '/chat': Chat,
  };
</script>

<div class="monitor-container">
  <div class="web-content">
    <Router {routes} />
  </div>

  <MonitorFrame
    hasError={$loginHasError || $passwordMismatch}
    ledState="all_off"
    {keyStates}
    onKeyClick={activateKey}
  />
</div>

<style>
  :global(html),
  :global(body) {
    margin: 0;
    width: 100%;
    min-height: 100%;
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
</style>
