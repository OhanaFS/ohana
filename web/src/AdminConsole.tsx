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
  // take from database.

  const theme = useMantineTheme();
  let [CurrentSSOGroups, setValue] = useState(props.groupList);
  const [Group, addGroup] = useState('');
  function add() {
    setValue((prevValue) => CurrentSSOGroups.concat(Group));
    setOpened(false);
  }
  const [openedModel, setOpened] = useState(false);
  const inputRef = useRef(null);
  const [checkboxList, setCheckedOne] = useState(['']);
  const title = "Add "+ props.addObjectLabel;
  const textField ="Name of the " +props.addObjectLabel;

  function deleteGroup() {
    checkboxList.forEach((element) => {
      setCheckedOne(checkboxList.filter((item) => item !== element));
      setValue(CurrentSSOGroups.filter((item) => item !== element))
    });
  }

  // if the user check the checkbox, it will add the item
  function update(index: string) {
    setCheckedOne((prevValue) => checkboxList.concat(index));
  }

    // if the user uncheck the checkbox, it will remove the item
  function remove(index: string) {
    setCheckedOne(checkboxList.filter((item) => item !== index));
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
    <tr style={{
    
    }}>
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
     
        <Checkbox
          style={{marginRight:'10px'}}
          onChange={(event) =>    
            [   
              event.currentTarget.checked ? update(items) : remove(items),
            ]
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


      <Modal
      centered 
        opened={openedModel}
        onClose={() => setOpened(false)} 
      >
        {
          <>
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
              submit 
            </Button>
            </div>
          </>
          
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

   
       <div  style={{
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


