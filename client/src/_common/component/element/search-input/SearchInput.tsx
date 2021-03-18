import { ChangeEventHandler, useEffect, useState } from "react";
import useSeachText from "_common/hook/useSearchText";

export type SearchInputProps = Omit<
  React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  >,
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
  const handleSearchChange = useSeachText(
    onChange,
    debounceDelay,
    discardDuplicates
  );

  useEffect(() => {
    setSearchString(value);
  }, [value]);

  const handleChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    const value = e.target.value;
    setSearchString(value);
    handleSearchChange(value);
  };

  return (
    <input
      className="input"
      onChange={handleChange}
      value={searchString}
      {...others}
    />
  );
};

export default SearchInput;
