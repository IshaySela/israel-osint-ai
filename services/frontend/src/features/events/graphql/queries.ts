import { gql } from '@apollo/client';

export const GET_LATEST_EVENTS = gql`
  query GetLatestEvents {
    latestEvents {
      raw_message
      summary
      timestamp_epoch
      locations {
        name
        lat
        lon
      }
    }
  }
`;

export const GET_EVENTS = gql`
  query GetEvents($fromMinutesAgo: Int!, $toMinutesAgo: Int!) {
    events(fromMinutesAgo: $fromMinutesAgo, toMinutesAgo: $toMinutesAgo) {
      raw_message
      summary
      timestamp_epoch
      locations {
        name
        lat
        lon
      }
    }
  }
`;
