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
  <a href="#/">Home</a>
  <a href="#/login">Login</a>
  <a href="#/chat">Chat</a>
</nav>

<Router {routes} />
