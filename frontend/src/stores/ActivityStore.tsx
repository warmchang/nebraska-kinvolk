import Link from '@mui/material/Link';
import React from 'react';
import { Link as RouterLink } from 'react-router';
import _ from 'underscore';

import API from '../api/API';
import { Activity } from '../api/apiDataTypes';
import { toLocaleDateString } from '../i18n/dateTime';
import Store from './BaseStore';

interface ActivityEntryClass {
  type: string;
  appName: string;
  groupName: string | null;
  channelName: string | null;
  description: string | React.ReactElement;
}
interface ActivityEntrySeverity {
  type: string;
  className: string;
  icon: string;
}

class ActivityStore extends Store {
  activity: { [key: string]: Activity[] } | null;
  interval: null | number;
  constructor(noRefresh?: boolean) {
    super();
    this.activity = null;

    if (noRefresh) {
      this.interval = null;
    } else {
      this.getActivity();
      this.interval = window.setInterval(() => {
        this.getActivity();
      }, 60 * 1000);
    }
  }

  stopRefreshing() {
    if (this.interval) {
      window.clearInterval(this.interval);
    }
    this.interval = null;
  }

  setActivity(entries: Activity[] | null) {
    this.activity = entries === null ? null : this.sortActivityByDate(entries);
  }

  getCachedActivity() {
    return this.activity;
  }

  getActivity() {
    API.getActivity()
      .then(response => {
        this.setActivity(response.totalCount === 0 ? [] : response.activities);
        this.emitChange();
      })
      .catch(error => {
        if (error.status === 404) {
          this.setActivity(null);
          this.emitChange();
        }
      });
  }

  sortActivityByDate(entries: Activity[]) {
    const sortedEntries: {
      [key: string]: Activity[];
    } = {};

    entries.forEach(entry => {
      const date = toLocaleDateString(entry.created_ts, undefined, {
        day: 'numeric',
        weekday: 'short',
        month: 'short',
        year: 'numeric',
      });
      if (_.has(sortedEntries, date)) {
        sortedEntries[date].push(entry);
      } else {
        sortedEntries[date] = [entry];
      }
    });

    return sortedEntries;
  }

  makeActivityEntryClass(classID: number, entry: Activity): ActivityEntryClass {
    const instancePath = `/apps/${entry.app_id}/groups/${entry.group_id}/instances/${entry.instance_id}?period=1d`;

    const classTypes: {
      [key: string]: ActivityEntryClass;
    } = {
      1: {
        type: 'activityPackageNotFound',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description:
          "An update request could not be processed because the group's channel is not linked to any package",
      },
      2: {
        type: 'activityRolloutStarted',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description: 'Version ' + entry.version + ' roll out started',
      },
      3: {
        type: 'activityRolloutFinished',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description: 'Version ' + entry.version + ' successfully rolled out',
      },
      4: {
        type: 'activityRolloutFailed',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description:
          'There was an error rolling out version ' +
          entry.version +
          " as the first update attempt failed. Group's updates have been disabled",
      },
      5: {
        type: 'activityInstanceUpdateFailed',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description: (
          <React.Fragment>
            Instance{' '}
            <Link component={RouterLink} to={instancePath} underline="hover">
              {entry.instance_id}
            </Link>{' '}
            reported an error while processing update to version {entry.version}
          </React.Fragment>
        ),
      },
      6: {
        type: 'activityChannelPackageUpdated',
        appName: entry.application_name,
        groupName: entry.group_name,
        channelName: entry.channel_name,
        description:
          'Channel ' + entry.channel_name + ' is now pointing to version ' + entry.version,
      },
    };

    const classDetails = classID ? classTypes[classID] : classTypes[1];
    return classDetails;
  }

  makeActivityEntrySeverity(severityID: number) {
    const severityType: {
      [key: string]: ActivityEntrySeverity;
    } = {
      1: {
        type: 'activitySuccess',
        className: 'success',
        icon: 'fa-check',
      },
      2: {
        type: 'activityInfo',
        className: 'info',
        icon: 'fa-info',
      },
      3: {
        type: 'activityWarning',
        className: 'warning',
        icon: 'fa-exclamation',
      },
      4: {
        type: 'activityError',
        className: 'error',
        icon: 'fa-close',
      },
    };

    const severityInfo = severityID ? severityType[severityID] : severityType[1];
    return severityInfo;
  }
}

export default ActivityStore;
