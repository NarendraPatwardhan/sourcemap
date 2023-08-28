// If in development mode, use the backend server at localhost:{API_PORT}
let API_URL: string;
if (import.meta.env.DEV) {
  const API_PORT = import.meta.env.VITE_API_PORT || 8080;
  API_URL = `http://localhost:${API_PORT}/api`;
} else {
  // Otherwise, use the backend server at /api
  API_URL = "/api";
}
// Export the API_URL
export { API_URL };
