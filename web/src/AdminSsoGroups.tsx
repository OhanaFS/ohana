import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminSsoGroups() {
  const SSOGroupList = ['Hr', 'Finance', 'IT', 'Management'];

  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
        <AdminConsole
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
