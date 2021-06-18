import { TFunction } from "i18next";
import React from "react";

const DEFAULT_MODAL_WIDTH = "80%";

export interface ModalConfig {
  width: number | string;
  title: React.ReactNode;
  saveLabel?: string;
  cancelLabel?: string;
  content: React.ReactNode;
  onSave: (event: undefined | React.MouseEvent<HTMLButtonElement, MouseEvent> | KeyboardEvent) => Promise<boolean>;
  onCancel: (
    event: "onPrevious" | "onNext" | undefined | React.MouseEvent<HTMLButtonElement, MouseEvent> | KeyboardEvent
  ) => Promise<boolean>;
  focused?: "save" | "cancel";
  showSaveButton?: boolean;
  showCancelButton?: boolean;
}

class ModalServiceImpl {
  private modalProviderHandler: ((config: ModalConfig) => void) | undefined = undefined;

  private closeTopModal: (() => void) | undefined = undefined;

  private t: TFunction | undefined;

  setModalProvider = (modalProviderHandler: (config: ModalConfig) => void, closeTopModal: () => void, t: TFunction) => {
    this.modalProviderHandler = modalProviderHandler;
    this.closeTopModal = closeTopModal;
    this.t = t;
  };

  open = (config: ModalConfig) =>
    this.modalProviderHandler &&
    this.modalProviderHandler({
      width: config.width || DEFAULT_MODAL_WIDTH,
      title: config.title || "",
      saveLabel: config.saveLabel,
      cancelLabel: config.cancelLabel,
      content: config.content,
      onSave: config.onSave || (() => Promise.resolve(true)),
      onCancel: config.onCancel || (() => Promise.resolve(true)),
      focused: config.focused,
      showSaveButton: config.showSaveButton,
      showCancelButton: config.showCancelButton,
    });

  closeTop = () => this.closeTopModal && this.closeTopModal();

  confirm = ({
    title,
    message,
    saveLabel,
    cancelLabel,
    width,
    focused,
  }: {
    title?: React.ReactNode;
    message: React.ReactNode;
    saveLabel?: string;
    cancelLabel?: string;
    width?: number;
    focused?: "save" | "cancel";
  }): Promise<boolean> =>
    new Promise(
      (resolve) =>
        this.modalProviderHandler &&
        this.modalProviderHandler({
          width: width || DEFAULT_MODAL_WIDTH,
          title: title || "",
          saveLabel: saveLabel || (this.t && this.t("Confirm")),
          cancelLabel: cancelLabel,
          content: message,
          onSave: async () => {
            resolve(true);
            return Promise.resolve(true);
          },
          onCancel: async () => {
            resolve(false);
            return Promise.resolve(true);
          },
          focused,
        })
    );

  message = ({
    title,
    message,

    width,
  }: {
    title?: React.ReactNode;
    message: React.ReactNode;

    width?: number;
  }): Promise<boolean> =>
    new Promise(
      (resolve) =>
        this.modalProviderHandler &&
        this.modalProviderHandler({
          width: width || DEFAULT_MODAL_WIDTH,
          title: title || "",
          showSaveButton: false,
          cancelLabel: this.t && this.t("Close-button"),
          content: message,
          onSave: async () => {
            resolve(true);
            return Promise.resolve(true);
          },
          onCancel: async () => {
            resolve(false);
            return Promise.resolve(true);
          },
          focused: "cancel",
        })
    );
}

const ModalService = new ModalServiceImpl();

export default ModalService;
