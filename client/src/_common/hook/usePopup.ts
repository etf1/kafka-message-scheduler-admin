import { useState, useRef, useEffect } from "react";
import { hideOnEscapeOrClickOutside } from "../service/DomUtil";
/**
 *
 * @param isPopupInitiallyVisisible if popup should be visible initially
 */
export default function usePopup<T>(
  isPopupInitiallyVisisible: boolean = false
) {
  const [popupVisible, setPopupVisible] = useState(isPopupInitiallyVisisible);
  const popupRef = useRef<T>(null);

  useEffect(() => {
    if (popupVisible) {
      return hideOnEscapeOrClickOutside(popupRef.current, () => {
        setPopupVisible(false);
      });
    }
  }, [popupVisible]);

  return { popupVisible, setPopupVisible, popupRef };
}
