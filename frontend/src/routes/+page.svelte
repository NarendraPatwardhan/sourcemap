<script lang="ts">
	import { apiURL } from '$lib';
	import type { Repository, Commit } from '$lib';
	import Commit from '$lib/Commit.svelte';

	let repoAddr = 'https://github.com/NarendraPatwardhan/sourcemap';
	let limit = 10;
	let excludeGlobs = '';
	let excludePaths = '';

	let repo: Repository;

	const fetchRepo = async () => {
		const body = JSON.stringify({
			address: repoAddr + '.git',
			limit,
			excludeGlobs,
			excludePaths
		});
		const res = await fetch(`${apiURL('sourcemap')}`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body
		});
		const parsed: Repository = await res.json();
		return parsed;
	};
</script>

<h1>Sourcemap</h1>

<form
	on:submit|preventDefault={() => {
		repo = fetchRepo();
	}}
>
	<label
		>Repository Address
		<input bind:value={repoAddr} placeholder="enter repository address" />
	</label>

	<label
		>Limit
		<input bind:value={limit} type="number" min="0" />
	</label>

	<label
		>Exclude Globs
		<input bind:value={excludeGlobs} placeholder="enter globs to exclude" />
	</label>

	<label
		>Exclude Paths
		<input bind:value={excludePaths} placeholder="enter paths to exclude" />
	</label>

	<button type="submit">Submit</button>
</form>

{#await repo}
	Loading Repository...
{:then fetchedRepo}
	{#if fetchedRepo}
		{#each fetchedRepo as commit}
			<p>{commit.hash.slice(0, 7)}: {commit.message}</p>
		{/each}
		<hr />
		<Commit {...fetchedRepo[0].data} />
	{/if}
{:catch someError}
	System error: {someError.message}.
{/await}
