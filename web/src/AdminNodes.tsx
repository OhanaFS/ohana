import { Group, Text, Accordion, Title, ThemeIcon } from '@mantine/core';
import AppBase from './components/AppBase';
import { ExclamationMark } from 'tabler-icons-react';

const charactersList = [
  {
    label: 'Bender Bending Rodríguez',
    ip: '192.168.1.1',
    warnings: 1,
    uptime: 356000,
    loadavg: [0.5, 0.5, 0.5],
  },

  {
    label: 'Carol Miller',
    ip: '192.168.1.23',
    warnings: 0,
    uptime: 350,
    loadavg: [0.5, 0.5, 0.5],
  },
  {
    label: 'Homer Simpson',
    ip: '192.168.1.3',
    warnings: 2,
    uptime: 30,
    loadavg: [0.5, 0.5, 0.5],
  },
  {
    label: 'Spongebob Squarepants',
    ip: '192.168.1.77',
    warnings: 0,
    uptime: 35600,
    loadavg: [0.5, 0.5, 0.5],
  },
];

interface AccordionLabelProps {
  label: string;
  warnings: number;
}

function AccordionLabel({ label, warnings }: AccordionLabelProps) {
  return (
    <Group noWrap>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
          alignItems: 'center',
          width: '100%',
        }}
      >
        <Text>{label}</Text>
        <Text color="red">{warnings > 0 ? warnings : ''}</Text>
      </div>
    </Group>
  );
}

function secondsToDhms(seconds: number) {
  seconds = Number(seconds);
  var d = Math.floor(seconds / (3600 * 24));
  var h = Math.floor((seconds % (3600 * 24)) / 3600);
  var m = Math.floor((seconds % 3600) / 60);
  var s = Math.floor(seconds % 60);

  var dDisplay = d > 0 ? d + (d == 1 ? ' day, ' : ' days, ') : '';
  var hDisplay = h > 0 ? h + (h == 1 ? ' hour, ' : ' hrs, ') : '';
  var mDisplay = m > 0 ? m + (m == 1 ? ' minute, ' : ' mins, ') : '';
  var sDisplay = s > 0 ? s + (s == 1 ? ' second' : ' secs') : '';
  return dDisplay + hDisplay + mDisplay + sDisplay;
}

export function AdminNodes() {
  const items = charactersList.map((item) => (
    <Accordion.Item
      label={<AccordionLabel {...item} />}
      key={item.label}
      icon={
        item.warnings > 0 ? (
          <ThemeIcon color="red" variant="light">
            <ExclamationMark size={30} />
          </ThemeIcon>
        ) : null
      }
    >
      <Text size="md">
        <Text weight={500} component="span">
          IP-Address:&nbsp;
        </Text>
        {item.ip}
      </Text>
      <Text size="md">
        <Text weight={500} component="span">
          Uptime:&nbsp;
        </Text>
        {secondsToDhms(item.uptime)}
      </Text>
      <Text size="md">
        <Text weight={500} component="span">
          Load Average:&nbsp;
        </Text>
        {item.loadavg.join(', ')}
      </Text>
    </Accordion.Item>
  ));

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
          }}
        >
          <div
            className="w-full md:w-9/12 p-2 md:p-4"
            style={{
              backgroundColor: 'white',
              borderRadius: '10px',
              border: 'none',
              overflow: 'hidden',
              paddingBottom: '40px',
            }}
          >
            <Title
              order={2}
              style={{
                paddingLeft: '16px',
                marginTop: '10px',
                marginBottom: '20px',
              }}
            >
              Nodes Summary
            </Title>
            <Accordion
              disableIconRotation
              initialItem={-1}
              iconPosition="right"
            >
              {items}
            </Accordion>
          </div>
        </div>
      </AppBase>
    </>
  );
}
