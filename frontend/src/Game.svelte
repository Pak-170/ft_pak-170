<script lang="ts">
  import { onMount } from 'svelte';
  import { Game, handleMouseMove } from './stores/game.ts';

  let canvas: HTMLCanvasElement;
  let wrapper: HTMLDivElement;
  let game: Game | undefined;

  function handleMouseMoveOuter(event: MouseEvent) {
    if (game) handleMouseMove(game, event, canvas);
  }

  function requestLock() {
    canvas.requestPointerLock();
  }

  onMount(() => {
    game = new Game(canvas, wrapper);
    const observer = new ResizeObserver(() => {
      if (!game) return;
      const width = wrapper.clientWidth;
      const height = wrapper.clientHeight;
      game.renderer.setSize(width, height);
    });
    observer.observe(wrapper);
    return () => observer.disconnect();
  });
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  bind:this={wrapper}
  onmousemove={handleMouseMoveOuter}
  onclick={requestLock}
>
  <canvas bind:this={canvas}></canvas>
</div>

<style>
  div {
    width: 100%;
    height: 100%;
    overflow: hidden;
    cursor: pointer;
  }
</style>
