import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client';
import config from './store/config';

const client = new ApolloClient({
  link: new HttpLink({ uri: `${config.BACKEND_URL}/graphql` }),
  cache: new InMemoryCache(),
});

export default client;