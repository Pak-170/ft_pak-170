<script lang="ts">
  import Canvas from './Canvas.svelte';

  let speed = $state(1);
</script>

<div>
  <div class="inner">
    <p>Speeddddd: {Intl.NumberFormat().format(speed)}</p>
    <button onclick={() => (speed *= 1.2)}> Increment </button>
  </div>
  <Canvas {speed} />
</div>

<style>
  div {
    width: auto;
    height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem;

    .inner {
      border-radius: 1rem;
      background-color: darkblue;
      padding-top: 1rem;
      height: auto;
    }
  }
</style>
