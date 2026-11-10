<script lang="ts">
  import { currentUser, disconnectWebSocket } from '../client';

  let loggingOut = $state(false);
  let errorMessage = $state('');

  async function logout() {
    if (loggingOut) return;

    loggingOut = true;
    errorMessage = '';
    try {
      const response = await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
      if (!response.ok) throw new Error();

      currentUser.set(null);
      disconnectWebSocket();
      window.location.hash = '#/login';
    } catch {
      errorMessage = 'COULD NOT LOG OUT';
      loggingOut = false;
    }
  }
</script>

<div class="profile-page">
  <div class="scanlines"></div>
  <section class="profile">
    <header>
      <span class="section-label">PAK SYSTEMS / ACCOUNT</span>
      <h1>USER PROFILE</h1>
    </header>
    <div class="session-status">
      <span class="status-dot"></span>
      SESSION ACTIVE
    </div>
    <div class="user-details">
      <div class="detail">
        <span class="label">USERNAME</span>
        <span class="value">{$currentUser?.username ?? '---'}</span>
      </div>
      <div class="detail">
        <span class="label">ALIAS</span>
        <span class="value">{$currentUser?.alias || 'NOT SET'}</span>
      </div>
    </div>
    {#if errorMessage}
      <p class="error">{errorMessage}</p>
    {/if}
    <footer>
      <span>END SESSION</span>
      <button onclick={logout} disabled={loggingOut}>
        {loggingOut ? 'DISCONNECTING...' : 'LOG OUT'}
      </button>
    </footer>
  </section>
</div>

<style>
  .profile-page {
    position: relative;
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    overflow: hidden;
    background: var(--black);
  }

  .profile-page :global(.scanlines) {
    position: absolute;
    width: 100%;
    height: 100%;
    z-index: 3;
  }

  .profile {
    position: relative;
    z-index: 2;
    width: min(100%, 520px);
    padding: 28px;
    background: var(--surface);
    border: 1px solid var(--terminal-border);
  }

  header {
    padding-bottom: 18px;
    border-bottom: 1px solid var(--green-dark);
  }

  .section-label,
  .label {
    color: var(--green-dim);
    letter-spacing: 2px;
  }

  h1 {
    margin: 10px 0 0;
    color: var(--green-bright);
    font-family: 'Press Start 2P', monospace;
    font-size: 14px;
    line-height: 1.8;
  }

  .session-status {
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 22px 0;
  }

  .status-dot {
    width: 8px;
    height: 8px;
    background: var(--green);
    box-shadow: 0 0 8px var(--green);
  }

  .user-details {
    display: flex;
    flex-direction: column;
    gap: 18px;
  }

  .detail {
    display: flex;
    flex-direction: column;
    gap: 5px;
    padding-bottom: 12px;
    border-bottom: 1px dashed var(--green-dark);
  }

  .value {
    color: var(--green-bright);
    overflow-wrap: anywhere;
  }

  footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 15px;
    margin-top: 25px;
    color: var(--green-dim);
  }

  button {
    padding: 8px 14px;
    color: var(--green);
    background: transparent;
    border: 1px solid var(--green-dark);
    cursor: pointer;
  }

  button:hover:not(:disabled) {
    color: var(--green-bright);
    border-color: var(--green-light);
  }

  button:disabled {
    opacity: 0.5;
    cursor: wait;
  }

  .error {
    color: var(--green-bright);
  }

  @media (max-width: 480px) {
    .profile {
      padding: 18px;
    }

    h1 {
      font-size: 11px;
    }
  }
</style>
