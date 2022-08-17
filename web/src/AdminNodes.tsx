import { Group, Text, Accordion, Title, ThemeIcon } from '@mantine/core';
import { IconExclamationMark } from '@tabler/icons';
import { useState } from 'react';
import {
  useQueryGetAlerts,
  useQueryGetserverStatuses,
  useQueryGetserverStatusesID,
} from './api/cluster';
import AppBase from './components/AppBase';
import { humanFileSize } from './shared/util';

interface AccordionLabelProps {
  name: string;
  warnings: number;
  errors: number;
}

function AccordionLabel({ name, warnings, errors }: AccordionLabelProps) {
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
        <Text>{name}</Text>
        <div
          style={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
          }}
        >
          {warnings > 0 ? (
            <ThemeIcon color="yellow" variant="light">
              <IconExclamationMark size={30} />
            </ThemeIcon>
          ) : null}
          <Text style={{ padding: '5px' }} color="black">
            {warnings > 0 ? warnings : ''}
          </Text>
          {errors > 0 ? (
            <ThemeIcon color="red" variant="light">
              <IconExclamationMark size={30} />
            </ThemeIcon>
          ) : null}
          <Text style={{ padding: '5px' }} color="red">
            {errors > 0 ? errors : ''}
          </Text>
        </div>
      </div>
    </Group>
  );
}

const serverStatusCodeConversion = (code: number) => {
  switch (code) {
    case 1:
      return 'Online';
      break;
    case 2:
      return 'Offline';
      break;
    case 3:
      return 'Starting';
      break;
    case 4:
      return 'Stopping';
      break;
    case 5:
      return 'Needs Attention';
      break;
    case 6:
      return 'Error!';
      break;
    case 7:
      return 'Offline Error!';
      break;

    default:
      return 'Unknown';
      break;
  }
};

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

interface NodeDetailsProps {
  name: string;
}

function NodeDetails({ name }: NodeDetailsProps) {
  const qServerStatusDetails = useQueryGetserverStatusesID(name);

  return (
    <Accordion.Item value={name}>
      <Accordion.Control>
        <AccordionLabel
          name={name}
          warnings={qServerStatusDetails.data?.warnings ?? 0}
          errors={qServerStatusDetails.data?.errors ?? 0}
        />
      </Accordion.Control>
      <Accordion.Panel>
        <Text size="md">
          <Text weight={500} component="span">
            Hostname:{' '}
          </Text>
          {qServerStatusDetails.data?.hostname}
        </Text>
        <Text size="md">
          <Text weight={500} component="span">
            Port:{' '}
          </Text>
          {qServerStatusDetails.data?.port}
        </Text>
        <Text size="md">
          <Text weight={500} component="span">
            Status:{' '}
          </Text>
          {serverStatusCodeConversion(qServerStatusDetails.data?.status ?? 0)}
        </Text>
        <Text size="md">
          <Text weight={500} component="span">
            Storage{' (free/total): '}
          </Text>
          {humanFileSize(qServerStatusDetails.data?.free_space ?? 0) +
            ' / ' +
            humanFileSize(
              (qServerStatusDetails.data?.used_space ?? 0) +
                (qServerStatusDetails.data?.free_space ?? 0)
            )}
        </Text>
        <Text size="md">
          <Text weight={500} component="span">
            Memory{' (free/total): '}
          </Text>
          {humanFileSize(qServerStatusDetails.data?.memory_free ?? 0) +
            ' / ' +
            humanFileSize(
              (qServerStatusDetails.data?.memory_free ?? 0) +
                (qServerStatusDetails.data?.memory_used ?? 0)
            )}
        </Text>
        <Text size="md">
          <Text weight={500} component="span">
            Uptime:&nbsp;
          </Text>
          {secondsToDhms(qServerStatusDetails.data?.uptime ?? 0)}
        </Text>
      </Accordion.Panel>
    </Accordion.Item>
  );
}

export function AdminNodes() {
  const qServerStatus = useQueryGetserverStatuses();
  const serversList = qServerStatus?.data ?? [];

  //console.log(useQueryGetAlerts().data?.log_entries);

  const items = serversList.map((item, i) => (
    <NodeDetails name={item.name} key={i} />
  ));

  return (
    <>
      <AppBase userType="admin">
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
              Nodes Connected: {serversList.length}
            </Title>
            <Accordion>{items}</Accordion>
          </div>
        </div>
      </AppBase>
    </>
  );
}
