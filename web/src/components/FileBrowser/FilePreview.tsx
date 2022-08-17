import { Code, Group, Image, Loader, Text } from '@mantine/core';
import { IconLock } from '@tabler/icons';
import { useQuery } from '@tanstack/react-query';
import { APIClient } from '../../api/api';
import { getFileDownloadURL, useQueryFileMetadata } from '../../api/file';
import {
  getSharingLinkURL,
  useQuerySharingLinkMetadata,
} from '../../api/sharing';

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

export type FilePreviewProps = { fileId: string } | { shareId: string };

const FilePreview = (props: FilePreviewProps) => {
  const isShare = 'shareId' in props;
  const qFile = useQueryFileMetadata(!isShare ? props.fileId : '');
  const qShare = useQuerySharingLinkMetadata(isShare ? props.shareId : '');

  const downloadUrl = !isShare
    ? getFileDownloadURL(props.fileId, { inline: true })
    : getSharingLinkURL(props.shareId, 'inline');
  const metadata = !isShare ? qFile.data : qShare.data;

  return metadata?.password_protected ? (
    <Group p="xl">
      <IconLock />
      <Text>Password protected</Text>
    </Group>
  ) : metadata?.mime_type.startsWith('image/') ? (
    <Image src={downloadUrl} />
  ) : metadata?.mime_type.startsWith('video/') ? (
    <video className="w-full" controls playsInline src={downloadUrl}></video>
  ) : metadata?.mime_type.startsWith('audio/') ? (
    <audio className="w-full" controls playsInline src={downloadUrl}></audio>
  ) : metadata?.mime_type.startsWith('text/') ? (
    <TextFilePreview url={downloadUrl} />
  ) : (
    <Text>{metadata?.file_name || ''}</Text>
  );
};

export default FilePreview;
