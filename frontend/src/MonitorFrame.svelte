<script lang="ts">
  import { onDestroy } from 'svelte';

  let { hasError = false } = $props();

  const normalImage = '00.led_off_Personalizzato.png';
  const errorImage = '03.led_red_Personalizzato.png';

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
</script>

<div class="monitor-frame">
  <img
    src={showErrorImage ? errorImage : normalImage}
    alt=""
  />
</div>

<style>
  .monitor-frame {
    position: fixed;
    inset: 0;
    width: 100%;
    height: 100%;
    z-index: 10000;
    pointer-events: none;
  }

  .monitor-frame img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: fill;
    pointer-events: none;
  }
</style>
