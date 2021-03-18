import { useRef, useEffect } from "react";
import { Subject, identity } from "rxjs";
import { debounceTime, distinctUntilChanged } from "rxjs/operators";

/**
 *
 * Small textInput hook helper to avoid too much onChange event triggers
 *
 * @param onChange handler that should be called on effective changes
 * @param debounceDelay call onChange only after a certain debounce delay. (default = 650)
 * @param discardDuplicates should we discard duplicate changes ? (default = true)
 * @returns
 */
function useSeachText(
  onChange: (value: string | undefined) => void,
  debounceDelay: number = 650,
  discardDuplicates: boolean = true
) {
  const searchSubject = useRef<Subject<string>>(new Subject());

  useEffect(() => {
    const searchResultObservable = searchSubject.current.pipe(
      debounceTime(debounceDelay),
      discardDuplicates ? distinctUntilChanged() : identity
    );
    const subscription = searchResultObservable.subscribe(onChange);
    return () => subscription.unsubscribe();
  }, [debounceDelay, discardDuplicates, searchSubject, onChange]);

  const handleSearchChange = (value: string) => {
    searchSubject.current.next(value);
  };

  return handleSearchChange;
}

export default useSeachText;
