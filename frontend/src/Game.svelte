<script lang="ts">
	import { onMount } from 'svelte';
	import { Game, handleMouseMove } from './stores/game.ts';

	type Vect2 = {x : number, y : number};
	let canvas: HTMLCanvasElement;
	let wrapper: HTMLDivElement;
	let game: Game | undefined;

	function handleMouseMoveOuter(event: MouseEvent) {
		if (game) handleMouseMove(game, event, canvas);
	}

	function requestLock() {
		canvas.requestPointerLock();
	}

	onMount(() =>{
		const socket = new WebSocket("/api/game");
		game = new Game(canvas, wrapper);

		socket.addEventListener("message", (event) =>
		{
				// console.log(JSON.parse(event.data));

			const data_backend = JSON.parse(event.data);

			let table_size: Vect2 = {x: data_backend.width, y: data_backend.height};
			let pos_p1: Vect2 = data_backend.pos_p1;
			let pos_p2: Vect2 = data_backend.pos_p2;
			let pos_ball: Vect2 = data_backend.pos_ball;
			let ball_radius: number = data_backend.ball_radius;

			game!.Disco.position.x = pos_ball.x;
			game!.Disco.position.x = pos_ball.x;
		})

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
