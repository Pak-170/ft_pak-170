<script lang="ts">
  let email = $state('');
  let password = $state('');
  let username = $state('');
  let register = $state(false);
  let confPassword = $state('');

  function submit(e: Event) {
    e.preventDefault();


    if (register) {
	  if (password !== confPassword)
	  	return;
      console.log('Register:', username, email, password);
    } else {
      console.log('Login:', email, password);
    }
  }
</script>

<div class="login-page">
  <div class="scanlines"></div>
  <div class="login-box">
    <div class="header">
      <div class="logo">◆</div>
      <div class="title">
        <p class="system">PAK SYSTEMS</p>
        <h1>{register ? 'NEW USER' : 'ACCESS TERMINAL'}</h1>
      </div>
    </div>
    <div class="status">
      <span class="dot"></span>
      <span>SYSTEM ONLINE</span>
    </div>
    <form onsubmit={submit}>
      <label>
	  	{#if register}
        <span>OPERATOR ID</span>
        <input
          type="text"
          placeholder="USERNAME"
          bind:value={username}
          autocomplete="nickname"
          required
        />
		{/if}
		<span>{register ? 'NETWORK ID' : 'IDENTIFICATION'}</span>
		<input
          type="email"
          placeholder="USER@NETWORK"
          bind:value={email}
          autocomplete="email"
          required
        />
      </label>
      <label>
        <span>ACCESS CODE</span>
        <input
          type="password"
          placeholder="••••••••"
          bind:value={password}
          autocomplete={register ? 'new-password' : 'current-password'}
          required
        />
		{#if register}
		<span>CONFIRM ACCESS CODE</span>
        <input
          type="password"
          placeholder="••••••••"
          bind:value={confPassword}
          autocomplete={register ? 'new-password' : 'current-password'}
          required
        />
		{#if confPassword.length > 0}
		  {#if password === confPassword}
		    <span class="passconfirm">ACCESS CODES MATCH</span>
		  {:else}
		    <span class="passconfirm">ACCESS CODES DO NOT MATCH</span>
		  {/if}
		{/if}
		{/if}
      </label>
      <button class="submit" type="submit">
        {register ? 'CREATE ACCOUNT' : 'LOG IN'}
      </button>
    </form>
    <button class="switch" type="button" onclick={() => (register = !register)}>
      {register ? '>> EXISTING USER // LOGIN' : '>> NEW USER // REGISTER'}
    </button>
    <div class="footer">
      <span>SECURE CONNECTION</span>
      <span>SYS. 1982</span>
    </div>
  </div>
</div>

<style lang="scss">
  .login-page {
    position: relative;
    width: 100%;
    height: 100dvh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    background: radial-gradient(
      circle at center,
      var(--surface) 0%,
      var(--black) 70%
    );
    color: var(--green);
    text-shadow: 0 0 4px var(--text-glow);
  }

  .login-box {
    position: relative;
    width: 430px;
    max-width: 100%;
    padding: 30px;
    background: linear-gradient(180deg, var(--surface-light), var(--surface));
    border: var(--terminal-border-width) solid var(--terminal-border);
    box-shadow:
      0 0 0 1px var(--black),
      8px 8px 0 var(--black),
      12px 12px 0 var(--green-dark),
      0 0 30px var(--terminal-glow);
    z-index: 2;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 18px;
    padding-bottom: 20px;
    border-bottom: 1px dashed var(--green-dark);
  }

  .logo {
    width: 52px;
    height: 52px;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--green);
    color: var(--black);
    font-family: 'Press Start 2P', monospace;
    box-shadow:
      4px 4px 0 var(--green-dark),
      0 0 15px var(--terminal-glow);
  }

  .system {
    margin: 0 0 6px;
    color: var(--green);
    letter-spacing: 3px;
  }

  h1 {
    margin: 0;
    color: var(--green-bright);
    font-family: 'Press Start 2P', monospace;
    font-size: 11px;
    line-height: 1.6;
    text-shadow: 0 0 6px var(--terminal-glow);
  }

  .status {
    display: flex;
    align-items: center;
    margin: 24px 0 28px;
    color: var(--green);
    letter-spacing: 2px;
  }

  .dot {
    width: 9px;
    height: 9px;
    margin-right: 9px;
    background: var(--green);
    box-shadow:
      0 0 5px var(--green),
      0 0 12px var(--terminal-glow);
    animation: status-blink 2s steps(1) infinite;
  }

  @keyframes status-blink {
    0%,
    85% {
      opacity: 1;
    }
    86%,
    100% {
      opacity: 0.35;
    }
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  label {
    display: flex;
    flex-direction: column;
    gap: 7px;
    color: var(--green);
    letter-spacing: 2px;
  }

  input {
    width: 100%;
    padding: 12px;
    background: var(--black);
    border: 1px solid var(--green-dark);
    outline: none;
    color: var(--green);
    box-shadow: inset 0 0 12px rgba(51, 255, 51, 0.025);
  }

  input:focus {
    border-color: var(--green);
    box-shadow:
      0 0 8px var(--terminal-glow),
      inset 0 0 12px rgba(51, 255, 51, 0.035);
  }

  input::placeholder {
    color: var(--green-dark);
  }

  .submit {
    width: 100%;
    margin-top: 5px;
    padding: 15px;
    background: var(--green);
    border: none;
    color: var(--black);
    font-family: 'Press Start 2P', monospace;
    font-size: 11px;
    cursor: pointer;
    box-shadow: 0 6px 0 var(--green-dark);
    transition:
      transform 80ms linear,
      box-shadow 80ms linear,
      background 80ms linear;
  }

  .submit:hover {
    background: var(--green-light);
    transform: translateY(3px);
    box-shadow: 0 3px 0 var(--green-dark);
  }

  .submit:active {
    transform: translateY(5px);
    box-shadow: none;
  }

  .switch {
    display: block;
    margin: 25px auto 0;
    padding: 0;
    background: none;
    border: none;
    color: var(--green-dark);
    cursor: pointer;
    transition: color 100ms linear;
  }

  .switch:hover {
    color: var(--green);
  }

  .footer {
    display: flex;
    justify-content: space-between;
    margin-top: 28px;
    padding-top: 12px;
    border-top: 1px dashed var(--green-dark);
    color: var(--green-dark);
    letter-spacing: 1px;
  }

  .passconfirm {
    display: flex;
	font-size: 17px;
    color: var(--green);
    letter-spacing: 1px;
  }

  @media (max-width: 550px) {
    .login-page {
      padding: 15px;
    }
    .login-box {
      width: 100%;
      padding: 24px;
      box-shadow:
        6px 6px 0 var(--black),
        9px 9px 0 var(--green-dark),
        0 0 25px var(--terminal-glow);
    }
  }
</style>
