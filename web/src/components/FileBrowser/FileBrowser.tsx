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
import { showNotification } from '@mantine/notifications';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  EntryType,
  getFileDownloadURL,
  useMutateCopyFile,
  useMutateDeleteFile,
  useMutateMoveFile,
  useMutateUpdateFileMetadata,
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

export const VFSBrowser: React.FC<VFSProps> = React.memo((props) => {
  const [fuOpened, setFuOpened] = useState(false);
  const [previewFileId, setPreviewFileId] = useState('');
  const [propertiesFileId, setPropertiesFileId] = useState('');
  const [clipboardIds, setClipboardsIds] = useState<string[]>([]);
  const params = useParams();
  const navigate = useNavigate();

  const qUser = useQueryUser();
  const homeFolderID: string = qUser.data?.home_folder_id || '';

  useEffect(() => {
    if (!params.id && homeFolderID) navigate(`/home/${homeFolderID}`);
  }, [params, homeFolderID]);

  const folderID = params.id || '';

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
          mDeleteFolder.mutateAsync(selectedItem.id).then(() =>
            showNotification({
              title: `${selectedItem.name} deleted`,
              message: 'Successfully Deleted',
            })
          );
        } else {
          mDeleteFile.mutateAsync(selectedItem.id).then(() =>
            showNotification({
              title: `${selectedItem.name} deleted`,
              message: 'Successfully Deleted',
            })
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
        setPreviewFileId(data.state.selectedFilesForAction[0].id);
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
  const thumbnailGenerator = useCallback(
    (file: FileData) =>
      file.thumbnailUrl ? `https://chonky.io${file.thumbnailUrl}` : null,
    []
  );

  const qFilesList = useQueryFolderContents(folderID);

  const mCreateFolder = useMutateCreateFolder();
  const mDeleteFolder = useMutateDeleteFolder();
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
      </div>
      <UploadFileModal
        onClose={() => setFuOpened(false)}
        opened={fuOpened}
        parentFolderId={folderID}
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
