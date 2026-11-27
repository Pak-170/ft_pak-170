<script lang="ts">
  import Chat from './Chat.svelte';
  import Game from './Game.svelte';
</script>

<div class="game-chat-layout">
  <!-- Area del Juego-->
  <main class="game-area">
    <Game />
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
    display: grid;
    grid-template-columns: minmax(0, 1fr) 30.24%;
    gap: 0;
    width: 100%;
    height: 100%;
    overflow: hidden;
  }

  .game-area {
    width: 100%;
    height: 100%;
    min-width: 0;
    overflow: hidden;
  }

  .crt-panel {
    position: relative;
    top: 5%;
    width: 100%;
    height: 80%;
    min-width: 0;
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
  }

  :global(:root:not(:fullscreen)) .game-chat-layout {
    grid-template-columns: minmax(0, 1fr) 36%;
  }
</style>
