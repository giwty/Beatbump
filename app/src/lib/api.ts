import {SERVER_DOMAIN} from "../env";


export const APIClient = {
    fetch: (url: string): Promise<Response> => {
        const headers: Record<string, string> = {}
        
        // add the headers to the options
        let uri = `${SERVER_DOMAIN}` + url;
        return fetch(uri, {headers: headers, credentials: 'same-origin'})
    }
};
