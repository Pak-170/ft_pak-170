<script lang="ts">
  type Props = {
    onClose: () => void;
  };

  let { onClose }: Props = $props();
</script>

<div class="settings-overlay">
  <dialog open class="settings-modal">
    <header>
      <div>
        <span class="section-label">PAK SYSTEMS</span>
        <h2>SETTINGS</h2>
      </div>
      <button onclick={onClose}>X</button>
    </header>
    <main>
      <p>CONFIGURATION OPTIONS</p>
    </main>
    <footer>PAK-170</footer>
  </dialog>
</div>

<style>
  .settings-overlay {
    position: absolute;
    inset: 0;
    z-index: 10;
    display: grid;
    place-items: center;
    padding: 20px;
    background: rgba(0, 0, 0, 0.75);
  }

  .settings-modal {
    width: min(100%, 600px);
    max-height: 100%;
    margin: 0;
    padding: 0;
    color: var(--green);
    background: var(--surface);
    border: 1px solid var(--terminal-border);
    box-shadow: 0 0 25px var(--terminal-glow);
    overflow: auto;
  }

  header,
  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid var(--terminal-border);
  }

  .section-label,
  footer {
    color: var(--green-dim);
  }

  h2 {
    margin: 8px 0 0;
    color: var(--green-bright);
    font-family: 'Press Start 2P', monospace;
    font-size: 13px;
  }

  button {
    padding: 5px 10px;
    color: var(--green);
    background: transparent;
    border: 1px solid var(--green-dark);
    cursor: pointer;
  }

  main {
    padding: 20px;
  }

  main p {
    color: var(--green-dim);
  }

  footer {
    border-top: 1px solid var(--terminal-border);
    border-bottom: none;
  }
</style>
