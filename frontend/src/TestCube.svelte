<script lang="ts">
    import * as THREE from 'three';
    import { onMount } from 'svelte';

    // Props de Svelte 5 con las runas
    let { speed = 0 } = $props();

    let canvas: HTMLCanvasElement;
    let wrapper: HTMLDivElement;

    // Estado reactivo: indica si el ratón está atrapado (Pointer Lock activo)
    let isLocked = $state(false);

    let mouse_x = 0;
    let mouse_y = 0;

    // Memoria de la rotación acumulada: la cara donde se queda no se resetea jamás
    let targetRotationX = 0;
    let targetRotationY = 0;

    let width = 0;
    let height = 0;

    function handleMouseMove(event: MouseEvent) {
        if (!wrapper) return;

        if (isLocked) {
            // Cuando el ratón está TRAPADO:
            // Usamos event.movementX/Y que nos da el desplazamiento relativo puro sin límites de pantalla.
            // Sensibilidad de rotación ajustable (0.005).
            targetRotationY += event.movementX * 0.005;
            targetRotationX += event.movementY * 0.005;
        } else {
            // Modo libre (flotante): calculamos respecto al centro del contenedor
            const rect = wrapper.getBoundingClientRect();
            mouse_x = event.clientX - (rect.left + rect.width / 2);
            mouse_y = event.clientY - (rect.top + rect.height / 2);
        }
    }

    // Toggle para atrapar/desatrapar al pinchar con el clic
    function handleClick() {
        if (!document.pointerLockElement) {
            // 1. Si no está atrapado, lo atrapamos en el contenedor/canvas
            wrapper.requestPointerLock();
        } else {
            // 2. Si ya estaba atrapado, al pinchar otra vez se desatrapa (soltamos il ratón)
            document.exitPointerLock();
        }
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

        // 1. Cargador de texturas
        const textureLoader = new THREE.TextureLoader();

        // 2. Creamos materiales iniciales con tus texturas
        const materials = Array.from({ length: 6 }, () => new THREE.MeshPhongMaterial({ color: 0xffffff }));

        // 3. Rutas de tus imágenes originales
        const texturePaths = [
            "cube_texture/00.png",
            "cube_texture/01.png",
            "cube_texture/02.png",
            "cube_texture/03.png",
            "cube_texture/04.png",
            "cube_texture/05.png"
        ];

        // La función async es tipo logica en C
        async function cargarTexturasProfesional() {
            // 1. Disparamos las 6 peticiones de red simultáneamente.
            let procesosDeCarga = texturePaths.map(path => textureLoader.loadAsync(path));

            // 2. La barrera: el código espera aquí hasta que las 6 descargas terminan.
            let texturas = await Promise.all(procesosDeCarga);

            // 3. Aplicamos las 6 imágenes
            for (let i = 0; i < 6; i++) {
                materials[i].map = texturas[i];
                materials[i].needsUpdate = true;
            }
        }
        cargarTexturasProfesional();

        // Pasar el array de materiales al Mesh
        const cube = new THREE.Mesh(geometry, materials);
        scene.add(cube);

        camera.position.z = 60;

        // Listener nativo del navegador para enterarnos cuando se atrapa o desatrapa (por clic o con tecla ESC)
        function handlePointerLockChange() {
            isLocked = document.pointerLockElement === wrapper;
        }

        document.addEventListener('pointerlockchange', handlePointerLockChange);

        function animate() {
            requestAnimationFrame(animate);
            
            if (isLocked) {
                // Modo Atrapado:
                // El cubo se queda quieto en su posición actual (no sigue a mouse_x/y) y solo rueda sobre sus ejes
                cube.rotation.x = targetRotationX;
                cube.rotation.y = targetRotationY;
            } else {
                // Modo Libre (flotante):
                // Retoma desde la cara donde se quedó (targetRotation) y añade inclinación relativa al movimiento
                cube.rotation.x = targetRotationX + 0.005 * mouse_y;
                cube.rotation.y = targetRotationY + 0.005 * mouse_x;

                // La posición vuelve a seguir al cursor suavemente
                cube.position.x = 0.05 * mouse_x;
                cube.position.y = 0.05 * -mouse_y;
            }

            renderer.render(scene, camera);
        }

        animate();

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
            document.removeEventListener('pointerlockchange', handlePointerLockChange);
        };
    });
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div 
    bind:this={wrapper} 
    onmousemove={handleMouseMove}
    onclick={handleClick}
>
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