import { Dictionary, isPrimitive } from "_common/type/utils";

const BASE_KEY = "kafka-msg-scheduler-admin-v0";

export function load<T>(
  key: string,
  defaultValue: T | undefined
): T | undefined {
  const store = sessionStorage.getItem(BASE_KEY);
  if (store) {
    try {
      const result: any = JSON.parse(window.atob(store));
      if (result && result[key]) {
        if (result[key] && result[key].__primitive__value === true) {
          return result[key].value as T;
        } else {
          return result[key] as T;
        }
      } else {
        return defaultValue;
      }
    } catch {
      return defaultValue;
    }
  } else {
    return defaultValue;
  }
}
export function save<T>(key: string, value: T) {
  const store = sessionStorage.getItem(BASE_KEY);

  const result: Dictionary = store ? JSON.parse(window.atob(store)) : {};
  if (isPrimitive(value) || value === undefined) {
    const storedValue = {
      ...result,
      [key]: { __primitive__value: true, value },
    };
    sessionStorage.setItem(BASE_KEY, window.btoa(JSON.stringify(storedValue)));
  } else {
    const storedValue = { ...result, [key]: value };
    sessionStorage.setItem(BASE_KEY, window.btoa(JSON.stringify(storedValue)));
  }
}

export function clear(keepKeyPredicat: (key: string) => boolean) {
  const store = sessionStorage.getItem(BASE_KEY);
  if (store) {
    let data: any = JSON.parse(window.atob(store));
    let result: Dictionary = {};
    Object.keys(data).forEach((key) => {
      if (keepKeyPredicat(key)) {
        result[key] = data[key];
      }
    });
    sessionStorage.setItem(BASE_KEY, window.btoa(JSON.stringify(result)));
  }
}
