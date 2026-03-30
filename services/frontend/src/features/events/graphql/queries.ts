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
