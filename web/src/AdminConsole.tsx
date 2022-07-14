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

    const userinput = prompt('Please enter ' +props.addObjectLabel)
    setValue((prevValue) => CurrentSSOGroups.concat(userinput));
  }

  const [checkedOne, setCheckedOne] = useState(['']);
  function deleteGroup() {
    checkedOne.forEach((element) => {
      setCheckedOne(checkedOne.filter((item) => item !== element));
      setValue(CurrentSSOGroups.filter((item) => item !== element));
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
          to="/insidessogroup"
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
        height: '80vh',
        justifyContent: 'center',
        
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

   
       <div  style={{
       display: 'flex',
      flexDirection:'row',
      justifyContent: 'space-between',
      
      }} >
     
            <Button
              variant="default"
              color="dark"
              size="md"
              onClick={() => add()}
              style={{ marginLeft:'15px' }}
            >
              Add {props.addObjectLabel}
            </Button>
          
            <Button
              variant="default"
              color="dark"
              size="md"
              style={{ marginRight:'15px' }}
              onClick={() => deleteGroup()}
            >
              Delete {props.deleteObjectLabel}
            </Button>
            </div>
          
      </div>
    </div>
  );
}


