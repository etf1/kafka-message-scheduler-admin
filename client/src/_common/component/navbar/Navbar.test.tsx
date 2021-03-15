import React from "react";
import { render, screen, fireEvent } from "@testing-library/react";
import Navbar from "./Navbar";
import { MemoryRouter } from "react-router";

test("renders Navbar main menu part", () => {
  const handleClick = jest.fn();

  const items = [
    { label: "Home", href: "/home" },
    { label: "ClickTest", onClick: handleClick },
    {
      label: "With Children",
      children: [{ label: "First", href: "/first" }, "-", { label: "Second", href: "/second" }],
    },
  ];

  const component = render(
    <MemoryRouter>
      <Navbar brand={"Home"} items={items} />
    </MemoryRouter>
  );
  const container = component.container;

  // checking brand part
  const brand = container.querySelector(".navbar-brand");
  expect(brand).toHaveTextContent("Home");

  const foundItems = container.querySelector(".navbar-start");

  // first item should be an anchor
  const anchor = foundItems?.children[0];
  expect(anchor).toHaveTextContent("Home");
  expect(anchor).toHaveAttribute("href");

  // second item should be clickable
  const clickable = foundItems?.children[1];
  expect(clickable).toHaveTextContent("ClickTest");
  clickable && fireEvent.click(clickable);
  expect(handleClick).toHaveBeenCalledTimes(1);

  // last item should be a dropdown menu
  const dropDown = foundItems?.children[foundItems?.children.length - 1];
  expect(dropDown).toHaveTextContent("With Children");
  expect(dropDown).toBeDefined();
  const ddChildren = dropDown?.children[1].children;
  expect(ddChildren).toBeDefined();
  expect(ddChildren?.length).toEqual(3);

  if (ddChildren) {
    // second child item should be an anchor
    const childAnchor1 = ddChildren[0];
    expect(childAnchor1).toHaveTextContent("First");
    expect(childAnchor1).toHaveAttribute("href");

    // third item should be a hr separator
    const hr = ddChildren[1];
    expect(hr?.tagName.toLowerCase()).toEqual("hr");

    // last child item should be an anchor
    const childAnchorLast = ddChildren[ddChildren.length - 1];
    expect(childAnchorLast).toHaveTextContent("Second");
    expect(childAnchorLast).toHaveAttribute("href");
  }
});

test("renders Navbar right menu part", () => {
  const handleClick = jest.fn();

  const items = [
    { label: "Home", href: "/home" },
    { label: "ClickTest", onClick: handleClick },
    {
      label: "With Children",
      children: [{ label: "First", href: "/first" }, "-", { label: "Second", href: "/second" }],
    },
  ];

  const component = render(
    <MemoryRouter>
      <Navbar brand={"Home"} rightItems={items} />
    </MemoryRouter>
  );
  const container = component.container;

  // checking brand part
  const brand = container.querySelector(".navbar-brand");
  expect(brand).toHaveTextContent("Home");

  const foundItems = container.querySelector(".navbar-end");

  // first item should be an anchor
  const anchor = foundItems?.children[0];
  expect(anchor).toHaveTextContent("Home");
  expect(anchor).toHaveAttribute("href");

  // second item should be clickable
  const clickable = foundItems?.children[1];
  expect(clickable).toHaveTextContent("ClickTest");
  clickable && fireEvent.click(clickable);
  expect(handleClick).toHaveBeenCalledTimes(1);

  // last item should be a dropdown menu
  const dropDown = foundItems?.children[foundItems?.children.length - 1];
  expect(dropDown).toHaveTextContent("With Children");
  expect(dropDown).toBeDefined();
  const ddChildren = dropDown?.children[1].children;
  expect(ddChildren).toBeDefined();
  expect(ddChildren?.length).toEqual(3);

  if (ddChildren) {
    // second child item should be an anchor
    const childAnchor1 = ddChildren[0];
    expect(childAnchor1).toHaveTextContent("First");
    expect(childAnchor1).toHaveAttribute("href");

    // third item should be a hr separator
    const hr = ddChildren[1];
    expect(hr?.tagName.toLowerCase()).toEqual("hr");

    // last child item should be an anchor
    const childAnchorLast = ddChildren[ddChildren.length - 1];
    expect(childAnchorLast).toHaveTextContent("Second");
    expect(childAnchorLast).toHaveAttribute("href");
  }
});
