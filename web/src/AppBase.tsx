import { AppShell, Navbar, Header, Title } from '@mantine/core';
import { MainLinks } from './MainLinks';
import { User } from './user';


type AppBaseProps = {
    children: React.ReactNode;
}
export default function AppBase(props: AppBaseProps) {
    return (
      <AppShell
      padding="md"
      navbar={<Navbar width={{ sm: 300 }} p="md">
        <Navbar.Section grow mt="md">
          <MainLinks />
        </Navbar.Section>
        <Navbar.Section>
          <User />
        </Navbar.Section> 
      </Navbar>}
      header={<Header height={60} p="xs">{<Title order={2}>Ohana</Title>}</Header>}
      styles={(theme) => ({
        main: { backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0] },
      })}
      >
        {props.children}
      </AppShell>
    );
}