import { ChangeEventHandler, useEffect, useRef, useState } from "react";
import { identity, Subject } from "rxjs";
import { debounceTime, distinctUntilChanged } from "rxjs/operators";

export type SearchInputProps = Omit<
  React.DetailedHTMLProps<React.InputHTMLAttributes<HTMLInputElement>, HTMLInputElement>,
  "value" | "onChange"
> & {
  value: string | undefined;
  onChange: (value: string | undefined) => void;
  debounceDelay?: number;
  discardDuplicates?: boolean;
};
const SearchInput: React.FC<SearchInputProps> = ({
  value,
  onChange,
  debounceDelay,
  discardDuplicates = true,
  ...others
}) => {
  const [searchString, setSearchString] = useState<string | undefined>(value);
  const searchSubject = useRef<Subject<string>>(new Subject());

  useEffect(() => {
    setSearchString(value);
  }, [value]);

  useEffect(() => {
    const searchResultObservable = searchSubject.current.pipe(
      debounceTime(debounceDelay || 650),
      discardDuplicates ? distinctUntilChanged() : identity
    );
    const subscription = searchResultObservable.subscribe(onChange);
    return () => subscription.unsubscribe();
  }, [debounceDelay, discardDuplicates, searchSubject, onChange]);

  const handleChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    const value = e.target.value;
    setSearchString(value);
    searchSubject.current.next(value);
  };

  return <input className="input" onChange={handleChange} value={searchString} {...others} />;
};

export default SearchInput;
