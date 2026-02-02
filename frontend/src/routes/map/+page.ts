import type { PageLoad } from "./$types";
import { PharmacyInfo } from "$lib/service/pharmacy-info";

export const load: PageLoad = async ({ params }) => {
    return {
        pharmacies: await PharmacyInfo.readPharmacies()
    }
}