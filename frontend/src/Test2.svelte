
<!-- docker compose up --build frontend -->
<script lang="ts">
	import * as THREE from 'three';
	import { onMount } from 'svelte';

	let { speed = 0 } = $props();

	let canvas: HTMLCanvasElement;
	let wrapper: HTMLDivElement;

	let width = 0;
	let height = 0;

	// Límites físicos directos del tablero del pad 1
	const PAD1_MIN_X = -4.1;
	const PAD1_MAX_X = 4.1;
	const PAD1_MIN_Z = 5.0; // Línea central
	const PAD1_MAX_Z = 9.5; // Fondo de pista (portería propia)

	// Límites físicos directos del tablero del pad 2
	const PAD2_MIN_X = -4.1;
	const PAD2_MAX_X = 4.1;
	const PAD2_MIN_Z = -9.5; // Línea central
	const PAD2_MAX_Z = -5; // Fondo de pista (portería propia)


	// Sensibilidad del movimiento (Ajusta este valor a tu gusto)
	const SENSITIVITY = 0.015;

	let PadP1: THREE.Mesh | null = null;
	let PadP2: THREE.Mesh | null = null;

	function handleMouseMove(event: MouseEvent)
	{
		// Solo movemos el pad si el puntero está bloqueado/atrapado dentro del canvas
		if (document.pointerLockElement === canvas && PadP1 && !event.ctrlKey)
		{
			// Sumamos el desplazamiento del ratón directamente a las coordenadas 3D del Pad
			let nextX1 = PadP1.position.x + event.movementX * SENSITIVITY;
			let nextZ1 = PadP1.position.z + event.movementY * SENSITIVITY;
			// Límites directos e infranqueables (Clamp)
			PadP1.position.x = Math.max(PAD1_MIN_X, Math.min(PAD1_MAX_X, nextX1));
			PadP1.position.z = Math.max(PAD1_MIN_Z, Math.min(PAD1_MAX_Z, nextZ1));
		}
		if (event.ctrlKey && PadP2)
		{
			// Sumamos el desplazamiento del ratón directamente a las coordenadas 3D del Pad
			let nextX2 = PadP2.position.x - event.movementX * SENSITIVITY;
			let nextZ2 = PadP2.position.z - event.movementY * SENSITIVITY;
			// Límites directos e infranqueables (Clamp)
			PadP2.position.x = Math.max(PAD2_MIN_X, Math.min(PAD2_MAX_X, nextX2));
			PadP2.position.z = Math.max(PAD2_MIN_Z, Math.min(PAD2_MAX_Z, nextZ2));
		}
	}

	// Flipa que con esto nos capturamos el raton!Al hacer clic en el lienzo del juego
	function requestLock()
	{
		canvas.requestPointerLock();
	}

	onMount(() => {
		width = wrapper.getBoundingClientRect().width;
		height = wrapper.getBoundingClientRect().height;

		// 1. ESCENA Y RENDERIZADOR
		const scene = new THREE.Scene();
		const camera = new THREE.PerspectiveCamera(60, width / height, 0.1, 1000);
		
		const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
		renderer.setSize(width, height);
		renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));

		// 2. ILUMINACIÓN
		const ambientLight = new THREE.AmbientLight(0xffffff, 0.3); 
		scene.add(ambientLight);
		
		const dirLight = new THREE.DirectionalLight(0xffffff, 0.9);
		dirLight.position.set(5, 12, 8);
		scene.add(dirLight);

		// 3. ESTRUCTURA DE LA MESA (GRUPO)
		const mesaGroup = new THREE.Group();

		// 3.1 Tablero de juego
		const geometryMesa = new THREE.BoxGeometry(10.5, 0.5, 20);
		const materialMesa = new THREE.MeshStandardMaterial({ color: 0xff0000, roughness: 0.3 }); 
		const mesa = new THREE.Mesh(geometryMesa, materialMesa);
		mesaGroup.add(mesa);

		// Material para los bordes metálicos
		const materialBorde = new THREE.MeshStandardMaterial({ color: 0x8a9597, metalness: 0.4, roughness: 0.5 });

		// 3.2 Paredes laterales
		const geometryParedLarga = new THREE.BoxGeometry(1, 1.5, 20);
		
		const paredIzq = new THREE.Mesh(geometryParedLarga, materialBorde);
		paredIzq.position.set(-5.7, 0.5, 0);
		
		const paredDer = new THREE.Mesh(geometryParedLarga, materialBorde);
		paredDer.position.set(5.7, 0.5, 0);

		// 3.3 Pad Player 1 (Azul)
		const geometryPadP1 = new THREE.BoxGeometry(2, 0.5, 0.7);
		const materialPadP1 = new THREE.MeshStandardMaterial({ color: 0x0000ff, roughness: 0.3 }); 
		PadP1 = new THREE.Mesh(geometryPadP1, materialPadP1);
		PadP1.position.set(0, 0.5, 8);

		// 3.3 Pad Player 2 (verde)
		const geometryPadP2 = new THREE.BoxGeometry(2, 0.5, 0.7);
		const materialPadP2 = new THREE.MeshStandardMaterial({ color: 0x00ff00, roughness: 0.3 }); 
		PadP2 = new THREE.Mesh(geometryPadP2, materialPadP2);
		PadP2.position.set(0, 0.5, -8);



		mesaGroup.add(paredIzq);
		mesaGroup.add(paredDer);
		mesaGroup.add(PadP1);
		mesaGroup.add(PadP2);

		scene.add(mesaGroup);

		// 4. ADAPTABILIDAD DE CÁMARA
		function adaptarCamara() {
			if (width < height) {
				camera.position.set(0, 10, 22); 
			} else {
				camera.position.set(0, 6, 14); 
			}
			camera.lookAt(0, 0, 0);
		}

		adaptarCamara();

		// 5. BUCLE DE ANIMACIÓN
		function animate()
		{
			requestAnimationFrame(animate);
			renderer.render(scene, camera);
		}
		animate();

		// 6. EVENTO RESPONSIVE
		function handleResize() {
			if (!wrapper) return;
			width = wrapper.getBoundingClientRect().width;
			height = wrapper.getBoundingClientRect().height;
			
			camera.aspect = width / height;
			adaptarCamara();
			camera.updateProjectionMatrix();
			
			renderer.setSize(width, height);
		}

		window.addEventListener('resize', handleResize);

		return () => {
			window.removeEventListener('resize', handleResize);
		};
	});
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={wrapper} onmousemove={handleMouseMove} onclick={requestLock}>
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