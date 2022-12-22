import { useEffect, useReducer, useState } from 'react';
import { uid } from 'react-uid';
import { GuildInfo, User } from '../api';
import { Card } from '../components/Card';
import { Flex } from '../components/Flex';
import { RouteContainer } from '../components/RouteContainer';
import { UserTile } from '../components/UserTile';
import { useApi } from '../hooks/useApi';
import { ReactComponent as IconCross } from '..//assets/cross.svg';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import styled from 'styled-components';
import { useSnackBar } from '../hooks/useSnackBar';
import { GuildTile } from '../components/GuildTile';

type Props = {};

const AdminControls = styled.div`
  display: flex;
  width: 100%;
  gap: 1em;
  white-space: nowrap;
  margin-bottom: 1em;

  > input {
    width: 100%;
  }
`;

const GuildControls = styled.div``;

const StyledGuildTile = styled(GuildTile)`
  width: 100%;
  background-color: ${(p) => p.theme.background3} !important;
`;

type AdminsReducerAction =
  | {
      type: 'add';
      payload: User;
    }
  | {
      type: 'remove';
      payload: string;
    }
  | {
      type: 'set';
      payload: User[];
    };

const adminsReducer = (state: User[], action: AdminsReducerAction) => {
  switch (action.type) {
    case 'set':
      return action.payload;
    case 'add':
      if (state.find((u) => u.id === action.payload.id)) return state;
      return [...state, action.payload];
    case 'remove':
      const i = state.findIndex((u) => u.id === action.payload);
      if (i === -1) return state;
      state.splice(i, 1);
      return state;
    default:
      return state;
  }
};

export const AdminRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const { show } = useSnackBar();
  const [admins, adminsDispatch] = useReducer(adminsReducer, []);
  const [userID, setUserID] = useState('');
  const [guilds, setGuilds] = useState<GuildInfo[]>();

  const _addAdmin = async () => {
    try {
      const user = await fetch((c) => c.setAdmin(userID));
      adminsDispatch({ type: 'add', payload: user });
      setUserID('');
      show(`${user.username} has been added as admin.`, 'success');
    } catch {}
  };

  const _removeAdmin = async (id: string) => {
    try {
      await fetch((c) => c.removeAdmin(id));
      adminsDispatch({ type: 'remove', payload: id });
      show('The user has been removed from admins.', 'success');
    } catch {}
  };

  const _removeGuild = async (id: string) => {
    try {
      await fetch((c) => c.removeGuild(id));
      setGuilds(guilds?.filter((g) => g.id != id));
      show('Yuri has left the guild.', 'success');
    } catch {}
  };

  useEffect(() => {
    fetch((c) => c.admins())
      .then((res) => adminsDispatch({ type: 'set', payload: res }))
      .catch();

    fetch((c) => c.guilds())
      .then((g) => setGuilds(g))
      .catch();
  }, []);

  return (
    <RouteContainer>
      <h1>Admin Area</h1>

      <Card margin="0 0 1.5em 0">
        <h2>Admins</h2>
        <AdminControls>
          <Input
            placeholder="User ID"
            value={userID}
            onInput={(e) => setUserID(e.currentTarget.value)}
          />
          <Button variant="green" disabled={!userID} onClick={_addAdmin}>
            Add Admin
          </Button>
        </AdminControls>
        {admins && (
          <Flex direction="column" gap="0.5em">
            {admins.map((a) => (
              <Flex key={uid(a)} gap="1em" vCenter>
                <UserTile user={a} />
                <Button variant="red" disabled={a.is_owner} onClick={() => _removeAdmin(a.id)}>
                  <IconCross />
                  Remove
                </Button>
              </Flex>
            ))}
          </Flex>
        )}
      </Card>

      <Card>
        <h2>Guilds</h2>
        <GuildControls>
          {guilds && (
            <Flex direction="column" gap="0.5em">
              {guilds.map((g) => (
                <Flex key={uid(g)} gap="1em" vCenter>
                  <StyledGuildTile guild={g} />
                  <Button variant="red" onClick={() => _removeGuild(g.id)}>
                    <IconCross />
                    Remove
                  </Button>
                </Flex>
              ))}
            </Flex>
          )}
        </GuildControls>
      </Card>
    </RouteContainer>
  );
};
