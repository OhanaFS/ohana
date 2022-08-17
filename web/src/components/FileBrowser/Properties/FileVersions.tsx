import { useState } from 'react';
import { ActionIcon, Button, Table } from '@mantine/core';
import { IconDownload, IconTrash } from '@tabler/icons';
import {
  useMutateDeleteFileVersion,
  useQueryFileVersionHistory,
} from '../../../api/file';
import UploadFileModal from '../UploadFileModal';

type FileVersionsProps = {
  fileId: string;
};
const FileVersions = (props: FileVersionsProps) => {
  const qFileVersionHistory = useQueryFileVersionHistory(props.fileId);
  const mDeleteFileVersion = useMutateDeleteFileVersion();

  const [isUploadOpen, setUploadOpen] = useState(false);
  return (
    <>
      <Button onClick={() => setUploadOpen(true)}>Upload New Version</Button>
      <UploadFileModal
        opened={isUploadOpen}
        onClose={() => setUploadOpen(false)}
        update
        updateFileId={props.fileId}
      />
      <Table>
        <thead>
          <tr>
            <th>Ver</th>
            <th>Last Modified</th>
            <th></th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {qFileVersionHistory.data
            ?.sort((a, b) => b.version_no - a.version_no)
            .map((version, key) => (
              <tr key={key}>
                <td>{version.version_no}</td>
                <td>{version.modified_time}</td>
                <td>
                  <ActionIcon
                    onClick={() =>
                      mDeleteFileVersion.mutate({
                        file_id: props.fileId,
                        version_id: version.version_no,
                      })
                    }
                  >
                    <IconTrash className="text-red-500" />
                  </ActionIcon>
                </td>
                <td>
                  <ActionIcon>
                    <IconDownload className="text-blue-500" />
                  </ActionIcon>
                </td>
              </tr>
            ))}
        </tbody>
      </Table>
    </>
  );
};

export default FileVersions;
