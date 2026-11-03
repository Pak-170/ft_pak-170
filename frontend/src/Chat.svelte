<script lang="ts">
  import { onMount, tick } from 'svelte';
  import {
    messages,
    currentUser,
    users,
    chats,
    connectWebSocket,
    sendMessage,
    fetchUsers,
    fetchChats,
    createChat as createChatRequest,
  } from '../client';

  type ChatUser = {
    ID: number;
    username: string;
    alias: string | null;
  };

  let message = '';
  let selectedChatID: number | null = null;
  let messagesContainer: HTMLDivElement;

  let showNewChat = false;
  let chatType: 'private' | 'group' = 'private';
  let search = '';
  let selectedUsers: ChatUser[] = [];
  let groupName = '';

  let loadingUsers = false;
  let creatingChat = false;
  let errorMessage = '';

  onMount(async () => {
    await connectWebSocket();

    try {
      await Promise.all([fetchUsers(), fetchChats()]);
    } catch (error) {
      console.error('[Chat] Could not load initial data:', error);
      errorMessage = 'Could not load users or chats.';
    }
  });

  $: username = $currentUser?.username ?? '';

  $: chatIDs = $chats.map((chat) => chat.ID);

  $: selectedMessages = $messages
    .filter((msg) => msg.ChatID === selectedChatID)
    .sort(
      (a, b) => new Date(a.SendTime).getTime() - new Date(b.SendTime).getTime(),
    );

  $: filteredUsers = $users.filter(
    (user) =>
      user.ID !== $currentUser?.ID &&
      `${user.username} ${user.alias ?? ''}`
        .toLowerCase()
        .includes(search.toLowerCase()),
  );

  $: selectedChat = $chats.find((chat) => chat.ID === selectedChatID);

  $: if (selectedChatID === null && chatIDs.length > 0) {
    selectedChatID = chatIDs[0];
  }

  async function submit() {
    if (!message.trim() || selectedChatID === null || !username) {
      return;
    }

    sendMessage(selectedChatID, username, message.trim());
    message = '';
  }

  async function selectChat(chatID: number) {
    selectedChatID = chatID;
    await scrollDown();
  }

  function toggleUser(user: ChatUser) {
    const exists = selectedUsers.some((u) => u.ID === user.ID);

    if (exists) {
      selectedUsers = selectedUsers.filter((u) => u.ID !== user.ID);
      return;
    }

    if (chatType === 'private') {
      selectedUsers = [user];
    } else {
      selectedUsers = [...selectedUsers, user];
    }
  }

  function openNewChat() {
    chatType = 'private';
    selectedUsers = [];
    search = '';
    groupName = '';
    errorMessage = '';
    showNewChat = true;
  }

  async function createNewChat() {
    if (creatingChat || selectedUsers.length === 0) {
      return;
    }

    if (chatType === 'group' && !groupName.trim()) {
      errorMessage = 'Enter a group name.';
      return;
    }

    creatingChat = true;
    errorMessage = '';

    try {
      const chat = await createChatRequest({
        infoType: 'chat',
        type: chatType === 'private' ? 'individual' : 'group',
        userIDs: selectedUsers.map((user) => user.ID),
        ...(chatType === 'group' ? { name: groupName.trim() } : {}),
      });

      selectedChatID = chat.ID;
      showNewChat = false;
      selectedUsers = [];
      groupName = '';
      search = '';

      await fetchChats();
      selectedChatID = chat.ID;

      await scrollDown();
    } catch (error) {
      console.error('[Chat] Could not create chat:', error);
      errorMessage = 'Could not create chat. Please try again.';
    } finally {
      creatingChat = false;
    }
  }

  async function scrollDown() {
    await tick();

    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  $: if (selectedMessages.length > 0) {
    scrollDown();
  }
</script>

<div class="monitor-content">
<div class="scanlines"></div>
  <div class="chat-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>Chats</h2>

        <button
          class="new-chat"
          onclick={openNewChat}
          aria-label="New Chat"
          title="New Chat"
        >
          +
        </button>
      </div>

      {#each $chats as chat (chat.ID)}
        <button
          class:active={selectedChatID === chat.ID}
          onclick={() => selectChat(chat.ID)}
        >
          {chat.Name || `Chat ${chat.ID}`}
        </button>
      {/each}
    </aside>

    <div class="chat">
      {#if selectedChatID !== null}
        <div class="header">
          <h1>{selectedChat?.Name || `Chat ${selectedChatID}`}</h1>
        </div>

        <div class="messages" bind:this={messagesContainer}>
          {#each selectedMessages as msg, i (`${msg.ChatID}-${msg.SendTime}-${i}`)}
            <div class="message" class:own-message={msg.Username === username}>
              <strong>{msg.Username}</strong>
              <div>{msg.Msg}</div>
              <small>
                {new Date(msg.SendTime).toLocaleTimeString()}
              </small>
            </div>
          {/each}
        </div>

        <form
          onsubmit={(e) => {
            e.preventDefault();
            submit();
          }}
        >
          <input
            aria-label="Write a message"
            placeholder="Write a message..."
            bind:value={message}
          />
        </form>
      {:else}
        <div class="header">
          <h1>Chat</h1>
        </div>

        <div class="messages">
          <div class="message">Select a chat or create a new one.</div>
        </div>
      {/if}
    </div>
  </div>
</div>

{#if showNewChat}
  <div
    class="modal-overlay"
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) showNewChat = false;
    }}
  >
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
    >
      <div class="modal-header">
        <h2 id="modal-title">New Chat</h2>

        <button
          type="button"
          onclick={() => (showNewChat = false)}
          aria-label="Close"
          title="Close"
        >
          x
        </button>
      </div>

      <div class="chat-type">
        <button
          type="button"
          class:active={chatType === 'private'}
          onclick={() => {
            chatType = 'private';
            selectedUsers = [];
            errorMessage = '';
          }}
        >
          Private
        </button>

        <button
          type="button"
          class:active={chatType === 'group'}
          onclick={() => {
            chatType = 'group';
            selectedUsers = [];
            errorMessage = '';
          }}
        >
          Group
        </button>
      </div>

      {#if chatType === 'group'}
        <input
          class="group-name"
          aria-label="Group name"
          placeholder="Group Name"
          bind:value={groupName}
        />
      {/if}

      <input
        class="user-search"
        aria-label="Search users"
        placeholder="Search users..."
        bind:value={search}
      />

      {#if errorMessage}
        <p role="alert">{errorMessage}</p>
      {/if}

      <div class="user-list">
        {#each filteredUsers as user (user.ID)}
          <button
            type="button"
            class:selected={selectedUsers.some((u) => u.ID === user.ID)}
            aria-pressed={selectedUsers.some((u) => u.ID === user.ID)}
            onclick={() => toggleUser(user)}
          >
            <span>
              {user.username.toUpperCase()}
            </span>

            <span>{user.alias || user.username}</span>

            <span aria-hidden="true">
              {selectedUsers.some((u) => u.ID === user.ID) ? '✓' : ''}
            </span>
          </button>
        {:else}
          <p class="empty-users">No users available.</p>
        {/each}
      </div>

      <div class="modal-actions">
        <button type="button" onclick={() => (showNewChat = false)}>
          Cancel
        </button>

        <button
          type="button"
          disabled={creatingChat ||
            selectedUsers.length === 0 ||
            (chatType === 'group' && !groupName.trim())}
          onclick={createNewChat}
        >
          {creatingChat ? 'Creating...' : 'Create Chat'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style lang="scss">
.scanlines {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 10;
}

.monitor-content {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  overflow: hidden;
}

  .modal-overlay {
    position: fixed;
    inset: 0;
    z-index: 900;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    background: rgba(0, 0, 0, 0.8);
  }

  .modal {
    display: flex;
    flex-direction: column;
    gap: 15px;
    width: 100%;
    max-width: 450px;
    max-height: 80dvh;
    padding: 25px;
    overflow-y: auto;
    background: var(--surface);
    border: 1px solid var(--green-light);
    box-shadow: 0 0 20px var(--green-dark);
  }

  .modal h2,
  .sidebar h2 {
    margin-top: 0;
    color: var(--green-bright);
  }

  .modal label {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .chat-type,
  .modal-actions {
    display: flex;
    gap: 10px;
  }

  .chat-type button {
    flex: 1;
    padding: 10px;
  }

  .user-list {
    display: flex;
    flex-direction: column;
    gap: 5px;
    max-height: 200px;
    overflow-y: auto;
    border: 1px solid var(--green-dark);
  }

  .user-list button {
    padding: 12px;
    text-align: left;
    border: 0;
    border-bottom: 1px solid var(--green-dark);
  }

  .user-list button:last-child {
    border-bottom: 0;
  }

  .modal-actions {
    justify-content: flex-end;
  }

  .modal-actions button {
    padding: 10px 15px;
  }

  .modal-actions button:last-child {
    border-color: var(--green-light);
    color: var(--green-bright);
  }

  .modal-actions button:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  .chat-layout {
    display: flex;
    width: 100%;
    height: 100%;
    min-width: 0;
    min-height: 0;
    overflow: hidden;
    background: var(--black);
    color: var(--green);
  }

  .sidebar {
    flex: 0 0 15%;
    min-width: 150px;
    padding: 20px;
    overflow-y: auto;
    border-right: 1px solid var(--green-dark);
    background: var(--surface);
  }

  .sidebar button {
    width: 100%;
    padding: 12px;
    margin-bottom: 8px;
    text-align: left;
  }

  .chat {
    display: flex;
    flex: 1;
    flex-direction: column;
    width: 0;
    min-width: 0;
    min-height: 0;
    height: 100%;
    padding: 30px;
    overflow: hidden;
    background: var(--black);
    color: var(--green);
  }

  .header {
    flex-shrink: 0;
    padding: 15px;
    background: var(--surface);
    border: 1px solid var(--terminal-border);
    box-shadow:
      inset 0 0 15px var(--scanline),
      0 0 8px var(--scanline-glow);
  }

  .messages {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;
    padding: 12px 16px;
    overflow-y: auto;
    overflow-x: hidden;
    border: 1px solid var(--green-dark);
    box-shadow: inset 0 0 30px rgba(51, 255, 51, 0.025);
    scrollbar-width: none;
  }

  .message {
    max-width: 70%;
    padding: 10px 15px;
    border-bottom: 1px dashed var(--green-dark);
    overflow-wrap: anywhere;
  }

  .message.own-message {
    align-self: flex-end;
  }

  .message:not(.own-message) {
    align-self: flex-start;
  }

  .message strong {
    display: block;
    margin-bottom: 4px;
    color: var(--green-bright);
    text-shadow: 0 0 5px var(--green-bright);
  }

  .message small {
    display: block;
    margin-top: 5px;
    color: var(--green-dark);
  }

  .chat form {
    flex-shrink: 0;
    margin-top: 15px;
  }

  .chat form input {
    padding: 15px;
    margin-bottom: 5dvh;
    box-shadow: inset 0 0 10px var(--scanline);
  }

  .chat form input:focus {
    box-shadow:
      0 0 8px var(--scanline-glow),
      inset 0 0 10px var(--text-glow);
  }

  .new-chat {
    border-color: var(--green-light);
    color: var(--green-bright);
    background-color: var(--black);
  }
</style>
