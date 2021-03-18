import { fireEvent, render, screen } from "@testing-library/react";
import { later } from "_common/service/FunUtil";
import SearchInput from "./SearchInput";

test("renders SearchInput", async () => {
  const value = "an initial value";
  const onChange = jest.fn();
  const component = render(
    <SearchInput
      debounceDelay={1000}
      placeholder="a placeholder"
      value={value}
      onChange={onChange}
    />
  );
  const container = component.container;

  const txt: HTMLInputElement = screen.getByDisplayValue(
    /an initial value/i
  ) as HTMLInputElement;
  expect(txt).toBeInTheDocument();
  fireEvent.change(txt, { target: { value: "new value" } });
  expect(txt.value).toBe("new value");
  fireEvent.change(txt, { target: { value: "new value 2" } });
  expect(txt.value).toBe("new value 2");

  // on change is called after a debounce debounceDelay
  expect(onChange).toBeCalledTimes(0);
  await later(700);
  expect(onChange).toBeCalledTimes(0);
  await later(300);
  expect(onChange).toBeCalledTimes(1);

  fireEvent.change(txt, { target: { value: "new value 3" } });
  expect(txt.value).toBe("new value 3");
  expect(onChange).toBeCalledTimes(1);
  await later(1000);
  expect(onChange).toBeCalledTimes(2);

  // if search value has not changed, no OnChange is triggered
  fireEvent.change(txt, { target: { value: "new value 4" } });
  expect(txt.value).toBe("new value 4");
  fireEvent.change(txt, { target: { value: "new value 3" } });
  expect(txt.value).toBe("new value 3");
  await later(1000);
  expect(onChange).toBeCalledTimes(2);
});
