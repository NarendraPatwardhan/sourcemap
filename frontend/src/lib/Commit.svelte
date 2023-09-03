<script lang="ts">
	import type { Data, Changes } from '$lib';
	export let name: string;
	export let path: string;
	export let size: string;
	export let changes: Changes;
  export let repr: string = "";
	export let children: Data[] = [];
</script>

{#if path == ''}
	<p>{name}: {size}</p>
{:else}
	<p>{path}: {size}</p>
{/if}
{#if children && children.length > 0}
	{#each children as child}
		<svelte:self {...child} />
	{/each}
{/if}
