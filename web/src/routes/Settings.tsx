import { useEffect, useReducer, useRef, useState } from 'react';
import { uid } from 'react-uid';
import styled, { useTheme } from 'styled-components';
import { ImportSoundsResult, OTAToken, Sound, TwitchState } from '../api';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { FileDrop } from '../components/FileDrop';
import { Flex } from '../components/Flex';
import { GuildTile } from '../components/GuildTile';
import { HoverToShow } from '../components/HoverToShow';
import { InfoPanel } from '../components/InfoPanel';
import { Input } from '../components/Input';
import { RouteContainer } from '../components/RouteContainer';
import { Select } from '../components/Select';
import { Smol } from '../components/Smol';
import { Spinner } from '../components/Spinner';
import { SplitContainer } from '../components/SplitContainer';
import { TagsInput } from '../components/TagsInput';
import { useApi } from '../hooks/useApi';
import { useLocalStorage } from '../hooks/useLocalStorage';
import { useSnackBar } from '../hooks/useSnackBar';
import { useSounds } from '../hooks/useSounds';
import { ApiClientInstance } from '../instances';
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

const ImportContainer = styled.div`
  margin-top: 2em;
  display: flex;
  flex-direction: column;
  gap: 1.5em;

  > h4 {
    margin: 0;
  }

  > div > button {
    margin-top: 1em;
  }
`;

const ErrContainer = styled.div`
  color: ${(p) => p.theme.red};
`;

const ControlButtons = styled.div`
  display: flex;
  gap: 1em;
  margin-top: 1em;

  > * {
    width: 100%;
  }
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

const twitchReducer = (
  state: TwitchState,
  action:
    | {
        type: 'set_username' | 'set_prefix';
        payload: string;
      }
    | {
        type: 'set_ratelimit_burst' | 'set_ratelimit_reset';
        payload: number;
      }
    | {
        type: 'set_filters_include' | 'set_filters_exclude' | 'set_blocked';
        payload: string[];
      }
    | {
        type: 'set_connected';
        payload: boolean;
      }
    | {
        type: 'set_state';
        payload: TwitchState;
      },
) => {
  switch (action.type) {
    case 'set_username':
      return { ...state, twitch_user_name: action.payload };
    case 'set_prefix':
      return { ...state, prefix: action.payload };
    case 'set_ratelimit_burst':
      return { ...state, ratelimit: { ...state.ratelimit, burst: action.payload } };
    case 'set_ratelimit_reset':
      return { ...state, ratelimit: { ...state.ratelimit, reset_seconds: action.payload } };
    case 'set_filters_include':
      return { ...state, filters: { ...state.filters, include: action.payload } };
    case 'set_filters_exclude':
      return { ...state, filters: { ...state.filters, exclude: action.payload } };
    case 'set_blocked':
      return { ...state, blocklist: action.payload };
    case 'set_connected':
      return { ...state, connected: action.payload };
    case 'set_state':
      return action.payload;
    default:
      return state;
  }
};

export const SettingsRoute: React.FC<Props> = ({}) => {
  const [connected, guild, filters, isAdmin] = useStore((s) => [
    s.connected,
    s.guild,
    s.filters,
    s.isAdmin,
  ]);
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
  const [downloadLock, setDownloadLock] = useLocalStorage<number>('yuri_downloadlock');
  const [importFile, setImportFile] = useState<File>();
  const [importProcessing, setImportProcessing] = useState(false);
  const [importResult, setImportResult] = useState<ImportSoundsResult>();
  const [twitchState, dispatchTwitchState] = useReducer(twitchReducer, {
    ratelimit: {},
    filters: {},
  } as TwitchState);

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

  const _onDownloadAll = () => {
    const url = ApiClientInstance.allSoundsDownloadUrl();
    const a = document.createElement('a');
    a.href = url;
    a.target = '_blank';
    a.click();
    setDownloadLock(Date.now() + 5 * 60 * 1000);
  };

  const _onImport = () => {
    if (!importFile) return;

    setImportFile(undefined);
    setImportProcessing(true);
    setImportResult(undefined);
    fetch((c) => c.soundsImport(importFile))
      .then(setImportResult)
      .catch()
      .finally(() => setImportProcessing(false));
  };

  const _onTwitchSave = () => {
    fetch((c) => c.setTwitchSettings(twitchState))
      .then(() => show('Twitch settings saved.', 'success'))
      .catch();
  };

  const _onTwitchJoin = () => {
    fetch((c) => c.joinTwitch(twitchState))
      .then(() => {
        dispatchTwitchState({ type: 'set_connected', payload: true });
        show(`Joined twitch channel "${twitchState.twitch_user_name}".`, 'success');
      })
      .catch();
  };

  const _onTwitchLeave = () => {
    fetch((c) => c.leaveTwitch())
      .then(() => {
        dispatchTwitchState({ type: 'set_connected', payload: false });
        show(`Left twitch channel "${twitchState.twitch_user_name}".`, 'success');
      })
      .catch();
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

    fetch((c) => c.twitchState())
      .then((res) => dispatchTwitchState({ type: 'set_state', payload: res }))
      .catch();

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

  const _downloadDisabled = (downloadLock ?? 0) > Date.now();

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

        <Card>
          <div>
            <h2>Sounds</h2>
            {_downloadDisabled && (
              <Smol>
                Sounds download can only be requested every 5 minutes.
                <br />
                <br />
              </Smol>
            )}
            <Button onClick={_onDownloadAll} disabled={_downloadDisabled}>
              Download all Sounds
            </Button>
          </div>
          {isAdmin && (
            <ImportContainer>
              <h4>Import Sounds</h4>
              {importProcessing || (
                <div>
                  <FileDrop file={importFile} onFileInput={setImportFile} />
                  <Button disabled={!importFile} onClick={_onImport}>
                    Import
                  </Button>
                </div>
              )}
              {importProcessing && (
                <Flex gap="1em" vCenter>
                  <Spinner />
                  <span>Importing ...</span>
                </Flex>
              )}
              {importResult && (
                <div>
                  <p>Successful Imports: {importResult.successful?.length ?? 0}</p>
                  <ErrContainer>
                    {importResult.failed?.map((err) => (
                      <p>
                        <strong>{err.uid}</strong>
                        <br />
                        <span>{err.error}</span>
                      </p>
                    ))}
                  </ErrContainer>
                </div>
              )}
            </ImportContainer>
          )}
        </Card>

        {twitchState.capable && (
          <Card>
            <h2>Twitch</h2>
            <Controls>
              <label>Twitch Channel to Join</label>
              <Input
                disabled={twitchState.connected}
                placeholder="zekrotja"
                value={twitchState.twitch_user_name}
                onInput={(e) =>
                  dispatchTwitchState({ type: 'set_username', payload: e.currentTarget.value })
                }
              />
              <label>Chat Command Prefix</label>
              <Input
                placeholder="!yuri"
                value={twitchState.prefix}
                onInput={(e) =>
                  dispatchTwitchState({ type: 'set_prefix', payload: e.currentTarget.value })
                }
              />
              <label>Rate Limit Burst Rate</label>
              <Input
                placeholder="!yuri"
                value={twitchState.ratelimit.burst}
                type="number"
                min="1"
                onInput={(e) =>
                  dispatchTwitchState({
                    type: 'set_ratelimit_burst',
                    payload: parseInt(e.currentTarget.value),
                  })
                }
              />
              <label>
                Rate Limit Reset <Smol>(in seconds)</Smol>
              </label>
              <Input
                placeholder="!yuri"
                value={twitchState.ratelimit.reset_seconds}
                type="number"
                min="1"
                onInput={(e) =>
                  dispatchTwitchState({
                    type: 'set_ratelimit_reset',
                    payload: parseInt(e.currentTarget.value),
                  })
                }
              />
              <label>Include Filter Tags</label>
              <TagsInput
                tags={twitchState.filters.include}
                onTagsChange={(payload) =>
                  dispatchTwitchState({ type: 'set_filters_include', payload })
                }
              />
              <label>Exclude Filter Tags</label>
              <TagsInput
                tags={twitchState.filters.exclude}
                onTagsChange={(payload) =>
                  dispatchTwitchState({ type: 'set_filters_exclude', payload })
                }
              />
              <label>User Blocklist</label>
              <TagsInput
                tags={twitchState.blocklist}
                onTagsChange={(payload) => dispatchTwitchState({ type: 'set_blocked', payload })}
              />
              <ControlButtons>
                <Button
                  variant={twitchState.connected ? 'orange' : 'blue'}
                  onClick={twitchState.connected ? _onTwitchLeave : _onTwitchJoin}>
                  {twitchState.connected ? 'Disconnect' : 'Connect'}
                </Button>
                <Button variant="green" onClick={_onTwitchSave}>
                  Save
                </Button>
              </ControlButtons>
            </Controls>
          </Card>
        )}
      </SplitContainer>
    </SettingsRouteContainer>
  );
};
