import { Modal, FileButton, Button } from '@mantine/core';
import { showNotification } from '@mantine/notifications';
import { Loader } from 'tabler-icons-react';
import {
  useMutateUpdateFile,
  useMutateUpdateFileMetadata,
  useMutateUploadFile,
  useQueryFileMetadata,
} from '../../api/file';

export type UploadFileModalProps = {
  onClose: () => void;
  opened: boolean;
} & (
  | {
      update: true;
      updateFileId: string;
    }
  | {
      update: false;
      parentFolderId: string;
    }
);

const UploadFileModal = (props: UploadFileModalProps) => {
  const mUploadFile = useMutateUploadFile();
  const mUpdateFile = useMutateUpdateFile();
  const mUpdateFileMeta = useMutateUpdateFileMetadata();
  //@ts-ignore
  const qFileMeta = useQueryFileMetadata(props.updateFileId || '');

  return (
    <Modal
      centered
      opened={props.opened}
      onClose={props.onClose}
      title="Upload a File"
    >
      <div className="flex">
        {mUploadFile.isLoading || mUpdateFile.isLoading ? (
          <Loader className="mr-5" />
        ) : null}
        <FileButton
          onChange={async (item) => {
            if (!item) {
              return;
            }

            if (props.update) {
              if (qFileMeta.data?.versioning_mode !== 2) {
                await mUpdateFileMeta.mutateAsync({
                  file_id: props.updateFileId,
                  versioning_mode: 2,
                });
              }
            }

            (props.update
              ? mUpdateFile.mutateAsync({
                  file_id: props.updateFileId,
                  file: item,
                  frag_count: 1,
                  parity_count: 1,
                })
              : mUploadFile.mutateAsync({
                  file: item,
                  folder_id: props.parentFolderId,
                  frag_count: 1,
                  parity_count: 1,
                })
            )
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
              disabled={mUploadFile.isLoading || mUpdateFile.isLoading}
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

export default UploadFileModal;
