<script lang="ts">
    import { onMount, tick } from 'svelte';
    import type  {Rate}   from '../types'
    let socket
    let btc, eth
    let message
    let btc_rate:Rate
    let eth_rate:Rate
    let data = []
    const column = ["symbol","ask","bid","high","last","low","volume"]
    onMount(() => {
        socket = new WebSocket('ws://localhost:8000/ws');
        eth = new WebSocket('ws://localhost:8000/realtime_eth_rate');
        btc = new WebSocket('ws://localhost:8000/realtime_btc_rate');
        socket.onopen = () => {
        console.log('socket connected');
    };
    btc.onmessage = (ev) => {
        btc_rate = JSON.parse(JSON.parse(ev.data))
        console.log('btc')
    }
    eth.onmessage = (ev) => {
        eth_rate = JSON.parse(JSON.parse(ev.data))
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
    <tr>
        {#if btc_rate && eth_rate}
            {#each column as sym}
            <td>{btc_rate[sym]}</td>
            <td>{eth_rate[sym]}</td>
            {/each}
        {/if}
    </tr>
</table>

<style global  lang="postcss">

</style>