<script lang="ts">
  import {
    errorNotifications,
    removeErrorNotification,
  } from './stores/notifications';
</script>

<div class="notifications">
  {#each $errorNotifications as notification (notification.id)}
    <div class="notification error" role="alert">
      <span>{notification.message}</span>
      <button
        type="button"
        onclick={() => removeErrorNotification(notification.id)}
        aria-label="Close notification"
      >
        x
      </button>
    </div>
  {/each}
</div>

<style>
  .notifications {
    position: fixed;
    top: 20px;
    right: 20px;
    z-index: 2000;
    display: flex;
    flex-direction: column;
    gap: 10px;
    max-width: 320px;
    pointer-events: none;
  }

  .notification {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 12px;
    pointer-events: auto;
    background: var(--surface);
    border: 1px solid var(--green-dark);
    color: var(--green-bright);
    box-shadow: 0 0 10px var(--scanline-glow);
    font-size: 14px;
  }

  .notification.error {
    border-color: #ff4444;
    color: #ff7777;
  }

  .notification button {
    width: auto;
    padding: 2px 6px;
    color: inherit;
    background: transparent;
    border: none;
  }
</style>
