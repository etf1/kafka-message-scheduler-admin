export function isVisible(elem: any): boolean {
  if (!elem || !(elem.offsetWidth || elem.offsetHeight || elem.getClientRects().length)) {
    return false;
  }
  const st = window.getComputedStyle(elem);
  return st.display !== "none" && st.visibility !== "hidden";
}

export function hideOnEscapeOrClickOutside(element: any, hideFunc: () => void) {
  const hideElement = () => {
    if (hideFunc) {
      hideFunc();
    } else {
      element.style.display = "none";
    }
  };

  const mouseListener = (event: MouseEvent) => {
    if (isVisible(element) && !element.contains(event.target)) {
      hideElement();
    }
  };
  const kbdListener = (event: KeyboardEvent) => {
    if (isVisible(element) && event.key === "Escape") {
      hideElement();
    }
  };

  document.addEventListener("keydown", kbdListener);
  document.addEventListener("click", mouseListener);

  return () => {
    document.removeEventListener("keydown", kbdListener);
    document.removeEventListener("click", mouseListener);
  };;
}
