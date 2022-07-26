import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminSsoGroupsInside() {
  const data: Array<any> = ['Tom', 'Peter', 'Raymond'];

  return (
    <>
      <AppBase userType="admin">
        <AdminConsole
          groupList={data}
          addObjectLabel="User"
          deleteObjectLabel="User"
          tableHeader={['List of Users inside this group']}
          tableBody={[]}
          caption="User Management Console"
          pointerEvents={false}
        ></AdminConsole>
      </AppBase>
    </>
  );
}
