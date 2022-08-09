import AppBase from './AppBase';
import {
  ChonkyActions,
  ChonkyIconName,
  defineFileAction,
  FileActionHandler,
  FileBrowserProps,
  FileData,
  FullFileBrowser,
} from 'chonky';
import {
  Modal,
  FileInput,
  FileButton,
  Button,
  Loader,
  Image,
  Drawer,
  Table,
} from '@mantine/core';
import { showNotification, cleanNotifications } from '@mantine/notifications';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import {
  EntryType,
  getFileDownloadURL,
  MetadataKeyMap,
  useMutateCopyFile,
  useMutateDeleteFile,
  useMutateMoveFile,
  useMutateUpdateFileMetadata,
  useMutateUploadFile,
  useQueryFileMetadata,
} from '../api/file';
import {
  useMutateCreateFolder,
  useMutateDeleteFolder,
  useQueryFolderContents,
} from '../api/folder';
import { useQueryUser } from '../api/auth';
import { useMediaQuery } from '@mantine/hooks';

export type VFSProps = Partial<FileBrowserProps>;

const RenameFiles = defineFileAction({
  id: 'rename_files',
  button: {
    name: 'Rename',
    toolbar: true,
    contextMenu: true,
    group: 'Actions',
    icon: ChonkyIconName.config,
  },
} as const);

const PasteFiles = defineFileAction({
  id: 'paste_files',
  button: {
    name: 'Paste',
    toolbar: true,
    contextMenu: true,
    group: 'Actions',
    icon: ChonkyIconName.paste,
  },
} as const);

const FileProperties = defineFileAction({
  id: 'file_properties',
  button: {
    name: 'Properties',
    toolbar: true,
    contextMenu: true,
    group: 'Actions',
    icon: ChonkyIconName.info,
  },
} as const);

const keyMap = {
  file_name: 'Name',
};

export const VFSBrowser: React.FC<VFSProps> = React.memo((props) => {
  const [fuOpened, setFuOpened] = useState(false);
  const [fileOpened, setFileOpened] = useState('');
  const [filePropertiesOpened, setFilePropertiesOpened] = useState('');
  const [clipboardIds, setClipboardsIds] = useState<string[]>([]);
  const params = useParams();
  const navigate = useNavigate();
  const smallScreen = useMediaQuery('(max-width: 600px)');

  const qUser = useQueryUser();
  const homeFolderID: string = qUser.data?.home_folder_id || '';

  useEffect(() => {
    if (!params.id && homeFolderID) navigate(`/home/${homeFolderID}`);
  }, [params, homeFolderID]);

  const folderID = params.id || '';

  const showNotificationFunc = (title: string, message: string) => {
    showNotification({
      title: title,
      message: message,
      onClose: () => cleanNotifications(),
    });
  };

  const handleFileAction: FileActionHandler = async (data) => {
    if (data.action === ChonkyActions.UploadFiles) {
      setFuOpened(true);
    } else if (data.action === ChonkyActions.CreateFolder) {
      let name = window.prompt('Enter new folder name: ');
      if (!name) {
        return;
      }
      mCreateFolder.mutate({
        folder_name: name,
        parent_folder_id: folderID,
      });
    } else if (data.id === ChonkyActions.DeleteFiles.id) {
      for (const selectedItem of data.state.selectedFilesForAction) {
        if (selectedItem.isDir) {
          mDeleteFolder
            .mutateAsync(selectedItem.id)
            .then(() =>
              showNotificationFunc(
                `${selectedItem.name} deleted`,
                'Successfully Deleted'
              )
            );
        } else {
          mDeleteFile
            .mutateAsync(selectedItem.id)
            .then(() =>
              showNotificationFunc(
                `${selectedItem.name} deleted`,
                'Successfully Deleted'
              )
            );
        }
      }
    } else if (data.id === ChonkyActions.DownloadFiles.id) {
      // @ts-ignore
      window.location = getFileDownloadURL(
        data.state.selectedFilesForAction[0].id
      );
    } else if (data.id === ChonkyActions.OpenFiles.id) {
      if (!data.state.selectedFilesForAction[0].isDir)
        setFileOpened(data.state.selectedFilesForAction[0].id);
      else navigate(`/home/${data.state.selectedFilesForAction[0].id}`);
    } else if ((data.id as string) === RenameFiles.id) {
      let newFileName = window.prompt(
        'Enter new file name',
        data.state.selectedFilesForAction[0].name
      );
      if (!newFileName) return;
      mUpdateFileMetadata.mutate({
        file_name: newFileName,
        file_id: data.state.selectedFilesForAction[0].id,
      });
    } else if (data.id === ChonkyActions.MoveFiles.id) {
      if (data.payload.draggedFile.isDir) return;
      mMoveFile.mutate({
        file_id: data.payload.draggedFile.id,
        folder_id: data.payload.destination.id,
      });
    } else if (data.id === ChonkyActions.CopyFiles.id) {
      console.log(data);
      setClipboardsIds(
        data.state.selectedFilesForAction.map((file) => file.id)
      );
    } else if ((data.id as string) === PasteFiles.id) {
      for (const item of clipboardIds) {
        await mCopyFile.mutateAsync({
          file_id: item,
          folder_id: folderID,
        });
      }
    } else if ((data.id as string) === FileProperties.id) {
      setFilePropertiesOpened(data.state.selectedFilesForAction[0].id);
    }
  };

  const fileActions = useMemo(
    () => [
      ChonkyActions.CreateFolder,
      ChonkyActions.DeleteFiles,
      ChonkyActions.UploadFiles,
      ChonkyActions.DownloadFiles,
      ChonkyActions.MoveFiles,
      ChonkyActions.CopyFiles,
      RenameFiles,
      PasteFiles,
      FileProperties,
    ],
    []
  );
  const thumbnailGenerator = useCallback(
    (file: FileData) =>
      file.thumbnailUrl ? `https://chonky.io${file.thumbnailUrl}` : null,
    []
  );

  const qFilesList = useQueryFolderContents(folderID);
  const qActiveFile = useQueryFileMetadata(fileOpened);
  const qPropsFile = useQueryFileMetadata(filePropertiesOpened);

  const mCreateFolder = useMutateCreateFolder();
  const mDeleteFolder = useMutateDeleteFolder();
  const mUploadFile = useMutateUploadFile();
  const mDeleteFile = useMutateDeleteFile();
  const mUpdateFileMetadata = useMutateUpdateFileMetadata();
  const mMoveFile = useMutateMoveFile();
  const mCopyFile = useMutateCopyFile();

  const ohanaFiles =
    qFilesList.data?.map((file) => ({
      id: file.file_id,
      name: file.file_name,
      isDir: file.entry_type === EntryType.Folder,
      modDate: file.modified_time,
      size: file.size,
    })) || [];

  const ohanaFolderChain = [{ id: homeFolderID, name: 'Home', isDir: true }];
  return (
    <AppBase userType="user">
      <div style={{ height: '100%' }}>
        <FullFileBrowser
          files={ohanaFiles}
          folderChain={ohanaFolderChain}
          fileActions={fileActions}
          onFileAction={handleFileAction}
          thumbnailGenerator={thumbnailGenerator}
          {...props}
        />
        {/* Upload file modal */}
      </div>
      <Modal
        centered
        opened={fuOpened}
        onClose={() => setFuOpened(false)}
        title="Upload a File"
      >
        <div className="flex">
          {mUploadFile.isLoading ? <Loader className="mr-5" /> : null}
          <FileButton
            onChange={(item) => {
              console.log('we going in');
              if (!item) {
                return;
              }
              mUploadFile
                .mutateAsync({
                  file: item,
                  folder_id: folderID,
                  frag_count: 1,
                  parity_count: 1,
                })
                .then(() => setFuOpened(false))
                .then(() =>
                  showNotificationFunc(
                    `${item.name} uploaded`,
                    'File Uploaded Successfully'
                  )
                );
            }}
          >
            {(props) => (
              <Button
                disabled={mUploadFile.isLoading}
                className="bg-cyan-500"
                color="cyan"
                {...props}
              >
                Upload a File
              </Button>
            )}
          </FileButton>
        </div>
      </Modal>
      <Modal
        centered
        opened={fileOpened !== ''}
        onClose={() => setFileOpened('')}
        title={qActiveFile.data?.file_name}
        size={smallScreen ? '100%' : '70%'}
      >
        <div className="flex">
          {qActiveFile.data?.mime_type.startsWith('image/') ? (
            <Image src={getFileDownloadURL(fileOpened)} />
          ) : qActiveFile.data?.mime_type.startsWith('video/') ? (
            <video
              controls
              playsInline
              src={getFileDownloadURL(fileOpened)}
            ></video>
          ) : null}
        </div>
        <Button
          component="a"
          href={getFileDownloadURL(fileOpened) + '?inline=1'}
          className="bg-blue-600 mt-5"
          color="blue"
        >
          Download
        </Button>
      </Modal>
      <Drawer
        opened={filePropertiesOpened !== ''}
        onClose={() => setFilePropertiesOpened('')}
        title={qPropsFile.data?.file_name}
        padding="lg"
        position="right"
        size="xl"
      >
        <Table>
          <thead>
            <tr>
              <th>Property</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            {Object.keys(qPropsFile.data || {})
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
                  <td>{(qPropsFile.data as any)[key]}</td>
                </tr>
              ))}
          </tbody>
        </Table>
      </Drawer>
    </AppBase>
  );
});
