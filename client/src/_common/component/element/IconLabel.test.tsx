import React from 'react';
import { render, screen } from '@testing-library/react';
import IconLabel from './IconLabel';

test('renders multiple Icons with Labels', () => {
  const component = render(<IconLabel data={[{icon:"save", label:"Hello"}, {icon:"pencil", label:"World"}]} />);
  const container = component.container;
  

  // checking labels  
  const hello = screen.getByText(/Hello/i);
  expect(hello).toBeInTheDocument();
  const world = screen.getByText(/World/i);
  expect(world).toBeInTheDocument();

  // checking icons  
  const icons = container.querySelectorAll('i');
  expect(icons.length).toEqual(2);
  expect(icons[0]).toHaveClass('fa-save');
  expect(icons[1]).toHaveClass('fa-pencil');

});
