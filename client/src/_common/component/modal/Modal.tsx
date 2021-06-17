import React, { useCallback, useEffect } from "react";
import clsx from "clsx";
import useFocus from "_common/hook/useFocus";
import { later } from "_common/service/FunUtil";
import { useTranslation } from "react-i18next";
import useKeyDown from "_common/hook/useKeyDown";

const MODAL_SPACING = 40;

type OnClickType = (
  event:
    | undefined
    | React.MouseEvent<HTMLButtonElement, MouseEvent>
    | KeyboardEvent
) => void;

type ModalProps = {
  width: number | string;
  height?: number | string;
  zIndex: number;
  isActive: boolean;
  title: React.ReactNode;
  saveLabel?: string;
  cancelLabel?: string;
  onSave: OnClickType;
  onCancel: OnClickType;
  isOnTop: boolean;
  focused?: "save" | "cancel";
  showSaveButton?: boolean;
  showCancelButton?: boolean;
  disabled?: boolean;
  children: React.ReactNode;
};

const Modal: React.FC<ModalProps> = ({
  width,
  height,
  zIndex,
  isActive,
  title,
  children,
  saveLabel,
  cancelLabel,
  isOnTop,
  onSave,
  onCancel,
  focused = "none",
  showSaveButton = true,
  showCancelButton = true,
  disabled = false
}) => {
  const { t } = useTranslation();
  const padding = zIndex * MODAL_SPACING;
  const [safeRef, setSaveFocus] = useFocus();
  const [cancelRef, setCancelFocus] = useFocus();
  const style = {
    width: typeof width === "string" ? width : (width as number) + padding,
    paddingLeft: padding,
    paddingTop: padding
  };
  useEffect(() => {
    if (focused === "save") {
      later().then(setSaveFocus);
    } else if (focused === "cancel") {
      later().then(setCancelFocus);
    }
    // eslint-disable-next-line
  }, [focused]);

  useKeyDown(
    useCallback((ke) => isOnTop && ke.key === "Escape", [isOnTop]),
    onCancel,
    true
  );

  return (
    <div className={clsx("modal", isActive && "is-active")}>
      <div className="modal-background" />
      <div className="modal-card" style={style}>
        <header className="modal-card-head" style={{ padding: 16 }}>
          <p className="modal-card-title" style={{ fontSize: "1.2rem" }}>
            {title}
          </p>
          <button
            className="delete"
            aria-label="close"
            onClick={onCancel}
            type="button"
          />
        </header>
        <section
          className="modal-card-body"
          style={{
            minHeight: height || "20vh",
            display: "grid",
            alignItems: "center"
          }}
        >
          {children}
        </section>
        <footer
          className="modal-card-foot"
          style={{ padding: 16, justifyContent: "flex-end" }}
        >
          {showSaveButton && (
            <button
              className="button is-success"
              onClick={onSave}
              ref={safeRef}
              disabled={disabled}
              type="button"
            >
              {saveLabel || t("Save")}
            </button>
          )}
          {showCancelButton && (
            <button
              className={clsx("button" /*, !showSaveButton?"is-success":"" */)}
              onClick={onCancel}
              ref={cancelRef}
              disabled={disabled}
              type="button"
            >
              {cancelLabel || t("Cancel")}
            </button>
          )}
        </footer>
      </div>
    </div>
  );
};

export default React.memo(Modal);
