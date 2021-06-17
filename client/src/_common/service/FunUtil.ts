import { Dictionary, isDictionary, isFunction, isString } from "_common/type/utils";

export function escapeRegExp(s: string): string {
  return s.replace(/([.*+?^=!:${}()|[\]/\\])/g, "\\$1");
}

export function replaceAll(
  str: string,
  toFind: string,
  toReplace: string
): string {
  return str.replace(new RegExp(escapeRegExp(toFind), "g"), toReplace);
}

export function truncate(str:string, length:number, ending:string = "...") {
  if (str.length > length) {
    return str.substring(0, length - ending.length) + ending;
  } else {
    return str;
  }
}

export const later = async (duration: number = 10) =>
  new Promise((resolve) => setTimeout(resolve, duration));
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function slsx(...args: any[]): Record<string, unknown> {
  if (args) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return args.reduce((accumulator: any, currentValue: any) => {
      if (currentValue === undefined || currentValue === null) {
        return accumulator;
      }
      if (isDictionary(currentValue)) {
        const validObject = Object.keys(currentValue)
          .filter(
            (k) => currentValue[k] !== null && currentValue[k] !== undefined
          )
          .reduce((obj, key) => {
            obj[key] = currentValue[key];
            return obj as Dictionary;
          }, {} as Dictionary);

        return { ...accumulator, ...validObject };
      } else if (isString(currentValue)) {
        const parts = currentValue.split(":", 1);
        accumulator[parts[0]] = parts[1];
        return accumulator;
      } else if (isFunction(currentValue)) {
        return { ...accumulator, ...slsx(currentValue()) };
      } else {
        return accumulator;
      }
    }, {});
  } else {
    return {};
  }
}
