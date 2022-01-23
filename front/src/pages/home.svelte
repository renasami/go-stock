<script>
    import { onMount, tick } from 'svelte';
    let socket
    let sockets
    let message
    let rateData = []
    let data = []
    const column = ["code","ask","heigh","low"]
    onMount(() => {
        socket = new WebSocket('ws://localhost:8000/ws');
        sockets = new WebSocket('ws://localhost:8000/wss');
        socket.onopen = () => {
        console.log('socket connected');
    };
    socket.onmessage = (event) => {
      if(!event.data) { return }
      message = event.data
    //   console.log(message);
      const d = JSON.parse(message)
      rateData = d.quotes
    };
    sockets.onmessage = (ev) => {
        console.log(ev.data);
    }
  });
</script>

<h1>Home</h1>

<table>
    <tr>

		{#each column as col}
			<th>{col}</th>
		{/each}
    </tr>
        {#each rateData as data}
            <!-- {#each data as d} -->
                <tr>
                    <td><b>{data.currencyPairCode}</b></td>
                    <td>{data.ask}</td>
                    <td>{data.high}</td>
                    <td>{data.low}</td>
                </tr>
            <!-- {/each} -->
        {/each}
</table>

<style global  lang="postcss">

</style>