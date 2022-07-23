
import { AdminConsole } from './AdminConsole';
import AppBase from './components/AppBase';

export function AdminSsoGroupsInside() {
  
  const data: Array<any>= ['Tom', 'Peter', 'Raymond'];
  
  return (
    <>
      <AppBase
        userType="admin"
        name="Alex Simmons"
        username="@alex"
        image="https://images.unsplash.com/photo-1496302662116-35cc4f36df92?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2070&q=80"
      >
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


