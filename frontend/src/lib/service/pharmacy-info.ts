import type { HttpError } from "$lib/http-error";

export class PharmacyInfo {
    id: number | undefined;
    chain: string | undefined;
    name: string | undefined;
    address: string | undefined;
    city: string | undefined
    county: string | undefined
    postalCode: number | undefined;
    phoneNumber: string | undefined;
    email: string | undefined;
    lat: number | undefined;
    lng: number | undefined;
}

const ESTONIA_BOUNDS = [[57.74100835592354, 21.73728235397422], [59.694266694887354, 28.30161313994147]]

export async function getPharmacies(): Promise<PharmacyInfo[]> {
    return await fetch(`/api/v1/pharmacies?sw=${ESTONIA_BOUNDS[0][0]},${ESTONIA_BOUNDS[0][1]}&ne=${ESTONIA_BOUNDS[1][0]},${ESTONIA_BOUNDS[1][1]}`)
        .then(async res => {
            if (res.status != 200) {
                let err: HttpError = await res.json();
                console.log(err);
                throw new Error(`Failed to fetch pharmacies: ${err.msg}`);
            }

            let data: PharmacyInfo[] = await res.json()
            return data
        })
        .catch(e => {
            console.log(e);
            return [];
        }) /* TODO: Implement a better error signaling system so that the user can see errors */
}