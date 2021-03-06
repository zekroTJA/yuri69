import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { uid } from 'react-uid';
import styled, { css } from 'styled-components';
import { PlaybackLogEntry, PlaybackStats, Sound, StateStats } from '../api';
import { Card } from '../components/Card';
import { LinkButton } from '../components/LinkButton';
import { RouteContainer } from '../components/RouteContainer';
import { SplitContainer } from '../components/SplitContainer';
import { useApi } from '../hooks/useApi';
import { formatDate } from '../util/date';

type Props = {};

const Table = styled.table<{ fw?: boolean; clickable?: boolean }>`
  text-align: left;
  cursor: default;
  border-collapse: collapse;
  ${(p) => p.fw && 'width: 100%;'}

  th {
    text-transform: uppercase;
    opacity: 0.8;
    font-weight: 300;
  }

  tr {
    transition: all 0.2s ease;

    ${(p) =>
      p.clickable &&
      css`
        cursor: pointer;
      `}

    &:hover {
      ${(p) =>
        p.clickable &&
        css`
          background-color: ${p.theme.background3};
        `}
    }

    > * {
      padding: 0.5em 1em 0.5em 0.5em;
      &:last-child {
        padding-right: 0.5em;
      }
    }

    &:first-child {
      background-color: transparent !important;
      cursor: default !important;
    }
  }
`;

export const StatsRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const [state, setState] = useState<StateStats>();
  const [log, setLog] = useState<PlaybackLogEntry[]>();
  const [limitLog, setLimitLog] = useState(10);
  const [counts, setCounts] = useState<PlaybackStats[]>();
  const [limitCounts, setLimitCounts] = useState(5);
  const nav = useNavigate();

  const _edit = (e: { ident: string }) => nav('/sounds/' + e.ident);

  useEffect(() => {
    fetch((c) => c.statsState())
      .then((res) => setState(res))
      .catch();
    fetch((c) => c.statsCount())
      .then((res) => setCounts(res))
      .catch();
  }, []);

  useEffect(() => {
    fetch((c) => c.statsLog('', '', '', limitLog))
      .then((res) => setLog(res))
      .catch();
  }, [limitLog]);

  return (
    <RouteContainer>
      <h1>Stats</h1>
      {state && (
        <Card margin="0 0 1em 0">
          <Table>
            <tbody>
              <tr>
                <th>Number of Sounds</th>
                <td>{state.n_sounds}</td>
              </tr>
              <tr>
                <th>Number of Plays</th>
                <td>{state.n_plays}</td>
              </tr>
            </tbody>
          </Table>
        </Card>
      )}
      <SplitContainer margin="0 0 1.5em 0">
        {log && (
          <Card>
            <Table fw clickable>
              <tbody>
                <tr>
                  <th>Sound UID</th>
                  <th>Count</th>
                </tr>
                {log.map((c) => (
                  <tr key={uid(c)} onClick={() => _edit(c)}>
                    <td>{c.ident}</td>
                    <td>{formatDate(c.timestamp)}</td>
                  </tr>
                ))}
                <tr>
                  <LinkButton onClick={() => setLimitLog(limitLog === 10 ? 100 : 10)}>
                    Show {limitLog === 10 ? 'more' : 'less'} ...
                  </LinkButton>
                </tr>
              </tbody>
            </Table>
          </Card>
        )}
        {counts && (
          <Card>
            <Table fw clickable>
              <tbody>
                <tr>
                  <th>Sound UID</th>
                  <th>Count</th>
                </tr>
                {counts.slice(0, limitCounts).map((c) => (
                  <tr key={uid(c)} onClick={() => _edit(c)}>
                    <td>{c.ident}</td>
                    <td>{c.count}</td>
                  </tr>
                ))}
                <tr>
                  <LinkButton onClick={() => setLimitCounts(limitCounts === 5 ? 50 : 5)}>
                    Show {limitCounts === 5 ? 'more' : 'less'} ...
                  </LinkButton>
                </tr>
              </tbody>
            </Table>
          </Card>
        )}
      </SplitContainer>
    </RouteContainer>
  );
};

