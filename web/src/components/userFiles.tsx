import AppBase from './AppBase';
import {
  ChonkyActions,
  ChonkyFileActionData,
  ChonkyIconName,
  defineFileAction,
  FileActionHandler,
  FileArray,
  FileBrowserProps,
  FileData,
  FileHelper,
  FullFileBrowser,
} from 'chonky';
import { Modal, FileInput, FileButton, Button, Loader } from '@mantine/core';
import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';

import { useNavigate, useParams } from 'react-router-dom';

import {
  EntryType,
  getFileDownloadURL,
  useMutateCopyFile,
  useMutateDeleteFile,
  useMutateMoveFile,
  useMutateUpdateFile,
  useMutateUpdateFileMetadata,
  useMutateUploadFile,
} from '../api/file';
import {
  useMutateCreateFolder,
  useMutateDeleteFolder,
  useQueryFolderContents,
  useQueryFolderContentsByPath,
} from '../api/folder';
import { IconUpload } from '@tabler/icons';
import { useQueryUser } from '../api/auth';

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

export const VFSBrowser: React.FC<VFSProps> = React.memo((props) => {
  const [fuOpened, setFuOpened] = useState(false);
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
          mDeleteFolder.mutate(selectedItem.id);
        } else {
          mDeleteFile.mutate(selectedItem.id);
        }
      }
    } else if (data.id === ChonkyActions.DownloadFiles.id) {
      // @ts-ignore
      window.location = getFileDownloadURL(
        data.state.selectedFilesForAction[0].id
      );
    } else if (data.id === ChonkyActions.OpenFiles.id) {
      navigate(`/home/${data.state.selectedFilesForAction[0].id}`);
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
                .then(() => setFuOpened(false));
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
    </AppBase>
  );
});
