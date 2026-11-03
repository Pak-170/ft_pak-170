<script lang="ts">
  import { onDestroy } from 'svelte';

  type LedState = 'all_off' | 'all_on' | 'led_1' | 'led_2' | 'led_3';

  type Props = {
    hasError?: boolean;
    ledState?: LedState;
    keyStates?: boolean[];
    onKeyClick?: (index: number) => void;
  };

  let {
    hasError = false,
    ledState = 'all_off',
    keyStates = [true, false, false, false],
    onKeyClick,
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

  let showErrorImage = $state(false);
  let interval: ReturnType<typeof setInterval> | undefined;

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

<div class="monitor-frame">
  <img class="layer" src={borderImage} alt="" />
  <img class="layer" src={currentLedImage} alt="" />
  {#each keyImages as key, i}
    <img
      class="layer"
      src={keyStates[i] ? key.on : key.off}
      alt=""
    />
  {/each}
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

  .layer {
    grid-area: 1 / 1;
    width: 100%;
    height: 100%;
    object-fit: fill;
    pointer-events: none;
  }

  .key-hitbox {
    position: absolute;
    z-index: 2;
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
