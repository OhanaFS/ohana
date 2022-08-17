import { Button, Group, Modal, Stack } from '@mantine/core';
import { useQueryFileMetadata } from '../../../api/file';
import GeneralAccess from './GeneralAccess';
import LocalAccess from './LocalAccess';

export type SharingModalProps = {
  fileId: string;
  opened: boolean;
  onClose: () => any;
};

const SharingModal = (props: SharingModalProps) => {
  const { fileId, opened, onClose } = props;

  const qFile = useQueryFileMetadata(fileId);

  return (
    <Modal
      centered
      opened={opened}
      onClose={onClose}
      size="lg"
      title={`Share "${qFile.data?.file_name}"`}
      overflow="outside"
      styles={(theme) => ({ title: { fontSize: theme.fontSizes.xl } })}
    >
      <Stack>
        <LocalAccess fileId={fileId} />
        <GeneralAccess fileId={fileId} />

        <Group position="right">
          <Button onClick={onClose}>Done</Button>
        </Group>
      </Stack>
    </Modal>
  );
};

export default SharingModal;
