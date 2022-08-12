import { ActionIcon, Button, Table } from '@mantine/core';
import React, { useState } from 'react';
import { Download, Trash } from 'tabler-icons-react';
import {
  MetadataKeyMap,
  useMutateDeleteFileVersion,
  useMutateUpdateFile,
  useQueryFileVersionHistory,
} from '../../../api/file';
import UploadFileModal from '../UploadFileModal';

type FileVersionsProps = {
  fileId: string;
};
const FileVersions = (props: FileVersionsProps) => {
  const qFileVersionHistory = useQueryFileVersionHistory(props.fileId);
  const mUpdateFile = useMutateUpdateFile();
  const mDeleteFileVersion = useMutateDeleteFileVersion();

  const [isUploadOpen, setUploadOpen] = useState(false);
  return (
    <>
      <Button onClick={() => setUploadOpen(true)}>Upload New Version</Button>
      <UploadFileModal
        opened={isUploadOpen}
        onClose={() => setUploadOpen(false)}
        update={true}
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
          {qFileVersionHistory.data?.map((version) => (
            <tr>
              <td>{version.version_no}</td>
              <td>{version.modified_time}</td>
              <td>
                <ActionIcon>
                  <Trash color="red" />
                </ActionIcon>
              </td>
              <td>
                <ActionIcon>
                  <Download color="blue" />
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
