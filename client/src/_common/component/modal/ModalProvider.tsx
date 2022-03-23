import React, { useEffect } from "react";
import useRefState from "_common/hook/useRefState";
import ModalService, { ModalConfig } from "./ModalService";
import Modal from "./Modal";
import { useTranslation } from "react-i18next";

/**
 *
 * Gère la création/affichage/destruction de modales,
 * sur injonction du ModalService.
 *
 */
const ModalProvider = (): JSX.Element => {
  const { t } = useTranslation();
  const [getActiveModals, setActiveModals] = useRefState<ModalConfig[]>([]);
  const [getButtonsDisabled, setButtonsDisabled] = useRefState(false);
  const closeTopModal = () => {
    const activeModals = getActiveModals();
    if (activeModals.length > 0) {
      setActiveModals((prevActiveModals) => {
        prevActiveModals.pop();
        return [...prevActiveModals];
      });
    }
  };
  const closeModal = (conf: ModalConfig) => {
    setActiveModals((prevActiveModals) => {
      const index = prevActiveModals.indexOf(conf);
      prevActiveModals.splice(index, 1);
      return [...prevActiveModals];
    });
  };

  useEffect(() => {
    return ModalService.setModalProvider(
      (conf: ModalConfig) => {
        setActiveModals((activeModals) => [...activeModals, conf]);
      },
      closeTopModal,
      t
    );
    // eslint-disable-next-line
  }, []);

  const activeModals = getActiveModals();

  const modals = activeModals.map((conf, i) => {
    const handleSave = async (
      e:
        | undefined
        | React.MouseEvent<HTMLButtonElement, MouseEvent>
        | KeyboardEvent
    ) => {
      setButtonsDisabled(true);
      const result = await conf.onSave(e);
      setButtonsDisabled(false);
      if (result) {
        closeModal(conf);
      }
    };

    const handleCancel = async (
      e:
        | undefined
        | React.MouseEvent<HTMLButtonElement, MouseEvent>
        | KeyboardEvent
    ) => {
      setButtonsDisabled(true);
      const result = await conf.onCancel(e);
      setButtonsDisabled(false);
      if (result) {
        closeModal(conf);
      }
    };

    return (
      <Modal
        /* eslint-disable-next-line react/no-array-index-key */
        key={`${i}${conf.title}.${conf.width}`}
        onSave={handleSave}
        onCancel={handleCancel}
        isActive
        zIndex={i}
        width={conf.width}
        title={conf.title}
        saveLabel={conf.saveLabel}
        cancelLabel={conf.cancelLabel}
        isOnTop={i === activeModals.length - 1}
        focused={conf.focused}
        disabled={getButtonsDisabled()}
        showSaveButton={conf.showSaveButton}
        showCancelButton={conf.showCancelButton}
      >
        {conf.content}
      </Modal>
    );
  });

  return <>{modals}</>;
};

export default ModalProvider;
