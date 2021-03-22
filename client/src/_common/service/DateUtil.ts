import {
  addDays,
  differenceInDays,
  endOfMonth,
  endOfWeek,
  format,
  isSameMonth,
  isToday,
  startOfMonth,
  startOfWeek,
  Locale,
} from "date-fns";

// source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
export function getDayLabelsOfWeek(locale: Locale) {
  // source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
  return new Array(7)
    .fill(startOfWeek(new Date(), { locale }))
    .map((d, i) => format(addDays(d, i), "EEE", { locale }));
}

export type DayOfMonth = {
  date: Date;
  isToday: boolean;
  isThisMonth: boolean;
};

// source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
export function getDaysOfMonth(
  visibleDate: Date,
  locale: Locale
): DayOfMonth[] {
  try {
    // first day of current month view
    // source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
    const start = startOfWeek(startOfMonth(visibleDate), { locale });

    // last day of current month view
    // source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
    const end = endOfWeek(endOfMonth(visibleDate), { locale });

    // source https://gist.github.com/stevensacks/79c60d0f8b1f8bc06b475438f59d687e
    const days = new Array(differenceInDays(end, start) + 1)
      .fill(start)
      .map((s, i) => {
        const date = addDays(s, i);
        return {
          date,
          isToday: isToday(date),
          isThisMonth: isSameMonth(visibleDate, date),
        };
      });

    return days;
  } catch (err) {
    return getDaysOfMonth(new Date(), locale);
  }
}
