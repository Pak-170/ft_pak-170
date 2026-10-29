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
      (a, b) =>
        new Date(a.SendTime).getTime() -
        new Date(b.SendTime).getTime(),
    );

  $: filteredUsers = $users.filter(
    (user) =>
      user.ID !== $currentUser?.ID &&
      `${user.username} ${user.alias ?? ''}`
        .toLowerCase()
        .includes(search.toLowerCase()),
  );

  $: selectedChat = $chats.find(
    (chat) => chat.ID === selectedChatID,
  );

  $: if (
    selectedChatID === null &&
    chatIDs.length > 0
  ) {
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
        type: chatType === 'private' ? 'individual' : 'group',
        userIDs: selectedUsers.map((user) => user.ID),
        ...(chatType === 'group'
          ? { name: groupName.trim() }
          : {}),
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

<div class="scanlines"></div>

<div class="monitor-content">
  <div class="chat-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>Chats</h2>
        <button
          class="new-chat"
          onclick={openNewChat}
          aria-label="New Chat"
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
            <div
              class="message"
              class:own-message={msg.Username === username}
            >
              <strong>{msg.Username}</strong>
              <div>{msg.Msg}</div>
              <small>{new Date(msg.SendTime).toLocaleTimeString()}</small>
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
            placeholder="Write a message..."
            bind:value={message}
          />
        </form>
      {:else}
        <div class="header">
          <h1>Chat</h1>
        </div>
        <div class="messages">
          <div class="message">
            Select a chat or create a new one.
          </div>
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
          onclick={() => (showNewChat = false)}
          aria-label="Close"
        >
          x
        </button>
      </div>

      <div class="chat-type">
        <button
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
          placeholder="Group Name"
          bind:value={groupName}
        />
      {/if}

      <input
        class="user-search"
        placeholder="Search users..."
        bind:value={search}
      />

      {#if errorMessage}
        <p role="alert">{errorMessage}</p>
      {/if}

      <div class="user-list">
        {#each filteredUsers as user (user.ID)}
          <button
            class:selected={selectedUsers.some((u) => u.ID === user.ID)}
            onclick={() => toggleUser(user)}
          >
            <span class="avatar">
              {user.username.slice(0, 1).toUpperCase()}
            </span>
            <span>{user.alias || user.username}</span>
            <span class="check">
              {selectedUsers.some((u) => u.ID === user.ID) ? '✓' : ''}
            </span>
          </button>
        {:else}
          <p class="empty-users">
            No users available.
          </p>
        {/each}
      </div>

      <div class="modal-actions">
        <button onclick={() => (showNewChat = false)}>
          Cancel
        </button>
        <button
          class="create-button"
          disabled={
            creatingChat ||
            selectedUsers.length === 0 ||
            (chatType === 'group' && !groupName.trim())
          }
          onclick={createNewChat}
        >
          {creatingChat ? 'Creating...' : 'Create Chat'}
        </button>
      </div>
    </div>
  </div>
{/if}

<style lang="scss">

.new-chat-button {
  border-color: var(--green-light);
  color: var(--green-bright);
}

.new-chat-button:hover {
  box-shadow: 0 0 8px var(--green-dark);
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

.modal h2 {
  margin: 0;
  color: var(--green-bright);
}

.modal label {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.chat-type {
  display: flex;
  gap: 10px;
}

.chat-type button {
  flex: 1;
  padding: 10px;
  background: var(--black);
  border: 1px solid var(--green-dark);
  color: var(--green);
  cursor: pointer;
}

.chat-type button.active,
.chat-type button:hover {
  border-color: var(--green-light);
  color: var(--green-bright);
  box-shadow: 0 0 8px var(--green-dark);
}

.modal input {
  width: 100%;
  padding: 12px;
  background: var(--black);
  border: 1px solid var(--green-dark);
  color: var(--green);
  outline: none;
}

input {
    outline: none;
    box-shadow: inset 0 0 12px rgba(51, 255, 51, 0.025);
  }

.modal input:focus {
  outline: none;
  border-color: var(--green-light);
  box-shadow: 0 0 8px var(--green-dark);
}

.modal input::placeholder {
  color: var(--green-dark);
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
  background: var(--black);
  border: none;
  border-bottom: 1px solid var(--green-dark);
  color: var(--green);
  cursor: pointer;
}

.user-list button:last-child {
  border-bottom: none;
}

.user-list button.selected,
.user-list button:hover {
  background: var(--surface-light);
  color: var(--green-bright);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.modal-actions button {
  padding: 10px 15px;
  background: var(--black);
  border: 1px solid var(--green-dark);
  color: var(--green);
  cursor: pointer;
}

.modal-actions button:last-child {
  border-color: var(--green-light);
  color: var(--green-bright);
}

.modal-actions button:hover:not(:disabled) {
  box-shadow: 0 0 8px var(--green-dark);
}

.modal-actions button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
  .message small {
    display: block;
    margin-top: 5px;
    color: var(--green-dark);
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
    min-width: 150px;
    width: 15%;
    flex-shrink: 0;
    padding: 20px;
    border-right: 1px solid var(--green-dark);
    background: var(--surface);
  }

  .sidebar h2 {
    margin-top: 0;
    color: var(--green-bright);
  }

  .sidebar button {
    width: 100%;
    padding: 12px;
    margin-bottom: 8px;
    text-align: left;
    background: var(--black);
    border: 1px solid var(--green-dark);
    color: var(--green);
    cursor: pointer;
  }

  .sidebar button.active,
  .sidebar button:hover {
    border-color: var(--green-light);
    color: var(--green-bright);
    box-shadow: 0 0 8px var(--green-dark);
  }

  .chat {
    flex: 1;
    min-width: 0;
    height: 100%;
    padding: 30px;
    display: flex;
    flex-direction: column;
    background: var(--black);
    color: var(--green);
  }

  .header {
    flex-shrink: 0;
    padding: 15px;
    background: var(--surface);
    border: 1px solid var(--terminal-border);
    box-shadow:
      inset 0 0 15px --scanline,
      0 0 8px --scanline-glow;
  }

  .messages {
    display: flex;
    flex-direction: column;
    flex: 1;
    min-height: 0;
    padding: 12px 16px;

    border: 1px solid var(--green-dark);
    overflow-y: auto;
    box-shadow: inset 0 0 30px rgba(51, 255, 51, 0.025);
    scrollbar-width: none;
  }

  .message.own-message {
    align-self: flex-end;
  }

  .message:not(.own-message) {
    align-self: flex-start;
  }

  .message {
    max-width: 70%;
    padding: 10px 15px;
    border-bottom: 1px dashed var(--green-dark);
    overflow-wrap: anywhere;
  }

  .message strong {
    display: block;
    margin-bottom: 4px;
    color: var(--green-bright);
    text-shadow: 0 0 5px rgba(182, 255, 182, 0.4);
  }

  form {
    flex-shrink: 0;
    margin-top: 15px;
  }

  form input {
    width: 100%;
    padding: 15px;
    margin-bottom: 5dvh;
    background: var(--black);
    border: 1px solid var(--green-dark);
    color: var(--green);
    outline: none;
    box-shadow: inset 0 0 10px rgba(51, 255, 51, 0.025);
  }

  form input:focus {
    border-color: var(--green-light);
    box-shadow:
      0 0 8px --scanline-glow,
      inset 0 0 10px --text-glow;
  }

  form input::placeholder {
    color: var(--green-dark);
  }
</style>
