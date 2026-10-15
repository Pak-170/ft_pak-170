<script>
  import Router from 'svelte-spa-router';
  import TestCube from './TestCube.svelte';
  import Test2 from './Test2.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';

  const routes = {
    '/': TestCube,
    '/login': Login,
    '/chat': Chat,
    '/test2': Test2,
  };
</script>

<nav>
  <a href="#/">TestCube</a>
  <a href="#/test2">Test2</a>
  <a href="#/login">Login</a>
  <a href="#/chat">Chat</a>
</nav>

<Router {routes} />
