import { HealthCheckResult } from "$lib/service/data/health";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { authenticationSession } from "$lib/service/auth-session";

export const load: PageLoad = async ({ fetch }) => {
    const health = await HealthCheckResult.readHealthCheck(fetch);
    if (!health?.initialized)
        redirect(307, "/mod/register");
    else if (authenticationSession.getSessionToken() == null)
        redirect(307, "/mod/login");

    return {
        health: health
    };
}