import React from 'react';
import {
  ActionIcon,
  Anchor,
  Group,
  ScrollArea,
  Table,
  Tooltip,
} from '@mantine/core';
import { IconExternalLink, IconInfoCircle } from '@tabler/icons';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { APIClient, typedError } from '../../api/api';
import { EntryType, FileMetadata } from '../../api/file';
import AppBase from '../AppBase';
import FilePreviewModal from './FilePreviewModal';
import FilePropertiesDrawer from './FilePropertiesDrawer';

type SharedFavoritesListProps = { list: 'shared' | 'favorites' };

const SharedFavoritesList = (props: SharedFavoritesListProps) => {
  const [previewFileId, setPreviewFileId] = React.useState('');
  const [propertiesFileId, setPropertiesFileId] = React.useState('');

  const qFilesList = useQuery([props.list], () =>
    (props.list === 'shared'
      ? APIClient.get<FileMetadata[]>('/api/v1/sharedWith')
      : APIClient.get<FileMetadata[]>('/api/v1/favorites')
    )
      .then((res) => res.data)
      .catch(typedError)
  );

  return (
    <AppBase userType="user">
      <ScrollArea>
        <Table>
          <thead>
            <tr>
              <th>File name</th>
              <th>Type</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {qFilesList.data?.map((file) => (
              <tr key={file.file_id}>
                <td>
                  <Anchor onClick={() => setPreviewFileId(file.file_id)}>
                    {file.file_name}
                  </Anchor>
                </td>
                <td>
                  {file.entry_type === EntryType.File
                    ? file.mime_type.split('/')[0]
                    : 'Folder'}
                </td>
                <td>
                  <Group>
                    <Link
                      to={`/home/${
                        file.entry_type === EntryType.Folder
                          ? file.file_id
                          : file.parent_folder_id
                      }`}
                    >
                      <Tooltip label="View in folder">
                        <ActionIcon>
                          <IconExternalLink />
                        </ActionIcon>
                      </Tooltip>
                    </Link>
                    <Tooltip label="Properties">
                      <ActionIcon
                        onClick={() => setPropertiesFileId(file.file_id)}
                      >
                        <IconInfoCircle />
                      </ActionIcon>
                    </Tooltip>
                  </Group>
                </td>
              </tr>
            ))}
          </tbody>
        </Table>
      </ScrollArea>
      <FilePreviewModal
        fileId={previewFileId}
        onClose={() => setPreviewFileId('')}
      />
      <FilePropertiesDrawer
        fileId={propertiesFileId}
        onClose={() => setPropertiesFileId('')}
      />
    </AppBase>
  );
};

export default SharedFavoritesList;
