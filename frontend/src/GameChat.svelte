<script lang="ts">
  import Chat from './Chat.svelte';
</script>

<div class="game-chat-layout">
  <!-- Area del Juego-->
  <main class="game-area">
    <img src="/GamePreview.gif" alt="Game background" />
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
    right: 0;
    inset: 0 30% 0 0;
    display: flex;
    align-items: center;
    justify-content: flex-start;
    overflow: hidden;
  }

  .game-area img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    object-position: left center;
    display: block;
  }

  .crt-panel {
    position: absolute;
    top: 5dvh;
    right: 0;
    width: 56dvh;
    height: 82dvh;
    overflow: hidden;
    z-index: 2;
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
