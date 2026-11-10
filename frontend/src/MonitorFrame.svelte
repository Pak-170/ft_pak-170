<script lang="ts">
  import { onDestroy, onMount } from 'svelte';
  import { loginLoading, loginSuccess } from './stores/loginStatus.ts';

  type Props = {
    hasError?: boolean;
    vertical?: boolean;
    transitioning?: boolean;
    settingsOpen?: boolean;
    onSettings?: () => void;
  };

  let {
    hasError = false,
    vertical = false,
    transitioning = false,
    settingsOpen = false,
    onSettings,
  }: Props = $props();

  const ledImages = {
    all_off: '/leds/all_off.png',
    all_on: '/leds/all_on.png',
    led_1: '/leds/led_1_on.png',
    led_2: '/leds/led_2_on.png',
    led_3: '/leds/led_3_on.png',
  };

  const keyImages = [
    { off: '/teclas/key_1_off.png', on: '/teclas/key_1_on.png' },
    { off: '/teclas/key_2_off.png', on: '/teclas/key_2_on.png' },
    { off: '/teclas/key_3_off.png', on: '/teclas/key_3_on.png' },
    { off: '/teclas/key_4_off.png', on: '/teclas/key_4_on.png' },
  ];

  let currentRoute = $state('/');
  let isFullscreen = $state(false);
  let showErrorImage = $state(false);
  let interval: ReturnType<typeof setInterval> | undefined;

  onMount(() => {
    function updateState() {
      currentRoute = window.location.hash.slice(1).split('?')[0] || '/';
      isFullscreen = document.fullscreenElement !== null;
    }
    updateState();
    window.addEventListener('hashchange', updateState);
    document.addEventListener('fullscreenchange', updateState);
    return () => {
      window.removeEventListener('hashchange', updateState);
      document.removeEventListener('fullscreenchange', updateState);
    };
  });
  const currentKeyStates = $derived([
    currentRoute === '/',
    currentRoute === '/login' || currentRoute === '/profile',
    settingsOpen,
    isFullscreen,
  ]);
  const currentLedImage = $derived(
    $loginLoading
      ? ledImages.led_2
      : $loginSuccess
        ? ledImages.led_1
        : hasError
          ? showErrorImage
            ? ledImages.led_3
            : ledImages.all_off
          : ledImages.all_off,
  );

  $effect(() => {
    if (hasError) {
      if (interval === undefined) {
        showErrorImage = true;
        interval = setInterval(() => {
          showErrorImage = !showErrorImage;
        }, 1000);
      }
    } else {
      if (interval !== undefined) {
        clearInterval(interval);
        interval = undefined;
      }
      showErrorImage = false;
    }
  });

  onDestroy(() => {
    if (interval !== undefined) {
      clearInterval(interval);
    }
  });

  function handleFullscreen() {
    if (document.fullscreenElement) {
      document.exitFullscreen();
    } else {
      document.documentElement.requestFullscreen();
    }
  }
</script>

<div class="monitor-frame" class:vertical class:transitioning>
  <img class="layer border-layer" src="/border.png" alt="" />

  <div class="keys-layer">
    {#each keyImages as key, i}
      <img class="layer" src={currentKeyStates[i] ? key.on : key.off} alt="" />
    {/each}
  </div>

  <img class="layer leds-layer" src={currentLedImage} alt="" />

  {#if transitioning}
    <video
      class="transition-video"
      src="/transition.mp4"
      autoplay
      muted
      playsinline
    ></video>
  {/if}

  <a class="key-hitbox key-1" href="#/" aria-label="Home"></a>

  <a class="key-hitbox key-2" href="#/login" aria-label="Login"></a>

  <button
    class="key-hitbox key-3"
    type="button"
    aria-label="Settings"
    aria-pressed={settingsOpen}
    onclick={onSettings}
  ></button>

  <button
    class="key-hitbox key-4"
    type="button"
    aria-label="Fullscreen"
    onclick={handleFullscreen}
  ></button>
</div>

<style>
  .monitor-frame {
    position: fixed;
    top: 50%;
    left: 50%;
    width: min(100vw, calc(100dvh * 16 / 9));
    aspect-ratio: 16 / 9;
    display: grid;
    place-items: center;
    transform: translate(-50%, -50%) rotate(0deg);
    transform-origin: center;
    z-index: 10000;
    pointer-events: none;
    overflow: visible;
    transition: none;
  }

  .monitor-frame.transitioning {
    transition:
      width 3s cubic-bezier(0.65, 0, 0.35, 1),
      transform 3s cubic-bezier(0.65, 0, 0.35, 1);
  }

  .monitor-frame.vertical {
    width: 100dvh;
    left: 50%;
    transform: translate(calc(-50% + 50vw - 28.125dvh), -50%) rotate(90deg);
  }

  .layer {
    grid-area: 1 / 1;
    width: 100%;
    height: 100%;
    object-fit: fill;
    pointer-events: none;
  }

  .border-layer {
    z-index: 1;
  }

  .keys-layer {
    grid-area: 1 / 1;
    width: 100%;
    height: 100%;
    display: grid;
    z-index: 2;
    pointer-events: none;
  }

  .leds-layer {
    z-index: 3;
  }

  .transition-video {
    position: absolute;
    left: 5%;
    top: 7.5%;
    width: 80%;
    height: 84%;
    object-fit: fill;
    z-index: 0;
    pointer-events: none;
  }

  .key-hitbox {
    position: absolute;
    z-index: 5;
    pointer-events: auto;
    border: none;
    cursor: pointer;
    background: transparent;
  }

  .key-1 {
    top: 34%;
    right: 3%;
    width: 6%;
    height: 9%;
  }

  .key-2 {
    top: 49%;
    right: 3%;
    width: 6%;
    height: 9%;
  }

  .key-3 {
    top: 65%;
    right: 3%;
    width: 6%;
    height: 9%;
  }

  .key-4 {
    top: 82%;
    right: 3%;
    width: 6%;
    height: 8%;
  }
</style>
