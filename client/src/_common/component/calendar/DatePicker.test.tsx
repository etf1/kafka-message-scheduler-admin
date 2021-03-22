import { fireEvent, render, screen } from "@testing-library/react";
import { later } from "_common/service/FunUtil";
import DatePicker from "./DatePicker";
import { fr, enGB, enUS } from "date-fns/locale";

test("renders DatePicker", async () => {
  const onChange = jest.fn();
  const value = new Date(2020, 5, 21);
  const component = render(
    <DatePicker
      placeholder="A placeholder"
      value={value}
      onChange={onChange}
      locale={fr}
      dateFormat={"dd/MM/yyyy"}
      todayLabel={"aujourdhui"}
    />
  );
  const container = component.container;

  const div = container.getElementsByClassName("dropdown-trigger")[0];
  expect(div).toBeInTheDocument();
  const input = screen.getByPlaceholderText(/A placeholder/i);
  expect(input).toBeInTheDocument();

  // open calendar popup
  fireEvent.click(div);
  const txt: HTMLInputElement = screen.getByText(
    /aujourdhui/i
  ) as HTMLInputElement;
  expect(txt).toBeInTheDocument();

  // click day button
  const btn29 = screen.getByText("29");
  expect(btn29).toBeInTheDocument();
  fireEvent.click(btn29);
  expect(onChange).toBeCalledTimes(1);
  const btn19 = screen.getByText("19");
  expect(btn19).toBeInTheDocument();
  fireEvent.click(btn19);
  expect(onChange).toBeCalledTimes(2);
});
