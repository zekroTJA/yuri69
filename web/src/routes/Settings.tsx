import { useEffect, useReducer, useRef, useState } from 'react';
import { uid } from 'react-uid';
import styled, { useTheme } from 'styled-components';
import { OTAToken, Sound } from '../api';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { Flex } from '../components/Flex';
import { GuildTile } from '../components/GuildTile';
import { HoverToShow } from '../components/HoverToShow';
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

const SettingsRouteContainer = styled(RouteContainer)``;

const OTAContainer = styled.div`
  display: flex;
  gap: 1em;
  margin-bottom: 1.5em;
`;

const timerReducer = (state: number, action: { type: 'decrease' | 'set'; payload?: number }) => {
  switch (action.type) {
    case 'decrease':
      if (state - 1 < 0) return 0;
      return state - 1;
    case 'set':
      return action.payload ?? state;
    default:
      return state;
  }
};

export const SettingsRoute: React.FC<Props> = ({}) => {
  const [connected, guild, filters] = useStore((s) => [s.connected, s.guild, s.filters]);
  const { sounds } = useSounds();
  const theme = useTheme();
  const fetch = useApi();
  const { show } = useSnackBar();
  const [tagsInclude, setTagsInclude] = useState<string[]>([]);
  const [tagsExclude, setTagsExclude] = useState<string[]>([]);
  const [fastTrigger, setFastTrigger] = useState('');
  const [otaToken, setOtaToken] = useState<OTAToken>();
  const [deadline, dispatchDeadline] = useReducer(timerReducer, 0);
  const timerRef = useRef<ReturnType<typeof setTimeout>>();
  const intervalRef = useRef<ReturnType<typeof setInterval>>();
  const [apiKey, setApiKey] = useState('');

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

  const _fetchOtaToken = async () => {
    try {
      const token = await fetch((c) => c.getOTAToken());
      setOtaToken(token);
      const refreshIn = new Date(token.deadline).getTime() - Date.now();
      clearTimeout(timerRef.current);
      if (refreshIn - 1000 > 0) {
        dispatchDeadline({ type: 'set', payload: Math.floor(refreshIn / 1000 - 1) });
        timerRef.current = setTimeout(() => _fetchOtaToken(), refreshIn - 1000);
      }
    } catch {}
  };

  const _generateApiKey = () => {
    fetch((c) => c.generateApiKey())
      .then((res) => {
        setApiKey(res.api_key);
        show('API key has been generated.', 'success');
      })
      .catch();
  };

  const _removeApiKey = () => {
    fetch((c) => c.removeApiKey())
      .then(() => {
        setApiKey('');
        show('API key has been removed.', 'success');
      })
      .catch();
  };

  const _copyApikeyToClipboard = () => {
    navigator.clipboard
      .writeText(apiKey)
      .then(() => show('API key has been copied to you clipboard.', 'success'))
      .catch((err) => show(`Failed copying API key to clipboard: ${err}`, 'error'));
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

    fetch((c) => c.apiKey(), 404)
      .then((res) => setApiKey(res.api_key))
      .catch();

    _fetchOtaToken();

    clearInterval(intervalRef.current);
    intervalRef.current = setInterval(() => dispatchDeadline({ type: 'decrease' }), 1000);

    return () => {
      clearInterval(intervalRef.current);
      clearTimeout(timerRef.current);
    };
  }, []);

  const _soundOptions: Sound[] = [
    { uid: '', display_name: '< unset >' },
    { uid: 'random', display_name: '< random >' },
    ...sounds,
  ];

  return (
    <SettingsRouteContainer>
      <h1>Settings</h1>
      <Card margin="0 0 1em 0">
        <h2>Authorization</h2>
        {otaToken && (
          <OTAContainer>
            <img src={otaToken.qrcode_data} />
            <div>
              <p>
                You can scan this QR code with your mobile device and use the Yuri web interface
                from there witout a login required!
              </p>
              <p>Code resets in {deadline} seconds.</p>
            </div>
          </OTAContainer>
        )}
        <h2>API Key</h2>
        {(apiKey && <HoverToShow>{apiKey}</HoverToShow>) || <p>No API key generated.</p>}
        <Flex gap="1em">
          <Button onClick={_generateApiKey}>{(apiKey && 'Regenerate') || 'Generate'}</Button>
          {apiKey && (
            <>
              <Button onClick={_copyApikeyToClipboard}>Copy to Clipboard</Button>
              <Button variant="red" onClick={_removeApiKey}>
                Delete
              </Button>
            </>
          )}
        </Flex>
      </Card>

      <SplitContainer margin="0 0 1.5em 0">
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
