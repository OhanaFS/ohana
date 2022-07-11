import {
  Table,
  Button,
  Text,
  Checkbox,
  useMantineTheme,
  ScrollArea,
} from '@mantine/core';

import { Link } from 'react-router-dom';
import { useState } from 'react';

export interface ConsoleDetails {
  groupList: Array<any>;
  addObjectLabel: string;
  deleteObjectLabel: string;
  tableHeader: Array<string>;
  tableBody: Array<string>;
  caption: string;
  pointerEvents: boolean;
  consoleWidth: number;
  consoleHeight: number;
}

export function AdminConsole(props: ConsoleDetails) {
  // take from database.

  const theme = useMantineTheme();
  let [CurrentSSOGroups, setValue] = useState(props.groupList);
  function add() {
    setValue((prevValue) => CurrentSSOGroups.concat('asd'));
  }

  const [checkedOne, setCheckedOne] = useState(['']);

  function deleteGroup() {
    checkedOne.forEach((element) => {
      setCheckedOne(checkedOne.filter((item) => item !== element));
      setValue(CurrentSSOGroups.filter((item) => item !== element));

      console.log(CurrentSSOGroups);
    });
  }
  function update(index: string) {
    setCheckedOne((prevValue) => checkedOne.concat(index));
  }
  function remove(index: string) {
    setCheckedOne(checkedOne.filter((item) => item !== index));
  }

  const ths = (
    <tr>
      <th
        style={{
          width: '80%',
          textAlign: 'left',
          fontWeight: '700',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {props.tableHeader}
      </th>
    </tr>
  );
  const rows = CurrentSSOGroups.map((items, index) => (
    <tr>
      <td
        width="80%"
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        <Text
          color="dark"
          style={{
            marginBottom: '20%',
            height: '50px',
            pointerEvents: props.pointerEvents ? 'auto' : 'none',
          }}
          component={Link}
          to="/Admin_ssogroups_inside"
          variant="link"
        >
          {items}
        </Text>
      </td>
      <td>
        <Checkbox
          onChange={(event) =>
            event.currentTarget.checked ? update(items) : remove(items)
          }
        ></Checkbox>
      </td>
    </tr>
  ));

  return (
    <div
      style={{
        display: 'flex',

        justifyContent: 'center',
        height: '80vh',
      }}
    >
      <div className="console">
        <ScrollArea style={{ height: '90%', width: '100%', marginTop: '1%' }}>
          <Table captionSide="top" verticalSpacing="sm" style={{}}>
            <caption
              style={{
                textAlign: 'center',
                fontWeight: '600',
                fontSize: '24px',
                color: 'black',
                marginTop: '2%',
              }}
            >
              {props.caption}{' '}
            </caption>
            <thead>{ths}</thead>

            <tbody>{rows}</tbody>
          </Table>
        </ScrollArea>

        <div style={{ position: 'relative' }}>
          <td
            style={{
              position: 'absolute',
              top: '5px',
              left: '16px',
            }}
          >
            {' '}
            <Button
              variant="default"
              color="dark"
              size="md"
              style={{}}
              onClick={() => add()}
            >
              Add {props.addObjectLabel}
            </Button>
          </td>

          <td
            style={{
              position: 'absolute',
              top: '5px',
              right: '16px',
            }}
          >
            {' '}
            <Button
              variant="default"
              color="dark"
              size="md"
              style={{}}
              onClick={() => deleteGroup()}
            >
              Delete {props.deleteObjectLabel}
            </Button>
          </td>
        </div>
      </div>
    </div>
  );
}


