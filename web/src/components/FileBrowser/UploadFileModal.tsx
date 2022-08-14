import { Modal, FileButton, Button, Loader } from '@mantine/core';
import { showNotification, updateNotification } from '@mantine/notifications';
import { IconCheck } from '@tabler/icons';
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
  const qFileMeta = useQueryFileMetadata(
    props.update ? props.updateFileId : ''
  );

  const handleUpload = async (files: File | File[]) => {
    if (!files) {
      return;
    }

    const loaderNotificationId = 'upload-loader';
    showNotification({
      id: loaderNotificationId,
      title: 'Uploading files...',
      message: 'Please wait',
      loading: true,
      autoClose: false,
      disallowClose: true,
    });

    let uploadedCount = 0;
    let totalCount = Array.isArray(files) ? files.length : 1;

    if (props.update) {
      // If updating, make sure versioning mode is set to 2.
      if (qFileMeta.data?.versioning_mode !== 2) {
        await mUpdateFileMeta.mutateAsync({
          file_id: props.updateFileId,
          versioning_mode: 2,
        });
      }

      // Assertion that `files` is not an array if updating
      if (Array.isArray(files)) return;

      await mUpdateFile
        .mutateAsync({
          file_id: props.updateFileId,
          file: files,
          frag_count: 1,
          parity_count: 1,
        })
        .then(() => {
          uploadedCount++;
          showNotification({
            title: `${files.name} uploaded`,
            message: 'File Uploaded Successfully',
          });
        })
        .catch((e) =>
          showNotification({
            title: `Error uploading ${files.name}`,
            message: JSON.stringify(e),
          })
        );
    } else {
      // Assertion that `files` is an array if uploading
      if (!Array.isArray(files)) return;

      for (const file of files) {
        updateNotification({
          id: loaderNotificationId,
          title: `Uploading files... ${uploadedCount + 1} / ${totalCount}`,
          message: file.name,
          loading: true,
          autoClose: false,
          disallowClose: true,
        });

        await mUploadFile
          .mutateAsync({
            file,
            folder_id: props.parentFolderId,
            frag_count: 1,
            parity_count: 1,
          })
          .then(() => {
            uploadedCount++;
          })
          .catch((e) =>
            showNotification({
              title: `Error uploading ${file.name}`,
              message: JSON.stringify(e),
            })
          );
      }
    }

    updateNotification({
      id: loaderNotificationId,
      title: 'Finished uploading',
      message: `Uploaded ${uploadedCount} / ${totalCount}`,
      color: 'teal',
      icon: <IconCheck size={16} />,
      autoClose: 5000,
    });

    props.onClose();
  };

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
        <FileButton multiple={!props.update} onChange={handleUpload}>
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
