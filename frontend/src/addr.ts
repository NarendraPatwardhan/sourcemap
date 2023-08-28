// If in development mode, use the backend server at localhost:{API_PORT}
let API_URL: string;
if (import.meta.env.DEV) {
  const API_PORT = import.meta.env.SOURCEMAP_API_PORT;
  if (!API_PORT) {
    throw new Error("SOURCEMAP_API_PORT is not set");
  }
  API_URL = `http://localhost:${API_PORT}/api`;
} else {
  // Otherwise, use the backend server at /api
  API_URL = "/api";
}
// Export the API_URL
export { API_URL };
