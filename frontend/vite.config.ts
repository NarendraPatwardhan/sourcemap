import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import dotenv from 'dotenv';

dotenv.config({ path: '../.env' });

const port = Number(process.env.SOURCEMAP_FRONTEND_PORT);
if (!port) throw new Error('SOURCEMAP_FRONTEND_PORT is not set');

export default defineConfig({
	plugins: [sveltekit()],
	envDir: '../',
	envPrefix: 'SOURCEMAP_',
	server: { port }
});
