import type { PageLoad } from "./$types";
import { getPharmacies } from "$lib/service/pharmacy-info";

export const load: PageLoad = async ({ params }) => {
    return {
        pharmacies: await getPharmacies()
    }
}