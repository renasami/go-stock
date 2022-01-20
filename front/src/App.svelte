<script lang="ts">
  import {Router,Route} from  "svelte-routing"
  import Home from "./pages/Home.svelte"  
</script>

<main>
  <Router>
    <Route path="" component={Home}/>
  </Router>
</main>

<style global  lang="postcss" >
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>