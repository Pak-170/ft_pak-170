<script lang="ts">
  import Chat from './Chat.svelte';
</script>

<div class="game-chat-layout">
  <!-- Area del Juego-->
  <main class="game-area">
    <div class="game-placeholder">
      <h1>GAME AREA</h1>
      <p>Poner El Juego Aqui</p>
    </div>
  </main>
  <!-- Area del Juego-->
  <aside class="crt-panel">
    <div class="crt-container">
      <div class="chat-screen">
        <Chat />
      </div>
    </div>
  </aside>
</div>

<style>
  .game-chat-layout {
    position: absolute;
    inset: 0;
    overflow: hidden;
  }

  .game-area {
    position: absolute;
    inset: 0 30% 0 0;
    display: grid;
    place-items: center;
  }

  .game-placeholder {
    color: white;
    text-align: center;
  }

  .crt-panel {
    position: absolute;
    top: 5dvh;
    right: 0;
    width: 56dvh;
    height: 82dvh;
    overflow: hidden;
  }

  .crt-container {
    position: relative;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .chat-screen {
    position: absolute;
    inset: 0;
    z-index: 1;
  }
</style>
