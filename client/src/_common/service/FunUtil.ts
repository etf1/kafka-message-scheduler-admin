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
  
 /**
  * 
  * @param args Styles Css sous forme d'objets de string (key:value) ou même de fonctions
  * @returns un objet de styles Css que l'on peut appliquer directement à un comoposant React
  */ 
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function slsx(...args: any[]): Record<string, unknown> {
  if (args) {
    return args.reduce((accumulator, currentValue) => {
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




/* Décoder un tableau d'octets depuis une chaîne en base64 */

function b64ToUint6 (nChr:number) {

  return nChr > 64 && nChr < 91 ?
      nChr - 65
    : nChr > 96 && nChr < 123 ?
      nChr - 71
    : nChr > 47 && nChr < 58 ?
      nChr + 4
    : nChr === 43 ?
      62
    : nChr === 47 ?
      63
    :
      0;

}
export function base64DecToArr (sBase64:string, nBlocksSize?:number) {

  var
    sB64Enc = sBase64.replace(/[^A-Za-z0-9+/]/g, ""), nInLen = sB64Enc.length,
    nOutLen = nBlocksSize ? Math.ceil((nInLen * 3 + 1 >> 2) / nBlocksSize) * nBlocksSize : nInLen * 3 + 1 >> 2, taBytes = new Uint8Array(nOutLen);

  for (var nMod3, nMod4, nUint24 = 0, nOutIdx = 0, nInIdx = 0; nInIdx < nInLen; nInIdx++) {
    nMod4 = nInIdx & 3;
    nUint24 |= b64ToUint6(sB64Enc.charCodeAt(nInIdx)) << 18 - 6 * nMod4;
    if (nMod4 === 3 || nInLen - nInIdx === 1) {
      for (nMod3 = 0; nMod3 < 3 && nOutIdx < nOutLen; nMod3++, nOutIdx++) {
        // eslint-disable-next-line
        taBytes[nOutIdx] = nUint24 >>> (16 >>> nMod3 & 24) & 255;
      }
      nUint24 = 0;

    }
  }

  return taBytes;
}


/* Tableau UTF-8 en DOMString et vice versa */

export function UTF8ArrToStr (aBytes:Uint8Array) {

  var sView = "";

  for (var nPart, nLen = aBytes.length, nIdx = 0; nIdx < nLen; nIdx++) {
    nPart = aBytes[nIdx];
    sView += String.fromCharCode(
      nPart > 251 && nPart < 254 && nIdx + 5 < nLen ? /* six bytes */
        /* (nPart - 252 << 32) n'est pas possible pour ECMAScript donc, on utilise un contournement... : */
        (nPart - 252) * 1073741824 + (aBytes[++nIdx] - 128 << 24) + (aBytes[++nIdx] - 128 << 18) + (aBytes[++nIdx] - 128 << 12) + (aBytes[++nIdx] - 128 << 6) + aBytes[++nIdx] - 128
      : nPart > 247 && nPart < 252 && nIdx + 4 < nLen ? /* five bytes */
        (nPart - 248 << 24) + (aBytes[++nIdx] - 128 << 18) + (aBytes[++nIdx] - 128 << 12) + (aBytes[++nIdx] - 128 << 6) + aBytes[++nIdx] - 128
      : nPart > 239 && nPart < 248 && nIdx + 3 < nLen ? /* four bytes */
        (nPart - 240 << 18) + (aBytes[++nIdx] - 128 << 12) + (aBytes[++nIdx] - 128 << 6) + aBytes[++nIdx] - 128
      : nPart > 223 && nPart < 240 && nIdx + 2 < nLen ? /* three bytes */
        (nPart - 224 << 12) + (aBytes[++nIdx] - 128 << 6) + aBytes[++nIdx] - 128
      : nPart > 191 && nPart < 224 && nIdx + 1 < nLen ? /* two bytes */
        (nPart - 192 << 6) + aBytes[++nIdx] - 128
      : /* nPart < 127 ? */ /* one byte */
        nPart
    );
  }

  return sView;

}
