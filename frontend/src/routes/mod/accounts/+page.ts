import { HealthCheckResult } from "$lib/service/data/health";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "../$types";
import { authenticationSession } from "$lib/service/auth-session";
import { UserProfile } from "$lib/service/data/users";

export const load: PageLoad = async ({ fetch }) => {
    const health = await HealthCheckResult.readHealthCheck(fetch);
    if (!health?.initialized)
        redirect(307, "/mod/register");
    else if (authenticationSession.getSessionToken() == null)
        redirect(307, "/mod/login");
    else if (!authenticationSession.isAdmin(authenticationSession.getSessionToken() || ""))
        redirect(307, "/mod");

    const users = await UserProfile.getAllUsers();

    return {
        health: health,
        accounts: users
    }
}