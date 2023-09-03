const apiURL = (relativePath: string): string => {
	let res: string;

	if (import.meta.env.DEV) {
		const apiPort = import.meta.env.SOURCEMAP_API_PORT;
		if (!apiPort) throw new Error('SOURCEMAP_API_PORT is not set');
		res = `http://localhost:${apiPort}/api/${relativePath}`;
	} else {
		res = `/api/${relativePath}`;
	}
	return res;
};

export { apiURL };

interface Commit {
	hash: string;
	author: string;
	message: string;
	timestamp: string;
	data: Data | null;
}

interface Data {
	name: string;
	path: string;
	size: string;
	changes: Changes;
	content?: string;
	repr?: string;
	children?: Data[] | null;
}

interface Changes {
	addition: number;
	deletion: number;
}

type Repository = Commit[];

export type { Changes, Commit, Data, Repository };
