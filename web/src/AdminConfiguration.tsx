import {
  Button,
  useMantineTheme,
  Textarea,
  Checkbox,
  Table,
} from '@mantine/core';

import AppBase from './components/AppBase';

export function AdminConfiguration() {
  const theme = useMantineTheme();
  return (
  
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
            alignItems: 'flex-start',
           
          }}

        >
          <div className="rotateKey">
              <Table captionSide="top" verticalSpacing="md">
                <caption
                  style={{
                    textAlign: 'center',
                    fontWeight: 600,
                    fontSize: '24px',
                    color: 'black',
                   
                  }}
                >
                  Rotate Key
                </caption>
                <tbody>
                  <tr>
                    <td
                      style={{
                        textAlign: 'left',
                        fontWeight: 400,
                        fontSize: '16px',
                        color: 'black',
                        border: 'none',
                      }}
                    >
                      {' '}
                      Specify the file/directory location and the system will
                      auto rotate the key
                    </td>
                  </tr>
                  <tr>
                    <td style={{ border: 'none' }}>
                      <Textarea
                        style={{}}
                        label="File location"
                        radius="xs"
                        size="md"
                      />
                    </td>
                  </tr>
                  <tr>
                    <td
                      style={{
                        border: 'none',
                        display: 'flex',
                        textAlign: 'left',
                        fontWeight: 400,
                        fontSize: '16px',
                        color: 'black',
                      }}
                    >
                      Master Key :{' '}
                      <Checkbox style={{ marginLeft: '10px' }}> </Checkbox>
                    </td>
                  </tr>
              

                <tr style={{ position: 'relative' }}>
               
                <Button
                  variant="default"
                  color="dark"
                  size="md"
                  style={{
                    position: 'absolute',
                    top: '18px',
                    right: '16px'
                  }}
                >
                  Rotate Key
                </Button>
                
                </tr>
                </tbody>
              </Table>

     
                </div>
        </div>
      </AppBase>
    
  );
}

