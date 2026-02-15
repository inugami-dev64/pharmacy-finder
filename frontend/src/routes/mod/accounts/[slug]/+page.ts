import { authenticationSession } from '$lib/service/auth-session';
import { LocalizedBackendError } from '$lib/service/data/error.js';
import { HealthCheckResult } from '$lib/service/data/health';
import { UserProfile } from '$lib/service/data/users.js';
import { redirect } from '@sveltejs/kit';

export async function load({ params }) {
    let profile: UserProfile | undefined = undefined;
    let error: LocalizedBackendError | undefined = undefined;

    // Check if the page should be redirected
    const health = await HealthCheckResult.readHealthCheck(fetch);
    if (!health?.initialized)
        redirect(307, "/mod/register");
    else if (authenticationSession.getSessionToken() == null)
        redirect(307, "/mod/login");

    // If the slug is "me", then query data about the current user
    const userId: string = params.slug;
    if (userId === "me") {
        try {
            profile = await UserProfile.getAuthenticatedUser();
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                error = e;
        }
    } else {
        try {
            profile = await UserProfile.getUserById(userId);
        } catch (e) {
            if (e instanceof LocalizedBackendError)
                error = e;
        }
    }

    return {
        profile: profile,
        error: error,
        isMe: userId === "me"
    }
}