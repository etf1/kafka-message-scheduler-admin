import React from "react";
import {
  isString,
  isNumber,
  isPrimitive,
  isArray,
  isFunction,
  isDictionary,
  sameKey
} from "_common/type/utils";

function isSelectedValueType(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any
): value is string | ReadonlyArray<string> | number | undefined {
  return isString(value) || isNumber(value) || isArray<string>(value);
}

function getSelectedValue<T>(
  value: T,
  keyFieldName: string | undefined
): string | ReadonlyArray<string> | number | undefined {
  if (isSelectedValueType(value)) {
    return value;
  } else if (
    keyFieldName &&
    isDictionary<string | ReadonlyArray<string> | number | undefined>(value)
  ) {
    return value[keyFieldName];
  }
}

export type SelectProps<T> = Omit<
  React.SelectHTMLAttributes<HTMLSelectElement>,
  "defaultValue" | "value" | "onBlur" | "onChange"
> & {
  value?: T | undefined;
  defaultValue?: T | undefined;
  options: T[];
  onChange?: (value: T | undefined) => void;
  onBlur?: (value: T | undefined) => void;
  keyField?: string;
  labelField?:
    | string
    | ((value: T | undefined, asString?: boolean) => string | undefined);
};

function Select<T>({
  options,
  onChange,
  onBlur,
  value,
  defaultValue,
  keyField = "key",
  labelField = "label",
  ...restProps
}: SelectProps<T>): JSX.Element {
  const handleChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    const key = event.target.value;
    const value = options.find(sameKey(keyField, key));
    onChange && onChange(value);
  };
  const handleBlur = (event: React.FocusEvent<HTMLSelectElement>) => {
    const key = event.target.value;
    const value = options.find(sameKey(keyField, key));
    onBlur && onBlur(value);
  };
  return (
    <div className="field">
      <div className="control">
        <div className="select is-fullwidth">
          <select
            defaultValue={getSelectedValue(defaultValue, keyField)}
            value={getSelectedValue(value, keyField)}
            onChange={handleChange}
            onBlur={handleBlur}
            {...restProps}
          >
            {options.map((option: T) => {
              if (isPrimitive(option)) {
                return (
                  <option
                    key={option + ""}
                    value={getSelectedValue(option, keyField)}
                  >
                    {option}
                  </option>
                );
              } else if (
                isDictionary<
                  string | ReadonlyArray<string> | number | undefined
                >(option)
              ) {
                return (
                  <option key={"" + option[keyField]} value={option[keyField]}>
                    {isFunction(labelField)
                      ? labelField(option)
                      : option[labelField]}
                  </option>
                );
              } else {
                return null;
              }
            })}
          </select>
        </div>
      </div>
    </div>
  );
}

export default Select;
