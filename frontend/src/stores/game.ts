import * as THREE from 'three';

export class Game {
  // gente, estos son los limites fisicos de la mesa para el Pad 1 (mi pad!)
  PAD1_MIN_X = -4.1;
  PAD1_MAX_X = 4.1;
  PAD1_MIN_Z = 5.0; // limite avance maximo
  PAD1_MAX_Z = 9.5; // La porteria propia

  // limites para el rival (Pad 2)
  PAD2_MIN_X = -4.1;
  PAD2_MAX_X = 4.1;
  PAD2_MIN_Z = -9.5; // Su fondo, porteria
  PAD2_MAX_Z = -5.0; // Su limite de avance

  // Sensibilidad del mouse:
  MOUSE_SENSIBILITY = 0.01;

  PadP1: THREE.Mesh;
  PadP2: THREE.Mesh; // es lo mismo que abajo!
  Disco: THREE.Mesh; // esto es para intaciar una variable de tipo (tenemos que darle la posibilidad de que sea undefined para instanciarla al principio)
  renderer: THREE.WebGLRenderer;
  // Fisicas del disco, mas rapido que el metro de Madrid a las 8 de la manana
  DISCO_RADIUS = 0.5;

  constructor(canvas: HTMLCanvasElement, wrapper: HTMLDivElement) {
    // adaptar juego a pagina abierta
    let width = wrapper.getBoundingClientRect().width;
    let height = wrapper.getBoundingClientRect().height;

    // 1. ESCENA Y RENDERIZADOR (el plato de juego)
    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(60, width / height, 0.1, 1000);

    this.renderer = new THREE.WebGLRenderer({
      canvas: canvas,
      antialias: true,
    });
    this.renderer.setSize(width, height);

    // 2. ILUMINACION
    const ambientLight = new THREE.AmbientLight(0xffffff, 2);
    scene.add(ambientLight);

    const dirLight = new THREE.DirectionalLight(0xffffff, 0.9);
    dirLight.position.set(5, 12, 8);
    scene.add(dirLight);

    // 3. ESTRUCTURA DE LA MESA (El grupo gigante con todo montado)
    const mesaGroup = new THREE.Group();

    // 3.1 La pista roja
    const geometryMesa = new THREE.BoxGeometry(10.5, 0.5, 20);
    const materialMesa = new THREE.MeshStandardMaterial({
      color: 0xff0000,
      roughness: 0.3,
    });
    const mesa = new THREE.Mesh(geometryMesa, materialMesa);
    mesaGroup.add(mesa);

    // Material para las vigas metalicas de los lados
    const materialBorde = new THREE.MeshStandardMaterial({
      color: 0x8a9597,
      metalness: 0.4,
      roughness: 0.5,
    });

    // 3.2 Las paredes de acero donde rebota el disco
    const geometryParedLarga = new THREE.BoxGeometry(1, 1.5, 20);

    const paredIzq = new THREE.Mesh(geometryParedLarga, materialBorde);
    paredIzq.position.set(-5.7, 0.5, 0);

    const paredDer = new THREE.Mesh(geometryParedLarga, materialBorde);
    paredDer.position.set(5.7, 0.5, 0);

    // 3.3 El Pad del Jugador 1 (Azul)
    const geometryPadP1 = new THREE.BoxGeometry(2, 0.5, 0.7);
    const materialPadP1 = new THREE.MeshStandardMaterial({
      color: 0x0000ff,
      roughness: 0.3,
    });
    this.PadP1 = new THREE.Mesh(geometryPadP1, materialPadP1);
    this.PadP1.position.set(0, 0.5, 8);

    // 3.4 El Pad del Jugador 2 (azul clarito)
    const geometryPadP2 = new THREE.BoxGeometry(2, 0.5, 0.7);
    const materialPadP2 = new THREE.MeshStandardMaterial({
      color: 0x00ffff,
      roughness: 0.3,
    });
    this.PadP2 = new THREE.Mesh(geometryPadP2, materialPadP2);
    this.PadP2.position.set(0, 0.5, -8);

    // 3.5 El Disco (amarillo brillante)
    const geometryDisco = new THREE.CylinderGeometry(
      this.DISCO_RADIUS,
      this.DISCO_RADIUS,
      0.2,
      32,
    );
    const materialDisco = new THREE.MeshStandardMaterial({
      color: 0xffff00,
      roughness: 0.2,
    });
    this.Disco = new THREE.Mesh(geometryDisco, materialDisco);
    this.Disco.position.set(0, 0.35, 0);

    // Añadir todos lo chachrros dentro del grupo qeu define la mesa
    mesaGroup.add(paredIzq);
    mesaGroup.add(paredDer);
    mesaGroup.add(this.PadP1);
    mesaGroup.add(this.PadP2);
    mesaGroup.add(this.Disco);
    // anadir todos los cacharrso agrupados en el MesaGroup a la escena para uqe sea todo visible
    scene.add(mesaGroup);

    camera.position.set(0, 10, 22); // Movil: alejamos un poco per favore

    camera.lookAt(0, 0, 0);

    render_game(this, scene, camera);
  }
}

export function handleMouseMove(
  game: Game,
  event: MouseEvent,
  canvas: HTMLCanvasElement,
) {
  // Mover el pad solo si tenemos el raton atrapado en el juego (Pointer Lock)
  if (document.pointerLockElement === canvas && !event.ctrlKey) {
    let nextX1 =
      game.PadP1.position.x + event.movementX * game.MOUSE_SENSIBILITY;
    let nextZ1 =
      game.PadP1.position.z + event.movementY * game.MOUSE_SENSIBILITY;
    // Clamp potente con Math.max/min para no salirnos de la mesa
    game.PadP1.position.x = Math.max(
      game.PAD1_MIN_X,
      Math.min(game.PAD1_MAX_X, nextX1),
    );
    game.PadP1.position.z = Math.max(
      game.PAD1_MIN_Z,
      Math.min(game.PAD1_MAX_Z, nextZ1),
    );
  }

  // Si aprietas CTRL, mueves el rival...
  // atencion tesdeadores ... esto requiere mucha coordinacion

  if (document.pointerLockElement === canvas && event.ctrlKey) {
    let nextX2 =
      game.PadP2.position.x - event.movementX * game.MOUSE_SENSIBILITY;
    let nextZ2 =
      game.PadP2.position.z - event.movementY * game.MOUSE_SENSIBILITY;
    game.PadP2.position.x = Math.max(
      game.PAD2_MIN_X,
      Math.min(game.PAD2_MAX_X, nextX2),
    );
    game.PadP2.position.z = Math.max(
      game.PAD2_MIN_Z,
      Math.min(game.PAD2_MAX_Z, nextZ2),
    );
  }
}

// 5. BUCLE DE ANIMACION Y FISICAS (El cuore del programa)
function render_game(game: Game, scene: THREE.Scene, camera: THREE.Camera) {
  requestAnimationFrame(() => {
    render_game(game, scene, camera);
  }); // Esa linea es el motor del bucle infinito de animacion (el render loop)

  if (game.renderer) game.renderer.render(scene, camera);
}
