<script lang="ts">
  import { onDestroy } from 'svelte';

  type LedState = 'all_off' | 'all_on' | 'led_1' | 'led_2' | 'led_3';

  type Props = {
    hasError?: boolean;
    ledState?: LedState;
    keyStates?: boolean[];
  };

  let {
    hasError = false,
    ledState = 'all_off',
    keyStates = [false, false, false, false],
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
</script>

<div class="monitor-frame">
  <img class="layer border" src={borderImage} alt="" />

  <img class="layer leds" src={currentLedImage} alt="" />

  {#each keyImages as key, i}
    <img class="layer key" src={keyStates[i] ? key.on : key.off} alt="" />
  {/each}
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
</style>
