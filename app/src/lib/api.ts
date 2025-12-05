import { SERVER_DOMAIN } from "../env";

export const APIClient = {
    fetch: (url: string): Promise<any> => {
        const headers: Record<string, string> = {}

        // add the headers to the options
        let uri = `${SERVER_DOMAIN}` + url;
        return fetch(uri, { headers: headers, credentials: 'same-origin' })
    },
    post: (url: string, body?: any): Promise<any> => {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json'
        }

        let uri = `${SERVER_DOMAIN}` + url;
        return fetch(uri, {
            method: 'POST',
            headers: headers,
            credentials: 'same-origin',
            body: JSON.stringify(body)
        })
    },
    del: (url: string): Promise<any> => {
        const headers: Record<string, string> = {}
        let uri = `${SERVER_DOMAIN}` + url;
        return fetch(uri, {
            method: 'DELETE',
            headers: headers,
            credentials: 'same-origin'
        })
    }
};
