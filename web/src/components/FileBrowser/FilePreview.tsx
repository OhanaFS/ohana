import { Code, Image, Loader, Text } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { APIClient } from '../../api/api';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';

type TextFilePreviewProps = {
  url: string;
};

const TextFilePreview = ({ url }: TextFilePreviewProps) => {
  const qFile = useQuery(['file', url], () =>
    APIClient.get<string>(url).then((res) => res.data)
  );

  if (qFile.isLoading) return <Loader />;

  return (
    <Code className="w-full overflow-x-scroll" px="md">
      <pre>{qFile.data || ''}</pre>
    </Code>
  );
};

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
  ) : qFile.data?.mime_type.startsWith('text/') ? (
    <TextFilePreview url={downloadUrl} />
  ) : (
    <Text>{qFile.data?.file_name || ''}</Text>
  );
};

export default FilePreview;
