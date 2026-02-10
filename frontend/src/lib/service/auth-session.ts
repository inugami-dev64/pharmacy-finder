const LOCAL_STORAGE_KEY = "SESSION";

class Session {
    token: string;
    expiry: number;

    constructor(token: string, expiry: number) {
        this.token = token;
        this.expiry = expiry;
    }
}

class JWTPayload {
    sub?: string;
    name?: string;
    admin?: boolean;
    iat?: number;
}

export class AuthenticationSession {
    /**
     * Decodes JWT payload and returns an instance to JWTPayload
     *
     * @param base64Payload specifies the base64 payload of the JWT
     * @returns an instance to JWTPayload if successful, undefined otherwise
     */
    private _decodePayload(base64Payload: string): JWTPayload|undefined {
        const base64Encoded = base64Payload.replace(/-/g, '+').replace(/_/g, '/');
        const padding = base64Payload.length % 4 === 0 ? '' : '='.repeat(4 - (base64Payload.length % 4));
        const base64WithPadding = base64Encoded + padding;

        return JSON.parse(atob(base64WithPadding))
    }

    /**
     * Returns true if provided session token belongs to an administrator
     * and false otherwise
     *
     * @param token specifies the token to admin verify
     * @returns true if provided session token belongs to an administrator, false otherwise
     */
    public isAdmin(token: string): boolean {
        const splitToken = token.split(".");
        if (splitToken.length === 3) {
            const payload = this._decodePayload(splitToken[1]);
            return payload?.admin || false;
        }

        return false;
    }

    /**
     * Gets currently authenticated user ID from provided token
     *
     * @param token specifies the session token to extract the user ID from
     * @returns a string value of the user ID or undefined if not found
     */
    public getUserId(token: string): string|undefined {
        const splitToken = token.split(".");
        if (splitToken.length === 3) {
            const payload = this._decodePayload(splitToken[1]);
            return payload?.sub;
        }

        return undefined;
    }

    /**
     * Sets the authentication token to LocalStorage
     *
     * @param token specifies the session token to persist
     * @param ttl specifies the token's TTL
     */
    public setSessionToken(token: string, ttl: number) {
        const splitToken = token.split(".");
        let expiry = new Date().getTime() + ttl;
        if (splitToken.length === 3) {
            const payload = this._decodePayload(splitToken[1]);
            if (payload != null && payload.iat != null)
                expiry = payload.iat + ttl;
        }

        const item = new Session(token, expiry);
        localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(item));
    }

    /**
     * Attempts to retrieve a session token from LocalStorage and return it
     *
     * @returns a retrieved JWT token or undefined if not found or expired
     */
    public getSessionToken(): string|undefined {
        const itemJson = localStorage.getItem(LOCAL_STORAGE_KEY);
        if (itemJson == null)
            return undefined;

        const item: Session = JSON.parse(itemJson);
        if (item == null || new Date().getTime() < item.expiry) {
            localStorage.removeItem(LOCAL_STORAGE_KEY);
            return undefined;
        }

        return item.token;
    }

    /**
     * Attempts to logout the user by deleting the JWT from LocalStorage
     */
    public logout() {
        localStorage.removeItem(LOCAL_STORAGE_KEY);
    }
}

export const authenticationSession = new AuthenticationSession();