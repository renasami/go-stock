<script lang="ts">
  import { onMount, tick } from 'svelte';
  let socket
  let message
  onMount(() => {
    socket = new WebSocket('ws://localhost:8000/ws');
    socket.onopen = () => {
      console.log('socket connected');
  };
  socket.onopen = () => {
      console.log('socket connected');
  };

  socket.onmessage = (event) => {
      if(!event.data) { return }
      message = JSON.parse(event.data);
      console.log(message)
    };
  });


</script>

<main>
	<button class="px-4 py-2 text-white bg-red-600 rounded shadow">ボタンテスト</button>
</main>

<style global  lang="postcss" >
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>