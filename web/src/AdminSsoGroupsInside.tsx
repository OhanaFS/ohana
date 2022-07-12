import {
  Table,
  Button,
  ScrollArea,
  Checkbox,
  useMantineTheme,
} from '@mantine/core';
import { useState } from 'react';
import AppBase from './components/AppBase';




export function AdminSsoGroupsInside() {
  
  const data: Array<any>= ['Tom', 'Peter', 'Raymond'];
  let [CurrentUser, setValue] = useState(data);
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
        List of Users inside this group
      </th>
    </tr>
  );
  const rows = CurrentUser.map((items) => (
    <tr>
      <td
        
        style={{
          textAlign: 'left',
          fontWeight: '400',
          fontSize: '16px',
          color: 'black',
        }}
      >
        {items}
      </td>
      <td>
        <Checkbox></Checkbox>{' '}
      </td>
    </tr>
  ));

  function addUser() {

    const userinput = prompt('Please enter user' );
    setValue((prevValue) => CurrentUser.concat(userinput));
    
  }

  

  const [checkedOne, setCheckedOne] = useState(['']);
  function deleteUser() {
    checkedOne.forEach((element) => {
      setCheckedOne(checkedOne.filter((item) => item !== element));
      setValue(data.filter((item) => item !== element));
    });
  }

  function remove(index: string) {
    setCheckedOne(checkedOne.filter((item) => item !== index));
  }
  const theme = useMantineTheme();

  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
      <div
        style={{
        display: 'flex',
        justifyContent: 'center',
        height: '80vh',
      }}
    >
      <div className="ssoGroupsDetails">
             
         

                <ScrollArea
                  style={{ height: '90%', width: '100%', marginTop: '1%' }}
                >
                  <Table
                    captionSide="top"
                    striped
                    highlightOnHover
                    verticalSpacing="sm"
                  >
                    <caption
                      style={{
                        textAlign: 'center',
                        fontWeight: '600',
                        fontSize: '24px',
                        color: 'black',
                      }}
                    >
                      User Management Console
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
              onClick={() => addUser()}
            >
              Add User
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
              onClick={() => deleteUser()}
            >
              Delete User
            </Button>
          </td>
        </div>
              
           
              </div>
          </div>
        
      </AppBase>
    </>
  );
}


