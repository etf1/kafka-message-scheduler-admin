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

export const later = async (duration: number) =>
  new Promise((resolve) => setTimeout(resolve, duration));
