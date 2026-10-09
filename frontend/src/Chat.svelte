<script lang="ts">
  import { onMount, tick } from 'svelte';

  let username = 'username';
  let message = '';
  let messages: { username: string; message: string }[] = [];
  let interval: ReturnType<typeof setInterval>;

  let messagesContainer: HTMLDivElement;
//Para Testes {
  function loadMessages() {
    const saved = localStorage.getItem('chat-messages');
    if (saved) {
      try {
        messages = JSON.parse(saved);
      } catch {
        messages = [];
      }
    }
  }

  onMount(() => {
    loadMessages();
    interval = setInterval(() => {
      loadMessages();
    }, 1000);
  });
// }
  const submit = async () => {
    if (!message.trim()) return;

    messages = [
      ...messages,
      {
        username,
        message,
      },
    ];
	//Para Testes {
    localStorage.setItem('chat-messages', JSON.stringify(messages, null, 2));
	// }

    message = '';
    await tick();

    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  };
</script>

<div class="chat">
  <div class="scanlines"></div>

  <div class="header">
    <input bind:value={username} />
  </div>

  <div class="messages" bind:this={messagesContainer}>
    {#each messages as msg}
      <div class="message">
        <strong>{msg.username}</strong>
        <div>{msg.message}</div>
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
</div>

<style lang="scss">
  .chat {
    width: 100%;
    height: 100vh;
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

  .header input {
    width: 100%;
    background: none;
    border: none;
    outline: none;
    color: var(--green);
  }

  .messages {
    flex: 1;
    min-height: 0;
    border: 1px solid var(--green-dark);
    overflow-y: auto;
    box-shadow: inset 0 0 30px rgba(51, 255, 51, 0.025);
  }

  .messages::-webkit-scrollbar {
    display: none;
  }

  .messages {
    scrollbar-width: none;
  }

  .message {
    padding: 15px;
    border-bottom: 1px dashed var(--green-dark);
    font-size: 20px;
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
