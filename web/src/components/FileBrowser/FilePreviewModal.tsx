import { Modal, Button } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';
import FilePreview from './FilePreview';

export type FilePreviewModalProps = {
  fileId: string;
  onClose: () => void;
};

const FilePreviewModal = (props: FilePreviewModalProps) => {
  const { fileId, onClose } = props;
  const qFile = useQueryFileMetadata(fileId);
  const downladUrl = getFileDownloadURL(fileId);
  const smallScreen = useMediaQuery('(max-width: 600px)');

  return (
    <Modal
      centered
      opened={!!fileId}
      onClose={onClose}
      title={qFile.data?.file_name}
      size={smallScreen ? '100%' : '70%'}
    >
      <div className="flex">
        <FilePreview fileId={fileId} />
      </div>
      <Button
        component="a"
        href={downladUrl}
        className="bg-blue-600 mt-5"
        color="blue"
      >
        Download
      </Button>
    </Modal>
  );
};

export default FilePreviewModal;
