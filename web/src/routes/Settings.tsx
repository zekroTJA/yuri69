import { useEffect, useState } from 'react';
import { uid } from 'react-uid';
import styled, { useTheme } from 'styled-components';
import { Sound } from '../api';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { GuildTile } from '../components/GuildTile';
import { InfoPanel } from '../components/InfoPanel';
import { RouteContainer } from '../components/RouteContainer';
import { Select } from '../components/Select';
import { Smol } from '../components/Smol';
import { SplitContainer } from '../components/SplitContainer';
import { TagsInput } from '../components/TagsInput';
import { useApi } from '../hooks/useApi';
import { useSnackBar } from '../hooks/useSnackBar';
import { useSounds } from '../hooks/useSounds';
import { useStore } from '../store';

type Props = {};

const Controls = styled.div`
  margin-top: 1.5em;

  > * {
    display: block;
    width: 100%;
  }

  > label {
    margin-top: 1em;
    margin-bottom: 0.4em;
  }

  > ${Button} {
    margin-top: 1em;
  }
`;

const SettingsRouteContainer = styled(RouteContainer)`
  /* > div {
    display: flex;
    gap: 1em;
    > ${Card} {
      width: 100%;
    }
  }

  @media screen and (max-width: 70em) {
    > div {
      flex-direction: column;
    }
  } */
`;

export const SettingsRoute: React.FC<Props> = ({}) => {
  const [connected, guild, filters] = useStore((s) => [s.connected, s.guild, s.filters]);
  const { sounds } = useSounds();
  const theme = useTheme();
  const fetch = useApi();
  const { show } = useSnackBar();
  const [tagsInclude, setTagsInclude] = useState<string[]>([]);
  const [tagsExclude, setTagsExclude] = useState<string[]>([]);
  const [fastTrigger, setFastTrigger] = useState('');

  const _applyGuild = async () => {
    try {
      await fetch((c) => c.guildsSetFilters({ include: tagsInclude, exclude: tagsExclude }));
      show('Guild preferences were sucessfully applied.', 'success');
    } catch {}
  };

  const _applyUser = async () => {
    try {
      await fetch((c) => c.usersSetFasttrigger(fastTrigger));
      show('Personal preferences were sucessfully applied.', 'success');
    } catch {}
  };

  useEffect(() => {
    if (filters?.include) setTagsInclude(filters.include);
    if (filters?.exclude) setTagsExclude(filters.exclude);
  }, [filters]);

  useEffect(() => {
    if (!connected) {
      setTagsExclude([]);
      setTagsInclude([]);
    }
  }, [connected]);

  useEffect(() => {
    fetch((c) => c.usersGetFasttrigger())
      .then((res) => setFastTrigger(res.fast_trigger))
      .catch();
  }, []);

  const _soundOptions: Sound[] = [
    { uid: '', display_name: '< unset >' },
    { uid: 'random', display_name: '< random >' },
    ...sounds,
  ];

  return (
    <SettingsRouteContainer>
      <h1>Settings</h1>
      <SplitContainer>
        <Card>
          <h2>Guild Preferences</h2>
          {(connected && guild && (
            <InfoPanel>
              <Smol>You are connected to</Smol>
              <GuildTile guild={guild} />
            </InfoPanel>
          )) || (
            <InfoPanel color={theme.orange}>
              <Smol>You are not connected to any guild.</Smol>
            </InfoPanel>
          )}
          <Controls>
            <label htmlFor="include-filters">Include Filters</label>
            <TagsInput
              disabled={!connected}
              id="include-filters"
              tags={tagsInclude}
              onTagsChange={setTagsInclude}
            />
            <label htmlFor="exclude-filters">Exclude Filters</label>
            <TagsInput
              disabled={!connected}
              id="exclude-filters"
              tags={tagsExclude}
              onTagsChange={setTagsExclude}
            />
            <Button disabled={!connected} variant="green" onClick={_applyGuild}>
              Apply
            </Button>
          </Controls>
        </Card>

        <Card>
          <h2>Personal Preferences</h2>
          <Controls>
            <label htmlFor="fast-trigger">Fast Trigger</label>
            <Select value={fastTrigger} onChange={(e) => setFastTrigger(e.currentTarget.value)}>
              {_soundOptions.map((s) => (
                <option key={uid(s)} value={s.uid}>
                  {s.display_name || s.uid}
                </option>
              ))}
            </Select>
            <Button variant="green" onClick={_applyUser}>
              Apply
            </Button>
          </Controls>
        </Card>
      </SplitContainer>
    </SettingsRouteContainer>
  );
};
