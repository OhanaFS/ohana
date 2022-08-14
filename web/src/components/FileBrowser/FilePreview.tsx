import { Image, Text } from '@mantine/core';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';

export type FilePreviewProps = {
  fileId: string;
};

const FilePreview = ({ fileId }: FilePreviewProps) => {
  const qFile = useQueryFileMetadata(fileId);
  const downladUrl = getFileDownloadURL(fileId, { inline: true });

  return qFile.data?.mime_type.startsWith('image/') ? (
    <Image src={downladUrl} />
  ) : qFile.data?.mime_type.startsWith('video/') ? (
    <video className="w-full" controls playsInline src={downladUrl}></video>
  ) : (
    <Text>{qFile.data?.file_name || ''}</Text>
  );
};

export default FilePreview;
