/**
 * LocalizedBackendError class can be thrown from fetch requests
 * to signal that a backend error has occurred
 */
export class LocalizedBackendError {
    msg: string;

    constructor(msg: string) {
        this.msg = msg;
    }
}