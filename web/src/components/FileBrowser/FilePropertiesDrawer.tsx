import { AspectRatio, Drawer, Table } from '@mantine/core';
import { MetadataKeyMap, useQueryFileMetadata } from '../../api/file';
import FilePreview from './FilePreview';

export type FilePropertiesDrawerProps = {
  fileId: string;
  onClose: () => void;
};

const FilePropertiesDrawer = (props: FilePropertiesDrawerProps) => {
  const { fileId, onClose } = props;
  const qFile = useQueryFileMetadata(fileId);

  return (
    <Drawer
      opened={fileId !== ''}
      onClose={onClose}
      title={qFile.data?.file_name}
      padding="lg"
      position="right"
      size="xl"
    >
      <AspectRatio ratio={16 / 9}>
        <FilePreview fileId={fileId} />
      </AspectRatio>
      <Table>
        <thead>
          <tr>
            <th>Property</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          {Object.keys(qFile.data || {})
            .map((key) => key as keyof typeof MetadataKeyMap)
            .filter((key) =>
              (
                [
                  'file_name',
                  'size',
                  'created_time',
                  'modified_time',
                  'version_no',
                ] as Array<keyof typeof MetadataKeyMap>
              ).includes(key)
            )
            .map((key) => (
              <tr key={key}>
                <td>{MetadataKeyMap[key]}</td>
                <td>{(qFile.data as any)[key]}</td>
              </tr>
            ))}
        </tbody>
      </Table>
    </Drawer>
  );
};

export default FilePropertiesDrawer;
