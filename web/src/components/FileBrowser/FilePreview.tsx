import { Image, Text } from '@mantine/core';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';

export type FilePreviewProps = {
  fileId: string;
};

const FilePreview = ({ fileId }: FilePreviewProps) => {
  const qFile = useQueryFileMetadata(fileId);
  const downloadUrl = getFileDownloadURL(fileId, { inline: true });

  return qFile.data?.mime_type.startsWith('image/') ? (
    <Image src={downloadUrl} />
  ) : qFile.data?.mime_type.startsWith('video/') ? (
    <video className="w-full" controls playsInline src={downloadUrl}></video>
  ) : qFile.data?.mime_type.startsWith('audio/') ? (
    <audio className="w-full" controls playsInline src={downloadUrl}></audio>
  ) : (
    <Text>{qFile.data?.file_name || ''}</Text>
  );
};

export default FilePreview;
