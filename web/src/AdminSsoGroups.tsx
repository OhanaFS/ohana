import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminSsoGroups() {
  const SSOGroupList = ['Hr', 'Finance', 'IT', 'Management'];

  return (
    <>
      <AppBase>
        <AdminConsole
          consoleWidth={80}
          consoleHeight={60}
          groupList={SSOGroupList}
          addObjectLabel="Group"
          deleteObjectLabel="Group"
          tableHeader={['Current SSO Groups']}
          tableBody={[]}
          caption="SSO Group Management Console"
          pointerEvents={true}
        ></AdminConsole>
      </AppBase>
    </>
  );
}
