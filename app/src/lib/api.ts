import {SERVER_DOMAIN} from "../env";


export const APIClient = {
    fetch: (url: string): Promise<Response> => {
        const headers: Record<string, string> = {}
        let companionUrl = localStorage.getItem('x-companion-base-url');
        let companionApiKey = localStorage.getItem('x-companion-api-key');
        if (companionUrl != undefined && companionUrl != "") {
            headers["x-companion-base-url"] = companionUrl || ""
        }
        if (companionApiKey != undefined && companionApiKey != "") {
            headers["x-companion-api-key"] = companionApiKey || ""
        }
      
        // add the headers to the options
        let uri = `${SERVER_DOMAIN}` + url;
        return fetch(uri, {headers: headers, credentials: 'same-origin'})
    }
};
