import { defineConfig } from "vite";
import react from "@vitejs/plugin-react-swc";
import dotenv from "dotenv";

dotenv.config({ path: "../.env" });

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  envDir: "../",
  server: { port: process.env.VITE_FRONTEND_PORT || 3000 },
});
