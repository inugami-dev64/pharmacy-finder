import { HealthCheckResult } from "$lib/service/data/health";
import { redirect } from "@sveltejs/kit";
import type { PageLoad } from "../$types";

export const load: PageLoad = async ({ fetch }) => {
    const health = await HealthCheckResult.readHealthCheck(fetch);
    if (health?.initialized)
        redirect(307, "/mod");

    return {
        health: health
    };
}