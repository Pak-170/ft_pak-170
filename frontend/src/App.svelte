<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  import Home from './GameChat.svelte';
  import Login from './Login.svelte';
  import Chat from './Chat.svelte';
  import Profile from './Profile.svelte';
  import SettingsModal from './SettingsModal.svelte';
  import MonitorFrame from './MonitorFrame.svelte';

  import {
    currentUser,
    notifications,
    chatPageActive,
    selectedChatID,
    markChatAsRead,
    dismissNotification,
    clearMessageNotifications,
  } from '../client';

  import { loginHasError, passwordMismatch } from './stores/loginStatus.ts';

  let settingsOpen = $state(false);
  let displayedRoute = $state('/');
  let isGameMode = $state(false);
  let targetGameMode = $state(false);
  let isTransitioning = $state(false);
  let authChecked = $state(false);
  let pendingRoute = '/';
  let transitionTimeout: ReturnType<typeof setTimeout> | undefined;

  function toggleSettings() {
    settingsOpen = !settingsOpen;
  }

  function getRoute() {
    return window.location.hash.slice(1).split('?')[0] || '/';
  }

  function updateChatPageActive(route: string) {
    const isChatPage = route === '/';
    chatPageActive.set(isChatPage);
    if (isChatPage) {
      clearMessageNotifications();
    }
  }

  function updateRoute() {
    const nextRoute = getRoute();

    if (!authChecked) return;
    if (nextRoute === '/login' && $currentUser) {
      window.location.hash = '#/profile';
      return;
    }
    if (nextRoute === '/profile' && !$currentUser) {
      window.location.hash = '#/login';
      return;
    }
    updateChatPageActive(nextRoute);
    if (nextRoute !== displayedRoute) {
      settingsOpen = false;
    }
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
  function openNotification(chatID: number) {
    selectedChatID.set(chatID);
    markChatAsRead(chatID);
    window.location.hash = '#/';
  }

  onMount(() => {
    window.addEventListener('hashchange', updateRoute);
    async function checkSession() {
      try {
        const response = await fetch('/api/me');
        if (response.ok) {
          const user = await response.json();
          currentUser.set(user);
        } else {
          currentUser.set(null);
        }
      } catch (error) {
        console.error('Session check failed:', error);
        currentUser.set(null);
      }
      authChecked = true;
      const route = getRoute();
      if (route === '/login' && $currentUser) {
        window.location.hash = '#/profile';
        return;
      }
      if (route === '/profile' && !$currentUser) {
        window.location.hash = '#/login';
        return;
      }
      updateChatPageActive(route);
      displayedRoute = route;
      isGameMode = route === '/';
      targetGameMode = isGameMode;
    }
    void checkSession();
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
    {#if authChecked}
      {#if displayedRoute === '/'}
        <Home />
      {:else if displayedRoute === '/login'}
        <Login />
      {:else if displayedRoute === '/profile'}
        <Profile />
      {:else if displayedRoute === '/chat'}
        <Chat />
      {/if}
    {/if}
    {#if settingsOpen}
      <SettingsModal onClose={() => (settingsOpen = false)} />
    {/if}
  </div>

  <MonitorFrame
    hasError={$loginHasError || $passwordMismatch}
    vertical={targetGameMode}
    transitioning={isTransitioning}
    {settingsOpen}
    onSettings={toggleSettings}
  />
</div>
<div class="notifications">
  {#each $notifications as notification (notification.id)}
    <div
      class="notification"
      class:error={notification.type === 'error'}
      class:message-notification={notification.type === 'message'}
    >
      {#if notification.type === 'message' && notification.chatID !== undefined}
        <button
          class="notification-message"
          onclick={() => openNotification(notification.chatID!)}
        >
          {notification.message}
        </button>
      {:else}
        <span class="notification-message">
          {notification.message}
        </span>
      {/if}
      <button
        class="notification-close"
        aria-label="Dismiss notification"
        onclick={() => dismissNotification(notification.id)}
      >
        ×
      </button>
    </div>
  {/each}
</div>

<style>
  .monitor-container {
    --unit: min(calc(100dvh - 32px), calc((100vw - 32px) / 1.5625));
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: min(calc(100vw - 32px), calc((100dvh - 32px) * 16 / 9));
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
    width: calc(var(--unit) * 1.5625);
    height: var(--unit);
    aspect-ratio: auto;
  }

  .web-content.game-mode {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
  }

  :global(:root:fullscreen) .monitor-container {
    --unit: min(100dvh, calc(100vw / 1.5625));
    top: 0;
    left: 0;
    width: 100vw;
    height: 100dvh;
    transform: none;
    aspect-ratio: auto;
  }

  :global(:root:fullscreen) .monitor-container.game-mode {
    width: 100vw;
    height: 100dvh;
  }

  .notifications {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 20000;
    display: flex;
    flex-direction: column;
    gap: 10px;
    width: min(360px, calc(100vw - 40px));
    pointer-events: none;
  }

  .notification {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px;
    border: 1px solid var(--green-dark);
    background: var(--surface);
    box-shadow: 0 0 12px var(--scanline-glow);
    pointer-events: auto;
  }

  .notification.error {
    border-color: #ff5555;
  }

  .notification.message-notification {
    border-color: var(--green-light);
  }

  .notification-message {
    flex: 1;
    min-width: 0;
    padding: 0;
    border: none;
    color: var(--green-bright);
    background: transparent;
    text-align: left;
    overflow-wrap: anywhere;
    cursor: pointer;
  }

  .notification.error .notification-message {
    color: #ff7777;
  }

  .notification-close {
    flex-shrink: 0;
    padding: 2px 6px;
    border: none;
    color: var(--green-bright);
    background: transparent;
    cursor: pointer;
  }
</style>
