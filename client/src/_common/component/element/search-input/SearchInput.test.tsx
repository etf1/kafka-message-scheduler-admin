import { fireEvent, render, screen } from '@testing-library/react';
import SearchInput from './SearchInput';

test('renders SearchInput', async () => {
  // Use fake timers to control RxJS debounceTime reliably
  // (real timers are flaky: 700+300ms ≈ 1000ms debounce but scheduling overhead causes races).
  vi.useFakeTimers();

  const value = 'an initial value';
  const onChange = vi.fn();
  render(
    <SearchInput
      debounceDelay={1000}
      placeholder='a placeholder'
      value={value}
      onChange={onChange}
    />,
  );

  const txt: HTMLInputElement = screen.getByDisplayValue(/an initial value/i) as HTMLInputElement;
  expect(txt).toBeInTheDocument();
  fireEvent.change(txt, { target: { value: 'new value' } });
  expect(txt.value).toBe('new value');
  fireEvent.change(txt, { target: { value: 'new value 2' } });
  expect(txt.value).toBe('new value 2');

  // on change is called after debounceDelay — not before
  expect(onChange).toBeCalledTimes(0);
  await vi.advanceTimersByTimeAsync(700);
  expect(onChange).toBeCalledTimes(0);
  await vi.advanceTimersByTimeAsync(300); // total 1000ms elapsed → debounce fires
  expect(onChange).toBeCalledTimes(1);

  fireEvent.change(txt, { target: { value: 'new value 3' } });
  expect(txt.value).toBe('new value 3');
  expect(onChange).toBeCalledTimes(1);
  await vi.advanceTimersByTimeAsync(1000);
  expect(onChange).toBeCalledTimes(2);

  // if search value has not changed (distinctUntilChanged), no onChange is triggered
  fireEvent.change(txt, { target: { value: 'new value 4' } });
  expect(txt.value).toBe('new value 4');
  fireEvent.change(txt, { target: { value: 'new value 3' } }); // revert to same as last emitted
  expect(txt.value).toBe('new value 3');
  await vi.advanceTimersByTimeAsync(1000);
  expect(onChange).toBeCalledTimes(2); // no extra call: value didn't change

  vi.useRealTimers();
});
