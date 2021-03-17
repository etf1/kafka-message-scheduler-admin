import i18n, { TFunction } from "i18next";
import { initReactI18next } from "react-i18next";

import LanguageDetector from "i18next-browser-languagedetector";
import Backend from "i18next-http-backend";
import { isNumber } from "_common/type/utils";

// don't want to use this?
// have a look at the Quick start guide
// for passing in lng and translations on init

i18n
  // load translation using http -> see /public/locales (i.e. https://github.com/i18next/react-i18next/tree/master/example/react/public/locales)
  // learn more: https://github.com/i18next/i18next-http-backend
  .use(Backend)
  // detect user language
  // learn more: https://github.com/i18next/i18next-browser-languageDetector
  .use(LanguageDetector)
  // pass the i18n instance to react-i18next.
  .use(initReactI18next)
  // init i18next
  // for all options read: https://www.i18next.com/overview/configuration-options
  .init({
    fallbackLng: "en",
    debug: false,

    interpolation: {
      escapeValue: false, // not needed for react as it escapes by default
    },
  });

function hasLength(value: any): value is { length: number } {
  return isNumber(value?.length);
}

export function pluralizeIf(
  count: number | undefined | { length: number },
  label: string,
  labels: string | undefined,
  t?: TFunction
) {
  const nb = hasLength(count) ? count.length : count;
  if (nb === undefined || nb === 0 || nb === 1 || labels === undefined) {
    return t ? t(label) : label;
  }
  return t ? t(labels) : labels;
}

export function getShortLanguageFromLS(): string | null {
  const locale = localStorage.getItem("i18nextLng");
  if (locale) {
    return locale;
  }
  return null;
}
export function changeLanguage(
  lng: "en-US" | "fr-FR" | string,
  callback?: ((error: any, t: TFunction) => void) | undefined
): Promise<TFunction> {
  return i18n.changeLanguage(lng, callback);
}
export default i18n;
