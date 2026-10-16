<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { messages, sendMessage } from '../client';

  let username = 'yuri';
  let message = '';
  let selectedChatID: number | null = null;
  let messagesContainer: HTMLDivElement;

  $: chatIDs = [...new Set(messages.map((msg) => msg.ChatID))];
  $: selectedMessages = messages.filter((msg) => msg.ChatID === selectedChatID);

  function submit() {
    if (!message.trim() || selectedChatID === null) return;

    sendMessage(selectedChatID, username, message.trim());
    message = '';
  }

  async function selectChat(chatID: number) {
    selectedChatID = chatID;

    await tick();

    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  onMount(async () => {
    await tick();

    if (chatIDs.length > 0) {
      selectedChatID = chatIDs[0];
    }
  });
</script>

<div class="scanlines"></div>

<div class="chat-layout">
  <aside class="sidebar">
    <h2>Chats</h2>

    {#each chatIDs as chatID}
      <button
        class:active={selectedChatID === chatID}
        onclick={() => selectChat(chatID)}
      >
        Chat {chatID}
      </button>
    {/each}
  </aside>

  <div class="chat">
    {#if selectedChatID !== null}
      <div class="header">
        <h1>Chat {selectedChatID}</h1>
      </div>

      <div class="messages" bind:this={messagesContainer}>
        {#each selectedMessages as msg, i (`${msg.ChatID}-${msg.SendTime}-${i}`)}
          <div class="message">
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
        <input placeholder="Write a message..." bind:value={message} />
      </form>
    {:else}
      <div class="header">
        <h1>Global Chat</h1>
      </div>

      <div class="messages">
        <div class="message">No chats available.</div>
      </div>
    {/if}
  </div>
</div>

<style lang="scss">
  .message small {
    display: block;
    margin-top: 5px;
    color: var(--green-dark);
  }

  .chat-layout {
    display: flex;
    width: 100%;
    height: 100vh;
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
    flex: 1;
    min-height: 0;
    border: 1px solid var(--green-dark);
    overflow-y: auto;
    box-shadow: inset 0 0 30px rgba(51, 255, 51, 0.025);
  }

  .messages {
    scrollbar-width: none;
  }

  .message {
    padding: 15px;
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
