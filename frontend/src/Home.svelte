<!-- He creado en home el canvas para Paolo, he puesto un cuadrado rojo solo para test-->
<script lang="ts">
  import * as THREE from 'three';
  import { onMount } from 'svelte';

  let { speed = 1 } = $props();

  const width = 500;
  const height = 500;

  let canvas: HTMLCanvasElement;

  onMount(() => {
    canvas.height = height;
    canvas.width = width;
    const scene = new THREE.Scene();

    const camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    renderer.setSize(width, height);

    const color = 0xffffff;
    const intensity = 3;
    const light = new THREE.DirectionalLight(color, intensity);
    light.position.set(-1, 2, 4);
    scene.add(light);

    const geometry = new THREE.BoxGeometry(10, 10, 10);
    const material = new THREE.MeshPhongMaterial({ color: '#f0ff00' });
    const cube = new THREE.Mesh(geometry, material);
    scene.add(cube);

    camera.position.z = 25;

    function animate() {
      requestAnimationFrame(animate);
      cube.rotation.x += 0.01 * speed;
      cube.rotation.y += 0.01 * speed;
      renderer.render(scene, camera);
    }

    animate();
  });
</script>

<canvas bind:this={canvas}></canvas>
