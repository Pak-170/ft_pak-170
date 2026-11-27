<script lang="ts">
  import { onMount } from 'svelte';
  import { currentUser, users, logout, fetchUsers } from '../client';

  type User = {
    ID: number;
    username: string;
    alias: string | null;
  };

  let loggingOut = $state(false);
  let errorMessage = $state('');
  let section = $state('profile');
  let friends = $state<User[]>([]);
  let selectedUser = $state<User | null>(null);
  let search = $state('');

  let filteredUsers = $derived(
    $users.filter(
      (user) =>
        user.ID !== $currentUser?.ID &&
        !friends.some((friend) => friend.ID === user.ID) &&
        user.username.toLowerCase().includes(search.toLowerCase()),
    ),
  );
  onMount(() => {
    fetchUsers().catch(() => {});
  });

  async function handleLogout() {
    if (loggingOut) return;

    loggingOut = true;
    errorMessage = '';
    try {
      await logout();
      window.location.hash = '#/login';
    } catch {
      errorMessage = 'COULD NOT LOG OUT';
      loggingOut = false;
    }
  }

  function addFriend(user: User) {
    friends = [...friends, user];
    selectedUser = user;
  }

  function removeFriend(user: User) {
    friends = friends.filter((friend) => friend.ID !== user.ID);
    selectedUser = null;
  }
</script>

<div class="profile-page">
  <div class="scanlines"></div>
  <section class="profile">
    <header>
      <span class="section-label">PAK SYSTEMS / SOCIAL</span>
      <h1>{section === 'profile' ? 'USER PROFILE' : 'FRIENDS'}</h1>
    </header>
    <nav>
      <button
        onclick={() => {
          section = 'profile';
          selectedUser = null;
        }}
      >
        PROFILE
      </button>
      <button
        onclick={() => {
          section = 'friends';
          selectedUser = null;
        }}
      >
        FRIENDS ({friends.length})
      </button>
    </nav>
    {#if section === 'profile'}
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
        <button onclick={handleLogout} disabled={loggingOut}>
          {loggingOut ? 'DISCONNECTING...' : 'LOG OUT'}
        </button>
      </footer>
    {:else if selectedUser}
      <button onclick={() => (selectedUser = null)}>← BACK</button>
      <div class="session-status">
        <span class="status-dot"></span>
        USER PROFILE
      </div>
      <div class="user-details">
        <div class="detail">
          <span class="label">USERNAME</span>
          <span class="value">{selectedUser.username}</span>
        </div>
        <div class="detail">
          <span class="label">ALIAS</span>
          <span class="value">{selectedUser.alias || 'NOT SET'}</span>
        </div>
      </div>
      {#if friends.some((friend) => friend.ID === selectedUser?.ID)}
        <button onclick={() => removeFriend(selectedUser!)}>
          REMOVE FRIEND
        </button>
      {:else}
        <button onclick={() => addFriend(selectedUser!)}> ADD FRIEND </button>
      {/if}
    {:else}
      <div class="social-section">
        <h2>MY FRIENDS</h2>
        {#each friends as friend}
          <button class="user-button" onclick={() => (selectedUser = friend)}>
            {friend.username}
          </button>
        {:else}
          <p>NO FRIENDS YET</p>
        {/each}
      </div>
      <div class="social-section">
        <h2>FIND USERS</h2>
        <input bind:value={search} placeholder="SEARCH USERNAME..." />
        {#each filteredUsers as user}
          <button class="user-button" onclick={() => (selectedUser = user)}>
            {user.username}
          </button>
        {:else}
          <p>NO USERS FOUND</p>
        {/each}
      </div>
    {/if}
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
    height: fit-content;
    max-height: 100%;
    padding: 28px;
    background: var(--surface);
    border: 1px solid var(--terminal-border);
    overflow-y: auto;
    scrollbar-width: none;
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

  nav {
    display: flex;
    gap: 10px;
    margin-top: 15px;
  }

  .social-section {
    margin-top: 22px;
  }

  .social-section h2 {
    color: var(--green-dim);
    font-size: 12px;
  }

  .social-section input {
    width: 100%;
    padding: 8px;
    margin-bottom: 10px;
    color: var(--green-bright);
    background: var(--black);
    border: 1px solid var(--green-dark);
  }

  .user-button {
    display: block;
    width: 100%;
    margin-bottom: 8px;
    text-align: left;
    background: var(--black);
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
