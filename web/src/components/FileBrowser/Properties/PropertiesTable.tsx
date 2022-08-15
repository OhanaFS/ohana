import { Table } from '@mantine/core';
import { FileMetadata } from '../../../api/file';
import { formatDateTime, humanFileSize } from '../../../shared/util';

export type PropertiesTableProps = {
  metadata: FileMetadata;
};

const Row = (props: { title: string; children: React.ReactNode }) => (
  <tr>
    <td className="font-bold whitespace-nowrap">{props.title}</td>
    <td>{props.children}</td>
  </tr>
);

const PropertiesTable = (props: PropertiesTableProps) => {
  return (
    <Table>
      <tbody>
        <Row title="File name">{props.metadata.file_name}</Row>
        <Row title="Size">{humanFileSize(props.metadata.size)}</Row>
        <Row title="Created at">
          {formatDateTime(props.metadata.created_time)}
        </Row>
        <Row title="Last modified">
          {formatDateTime(props.metadata.modified_time)}
        </Row>
      </tbody>
    </Table>
  );
};

export default PropertiesTable;
