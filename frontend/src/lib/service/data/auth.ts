import { HttpError } from "$lib/http-error";
import { LocalizedBackendError } from "./error";
import type { FetchFunc } from "./fetch";
import { UserProfile } from "./users";

/**
 * AuthenticatedUser class represents the successful authentication
 * response from the server
 */
export class AuthenticatedUser extends UserProfile {
    session?: {
        token: string;
        validFor?: number;
    }
}

export class LoginForm {
    username?: string;
    password?: string;

    /**
     * Login to moderation platform using credentials in the LoginForm instance
     *
     * @param fetch specifies the fetch function to use for making a request
     * @returns a Promise to AuthenticatedUser
     */
    public async login(fetch: FetchFunc): Promise<AuthenticatedUser> {
        return await fetch(`/api/v1/mod/auth/login`, {
            method: "POST",
            body: JSON.stringify(this),
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("User with provided"))
                        throw new LocalizedBackendError("mod.login.errors.userNotFound");
                    else if (err.msg?.startsWith("Invalid password"))
                        throw new LocalizedBackendError("mod.login.errors.invalidPassword");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                let data: AuthenticatedUser = await res.json();
                return data;
            })
    }
}

export class RegistrationForm {
    username: string | undefined;
    email: string | undefined;
    password: string | undefined;
    firstName: string | undefined;
    lastName: string | undefined;

    /**
     * Register a new administrator account using credentials in the
     * RegistrationForm instance and authenticate with it
     *
     * @param fetch specifies the fetch function to use for making requests
     * @returns a Promise to AuthenticatedUser object
     */
    public async registerAdministrator(fetch: FetchFunc): Promise<AuthenticatedUser> {
        return await fetch(`/api/v1/mod/auth/register`, {
            method: "POST",
            body: JSON.stringify(this),
            headers: {
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    let err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Username or email"))
                        throw new LocalizedBackendError("mod.register.errors.usernameOrEmailExists");
                    else if (err.msg?.startsWith("Password is too long"))
                        throw new LocalizedBackendError("mod.register.errors.passwordTooLong");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                let data: AuthenticatedUser = await res.json();
                return data;
            })
    }
}