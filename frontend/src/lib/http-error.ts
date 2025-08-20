/**
 * HttpError struct represents the unified error struct that backend responds with
 */
export class HttpError {
    code: number | undefined;
    ts: number | undefined;
    msg: string | undefined;
}