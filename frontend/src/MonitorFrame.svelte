<script lang="ts">
  import { onDestroy, onMount } from 'svelte';

  type LedState = 'all_off' | 'all_on' | 'led_1' | 'led_2' | 'led_3';

  type Props = {
    hasError?: boolean;
    ledState?: LedState;
    keyStates?: boolean[];
    onKeyClick?: (index: number) => void;
    vertical?: boolean;
  };

  let {
    hasError = false,
    ledState = 'all_off',
    onKeyClick,
    vertical = false,
  }: Props = $props();

  const borderImage = '/border.png';

  const ledImages: Record<LedState, string> = {
    all_off: '/leds/all_off.png',
    all_on: '/leds/all_on.png',
    led_1: '/leds/led_1_on.png',
    led_2: '/leds/led_2_on.png',
    led_3: '/leds/led_3_on.png',
  };

  const keyImages = [
    {
      off: '/teclas/key_1_off.png',
      on: '/teclas/key_1_on.png',
    },
    {
      off: '/teclas/key_2_off.png',
      on: '/teclas/key_2_on.png',
    },
    {
      off: '/teclas/key_3_off.png',
      on: '/teclas/key_3_on.png',
    },
    {
      off: '/teclas/key_4_off.png',
      on: '/teclas/key_4_on.png',
    },
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
    currentRoute === '/login',
    currentRoute === '/chat',
    isFullscreen,
  ]);

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

  const currentLedImage = $derived(
    hasError
      ? showErrorImage
        ? ledImages.led_3
        : ledImages.all_off
      : ledImages[ledState],
  );
  function handleFullscreen() {
    if (document.fullscreenElement) {
      document.exitFullscreen();
    } else {
      document.documentElement.requestFullscreen();
    }
    onKeyClick?.(3);
  }
</script>

<div class="monitor-frame" class:vertical>
  <img class="layer border-layer" src={borderImage} alt="" />

  <div class="keys-layer">
    {#each keyImages as key, i}
      <img class="layer" src={currentKeyStates[i] ? key.on : key.off} alt="" />
    {/each}
  </div>

  <img class="layer leds-layer" src={currentLedImage} alt="" />

  <a
    class="key-hitbox key-1"
    href="#/"
    aria-label="Home"
    onclick={() => onKeyClick?.(0)}
  ></a>
  <a
    class="key-hitbox key-2"
    href="#/login"
    aria-label="Login"
    onclick={() => onKeyClick?.(1)}
  ></a>
  <a
    class="key-hitbox key-3"
    href="#/chat"
    aria-label="Chat"
    onclick={() => onKeyClick?.(2)}
  ></a>
  <button
    class="key-hitbox key-4"
    type="button"
    aria-label="Fullscreen"
    onclick={handleFullscreen}
  ></button>
</div>

<style>
  .monitor-frame {
    position: absolute;
    inset: 0;
    width: 100%;
    height: 100%;
    display: grid;
    place-items: center;
    z-index: 10000;
    pointer-events: none;
    overflow: hidden;
  }

  .monitor-frame.vertical {
    width: 177.7778%;
    height: auto;
    aspect-ratio: 16 / 9;
    inset: auto;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) rotate(90deg);
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

  .key-hitbox {
    position: absolute;
    z-index: 4;
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
