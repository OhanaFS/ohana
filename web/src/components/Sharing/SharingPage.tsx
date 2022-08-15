import {
  Alert,
  Button,
  Center,
  Container,
  Loader,
  Title,
  Stack,
} from '@mantine/core';
import { IconAlertCircle, IconDownload } from '@tabler/icons';
import { useParams } from 'react-router-dom';
import {
  getSharingLinkURL,
  useQuerySharingLinkMetadata,
} from '../../api/sharing';
import FilePreview from '../FileBrowser/FilePreview';
import PropertiesTable from '../FileBrowser/Properties/PropertiesTable';

const SharingPage = () => {
  const params = useParams();
  const shareId = params.id || '';
  const qMeta = useQuerySharingLinkMetadata(shareId);
  const urlDownload = getSharingLinkURL(shareId, 'download');

  if (qMeta.isError)
    return (
      <Container py="xl">
        <Alert icon={<IconAlertCircle size={16} />} title="Error!" color="red">
          We're unable to find the requested file. It may have been deleted, or
          the link may have been revoked.
        </Alert>
      </Container>
    );

  if (qMeta.isLoading || !qMeta.data)
    return (
      <Container py="xl">
        <Center p="xl">
          <Loader />
        </Center>
      </Container>
    );

  return (
    <Container py="xl">
      <Stack spacing="md">
        <Title>File share</Title>
        <FilePreview shareId={shareId} />
        <Button
          component="a"
          href={urlDownload}
          target="_blank"
          leftIcon={<IconDownload />}
        >
          Download
        </Button>
        <PropertiesTable metadata={qMeta.data} />
      </Stack>
    </Container>
  );
};

export default SharingPage;
