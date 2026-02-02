import type { PageLoad } from "./$types";
import { PharmacyInfo } from "$lib/service/data/pharmacy-info";

export const load: PageLoad = async ({ fetch }) => {
    return {
        pharmacies: await PharmacyInfo.readPharmacies(fetch)
    }
}