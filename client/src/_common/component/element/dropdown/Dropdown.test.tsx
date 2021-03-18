import { render, screen } from "@testing-library/react";
import Dropdown from "./Dropdown";

type TestOption = {
  label: string;
};

test("renders DropDown with selected Option and 2 options", () => {
  const value: TestOption = { label: "An option 2" };
  const options: TestOption[] = [{ label: "An option 1" }, value];
  const onChange = jest.fn();
  const component = render(
    <Dropdown
      placeholder="a placeholder"
      value={value}
      options={options}
      getKey={(option) => option.label}
      renderOption={(option) => <span>{option.label}</span>}
      onChange={onChange}
    />
  );
  const container = component.container;

  // checking selected option
  const selected = container.querySelector(".dropdown-trigger");
  expect(selected).toHaveTextContent("An option 2");
  // checking options
  const foundItems = container.querySelectorAll(".dropdown-item");
  expect(foundItems.length).toEqual(options.length);
  expect(foundItems[0]).toHaveTextContent("An option 1");
  expect(foundItems[1]).toHaveTextContent("An option 2");
});

test("renders DropDown without pre-selection and 3 options", () => {
  const value = undefined;
  const options: TestOption[] = [{ label: "An option 1" }, { label: "An option 2" }, { label: "An option 3" }];
  const onChange = jest.fn();
  const component = render(
    <Dropdown
      placeholder="a placeholder"
      value={value}
      options={options}
      getKey={(option) => option.label}
      renderOption={(option) => <span>{option.label}</span>}
      onChange={onChange}
    />
  );
  const container = component.container;

  // checking placholder
  const selected = container.querySelector(".dropdown-trigger");
  expect(selected).toHaveTextContent("a placeholder");
  // checking options
  const foundItems = container.querySelectorAll(".dropdown-item");
  expect(foundItems.length).toEqual(options.length);
  expect(foundItems[0]).toHaveTextContent("An option 1");
  expect(foundItems[1]).toHaveTextContent("An option 2");
  expect(foundItems[2]).toHaveTextContent("An option 3");
});
