<!-- He creado en home el canvas para Paolo, he puesto un cuadrado rojo solo para test-->
<script lang="ts">
  import * as THREE from 'three';
  import { onMount } from 'svelte';

  let { speed = 0 } = $props();

  let canvas: HTMLCanvasElement;
  let wrapper: HTMLDivElement;

  let mouse_y = 0;

  function mouseMoveHandler(event: MouseEvent) {
  }

  onMount(() => {
  canvas.addEventListener("mousemove", (event) => {
    mouse_x = event.clientX
    mouse_y = event.clientY
  })
  let mouse_x = 0;
    const width = wrapper.getBoundingClientRect().width;
    const height = wrapper.getBoundingClientRect().height;

    const scene = new THREE.Scene();

    const camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);

    const renderer = new THREE.WebGLRenderer({ canvas, antialias: true });
    renderer.setSize(width, height);

    const color = 0xffffff;
    const intensity = 3;
    const light = new THREE.DirectionalLight(color, intensity);
    light.position.set(-1, 2, 4);
    scene.add(light);

    const geometry = new THREE.BoxGeometry(10, 1, 10);
    const material = new THREE.MeshPhongMaterial({ color: '#ff0000' });
    const cube = new THREE.Mesh(geometry, material);
    scene.add(cube);

    camera.position.z = 25;

    function animate() {
      requestAnimationFrame(animate);
      cube.rotation.x = 0.01 * mouse_x;
      cube.rotation.y = 0.01 * mouse_y;
      // cube.position.x = 0.01 * mouse_x;
      // cube.position.y = 0.01 * - mouse_y;
      renderer.render(scene, camera);
    }

    animate();
  });
</script>

<div bind:this={wrapper}>
  <canvas bind:this={canvas} ></canvas>
</div>

<style type="scss">
  div {
    width: 100vw;
    height: 100vh;
  }
</style>
