export const apiBaseUrl = import.meta.env.VITE_API_URL as string;

if (!apiBaseUrl) {
  console.warn("apiBaseUrl environment variable was not set");
}
