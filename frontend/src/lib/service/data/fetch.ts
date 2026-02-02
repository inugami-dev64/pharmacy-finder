export type FetchFunc = {
    (input: RequestInfo | URL, init?: RequestInit): Promise<Response>;
    (input: string | URL | Request, init?: RequestInit): Promise<Response>;
};