import { useEffect, useState } from 'react';
import { uid } from 'react-uid';
import { User } from '../api';
import { Card } from '../components/Card';
import { Flex } from '../components/Flex';
import { RouteContainer } from '../components/RouteContainer';
import { UserTile } from '../components/UserTile';
import { useApi } from '../hooks/useApi';
import { ReactComponent as IconCross } from '../../assets/cross.svg';
import { Button } from '../components/Button';
import { Input } from '../components/Input';
import styled from 'styled-components';

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

export const AdminRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const [admins, setAdmins] = useState<User[]>();

  useEffect(() => {
    fetch((c) => c.admins())
      .then(setAdmins)
      .catch();
  }, []);

  return (
    <RouteContainer>
      <h1>Admin Area</h1>
      <Card>
        <h2>Admins</h2>
        <AdminControls>
          <Input placeholder="User ID" />
          <Button>Add Admin</Button>
        </AdminControls>
        {admins && (
          <Flex direction="column" gap="0.5em">
            {admins.map((a) => (
              <Flex key={uid(a)} gap="1em" vCenter>
                <UserTile user={a} />
                <Button variant="red" disabled={a.is_owner}>
                  <IconCross />
                  Remove
                </Button>
              </Flex>
            ))}
          </Flex>
        )}
      </Card>
    </RouteContainer>
  );
};
