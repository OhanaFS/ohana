import {
  Table,
  Button,
  Text,
  Checkbox,
  useMantineTheme,
  ScrollArea,
  Modal,
  Group,
  Textarea,
} from '@mantine/core';

import { Link } from 'react-router-dom';
import { useRef, useState } from 'react';
import { group } from 'console';

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

  const [CurrentSSOGroups, setValue] = useState(props.groupList);
  const [Group, addGroup] = useState('');
  function add() {  
    setValue(CurrentSSOGroups.concat(Group));
    setOpened(false);
  }
  const [openedAddUserModel, setOpened] = useState(false);
  const title = "Add "+ props.addObjectLabel;
  const textField ="Name of the " +props.addObjectLabel;

  const deleteGroup = (index:any) => {
    setValue(CurrentSSOGroups.filter((v, i) => i !== index));
}
  
  const ths = (
    <tr>
      <th
        style={{
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
    <tr key={index}>
     <td style={{
        display:'flex',
        justifyContent:'space-between',
      }}>
        <Text
          color="dark"
          style={{
            marginLeft:'10px',
            pointerEvents: props.pointerEvents ? 'auto' : 'none',
          }}
          component={Link}
          to="/insidessogroup"
          variant="link"
        >
          {items}
        </Text>
        <Button
              variant="default"
              color="dark"
              size="md"
              style={{ marginRight:'15px' }}
              onClick={() => [         
               deleteGroup(index)            
              ]}
            >
              Delete 
            </Button>
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
      <Modal
         centered 
         opened={openedAddUserModel}
         onClose={() => setOpened(false)} 
      >
        {         
      <div style ={{    
             display: 'flex',
             height: '25vh',
             flexDirection:'column',
            }}>
            <Textarea
            placeholder={textField}
            label= {title}
            size="md"
            name="password"
            id="password"
            required
            onChange={(event) => {addGroup(event.target.value)}}    
            />  
            <Button
              variant="default"
              color="dark"
              size="md"
              onClick={() => add()}
              style={{ marginLeft:'15px', alignSelf:'flex-end',marginTop:'20px' }}
            >
              Submit 
            </Button>
            </div>
        }
      </Modal>

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

   
       <div style={{
             display: 'flex',
             flexDirection:'row',
            justifyContent: 'space-between',
      }} >
            <Button
              variant="default"
              color="dark"
              size="md"
              onClick={() => setOpened(true)}
              style={{ marginLeft:'15px' }}
            >
              Add {props.addObjectLabel}
            </Button>   
            </div>       
      </div>
    </div>
  );
}


