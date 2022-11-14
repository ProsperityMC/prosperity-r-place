export function getEnv(key: string) {
  key = key.toUpperCase();
  return window.CONFIG[key] ?? import.meta.env["VITE_" + key];
}
