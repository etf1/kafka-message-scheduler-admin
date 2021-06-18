const BASE_KEY = "kafka-msg-scheduler-admin";

export function load<T>(
  key: string,
  defaultValue: T | undefined
): T | undefined {
  const value = sessionStorage.getItem(BASE_KEY + "." + key);
  if (value) {
    try {
      return JSON.parse(window.atob(value)) as T;
    } catch {
      return defaultValue;
    }
  } else {
    return defaultValue;
  }
}
export function save<T>(key: string, value: T) {
  sessionStorage.setItem(
    BASE_KEY + "." + key,
    window.btoa(JSON.stringify(value))
  );
}
