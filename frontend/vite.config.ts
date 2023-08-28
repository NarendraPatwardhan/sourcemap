import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import dotenv from "dotenv";

dotenv.config({ path: "../.env" });
const port = process.env.SOURCEMAP_FRONTEND_PORT;

if (!port) {
  throw new Error("SOURCEMAP_FRONTEND_PORT is not set");
}

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  envDir: "../",
  envPrefix: "SOURCEMAP_",
  server: { port: port },
});
