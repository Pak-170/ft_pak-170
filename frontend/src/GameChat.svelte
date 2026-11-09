<script lang="ts">
  import MonitorFrame from './MonitorFrame.svelte';
  import Chat from './Chat.svelte';
</script>

<div class="game-chat-layout">
  <!-- Area del Juego -->
  <main class="game-area">
    <div class="game-placeholder">
      <h1>GAME AREA</h1>
      <p>Poner El Juego Aqui</p>
    </div>
  </main>
  <!-- Area del Juego -->
  <aside class="crt-panel">
    <div class="crt-container">
      <div class="chat-screen">
        <Chat />
      </div>
      <MonitorFrame
        vertical={true}
        ledState="all_off"
        keyStates={[false, false, false, false]}
      />
    </div>
  </aside>
</div>

<style>
  .game-chat-layout {
    position: fixed;
    display: grid;
    grid-template-columns: 70% 30%;
    width: 100%;
    height: 100dvh;
  }

  .game-area {
    display: grid;
    place-items: center;
  }

  .game-placeholder {
    color: white;
    text-align: center;
  }

  .crt-panel {
    display: grid;
    place-items: center;
  }

  .crt-container {
    position: relative;
    width: min(100%, 56dvh);
    aspect-ratio: 9 / 16;
  }

  .chat-screen {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 1;
  }
</style>
