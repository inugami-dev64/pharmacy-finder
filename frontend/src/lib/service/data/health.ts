import type { FetchFunc } from "./fetch";

export class PharmacyScraperResult {
    chain?: string;
    lastSuccessfulScrapeTimestamp?: number;
    timestamp?: number;
    success?: boolean;
}

export class HealthCheckResult {
    dbUp?: boolean;
    initialized?: boolean;
    lastScrapeResults?: Array<PharmacyScraperResult>

    /**
     * Retrieve a health check result from backend
     *
     * @returns a promise to HealthCheckResult
     */
    public static async readHealthCheck(fetch: FetchFunc): Promise<HealthCheckResult|undefined> {
        return await fetch(`/api/v1/health`)
            .then(async res => {
                if (res.status != 200) {
                    console.error((await res.bytes()).toString());
                    throw new Error(`Failed to fetch health check report`);
                }

                let data: HealthCheckResult = await res.json()
                return data
            })
            .catch(e => {
                console.error(e);
                return undefined;
            })
    }
}