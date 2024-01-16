<script lang="ts">
	// import { page } from '$app/stores';
	export let data;
</script>

<!-- <pre>{JSON.stringify($page, null, 2)}</pre> -->

<ul>
	{#each data.services as service}
		<li>{service.name} - {service.status}</li>
	{/each}
</ul>
