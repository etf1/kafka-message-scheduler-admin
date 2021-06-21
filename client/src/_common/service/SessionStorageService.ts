import { isPrimitive } from "_common/type/utils";

const BASE_KEY = "kafka-msg-scheduler-admin";

export function load<T>(key: string, defaultValue: T | undefined): T | undefined {
  const value = sessionStorage.getItem(BASE_KEY + "." + key);
  if (value) {
    try {
      const result: any = JSON.parse(window.atob(value));
      if (result && result.__primitive__value === true) {
        return result.value as T;
      } else {
        return result as T;
      }
    } catch {
      return defaultValue;
    }
  } else {
    return defaultValue;
  }
}
export function save<T>(key: string, value: T) {
  if (isPrimitive(value) || value === undefined) {
    sessionStorage.setItem(BASE_KEY + "." + key, window.btoa(JSON.stringify({ __primitive__value: true, value })));
  } else {
    sessionStorage.setItem(BASE_KEY + "." + key, window.btoa(JSON.stringify(value)));
  }
}
