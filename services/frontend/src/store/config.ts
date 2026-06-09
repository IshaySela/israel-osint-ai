
const getBackendUrl = () => {
    let url = import.meta.env.VITE_BACKEND_URL;
    
    if (url === undefined && import.meta.env.MODE === 'development') {
        url = 'http://localhost:5000'; // Default to localhost in development if not set
    }
    else if(url === undefined) {
        throw new Error("VITE_BACKEND_URL is not defined in environment variables");
    }

    return url;
}

const config = {
    BACKEND_URL: getBackendUrl()
} as const

export default config;