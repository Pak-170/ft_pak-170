<script>
  import Router from 'svelte-spa-router';
  import Home from './Home.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';

  const routes = {
    '/': Home,
    '/login': Login,
    '/chat': Chat,
  };
</script>

<nav>
  <a href="#/">HOME</a>
  <a href="#/login">LOGIN</a>
  <a href="#/chat">CHAT</a>
</nav>

<Router {routes} />
