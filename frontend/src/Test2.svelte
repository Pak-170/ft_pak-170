<!-- docker compose up --build frontend -->
<script lang="ts">
	import { onMount } from 'svelte';
    import { handleMouseMove, initGame } from './game/game';

	let canvas: HTMLCanvasElement;
	let wrapper: HTMLDivElement;

	function handleMouseMoveOuter(event: MouseEvent)
	{
		handleMouseMove(event, canvas);
	}
	// Click en el lienzo y fuera la frecha del raton. Magic!
	function requestLock()
	{
		canvas.requestPointerLock();
	}

	onMount(() =>
	{	// adaptar juego a pagina abierta
		initGame(canvas, wrapper);
	});
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={wrapper} onmousemove={handleMouseMoveOuter} onclick={requestLock}>
  <canvas bind:this={canvas}></canvas>
</div>

<style type="scss">
  div {
	width: 100vw;
	height: 100vh;
	overflow: hidden;
	cursor: pointer;
  }
</style>