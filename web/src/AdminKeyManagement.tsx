import { useState } from 'react';
import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminKeyManagement() {
  // data to be pass in
  let [data, setValue] = useState([
    '128c1d5d-2359-4ba1-8739-2cd30d694d67',
    '128c1d5d-2359-4ba1-8739-2cd30d69sds67',
  ]);

  return (
    <>
      <AppBase userType="admin">
        <AdminConsole
          groupList={data}
          addObjectLabel="Key"
          deleteObjectLabel="Key"
          tableHeader={['Key ID']}
          tableBody={[]}
          caption="API Key Management Console"
          pointerEvents={false}
        ></AdminConsole>
      </AppBase>
    </>
  );
}
