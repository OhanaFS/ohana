import { Image } from '@mantine/core';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';

export type FilePreviewProps = {
  fileId: string;
};

const FilePreview = ({ fileId }: FilePreviewProps) => {
  const qFile = useQueryFileMetadata(fileId);
  const downladUrl = getFileDownloadURL(fileId);

  return qFile.data?.mime_type.startsWith('image/') ? (
    <Image src={downladUrl} />
  ) : qFile.data?.mime_type.startsWith('video/') ? (
    <video controls playsInline src={downladUrl}></video>
  ) : null;
};

export default FilePreview;
