import AppBase from '../AppBase';
import {
  ChonkyActions,
  ChonkyIconName,
  defineFileAction,
  FileActionHandler,
  FileBrowserProps,
  FileData,
  FullFileBrowser,
} from 'chonky';
import { showNotification, updateNotification } from '@mantine/notifications';
import React, { useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  EntryType,
  getFileDownloadURL,
  isUserHome,
  useMutateCopyFile,
  useMutateDeleteFile,
  useMutateMoveFile,
  useMutateUpdateFileMetadata,
  useQueryFolderPathById,
} from '../../api/file';
import {
  useMutateCreateFolder,
  useMutateDeleteFolder,
  useQueryFolderContents,
} from '../../api/folder';
import { useQueryUser } from '../../api/auth';
import UploadFileModal from './UploadFileModal';
import FilePreviewModal from './FilePreviewModal';
import FilePropertiesDrawer from './FilePropertiesDrawer';
import { IconCheck } from '@tabler/icons';

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
    name: 'More',
    toolbar: true,
    contextMenu: true,
    group: 'Actions',
    icon: ChonkyIconName.info,
  },
} as const);

export const VFSBrowser: React.FC<VFSProps> = React.memo((props) => {
  const [fuOpened, setFuOpened] = useState(false);
  const [previewFileId, setPreviewFileId] = useState('');
  const [propertiesFileId, setPropertiesFileId] = useState('');
  const [clipboardIds, setClipboardsIds] = useState<string[]>([]);
  const params = useParams();
  const navigate = useNavigate();

  const qUser = useQueryUser();
  const homeFolderId: string = qUser.data?.home_folder_id || '';
  const currentFolderId = params.id || '';

  const qFilesList = useQueryFolderContents(currentFolderId);
  const qFolderChain = useQueryFolderPathById(currentFolderId);

  const mCreateFolder = useMutateCreateFolder();
  const mDeleteFolder = useMutateDeleteFolder();
  const mDeleteFile = useMutateDeleteFile();
  const mUpdateFileMetadata = useMutateUpdateFileMetadata();
  const mMoveFile = useMutateMoveFile();
  const mCopyFile = useMutateCopyFile();

  useEffect(() => {
    if (!params.id && homeFolderId) navigate(`/home/${homeFolderId}`);
  }, [params, homeFolderId]);

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
        parent_folder_id: currentFolderId,
      });
    } else if (data.id === ChonkyActions.DeleteFiles.id) {
      const loaderNotificationId = 'delete-loader';
      const totalCount = data.state.selectedFilesForAction.length;
      const showLoader = totalCount > 3;
      let deletedCount = 0;

      if (showLoader)
        showNotification({
          id: loaderNotificationId,
          title: 'Deleting files...',
          message: 'Please wait',
          loading: true,
          autoClose: false,
          disallowClose: true,
        });

      for (const selectedItem of data.state.selectedFilesForAction) {
        if (showLoader)
          updateNotification({
            id: loaderNotificationId,
            title: `Deleting files... ${deletedCount + 1} / ${totalCount}`,
            message: selectedItem.name,
            loading: true,
            autoClose: false,
            disallowClose: true,
          });

        if (selectedItem.isDir) {
          await mDeleteFolder
            .mutateAsync(selectedItem.id)
            .then(() => {
              deletedCount++;
              if (!showLoader)
                showNotification({
                  title: `${selectedItem.name} deleted`,
                  message: 'Successfully Deleted',
                });
            })
            .catch((e) =>
              showNotification({
                title: `Error deleting ${selectedItem.name}`,
                message: JSON.stringify(e),
              })
            );
        } else {
          await mDeleteFile
            .mutateAsync(selectedItem.id)
            .then(() => {
              deletedCount++;
              if (!showLoader)
                showNotification({
                  title: `${selectedItem.name} deleted`,
                  message: 'Successfully Deleted',
                });
            })
            .catch((e) =>
              showNotification({
                title: `Error deleting ${selectedItem.name}`,
                message: JSON.stringify(e),
              })
            );
        }
      }

      if (showLoader)
        updateNotification({
          id: loaderNotificationId,
          title: 'Finished deleting',
          message: `Deleted ${deletedCount} / ${totalCount}`,
          color: 'teal',
          icon: <IconCheck size={16} />,
          autoClose: 5000,
        });
    } else if (data.id === ChonkyActions.DownloadFiles.id) {
      if (!data.state.selectedFilesForAction[0].isDir)
        window.location.assign(
          getFileDownloadURL(data.state.selectedFilesForAction[0].id)
        );
    } else if (data.id === ChonkyActions.OpenFiles.id) {
      if (!data.payload.targetFile?.isDir)
        setPreviewFileId(data.payload.targetFile?.id || '');
      else navigate(`/home/${data.payload.targetFile?.id || ''}`);
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
          folder_id: currentFolderId,
        });
      }
    } else if ((data.id as string) === FileProperties.id) {
      setPropertiesFileId(data.state.selectedFilesForAction[0].id);
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

  const ohanaFiles = useMemo(
    () =>
      qFilesList.data?.map?.(
        (file) =>
          ({
            id: file.file_id,
            name: file.file_name,
            isDir: file.entry_type === EntryType.Folder,
            modDate: file.modified_time,
            size: file.size,
            thumbnailUrl:
              file.entry_type === EntryType.File &&
              file.mime_type.startsWith('image/')
                ? getFileDownloadURL(file.file_id, { inline: true })
                : undefined,
          } as FileData)
      ) || [],
    [qFilesList.data]
  );

  const folderChain = useMemo(
    () =>
      (qFolderChain.data ?? [])
        .slice()
        .reverse()
        .map((folder) => ({
          id: folder.file_id,
          name: folder.file_id === homeFolderId ? 'Home' : folder.file_name,
          isDir: true,
        })),
    [homeFolderId, qFolderChain.data]
  );

  return (
    <AppBase userType="user">
      <div style={{ height: '100%' }}>
        <FullFileBrowser
          files={ohanaFiles}
          folderChain={folderChain}
          fileActions={fileActions}
          onFileAction={handleFileAction}
          clearSelectionOnOutsideClick
          {...props}
        />
      </div>
      <UploadFileModal
        onClose={() => setFuOpened(false)}
        opened={fuOpened}
        parentFolderId={currentFolderId}
        update={false}
      />
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
});
