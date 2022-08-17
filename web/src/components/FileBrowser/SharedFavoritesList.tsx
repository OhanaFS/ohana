import { Card, ScrollArea, SimpleGrid } from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { APIClient, typedError } from '../../api/api';
import { FileMetadata } from '../../api/file';
import AppBase from '../AppBase';

type SharedFavoritesListProps = { list: 'shared' | 'favorites' };

const SharedFavoritesList = (props: SharedFavoritesListProps) => {
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
        <SimpleGrid
          cols={4}
          spacing="lg"
          breakpoints={[
            { maxWidth: 980, cols: 3, spacing: 'md' },
            { maxWidth: 755, cols: 2, spacing: 'sm' },
            { maxWidth: 600, cols: 1, spacing: 'sm' },
          ]}
        >
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
          <Card>a</Card>
        </SimpleGrid>
      </ScrollArea>
    </AppBase>
  );
};

export default SharedFavoritesList;
