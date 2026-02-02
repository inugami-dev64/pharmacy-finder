import type { HttpError } from "$lib/http-error";
import { authenticationSession } from "../auth-session";
import type { RegistrationForm } from "./auth";
import { LocalizedBackendError } from "./error";

class CurrentUserModificationForm {
    email?: string;
    firstName?: string;
    lastName?: string;
    password?: string;
    administrator?: boolean;
    currentPassword?: string;
}

export class UserProfile {
    id?: string;
    username?: string;
    email?: string;
    firstName?: string;
    lastName?: string;
    registrationTs?: number;
    lastLoginTs?: number;
    administrator?: boolean;

    /**
     * Queries for all moderator users for pharmacy finder
     * NOTE: Requires administrator privileges
     *
     * @returns a Promise to an array of UserProfile objects
     */
    protected static async getAllUsers(): Promise<UserProfile[]> {
        const token = authenticationSession.getSessionToken();
        return await fetch(`/api/v1/mod/users`, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + token
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                return await res.json();
            });
    }

    /**
     * Queries user profile data for currently authenticated user
     * NOTE: Requires authentication
     *
     * @returns a Promise to a UserProfile object
     */
    protected static async getAuthenticatedUser(): Promise<UserProfile> {
        const token = authenticationSession.getSessionToken();
        return await fetch(`/api/v1/mod/users/me`, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + token
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                return await res.json();
            })
    }

    /**
     * Creates a new moderator user from provided RegistrationForm
     * and returns a UserProfile instance for created user.
     * NOTE: Requires administrator privileges
     *
     * @param form specifies the RegistrationForm object to use
     * @returns a Promise to UserProfile object which represents the user just created
     */
    protected static async createNewModeratorUser(form: RegistrationForm): Promise<UserProfile> {
        const token = authenticationSession.getSessionToken();
        return await fetch(`/api/v1/mod/users`, {
            method: "POST",
            body: JSON.stringify(form),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 201) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Username or email"))
                        throw new LocalizedBackendError("mod.register.errors.usernameOrEmailExists");
                    else if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                return await res.json();
            })
    }

    /**
     * Updates the currently authenticated user.
     * NOTE: Requires authentication
     *
     * @param currentPassword specifies the password for current user
     * @param newPassword optionally specifies the new password for the user
     */
    protected async updateCurrentlyAuthenticatedUser(currentPassword: string, newPassword?: string): Promise<void> {
        const form = new CurrentUserModificationForm();
        form.email = this.email;
        form.firstName = this.firstName;
        form.lastName = this.lastName;
        form.password = newPassword;
        form.administrator = this.administrator;
        form.currentPassword = currentPassword;

        const token = authenticationSession.getSessionToken();
        await fetch(`/api/v1/mod/users/me`, {
            method: "PATCH",
            body: JSON.stringify(form),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Invalid password"))
                        throw new LocalizedBackendError("mod.login.errors.invalidPassword");
                    else if (err.msg?.startsWith("Password is too long"))
                        throw new LocalizedBackendError("mod.register.errors.passwordTooLong");
                    else if (err.msg?.startsWith("Email address is already"))
                        throw new LocalizedBackendError("mod.users.errors.emailInUse");
                    else if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: UserProfile = await res.json();
                Object.assign(this, data);
            })
    }

    /**
     * Update any pharmacy finder moderator user.
     * NOTE: Requires administrator privileges
     *
     * @param newPassword optionally specifies the new password for given user
     */
    protected async updateUser(newPassword?: string): Promise<void> {
        const form = new CurrentUserModificationForm();
        form.email = this.email;
        form.firstName = this.firstName;
        form.lastName = this.lastName;
        form.password = newPassword;
        form.administrator = this.administrator;

        const token = authenticationSession.getSessionToken();
        await fetch(`/api/v1/mod/users/${this.id}`, {
            method: "PATCH",
            body: JSON.stringify(form),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Malformed ID"))
                        throw new LocalizedBackendError("mod.users.errors.invalidUserId");
                    else if (err.msg?.startsWith("User with ID"))
                        throw new LocalizedBackendError("mod.users.errors.userNotFound");
                    else if (err.msg?.startsWith("Password is too long"))
                        throw new LocalizedBackendError("mod.register.errors.passwordTooLong");
                    else if (err.msg?.startsWith("Email address is already"))
                        throw new LocalizedBackendError("mod.users.errors.emailInUse");
                    else if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: UserProfile = await res.json();
                Object.assign(this, data);
            })
    }

    /**
     * Deletes currently authenticated user.
     * NOTE: Requires authentication
     *
     * @param currentPassword specifies the current password of the authenticated user
     */
    protected async deleteCurrentlyAuthenticatedUser(currentPassword: string): Promise<void> {
        const form = new CurrentUserModificationForm();
        form.currentPassword = currentPassword;

        const token = authenticationSession.getSessionToken();
        await fetch(`/api/v1/mod/users/me`, {
            method: "DELETE",
            body: JSON.stringify(form),
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Invalid Password"))
                        throw new LocalizedBackendError("mod.login.errors.invalidPassword");
                    else if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: UserProfile = await res.json();
                Object.assign(this, data);
            })
    }

    /**
     * Deletes someone else's user account from the system.
     * NOTE: Requires administrator privileges
     */
    protected async deleteUser(): Promise<void> {
        const token = authenticationSession.getSessionToken();
        await fetch(`/api/v1/mod/users/${this.id}`, {
            method: "DELETE",
            headers: {
                "Authorization": "Bearer " + token
            }
        })
            .then(async res => {
                if (res.status != 200) {
                    const err: HttpError = await res.json();
                    console.error(err);
                    if (err.msg?.startsWith("Malformed ID"))
                        throw new LocalizedBackendError("mod.users.errors.invalidUserId");
                    else if (err.msg?.startsWith("User with ID"))
                        throw new LocalizedBackendError("mod.users.errors.userNotFound");
                    else if (err.msg?.startsWith("Forbidden"))
                        throw new LocalizedBackendError("universalErrors.forbidden");
                    else throw new LocalizedBackendError("universalErrors.internalServerError");
                }

                const data: UserProfile = await res.json();
                Object.assign(this, data);
            })
    }
}