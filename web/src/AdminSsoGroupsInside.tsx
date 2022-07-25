import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminSsoGroupsInside() {
  const data: Array<any> = ['Tom', 'Peter', 'Raymond'];

  return (
    <>
      <AppBase>
        <AdminConsole
          consoleWidth={80}
          consoleHeight={60}
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
