import { gql } from '@apollo/client';

export const GET_LATEST_EVENTS = gql`
  query GetLatestEvents {
    latestEvents {
      raw_message
      summary
      timestamp
      locations {
        name
        lat
        lon
      }
    }
  }
`;

export const GET_EVENTS = gql`
  query GetEvents($fromHoursAgo: Int!, $toHoursAgo: Int!) {
    events(fromHoursAgo: $fromHoursAgo, toHoursAgo: $toHoursAgo) {
      raw_message
      summary
      timestamp
      locations {
        name
        lat
        lon
      }
    }
  }
`;
