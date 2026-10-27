
<script>
  import Router from 'svelte-spa-router';
  import Home from './Home.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';
  import MonitorFrame from './MonitorFrame.svelte';
  import { loginHasError, passwordMismatch } from './stores/loginStatus.js';

  const routes = {
    '/': Home,
    '/login': Login,
    '/chat': Chat,
  };
</script>

<div class="monitor-content">
  <nav>
    <a href="#/">HOME</a>
    <a href="#/login">LOGIN</a>
    <a href="#/chat">CHAT</a>
  </nav>

  <main class="page-content">
    <Router {routes} />
  </main>
</div>

<MonitorFrame hasError={$loginHasError || $passwordMismatch} />

<style>
  .monitor-content {
    position: fixed;
    top: 8%;
    left: 6%;
    width: 79%;
    height: 81%;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    z-index: 1;
  }

  nav {
    position: relative;
    z-index: 2;
    display: flex;
    flex-shrink: 0;
    flex-wrap: wrap;
    gap: 10px;
    padding: 10px 20px;
    background: var(--surface);
    border-bottom: 1px solid var(--terminal-border, var(--green-dark));
  }

  nav a {
    color: var(--green);
    text-decoration: none;
    padding: 5px 10px;
    border: 1px solid transparent;
  }

  nav a:hover {
    border-color: var(--green);
    text-shadow: 0 0 5px var(--green);
  }

  .page-content {
    flex: 1;
    min-width: 0;
    min-height: 0;
    overflow: auto;
    position: relative;
  }

  @media (max-width: 480px) {
    nav {
      gap: 4px;
      padding: 6px;
    }

    nav a {
      padding: 4px 6px;
    }
  }
</style>
