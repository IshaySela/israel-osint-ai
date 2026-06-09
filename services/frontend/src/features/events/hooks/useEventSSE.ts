import { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import { addEvent } from '../store/eventSlice';
import config from '../../../store/config';

export const useEventSSE = (url: string = `${config.BACKEND_URL}/events-stream`) => {
  const dispatch = useDispatch();

  useEffect(() => {
    console.info(`Connecting to SSE at ${url}...`);
    const eventSource = new EventSource(url);

    eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        console.debug('Received event via SSE:', data);
        dispatch(addEvent(data));
      } catch (err) {
        console.error('Failed to parse SSE message:', err);
      }
    };

    eventSource.onerror = (err) => {
      console.error('SSE connection error:', err);
      // EventSource will automatically retry by default
    };

    return () => {
      console.info('Closing SSE connection...');
      eventSource.close();
    };
  }, [url, dispatch]);
};
