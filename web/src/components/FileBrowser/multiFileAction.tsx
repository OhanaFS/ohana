import { showNotification, updateNotification } from '@mantine/notifications';
import { IconCheck } from '@tabler/icons';

export type MultiFileActionOptions<T> = {
  notifications: {
    loadingTitle: (
      successCount: number,
      failCount: number,
      total: number
    ) => string;
    doneTitle: string;
    errorTitle: (item: T, error: any) => string;
    itemName: (item: T) => string;
  };
  items: T[];
  handler: (item: T) => Promise<any>;
};

export const handleMultiFileAction = async function <T>(
  opts: MultiFileActionOptions<T>
) {
  const notificationId = Math.random().toString(36);
  const totalCount = opts.items.length;
  let successCount = 0;
  let failCount = 0;

  showNotification({
    id: notificationId,
    title: opts.notifications.loadingTitle(successCount, failCount, totalCount),
    message: 'Please wait',
    loading: true,
    autoClose: false,
    disallowClose: true,
  });

  for (const item of opts.items) {
    updateNotification({
      id: notificationId,
      title: opts.notifications.loadingTitle(
        successCount,
        failCount,
        totalCount
      ),
      message: opts.notifications.itemName(item),
      loading: true,
      autoClose: false,
      disallowClose: true,
    });

    await opts
      .handler(item)
      .then(() => {
        successCount++;
      })
      .catch((e) => {
        failCount++;
        showNotification({
          title: opts.notifications.errorTitle(item, e),
          message: JSON.stringify(e),
        });
      });
  }

  updateNotification({
    id: notificationId,
    title: opts.notifications.doneTitle,
    message: `${successCount} of ${totalCount} files`,
    color: 'teal',
    icon: <IconCheck size={16} />,
    autoClose: 5000,
  });
};
