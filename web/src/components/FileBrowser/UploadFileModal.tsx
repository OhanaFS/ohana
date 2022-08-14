import { Modal, FileButton, Button, Loader } from '@mantine/core';
import {
  useMutateUpdateFile,
  useMutateUpdateFileMetadata,
  useMutateUploadFile,
  useQueryFileMetadata,
} from '../../api/file';
import { handleMultiFileAction } from './multiFileAction';

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

  const handleUpload = async (filesParam: File | File[]) => {
    if (!filesParam) {
      return;
    }

    const files = Array.isArray(filesParam) ? filesParam : [filesParam];

    if (props.update) {
      // If updating, make sure versioning mode is set to 2.
      if (qFileMeta.data?.versioning_mode !== 2) {
        await mUpdateFileMeta.mutateAsync({
          file_id: props.updateFileId,
          versioning_mode: 2,
        });
      }
    }

    // Perform uploads
    await handleMultiFileAction({
      notifications: {
        loadingTitle: (success, _, total) =>
          `Uploading files... ${success + 1} / ${total}`,
        doneTitle: 'Finished uploading files',
        errorTitle: (item, _) => `Error uploading ${item.name}`,
        itemName: (item) => item.name,
      },
      items: files,
      handler: (file) =>
        props.update
          ? mUpdateFile.mutateAsync({
              file_id: props.updateFileId,
              file,
              frag_count: 1,
              parity_count: 1,
            })
          : mUploadFile.mutateAsync({
              file,
              folder_id: props.parentFolderId,
              frag_count: 1,
              parity_count: 1,
            }),
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
