import { HealthCheckResult } from "$lib/service/data/health";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "../$types";
import { authenticationSession } from "$lib/service/auth-session";

export const load: PageLoad = async ({ fetch }) => {
    const health = await HealthCheckResult.readHealthCheck(fetch);
    const token = authenticationSession.getSessionToken();
    if (health?.initialized && token != null)
        redirect(307, "/mod");
    else if (!health?.initialized)
        redirect(307, "/mod/register");

    return {
        health: health
    };
}