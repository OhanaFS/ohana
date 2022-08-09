import { Modal, FileButton, Button } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { Loader } from 'tabler-icons-react';
import { useMutateUploadFile } from '../../api/file';

export type UploadFileModalProps = {
  onClose: () => void;
  opened: boolean;
  parentFolderId: string;
};

export const UploadFileModal = (props: UploadFileModalProps) => {
  const mUploadFile = useMutateUploadFile();

  return (
    <Modal
      centered
      opened={props.opened}
      onClose={props.onClose}
      title="Upload a File"
    >
      <div className="flex">
        {mUploadFile.isLoading ? <Loader className="mr-5" /> : null}
        <FileButton
          onChange={(item) => {
            console.log('we going in');
            if (!item) {
              return;
            }
            mUploadFile
              .mutateAsync({
                file: item,
                folder_id: props.parentFolderId,
                frag_count: 1,
                parity_count: 1,
              })
              .then(() => props.onClose())
              .then(() =>
                showNotification({
                  title: `${item.name} uploaded`,
                  message: 'File Uploaded Successfully',
                })
              );
          }}
        >
          {(props) => (
            <Button
              disabled={mUploadFile.isLoading}
              className="bg-cyan-500"
              color="cyan"
              {...props}
            >
              Upload a File
            </Button>
          )}
        </FileButton>
      </div>
    </Modal>
  );
};
