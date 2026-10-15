<script lang="ts">
  import * as THREE from 'three';
  import { onMount } from 'svelte';

  let { speed = 0 } = $props();

  let canvas: HTMLCanvasElement;
  let wrapper: HTMLDivElement;

  let mouse_x = 0;
  let mouse_y = 0;

  let width = 0;
  let height = 0;

  function handleMouseMove(event: MouseEvent) {
	if (!wrapper) return;
	const rect = wrapper.getBoundingClientRect();
	// Coordenadas relativas al centro del wrapper/canvas
	mouse_x = event.clientX - (rect.left + rect.width / 2);
	mouse_y = event.clientY - (rect.top + rect.height / 2);
  }

  onMount(() => {
	width = wrapper.getBoundingClientRect().width;
	height = wrapper.getBoundingClientRect().height;

	const scene = new THREE.Scene();
	const camera = new THREE.PerspectiveCamera(40, width / height, 0.1, 1000);

	const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
	renderer.setSize(width, height);
	renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2)); // Mejora la nitidez en pantallas Retina

	const color = 0xffffff;
	const intensity = 3;
	const light = new THREE.DirectionalLight(color, intensity);
	light.position.set(0.1, 0.1, 0.1);
	scene.add(light);
	
	const ambientLight = new THREE.AmbientLight(0xffffff, 0.2); 
	scene.add(ambientLight);

	const geometry = new THREE.BoxGeometry(10, 10, 10);
	const material = new THREE.MeshPhongMaterial({ color: '#f0ff00' });


// 1. Cargador de texturas
	const textureLoader = new THREE.TextureLoader();

// 2. Array de 6 materiales
// 2. Creamos materiales iniciales (con un color para evitar el negro mientras cargan)
		const materials =
		[
			new THREE.MeshPhongMaterial({ color: 0xffffff }),
			new THREE.MeshPhongMaterial({ color: 0xffffff }),
			new THREE.MeshPhongMaterial({ color: 0xffffff }),
			new THREE.MeshPhongMaterial({ color: 0xffffff }),
			new THREE.MeshPhongMaterial({ color: 0xffffff }),
			new THREE.MeshPhongMaterial({ color: 0xffffff })
		];

		// 3. Rutas de las imágenes
		const texturePaths =
		[
			'cube_texture/00.png',
			'cube_texture/01.png',
			'cube_texture/02.png',
			'cube_texture/03.png',
			'cube_texture/04.png',
			'cube_texture/05.png'
		];

// 4. no entiendo esto:
		texturePaths.forEach((path, index) =>{
			textureLoader.load(path, (texture) => {
				materials[index].map = texture;
				materials[index].needsUpdate = true;
			});
		});

		// 3. Le pasamos el array de materiales al Mesh
	const cube = new THREE.Mesh(geometry, materials);
	scene.add(cube);

	camera.position.z = 60;

	function animate() {
		requestAnimationFrame(animate);
		
		// Interpolación suave opcional o mapeo directo
		cube.rotation.x = 0.005 * mouse_y;
		cube.rotation.y = 0.005 * mouse_x;
		cube.position.x = 0.05 * mouse_x;
		cube.position.y = 0.05 * - mouse_y;
	
		renderer.render(scene, camera);
	}

	animate();

	// Opcional: Manejar el redimensionado de ventana
	function handleResize() {
	  if (!wrapper) return;
	  width = wrapper.getBoundingClientRect().width;
	  height = wrapper.getBoundingClientRect().height;
	  camera.aspect = width / height;
	  camera.updateProjectionMatrix();
	  renderer.setSize(width, height);
	}

	window.addEventListener('resize', handleResize);

	return () => {
	  window.removeEventListener('resize', handleResize);
	};
  });
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={wrapper} onmousemove={handleMouseMove}>
  <canvas bind:this={canvas}></canvas>
</div>

<style type="scss">
  div {
	width: 100vw;
	height: 100vh;
	overflow: hidden;
  }
</style>