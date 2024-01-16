import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const results = await fetch('http://localhost:8080/health');
	const data = await results.json();

	return {
		services: data.services
	};
};
